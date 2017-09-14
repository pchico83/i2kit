package deploy

import (
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pchico83/i2kit/cf"
	"github.com/spf13/cobra"
)

//NewDeploy deploys a i2kit application
func NewDeploy(name, i2kitPath string, awsConfig *aws.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a i2kit application",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc := cloudformation.New(session.New(), awsConfig)
			bytes, err := ioutil.ReadFile("./cf_samples/elb-asg.json")
			if err != nil {
				return err
			}
			cfTemplate := string(bytes)
			inStack := &cloudformation.CreateStackInput{
				Capabilities: []*string{aws.String("CAPABILITY_NAMED_IAM")},
				StackName:    &name,
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
			return cf.Watch(svc, stack.StackId, cloudformation.ResourceStatusCreateInProgress, 0)
		},
	}
	return cmd
}
