package cmd

import (
	"os"

	"github.com/pchico83/i2kit/cli/providers/aws"
	"github.com/pchico83/i2kit/cli/schemas/service"
	"github.com/spf13/cobra"
)

//Destroy destroys an i2kit service
func Destroy() *cobra.Command {
	var file string
	var space string
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy an i2kit service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if dryRun {
				return nil
			}
			if err := initManifest(file); err != nil {
				return err
			}

			reader, err := os.Open(file)
			if err != nil {
				return err
			}

			s, err := service.Read(reader)
			if err != nil {
				return err
			}
			config, err := getAWSConfig()
			if err != nil {
				return err
			}
			return aws.Destroy(s, space, config)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "service.yml", "Service yml file to be deployed")
	cmd.Flags().StringVarP(&space, "space", "s", "", "subdomains for dns search configuration")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "True to fake i2kit deployments")
	return cmd
}
