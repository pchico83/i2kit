package aws

import (
	"fmt"

	"github.com/pchico83/i2kit/cli/providers/aws/cf"
	"github.com/pchico83/i2kit/cli/schemas/service"

	"github.com/aws/aws-sdk-go/aws"
	log "github.com/sirupsen/logrus"
)

//Destroy destroys a AWS Cloud Formation stack
func Destroy(s *service.Service, space string, config *aws.Config) error {
	stackName := s.Name
	if space != "" {
		stackName = fmt.Sprintf("%s-%s", s.Name, space)
	}
	log.Infof("Destroying the stack '%s'...", stackName)
	stack, err := cf.Get(stackName, config)
	if stack == nil {
		log.Infof("Stack '%s' doesn't exist.", stackName)
		return nil
	}
	consumed := cf.NumEvents(stackName, config)
	if err = cf.Delete(stackName, config); err != nil {
		return err
	}
	return cf.Watch(*stack.StackId, consumed, config)
}
