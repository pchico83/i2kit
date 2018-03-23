package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pchico83/i2kit/cli/schemas/environment"
)

//GetVPC returns the vpc-id of the environment subnets
func GetVPC(e *environment.Environment, config *aws.Config) (string, error) {
	svc := ec2.New(session.New(), config)
	dsi := &ec2.DescribeSubnetsInput{
		SubnetIds: e.Provider.Subnets,
	}
	dso, err := svc.DescribeSubnets(dsi)
	if err != nil {
		return "", err
	}
	vpc := ""
	for _, s := range dso.Subnets {
		if vpc == "" {
			vpc = *s.VpcId
		}
		if vpc != *s.VpcId {
			return "", fmt.Errorf("Environment subnets belong to different VPCs")
		}
	}
	return vpc, nil
}
