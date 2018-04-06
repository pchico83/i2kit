package cloudwatch

import (
	"fmt"
	logger "log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

//RetrieveLogs prints the logs asociatted to a given log stream
func RetrieveLogs(s *service.Service, e *environment.Environment, startTime *int64, config *aws.Config, log *logger.Logger) error {
	svc := cloudwatchlogs.New(session.New(), config)
	input := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName: aws.String(fmt.Sprintf("i2kit-%s", s.GetFullName(e, "-"))),
		StartTime:    startTime,
	}
	logEvents, err := svc.FilterLogEvents(input)
	if err != nil {
		return nil
	}
	for _, event := range logEvents.Events {
		log.Printf("%s: %q", *event.LogStreamName, *event.Message)
		if *event.Timestamp >= *startTime {
			*startTime = *event.Timestamp + 1
		}
	}
	return nil
}
