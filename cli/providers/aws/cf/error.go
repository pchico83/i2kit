package cf

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
)

//CheckError checks the message of a AWS Cloud Formation error
func CheckError(err error, code, message string) bool {
	if awsErr, ok := err.(awserr.Error); ok {
		if code != "" && awsErr.Code() == code {
			return true
		}
		if message != "" && awsErr.Message() == message {
			return true
		}
	}
	return false
}
