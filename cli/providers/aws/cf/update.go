package cf

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

//Update updates a AWS Cloud Formation stack
func Update(stackID, template string, config *aws.Config) (bool, error) {
	stack := &cloudformation.UpdateStackInput{
		Capabilities: []*string{aws.String("CAPABILITY_IAM")},
		StackName:    aws.String(stackID),
		TemplateBody: aws.String(template),
	}
	svc := cloudformation.New(session.New(), config)
	_, err := svc.UpdateStack(stack)
	if err != nil {
		if CheckError(err, "", "No updates are to be performed.") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
