package cf

import (
	"fmt"
	"strings"
	"time"

	logger "log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pchico83/i2kit/cli/providers/aws/cloudwatch"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

//Watch waits for a AWS Cloud Formation stack state
func Watch(name string, consumed int, s *service.Service, startTime *int64, config *aws.Config, log *logger.Logger) error {
	svc := cloudformation.New(session.New(), config)
	errors := 0
	endsWithInProgress := true
	for endsWithInProgress {
		time.Sleep(5 * time.Second)
		response, err := svc.DescribeStacks(
			&cloudformation.DescribeStacksInput{
				StackName: aws.String(name),
			},
		)
		if err != nil {
			errors++
			if errors >= 3 {
				return err
			}
			continue
		}
		errors = 0
		events, err := svc.DescribeStackEvents(
			&cloudformation.DescribeStackEventsInput{
				StackName: aws.String(name),
			},
		)
		if err == nil {
			for index := len(events.StackEvents) - consumed - 1; index >= 0; index-- {
				if events.StackEvents[index].ResourceStatusReason != nil {
					log.Printf("%s (%s) %s %s",
						*events.StackEvents[index].LogicalResourceId,
						*events.StackEvents[index].ResourceType,
						*events.StackEvents[index].ResourceStatus,
						*events.StackEvents[index].ResourceStatusReason,
					)
				} else {
					log.Printf("%s (%s)  %s",
						*events.StackEvents[index].LogicalResourceId,
						*events.StackEvents[index].ResourceType,
						*events.StackEvents[index].ResourceStatus,
					)
				}
			}
			consumed = len(events.StackEvents)
		}
		if *startTime != -1 {
			if err := cloudwatch.RetrieveLogs(s, startTime, config, log); err != nil {
				return err
			}
		}
		endsWithInProgress = strings.HasSuffix(*response.Stacks[0].StackStatus, "IN_PROGRESS")
		if !endsWithInProgress && (strings.HasPrefix(*response.Stacks[0].StackStatus, "ROLLBACK") ||
			!strings.HasSuffix(*response.Stacks[0].StackStatus, "COMPLETE")) {
			return fmt.Errorf("There was an error and the stack was rollbacked")
		}
	}
	return nil
}
