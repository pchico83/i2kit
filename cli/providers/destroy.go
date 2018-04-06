package providers

import (
	logger "log"
	"strings"

	"github.com/pchico83/i2kit/cli/providers/aws"
	"github.com/pchico83/i2kit/cli/providers/aws/cf"
	"github.com/pchico83/i2kit/cli/providers/aws/route53"
	"github.com/pchico83/i2kit/cli/providers/k8"
	K8Service "github.com/pchico83/i2kit/cli/providers/k8/service"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

//Destroy destroys a given service in a given environment
func Destroy(s *service.Service, e *environment.Environment, log *logger.Logger) error {
	if e.Provider == nil {
		log.Printf("Service '%s' dry-run destroy was successful", s.Name)
		return nil
	}

	log.Printf("Destroying the service '%s'...", s.GetFullName(e, "-"))
	if err := s.Validate(); err != nil {
		return err
	}
	if err := e.Validate(); err != nil {
		return err
	}

	if e.Provider.HostedZone == "" && len(s.GetPorts()) > 0 {
		var target string
		var err error
		switch e.Provider.GetType() {
		case environment.AWS:
			stackName := s.GetFullName(e, "-")
			target, err = cf.GetOutput(stackName, "elbURL", e.Provider.GetConfig())
		case environment.K8:
			target, err = K8Service.GetEndpoint(s, e)
			if err != nil && strings.Contains(err.Error(), "not found") {
				err = nil
			}
		}
		if err != nil {
			return err
		}
		if target != "" {
			route53.Destroy(s, e, target)
		}
	}
	switch e.Provider.GetType() {
	case environment.AWS:
		return aws.Destroy(s, e, log)
	case environment.K8:
		return k8.Destroy(s, e, log)
	}
	log.Print("Done!")
	return nil
}
