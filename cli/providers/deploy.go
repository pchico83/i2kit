package providers

import (
	logger "log"

	"github.com/pchico83/i2kit/cli/providers/aws"
	"github.com/pchico83/i2kit/cli/providers/aws/cf"
	"github.com/pchico83/i2kit/cli/providers/aws/route53"
	"github.com/pchico83/i2kit/cli/providers/k8"
	K8Service "github.com/pchico83/i2kit/cli/providers/k8/service"

	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

//Deploy deploys a given service in a given environment
func Deploy(s *service.Service, e *environment.Environment, log *logger.Logger) error {
	if e.Provider == nil {
		log.Printf("Service '%s' dry-run deployment was successful", s.Name)
		return nil
	}

	log.Printf("Deploying the service '%s'...", s.GetFullName(e, "-"))
	if err := s.Validate(); err != nil {
		return err
	}
	if err := e.Validate(); err != nil {
		return err
	}
	switch e.Provider.GetType() {
	case environment.AWS:
		if err := aws.Deploy(s, e, log); err != nil {
			return err
		}
	case environment.K8:
		if err := k8.Deploy(s, e, log); err != nil {
			return err
		}
	}
	if e.Provider.HostedZone != "" || len(s.GetPorts()) == 0 {
		return nil
	}
	var target string
	var err error
	switch e.Provider.GetType() {
	case environment.AWS:
		stackName := s.GetFullName(e, "-")
		target, err = cf.GetOutput(stackName, "elbURL", e.Provider.GetConfig())
	case environment.K8:
		target, err = K8Service.GetEndpoint(s, e)
	}
	if err != nil {
		return err
	}
	if target == "" {
		return nil
	}
	if err := route53.Create(s, e, target); err != nil {
		return err
	}
	log.Print("Done!")
	return nil

}
