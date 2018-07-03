package aws

import (
	"fmt"
	logger "log"
	"strings"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pchico83/i2kit/cli/providers/aws/cf"
	"github.com/pchico83/i2kit/cli/providers/aws/cf/templates"
	"github.com/pchico83/i2kit/cli/providers/aws/ec2"
	"github.com/pchico83/i2kit/cli/providers/aws/elb"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

const (
	createFailed         = "CREATE_FAILED"
	deleteFailed         = "DELETE_FAILED"
	rollbackComplete     = "ROLLBACK_COMPLETE"
	rollbackFailed       = "ROLLBACK_FAILED"
	updateRollbackFailed = "UPDATE_ROLLBACK_FAILED"
)

//Deploy deploys a AWS Cloud Formation stack
func Deploy(s *service.Service, e *environment.Environment, log *logger.Logger) error {
	if err := deployVPC(e, log); err != nil {
		return err
	}
	return deployService(s, e, log)
}

func deployVPC(e *environment.Environment, log *logger.Logger) error {
	consumed := 0
	stackName := e.Name
	log.Printf("Configuring environment '%s' infrastructure...", stackName)
	config := e.Provider.GetConfig()
	stack, err := cf.Get(stackName, config)
	if err != nil {
		return err
	}
	if mustDelete(stack) {
		log.Printf("Destroying previous environment '%s' in '%s' state...", stackName, *stack.StackStatus)
		if err = cf.Delete(stackName, config); err != nil {
			return err
		}
		startTime := new(int64)
		*startTime = -1
		if err = cf.Watch(*stack.StackId, consumed, nil, e, startTime, config, log); err != nil {
			return err
		}
		stack = nil
	}
	if stack != nil && *stack.StackStatus == updateRollbackFailed {
		log.Printf("Environment %s is in 'UPDATE_ROLLBACK_FAILED' state. Fix the stack in the AWS console or destroy this project to unlock it.", stackName)
		return fmt.Errorf("'UPDATE_ROLLBACK_FAILED' state")
	}

	if stack != nil && strings.HasSuffix(*stack.StackStatus, "IN_PROGRESS") {
		startTime := new(int64)
		*startTime = -1
		if err = cf.Watch(*stack.StackId, consumed, nil, e, startTime, config, log); err != nil {
			return err
		}
	}

	var stackID string
	if stack != nil {
		stackID = *stack.StackId
	}
	template, err := templates.VPC(e, config)
	if err != nil {
		return err
	}
	if err := deployTemplate(stackName, stackID, template, nil, e, log); err != nil {
		return err
	}
	if stack != nil {
		log.Printf("Environment '%s' infrastructure was already configured.", stackName)
	} else {
		log.Printf("Environment '%s' infrastructure is now configured.", stackName)
	}

	vpc, _ := cf.GetOutput(stackName, "VPC", config)
	e.Provider.VPC = vpc
	subnets, _ := cf.GetOutput(stackName, "Subnets", config)
	for _, subnet := range strings.Split(subnets, ",") {
		e.Provider.Subnets = append(e.Provider.Subnets, &subnet)
	}
	securtyGroup, _ := cf.GetOutput(stackName, "SecurityGroup", config)
	e.Provider.SecurityGroup = securtyGroup

	if e.Provider.Keypair == "" {
		if err := ec2.CreateKeypair(e, config); err != nil {
			return err
		}
		e.Provider.Keypair = fmt.Sprintf("i2kit-%s", e.Name)
	}

	return nil
}

func deployService(s *service.Service, e *environment.Environment, log *logger.Logger) error {
	stackName := s.GetFullName(e, "-")
	config := e.Provider.GetConfig()
	stack, err := cf.Get(stackName, config)
	if err != nil {
		return err
	}
	if mustDelete(stack) {
		if err = Destroy(s, e, log); err != nil {
			return err
		}
		log.Printf("Destroying previous stack '%s' in '%s' state...", stackName, *stack.StackStatus)
		stack = nil
	}
	if stack != nil && *stack.StackStatus == updateRollbackFailed {
		log.Printf("Stack %s is in 'UPDATE_ROLLBACK_FAILED' state. Fix the stack in the AWS console or destroy this service to unlock it.", stackName)
		return fmt.Errorf("'UPDATE_ROLLBACK_FAILED' state")
	}

	var stackID string
	if stack != nil {
		stackID = *stack.StackId
	}
	template, err := templates.Service(s, e, config)
	if err != nil {
		return err
	}
	if err := deployTemplate(stackName, stackID, template, s, e, log); err != nil {
		return err
	}

	elbName, _ := cf.GetOutput(stackName, "elbName", config)
	if elbName != "" {
		if err := elb.Wait(s, elbName, config, log); err != nil {
			return err
		}
	}
	return nil
}

func deployTemplate(stackName, stackID, template string, s *service.Service, e *environment.Environment, log *logger.Logger) error {
	consumed := 0
	config := e.Provider.GetConfig()
	if stackID == "" {
		log.Printf("Creating stack '%s'...", stackName)
		var err error
		if stackID, err = cf.Create(stackName, template, config); err != nil {
			return err
		}
	} else {
		consumed = cf.NumEvents(stackID, config)
		log.Printf("Updating the stack '%s'...", stackName)
		var updated bool
		updated, err := cf.Update(stackID, template, config)
		if err != nil {
			return err
		}
		if !updated {
			log.Printf("No updates are to be performed.")
		}
	}
	stack, err := cf.Get(stackID, config)
	if err != nil {
		return err
	}
	startTime := new(int64)
	*startTime = 0
	if stack.LastUpdatedTime != nil {
		*startTime = stack.LastUpdatedTime.Unix() * 1000
	}
	return cf.Watch(stackID, consumed, s, e, startTime, config, log)
}

func mustDelete(stack *cloudformation.Stack) bool {
	if stack == nil {
		return false
	}
	switch *stack.StackStatus {
	case createFailed:
		return true
	case deleteFailed:
		return true
	case rollbackComplete:
		return true
	case rollbackFailed:
		return true
	default:
		return false
	}
}
