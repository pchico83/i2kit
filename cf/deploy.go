package cf

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pchico83/i2kit/k8"
	"github.com/pchico83/i2kit/linuxkit"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

//NewDeploy deploys a k8 object
func NewDeploy(k8path string, awsConfig *aws.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a k8 object",
		RunE: func(cmd *cobra.Command, args []string) error {
			deployment, err := k8.Read(k8path)
			if err != nil {
				return err
			}
			mobyTemplate, err := linuxkit.GetTemplate(deployment)
			if err != nil {
				return err
			}
			// TODO: for debugging purposes only, will be implemented when added a "verbose" flag
			d, err := yaml.Marshal(&mobyTemplate)
			if err != nil {
				return fmt.Errorf("error marshalling: %v", err)
			}
			fmt.Printf("--- moby template:\n%s\n", string(d))
			// END TODO
			ami, err := linuxkit.Export(mobyTemplate, deployment.GetName())
			if err != nil {
				return err
			}
			cfTemplate, err := Translate(deployment, ami, false)
			if err != nil {
				return err
			}
			cfTemplateString := string(cfTemplate)
			svc := cloudformation.New(session.New(), awsConfig)
			deploymentName := deployment.GetName()
			inStack := &cloudformation.CreateStackInput{
				Capabilities: []*string{aws.String("CAPABILITY_NAMED_IAM")},
				StackName:    &deploymentName,
				TemplateBody: &cfTemplateString,
				Tags: []*cloudformation.Tag{
					&cloudformation.Tag{
						Key:   aws.String("i2kit"),
						Value: aws.String("alpha"),
					},
				},
			}
			stack, err := svc.CreateStack(inStack)
			if err != nil {
				return err
			}
			return Watch(svc, stack.StackId, cloudformation.ResourceStatusCreateInProgress, 0)
		},
	}
	return cmd
}
