package cmd

import (
	"fmt"
	"os"

	"github.com/pchico83/i2kit/cli/providers/aws"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
	"github.com/spf13/cobra"
)

//Deploy deploys an i2kit service
func Deploy() *cobra.Command {
	var servicePath string
	var environmentPath string
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an i2kit service",
		RunE: func(cmd *cobra.Command, args []string) error {
			reader, err := os.Open(servicePath)
			if err != nil {
				return err
			}
			if err = service.Validate(reader); err != nil {
				return err
			}
			reader, err = os.Open(servicePath)
			if err != nil {
				return err
			}
			s, err := service.Read(reader)
			if err != nil {
				return err
			}

			reader, err = os.Open(environmentPath)
			if err != nil {
				return err
			}
			if err = environment.Validate(reader); err != nil {
				return err
			}
			reader, err = os.Open(environmentPath)
			if err != nil {
				return err
			}
			e, err := environment.Read(reader)
			if err != nil {
				return err
			}

			if e.Provider == nil {
				fmt.Println("Service dry-run deployed succesfully. Did you define an 'environment.yml' file?")
				return nil
			}
			return aws.Deploy(s, e)
		},
	}
	cmd.Flags().StringVarP(&servicePath, "service", "s", "service.yml", "Service yml file to be deployed")
	cmd.Flags().StringVarP(&environmentPath, "environment", "e", "environment.yml", "Environment yml file used for deployment")
	return cmd
}
