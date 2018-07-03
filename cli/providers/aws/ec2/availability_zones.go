package ec2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pchico83/i2kit/cli/schemas/environment"
)

//AvailabilityZones returns availability zones where subnets can be created
func AvailabilityZones(e *environment.Environment, config *aws.Config) ([]string, error) {
	svc := ec2.New(session.New(), config)
	dazi := &ec2.DescribeAvailabilityZonesInput{}
	dazo, err := svc.DescribeAvailabilityZones(dazi)
	if err != nil {
		return nil, err
	}
	var azs []string
	for _, az := range dazo.AvailabilityZones {
		azs = append(azs, *az.ZoneName)
	}
	return azs, nil
}
