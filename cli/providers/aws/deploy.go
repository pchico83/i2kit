package aws

import (
	logger "log"
	"strings"

	"github.com/pchico83/i2kit/cli/providers/aws/cf"
	"github.com/pchico83/i2kit/cli/providers/aws/elb"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

//Deploy deploys a AWS Cloud Formation stack
func Deploy(s *service.Service, e *environment.Environment, log *logger.Logger) error {
	consumed := 0
	config := e.Provider.GetConfig()
	stack, err := cf.Get(s.Name, config)
	if err != nil {
		return err
	}
	if stack != nil && (*stack.StackStatus == "ROLLBACK_COMPLETE" || strings.HasSuffix(*stack.StackStatus, "_FAILED")) {
		if err = Destroy(s, e, log); err != nil {
			return err
		}
		log.Printf("Destroying previous stack '%s' in '%s' state...", s.Name, *stack.StackStatus)
		stack = nil
	}

	var stackID string
	var template string
	if stack == nil {
		template, err = Translate(s, e, config)
		if err != nil {
			return err
		}
		log.Printf("Creating stack '%s'...", s.Name)
		if stackID, err = cf.Create(s.Name, template, config); err != nil {
			return err
		}
	} else {
		template, err = Translate(s, e, config)
		if err != nil {
			return err
		}
		stackID = *stack.StackId
		consumed = cf.NumEvents(stackID, config)
		log.Printf("Updating the stack '%s'...", stackID)
		var updated bool
		updated, err = cf.Update(stackID, template, config)
		if err != nil {
			return err
		}
		if !updated {
			log.Printf("No updates are to be performed.")
		}
	}
	stack, err = cf.Get(stackID, config)
	if err != nil {
		return err
	}
	startTime := new(int64)
	*startTime = 0
	if stack.LastUpdatedTime != nil {
		*startTime = stack.LastUpdatedTime.Unix() * 1000
	}
	if err = cf.Watch(stackID, consumed, s, startTime, config, log); err != nil {
		return err
	}
	elbName, _ := cf.GetOutput(stackID, "elbName", config)
	if elbName != "" {
		return elb.Wait(s, elbName, config, log)
	}
	return nil
}
