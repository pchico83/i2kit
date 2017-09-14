package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/pchico83/i2kit/cf"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "i2kit COMMAND [ARG...]",
		Short: "Manage i2kit applications",
	}
	awsCredentials := os.Getenv("AWS_CREDENTIALS")
	if awsCredentials == "" {
		awsCredentials = "/run/secrets/aws-credentials"
	}
	if _, err := os.Stat(awsCredentials); err != nil {
		err := fmt.Errorf("Variable 'AWS_CREDENTIALS' must point to an existing file")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	awsConfig := &aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: credentials.NewSharedCredentials(awsCredentials, "default"),
	}
	name := "test"
	i2kitPath := "./example.yml"
	cmd.AddCommand(
		cf.NewDeploy(name, i2kitPath, awsConfig),
		cf.NewDestroy(name, awsConfig),
	)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
