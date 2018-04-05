package providers

import (
	logger "log"

	"github.com/pchico83/i2kit/cli/providers/aws"
	"github.com/pchico83/i2kit/cli/providers/aws/route53"

	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

//Deploy deploys a given service in a given environment
func Deploy(s *service.Service, e *environment.Environment, log *logger.Logger) error {
	if e.Provider == nil {
		log.Printf("Service '%s' dry-run deployment was successful", s.Name)
		return nil
	}

	if err := s.Validate(); err != nil {
		return err
	}
	if err := e.Validate(); err != nil {
		return err
	}
	if err := aws.Deploy(s, e, log); err != nil {
		return err
	}
	if e.Provider.HostedZone != "" {
		return nil
	}
	return route53.Create(s, e)
}
