package cf

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	log "github.com/sirupsen/logrus"
)

//Watch waits for a AWS Cloud Formation stack state
func Watch(name string, consumed int, config *aws.Config) error {
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
					log.Infof("%s: %s (%s) %s %s",
						events.StackEvents[index].Timestamp.Local().Format("2006-01-02 15:04:05"),
						*events.StackEvents[index].LogicalResourceId,
						*events.StackEvents[index].ResourceType,
						*events.StackEvents[index].ResourceStatus,
						*events.StackEvents[index].ResourceStatusReason,
					)
				} else {
					log.Infof("%s: %s (%s)  %s",
						events.StackEvents[index].Timestamp.Local().Format("2006-01-02 15:04:05"),
						*events.StackEvents[index].LogicalResourceId,
						*events.StackEvents[index].ResourceType,
						*events.StackEvents[index].ResourceStatus,
					)
				}
			}
			consumed = len(events.StackEvents)
		}
		endsWithInProgress = strings.HasSuffix(*response.Stacks[0].StackStatus, "IN_PROGRESS")
		if !endsWithInProgress && strings.HasPrefix(*response.Stacks[0].StackStatus, "UPDATE_ROLLBACK") {
			log.Fatal("There was an error and the stack was rollbacked")
		}
	}
	return nil
}
