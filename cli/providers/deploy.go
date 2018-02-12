package providers

import (
	"fmt"

	"github.com/pchico83/i2kit/cli/providers/aws"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

//Deploy deploys a given service in a given environment
func Deploy(s *service.Service, e *environment.Environment) error {
	if e.Provider == nil {
		fmt.Printf("Service '%s' dry-run deployment was successful\n", s.Name)
		return nil
	}
	return aws.Deploy(s, e)
}
