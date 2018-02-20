package cf

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

//Create creates a AWS Cloud Formation stack
func Create(name, template string, config *aws.Config) (string, error) {
	stack := &cloudformation.CreateStackInput{
		Capabilities: []*string{aws.String("CAPABILITY_IAM")},
		StackName:    aws.String(name),
		TemplateBody: aws.String(template),
		Tags: []*cloudformation.Tag{
			&cloudformation.Tag{
				Key:   aws.String("i2kit"),
				Value: aws.String("alpha"),
			},
			&cloudformation.Tag{
				Key:   aws.String("Name"),
				Value: aws.String(name),
			},
		},
	}
	svc := cloudformation.New(session.New(), config)
	response, err := svc.CreateStack(stack)
	if err != nil {
		return "", err
	}
	return *response.StackId, nil
}
