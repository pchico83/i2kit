package cf

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pchico83/i2kit/k8"
	"github.com/pchico83/i2kit/linuxkit"
	"github.com/spf13/cobra"
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
			linuxkitPath, err := linuxkit.GetTemplate(deployment)
			if err != nil {
				return err
			}
			defer os.Remove(linuxkitPath)
			ami, err := linuxkit.Export(linuxkitPath)
			fmt.Println(ami) //alberto: remove this line
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}
			svc := cloudformation.New(session.New(), awsConfig)
			bytes, err := ioutil.ReadFile("./cf_samples/elb-asg.json")
			if err != nil {
				return err
			}
			cfTemplate := string(bytes)
			inStack := &cloudformation.CreateStackInput{
				Capabilities: []*string{aws.String("CAPABILITY_NAMED_IAM")},
				StackName:    &deployment.Metadata.Name,
				TemplateBody: &cfTemplate,
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
