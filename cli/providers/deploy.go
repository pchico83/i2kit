package providers

import (
	logger "log"

	"github.com/pchico83/i2kit/cli/providers/aws"

	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

//Deploy deploys a given service in a given environment
func Deploy(s *service.Service, e *environment.Environment, log *logger.Logger) error {
	if e.Provider == nil {
		log.Printf("Service '%s' dry-run deployment was successful", s.Name)
		return nil
	}
	return aws.Deploy(s, e, log)
}
