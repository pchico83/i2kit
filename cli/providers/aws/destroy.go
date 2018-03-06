package aws

import (
	logger "log"

	"github.com/pchico83/i2kit/cli/providers/aws/cf"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

//Destroy destroys a AWS Cloud Formation stack
func Destroy(s *service.Service, e *environment.Environment, log *logger.Logger) error {
	config, err := getAWSConfig(e)
	if err != nil {
		return err
	}
	log.Printf("Destroying the stack '%s'...", s.Name)
	stack, err := cf.Get(s.Name, config)
	if stack == nil {
		log.Printf("Stack '%s' doesn't exist.", s.Name)
		return nil
	}
	consumed := cf.NumEvents(s.Name, config)
	if err = cf.Delete(s.Name, config); err != nil {
		return err
	}
	startTime := new(int64)
	*startTime = -1
	return cf.Watch(*stack.StackId, consumed, s, startTime, config, log)
}
