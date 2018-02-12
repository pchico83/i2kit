package cmd

import (
	"io/ioutil"

	"github.com/pchico83/i2kit/cli/providers"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

//Deploy deploys an i2kit service
func Deploy() *cobra.Command {
	var servicePath string
	var environmentPath string
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an i2kit service",
		RunE: func(cmd *cobra.Command, args []string) error {
			readBytes, err := ioutil.ReadFile(servicePath)
			if err != nil {
				return err
			}
			var s service.Service
			err = yaml.Unmarshal(readBytes, &s)
			if err != nil {
				return err
			}
			if err = s.Validate(); err != nil {
				return err
			}

			readBytes, err = ioutil.ReadFile(environmentPath)
			if err != nil {
				return err
			}
			var e environment.Environment
			err = yaml.Unmarshal(readBytes, &e)
			if err != nil {
				return err
			}
			if err = e.Validate(); err != nil {
				return err
			}

			return providers.Deploy(&s, &e)
		},
	}
	cmd.Flags().StringVarP(&servicePath, "service", "s", "service.yml", "Service yml file to be deployed")
	cmd.Flags().StringVarP(&environmentPath, "environment", "e", "environment.yml", "Environment yml file used for deployment")
	return cmd
}
