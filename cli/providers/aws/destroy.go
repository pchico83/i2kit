package aws

import (
	"github.com/pchico83/i2kit/cli/providers/aws/cf"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"

	log "github.com/sirupsen/logrus"
)

//Destroy destroys a AWS Cloud Formation stack
func Destroy(s *service.Service, e *environment.Environment) error {
	config, err := getAWSConfig(e)
	if err != nil {
		return err
	}
	log.Infof("Destroying the stack '%s'...", s.Name)
	stack, err := cf.Get(s.Name, config)
	if stack == nil {
		log.Infof("Stack '%s' doesn't exist.", s.Name)
		return nil
	}
	consumed := cf.NumEvents(s.Name, config)
	if err = cf.Delete(s.Name, config); err != nil {
		return err
	}
	return cf.Watch(*stack.StackId, consumed, config)
}
