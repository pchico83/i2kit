package cmd

import (
	"github.com/pchico83/i2kit/cli/aws"
	"github.com/pchico83/i2kit/cli/service"
	"github.com/spf13/cobra"
)

//Deploy deploys an i2kit service
func Deploy() *cobra.Command {
	var file string
	var space string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an i2kit service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateEnvironment(); err != nil {
				return err
			}
			if err := initManifest(file); err != nil {
				return err
			}
			s, err := service.Read(file)
			if err != nil {
				return err
			}
			config, err := getAWSConfig()
			if err != nil {
				return err
			}
			return aws.Deploy(s, space, config, dryRun)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "i2kit.yml", "Yml file to be deployed")
	cmd.Flags().StringVarP(&space, "space", "s", "", "subdomains for dns search configuration")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "True to fake i2kit deployments")
	return cmd
}
