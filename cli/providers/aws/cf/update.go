package cf

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	log "github.com/sirupsen/logrus"
)

//Update updates a AWS Cloud Formation stack
func Update(stackID, template string, config *aws.Config) error {
	log.Infof("Updating the stack '%s'...", stackID)
	stack := &cloudformation.UpdateStackInput{
		Capabilities: []*string{aws.String("CAPABILITY_IAM")},
		StackName:    aws.String(stackID),
		TemplateBody: aws.String(template),
	}
	svc := cloudformation.New(session.New(), config)
	_, err := svc.UpdateStack(stack)
	if err != nil {
		if CheckError(err, "", "No updates are to be performed.") {
			log.Info("No updates are to be performed.")
			return nil
		}
		return err
	}
	return nil
}
