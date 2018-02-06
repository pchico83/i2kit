package aws

import (
	"fmt"

	"github.com/pchico83/i2kit/cli/providers/aws/cf"
	"github.com/pchico83/i2kit/cli/providers/aws/elb"
	"github.com/pchico83/i2kit/cli/schemas/service"

	"github.com/aws/aws-sdk-go/aws"
	log "github.com/sirupsen/logrus"
)

//Deploy deploys a AWS Cloud Formation stack
func Deploy(s *service.Service, space string, config *aws.Config, dryRun bool) error {
	stackName := s.Name
	if space != "" {
		stackName = fmt.Sprintf("%s-%s", s.Name, space)
	}
	consumed := 0
	stack, err := cf.Get(stackName, config)
	if err != nil {
		return err
	}
	if stack != nil && *stack.StackStatus == "ROLLBACK_COMPLETE" {
		if err = Destroy(s, space, config); err != nil {
			return err
		}
		log.Infof("Destroying previous stack '%s' in 'ROLLBACK_COMPLETE' state...", stackName)
		stack = nil
	}
	var stackID string
	var template string
	if stack == nil {
		template, err = Translate(s, space, "NO")
		if err != nil {
			return err
		}
		if dryRun {
			fmt.Println(template)
			return nil
		}
		if stackID, err = cf.Create(stackName, template, config); err != nil {
			return err
		}
	} else {
		template, err = Translate(s, space, "YES")
		if err != nil {
			return err
		}
		if dryRun {
			fmt.Println(template)
			return nil
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
