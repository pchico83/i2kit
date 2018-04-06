package aws

import (
	logger "log"

	"github.com/pchico83/i2kit/cli/providers/aws/cf"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

//Destroy destroys a AWS Cloud Formation stack
func Destroy(s *service.Service, e *environment.Environment, log *logger.Logger) error {
	stackName := s.GetFullName(e, "-")
	log.Printf("Destroying the stack '%s'...", stackName)
	config := e.Provider.GetConfig()
	stack, err := cf.Get(stackName, config)
	if err != nil {
		return err
	}
	if stack == nil {
		log.Printf("Stack '%s' doesn't exist.", stackName)
		return nil
	}
	consumed := cf.NumEvents(stackName, config)
	if err = cf.Delete(stackName, config); err != nil {
		return err
	}
	startTime := new(int64)
	*startTime = -1
	return cf.Watch(*stack.StackId, consumed, s, e, startTime, config, log)
}
