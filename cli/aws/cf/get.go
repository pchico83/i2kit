package cf

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

//Get gets a AWS Cloud Formation stack
func Get(name string, config *aws.Config) (*cloudformation.Stack, error) {
	svc := cloudformation.New(session.New(), config)
	response, err := svc.DescribeStacks(
		&cloudformation.DescribeStacksInput{
			StackName: aws.String(name),
		},
	)
	if err != nil {
		if CheckError(err, "ValidationError", "") {
			return nil, nil
		}
		return nil, err
	}
	if response == nil || len((*response).Stacks) == 0 {
		return nil, nil
	}
	return response.Stacks[0], nil
}
