package k8

import (
	logger "log"
	"os"

	k8Deployment "github.com/pchico83/i2kit/cli/providers/k8/deployment"
	k8service "github.com/pchico83/i2kit/cli/providers/k8/service"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

//Destroy destroys a k8 deployment
func Destroy(s *service.Service, e *environment.Environment, log *logger.Logger) error {
	log.Printf("Destroying the k8 deployment '%s'...", s.GetFullName(e, "-"))
	c, tmpfile, err := e.Provider.GetConfigFile()
	if tmpfile != "" {
		defer os.Remove(tmpfile)
	}
	if err != nil {
		return err
	}

	err = k8Deployment.Destroy(s, e, c, log)
	if err != nil {
		return err
	}
	if len(s.GetPorts()) > 0 && s.Public {
		err = k8service.Destroy(s, e, c, log)
		if err != nil {
			return err
		}
	}
	return nil
}
