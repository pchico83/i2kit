package cf

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pchico83/i2kit/k8"
	"github.com/spf13/cobra"
)

//NewDestroy destroys a i2kit application
func NewDestroy(k8path string, awsConfig *aws.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy a k8 object",
		RunE: func(cmd *cobra.Command, args []string) error {
			deployment, err := k8.Read(k8path)
			if err != nil {
				return err
			}
			svc := cloudformation.New(session.New(), awsConfig)
			response, err := svc.DescribeStacks(
				&cloudformation.DescribeStacksInput{
					StackName: aws.String(deployment.Metadata.Name),
				},
			)
			if err != nil {
				return err
			}
			if len(response.Stacks) == 0 {
				fmt.Printf("Stack '%s' doesn't exist.\n", deployment.Metadata.Name)
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
			return Watch(svc, stackID, cloudformation.ResourceStatusDeleteInProgress, consumed)
		},
	}
	return cmd
}
