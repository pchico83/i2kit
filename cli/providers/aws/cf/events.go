package cf

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

//NumEvents get the current number of events of a AWS Cloud Formation stack
func NumEvents(name string, config *aws.Config) int {
	svc := cloudformation.New(session.New(), config)
	response, err := svc.DescribeStackEvents(
		&cloudformation.DescribeStackEventsInput{
			StackName: aws.String(name),
		},
	)
	if err == nil {
		return len(response.StackEvents)
	}
	return 0
}
