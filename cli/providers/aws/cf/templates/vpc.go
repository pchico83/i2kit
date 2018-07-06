package templates

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	gocf "github.com/crewjam/go-cloudformation"
	"github.com/pchico83/i2kit/cli/providers/aws/ec2"
	"github.com/pchico83/i2kit/cli/schemas/environment"
)

// VPC returns a AWS CloudFormation template representing an environment
func VPC(e *environment.Environment, config *aws.Config) (string, error) {
	t := gocf.NewTemplate()

	logGroupIAM(t, e)

	vpc := &gocf.EC2VPC{
		CidrBlock:          gocf.String("10.192.0.0/16"),
		EnableDnsSupport:   gocf.Bool(true),
		EnableDnsHostnames: gocf.Bool(true),
	}
	t.AddResource("VPC", vpc)

	ig := &gocf.EC2InternetGateway{}
	t.AddResource("InternetGateway", ig)

	iga := &gocf.EC2VPCGatewayAttachment{
		InternetGatewayId: gocf.Ref("InternetGateway").String(),
		VpcId:             gocf.Ref("VPC").String(),
	}
	t.AddResource("InternetGatewayAttachment", iga)

	rt := &gocf.EC2RouteTable{
		VpcId: gocf.Ref("VPC").String(),
	}
	t.AddResource("RouteTable", rt)

	r := &gocf.EC2Route{
		RouteTableId:         gocf.Ref("RouteTable").String(),
		DestinationCidrBlock: gocf.String("0.0.0.0/0"),
		GatewayId:            gocf.Ref("InternetGateway").String(),
	}
	t.AddResource("Route", r)

	azs, err := ec2.AvailabilityZones(e, config)
	if err != nil {
		return "", err
	}
	for i, az := range azs {
		cidr := fmt.Sprintf("10.192.1%d.0/24", i)
		subnetName := fmt.Sprintf("Subnet%d", i)
		subnetRouteTableAssociationName := fmt.Sprintf("Subnet%dRouteTableAssociation", i)
		subnet := &gocf.EC2Subnet{
			VpcId:               gocf.Ref("VPC").String(),
			AvailabilityZone:    gocf.String(az),
			CidrBlock:           gocf.String(cidr),
			MapPublicIpOnLaunch: gocf.Bool(true),
		}
		t.AddResource(subnetName, subnet)
		srta := &gocf.EC2SubnetRouteTableAssociation{
			RouteTableId: gocf.Ref("RouteTable").String(),
			SubnetId:     gocf.Ref(subnetName).String(),
		}
		t.AddResource(subnetRouteTableAssociationName, srta)
	}

	sg := &gocf.EC2SecurityGroup{
		GroupDescription: gocf.String(fmt.Sprintf("Security Group for internal traffic in project %s", e.Name)),
		VpcId:            gocf.Ref("VPC").String(),
	}
	if e.Provider.Keypair != "" {
		sg.SecurityGroupIngress = &gocf.EC2SecurityGroupRuleList{
			gocf.EC2SecurityGroupRule{
				IpProtocol: gocf.String("tcp"),
				CidrIp:     gocf.String("0.0.0.0/0"),
				FromPort:   gocf.Integer(22),
				ToPort:     gocf.Integer(22),
			},
		}
	}
	t.AddResource("SecurityGroup", sg)

	sgi := &gocf.EC2SecurityGroupIngress{
		GroupId:               gocf.Ref("SecurityGroup").String(),
		IpProtocol:            gocf.String("-1"),
		SourceSecurityGroupId: gocf.Ref("SecurityGroup").String(),
	}
	t.AddResource("SecurityGroupIngress", sgi)

	t.Outputs["VPC"] = &gocf.Output{
		Description: "VPC id",
		Value:       gocf.Ref("VPC"),
	}
	var subnetRefs []gocf.Stringable
	for i := range azs {
		subnetName := fmt.Sprintf("Subnet%d", i)
		subnetRefs = append(subnetRefs, gocf.Ref(subnetName))
	}
	t.Outputs["Subnets"] = &gocf.Output{
		Description: "List of subnet ids",
		Value:       gocf.Join(",", subnetRefs...),
	}
	t.Outputs["SecurityGroup"] = &gocf.Output{
		Description: "Security Group id",
		Value:       gocf.Ref("SecurityGroup"),
	}

	marshalledTemplate, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(marshalledTemplate), nil
}
