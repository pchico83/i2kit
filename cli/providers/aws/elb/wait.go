package elb

import (
	"fmt"
	"time"

	logger "log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

//Wait waits for the instances in a ELB to be healthy
func Wait(s *service.Service, name string, config *aws.Config, log *logger.Logger) error {
	registeredInstances, err := waitForRegistered(s, name, config, log)
	if err != nil {
		return err
	}
	log.Printf("Waiting for %s to be healthy in the '%s' load balancer...", registeredInstances, name)
	for _, instanceID := range registeredInstances {
		if err := waitForHealthy(name, instanceID, config, log); err != nil {
			return err
		}
	}
	log.Printf("%v are healthy in the '%s' load balancer.", registeredInstances, name)
	return nil
}

func waitForRegistered(s *service.Service, name string, config *aws.Config, log *logger.Logger) ([]string, error) {
	log.Printf("Waiting for instances to be registered in the '%s' load balancer...", name)
	retries := 0
	for retries < 30 {
		registeredInstances := Instances(name, config)
		log.Printf("%v instances registered in the '%s' load balancer.", registeredInstances, name)
		if len(registeredInstances) == s.Replicas {
			return registeredInstances, nil
		}
		time.Sleep(10 * time.Second)
	}
	return nil, fmt.Errorf("Instances do not register in the load balancer after 5 minutes")
}

func waitForHealthy(name, instanceID string, config *aws.Config, log *logger.Logger) error {
	svc := elb.New(session.New(), config)
	dii := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: aws.String(name),
	}
	retries := 0
	startTime := new(int64)
	*startTime = 0
	for retries < 90 {
		dio, err := svc.DescribeInstanceHealth(dii)
		if err == nil {
			for _, instance := range dio.InstanceStates {
				if *instance.InstanceId == instanceID && *instance.State == "InService" {
					log.Printf("%s is healthy in the '%s' load balancer.", instanceID, name)
					return nil
				}
			}
		}
		time.Sleep(10 * time.Second)
		retries++
	}
	return fmt.Errorf("'%s' is unhealthy after 15 minutes", instanceID)
}
