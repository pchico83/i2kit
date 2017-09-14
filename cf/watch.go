package cf

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation"
)

//Watch waits for a stack state
func Watch(svc *cloudformation.CloudFormation, stackID *string, watchedStatus string, consumed int) error {
	errors := 0
	status := watchedStatus
	for status == watchedStatus {
		time.Sleep(10 * time.Second)
		response, err := svc.DescribeStacks(
			&cloudformation.DescribeStacksInput{
				StackName: stackID,
			},
		)
		if err != nil {
			errors++
			fmt.Fprintln(os.Stderr, err)
			if errors >= 3 {
				return err
			}
			continue
		}
		errors = 0
		events, err := svc.DescribeStackEvents(
			&cloudformation.DescribeStackEventsInput{
				StackName: stackID,
			},
		)
		if err == nil {
			for index := len(events.StackEvents) - consumed - 1; index >= 0; index-- {
				if events.StackEvents[index].ResourceStatusReason != nil {
					fmt.Printf("%s: %s (%s) %s %s\n",
						events.StackEvents[index].Timestamp.Local().Format("2006-01-02 15:04:05"),
						*events.StackEvents[index].LogicalResourceId,
						*events.StackEvents[index].ResourceType,
						*events.StackEvents[index].ResourceStatus,
						*events.StackEvents[index].ResourceStatusReason,
					)
				} else {
					fmt.Printf("%s: %s (%s)  %s\n",
						events.StackEvents[index].Timestamp.Local().Format("2006-01-02 15:04:05"),
						*events.StackEvents[index].LogicalResourceId,
						*events.StackEvents[index].ResourceType,
						*events.StackEvents[index].ResourceStatus,
					)
				}
			}
			consumed = len(events.StackEvents)
		}
		status = *response.Stacks[0].StackStatus
	}
	return nil
}
