package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	log "github.com/sirupsen/logrus"
)

func initManifest(file string) error {
	manifest := os.Getenv("I2KIT_MANIFEST")
	if manifest != "" {
		decoded, err := base64.StdEncoding.DecodeString(manifest)
		if err != nil {
			return fmt.Errorf("Compose decode error: %s\n", err.Error())
		}
		if err = ioutil.WriteFile(file, decoded, 0644); err != nil {
			return fmt.Errorf("Compose write file error: %s\n", err.Error())
		}
	}
	return nil
}

func validateEnvironment() error {
	value := os.Getenv("I2KIT_HOSTED_ZONE")
	if value == "" {
		log.Infof("Variable 'I2KIT_HOSTED_ZONE' not defined (for example 'i2kit.com.'), CNAMEs won't be created")
	}
	value = os.Getenv("I2KIT_SECURITY_GROUP")
	if value == "" {
		return fmt.Errorf("Variable 'I2KIT_SECURITY_GROUP' must be defined\n")
	}
	value = os.Getenv("I2KIT_SUBNET")
	if value == "" {
		return fmt.Errorf("Variable 'I2KIT_SUBNET' must be defined\n")
	}
	value = os.Getenv("I2KIT_KEYPAIR")
	if value == "" {
		return fmt.Errorf("Variable 'I2KIT_KEYPAIR' must be defined\n")
	}
	return nil
}

func getAWSConfig() (*aws.Config, error) {
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if awsAccessKeyID == "" || awsSecretAccessKey == "" {
		return nil, fmt.Errorf("AWS credentials must be provided in 'AWS_ACCESS_KEY_ID' and 'AWS_SECRET_ACCESS_KEY'")
	}
	awsRegion := os.Getenv("I2KIT_REGION")
	awsConfig := &aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewEnvCredentials(),
	}
	return awsConfig, nil
}
