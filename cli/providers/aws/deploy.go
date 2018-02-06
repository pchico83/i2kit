package aws

import (
	"github.com/pchico83/i2kit/cli/providers/aws/cf"
	"github.com/pchico83/i2kit/cli/providers/aws/elb"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"

	log "github.com/sirupsen/logrus"
)

//Deploy deploys a AWS Cloud Formation stack
func Deploy(s *service.Service, e *environment.Environment) error {
	consumed := 0
	config, err := getAWSConfig(e)
	if err != nil {
		return err
	}
	stack, err := cf.Get(s.Name, config)
	if err != nil {
		return err
	}
	if stack != nil && *stack.StackStatus == "ROLLBACK_COMPLETE" {
		if err = Destroy(s, e); err != nil {
			return err
		}
		log.Infof("Destroying previous stack '%s' in 'ROLLBACK_COMPLETE' state...", s.Name)
		stack = nil
	}
	var stackID string
	var template string
	if stack == nil {
		template, err = Translate(s, e, "NO")
		if err != nil {
			return err
		}
		if stackID, err = cf.Create(s.Name, template, config); err != nil {
			return err
		}
	} else {
		template, err = Translate(s, e, "YES")
		if err != nil {
			return err
		}
		stackID = *stack.StackId
		consumed = cf.NumEvents(stackID, config)
		if err = cf.Update(stackID, template, config); err != nil {
			return err
		}
	}
	if err = cf.Watch(stackID, consumed, config); err != nil {
		return err
	}
	stack, err = cf.Get(stackID, config)
	if err != nil {
		return err
	}
	for _, o := range stack.Outputs {
		if *o.OutputKey == "elbName" {
			return elb.Wait(*o.OutputValue, config)
		}
	}
	return nil
}
