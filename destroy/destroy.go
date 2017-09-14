package destroy

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pchico83/i2kit/cf"
	"github.com/spf13/cobra"
)

//NewDestroy destroys a i2kit application
func NewDestroy(name string, awsConfig *aws.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy a i2kit application",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc := cloudformation.New(session.New(), awsConfig)
			response, err := svc.DescribeStacks(
				&cloudformation.DescribeStacksInput{
					StackName: aws.String(name),
				},
			)
			if err != nil {
				return err
			}
			if len(response.Stacks) == 0 {
				fmt.Printf("Stack '%s' doesn't exist.\n", name)
				return nil
			}
			stackID := response.Stacks[0].StackId
			events, err := svc.DescribeStackEvents(
				&cloudformation.DescribeStackEventsInput{
					StackName: stackID,
				},
			)
			consumed := 0
			if err == nil {
				consumed = len(events.StackEvents)
			}
			_, err = svc.DeleteStack(
				&cloudformation.DeleteStackInput{
					StackName: stackID,
				},
			)
			if err != nil {
				return err
			}
			return cf.Watch(svc, stackID, cloudformation.ResourceStatusDeleteInProgress, consumed)
		},
	}
	return cmd
}
