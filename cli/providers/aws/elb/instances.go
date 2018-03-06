package elb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
)

//Instances returns the instance ids of the instances registered in the ELB
func Instances(name string, config *aws.Config) []string {
	instanceIDs := []string{}
	svc := elb.New(session.New(), config)
	dii := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: aws.String(name),
	}
	dio, err := svc.DescribeInstanceHealth(dii)
	if err != nil {
		return instanceIDs
	}
	for _, instance := range dio.InstanceStates {
		instanceIDs = append(instanceIDs, *instance.InstanceId)
	}
	return instanceIDs
}
