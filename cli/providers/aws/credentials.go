package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/pchico83/i2kit/cli/schemas/environment"
)

func getAWSConfig(e *environment.Environment) (*aws.Config, error) {
	awsConfig := &aws.Config{
		Region:      aws.String(e.Provider.Region),
		Credentials: credentials.NewStaticCredentials(e.Provider.AccessKey, e.Provider.SecretKey, ""),
	}
	return awsConfig, nil
}
