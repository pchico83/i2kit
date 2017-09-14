package cf

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

//Delete deletes a AWS Cloud Formation stack
func Delete(name string, config *aws.Config) error {
	svc := cloudformation.New(session.New(), config)
	_, err := svc.DeleteStack(
		&cloudformation.DeleteStackInput{
			StackName: aws.String(name),
		},
	)
	return err
}
