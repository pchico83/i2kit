package cf

import (
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/spf13/cobra"
	"k8s.io/client-go/pkg/api"
)

//NewDeploy deploys a k8 object
func NewDeploy(name, k8Path string, awsConfig *aws.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a k8 object",
		RunE: func(cmd *cobra.Command, args []string) error {
			deploymentBytes, err := ioutil.ReadFile(k8Path)
			if err != nil {
				return err
			}
			decode := api.Codecs.UniversalDeserializer().Decode
			_, _, err = decode(deploymentBytes, nil, nil)
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
			return Watch(svc, stack.StackId, cloudformation.ResourceStatusCreateInProgress, 0)
		},
	}
	return cmd
}
