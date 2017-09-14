package elb

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	log "github.com/sirupsen/logrus"
)

//Wait waits for the instances in a ELB to be registered
func Wait(name string, config *aws.Config) error {
	log.Infof("Waiting for instances to be registered in the '%s' load balancer...", name)
	svc := elb.New(session.New(), config)
	dii := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: aws.String(name),
	}
	registered := 0
	current := 0
	outOfService := true
	limit := 60 * 20
	for outOfService && limit >= 0 {
		time.Sleep(5 * time.Second)
		limit -= 5
		dio, err := svc.DescribeInstanceHealth(dii)
		if err != nil {
			return err
		}
		current = 0
		outOfService = false
		for _, instance := range dio.InstanceStates {
			if *instance.State == "InService" {
				current++
			} else {
				outOfService = true
			}
		}
		if current != registered {
			registered = current
			log.Infof("%d instances registered in the '%s' load balancer", registered, name)
		}
	}
	if limit <= 0 {
		return fmt.Errorf("Instances failed to register in the load balancer after 20 minutes")
	}
	log.Infof("All instances are now registered in the '%s' load balancer", name)
	return nil
}
