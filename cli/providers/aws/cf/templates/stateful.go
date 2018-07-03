package templates

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	gocf "github.com/crewjam/go-cloudformation"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

func stateful(t *gocf.Template, s *service.Service, e *environment.Environment, ami, vpc, encodedCompose string) error {
	rand.Seed(time.Now().Unix())
	subnet := *e.Provider.Subnets[rand.Intn(len(e.Provider.Subnets))]
	securityGroups := &gocf.StringListExpr{
		Literal: []*gocf.StringExpr{gocf.String(e.Provider.SecurityGroup)}}
	if s.Public {
		securityGroup(t, s, e, vpc)
		securityGroups.Literal = append(securityGroups.Literal, gocf.Ref("SecurityGroup").String())
	}

	ec2Instance := &gocf.EC2Instance{
		IamInstanceProfile: gocf.Ref("InstanceProfile").String(),
		ImageId:            gocf.String(ami),
		InstanceType:       gocf.String(s.GetInstanceType(e)),
		KeyName:            gocf.String(e.Provider.Keypair),
		SecurityGroupIds:   securityGroups,
		SubnetId:           gocf.String(subnet),
		UserData:           gocf.String(userData(s.GetFullName(e, "-"), encodedCompose, e)),
	}
	t.AddResource("EC2Instance", ec2Instance)
	elasticIP := &gocf.EC2EIP{
		Domain:     gocf.String("vpc"),
		InstanceId: gocf.Ref("EC2Instance").String(),
	}
	t.AddResource("EIP", elasticIP)
	return nil
}

func securityGroup(t *gocf.Template, s *service.Service, e *environment.Environment, vpc string) {
	instanceIngressRules := gocf.EC2SecurityGroupRuleList{}
	for _, port := range s.GetPorts() {
		portNumber, _ := strconv.ParseInt(port.Port, 10, 64)
		instanceIngressRules = append(
			instanceIngressRules,
			gocf.EC2SecurityGroupRule{
				CidrIp:     gocf.String("0.0.0.0/0"),
				IpProtocol: gocf.String("tcp"),
				FromPort:   gocf.Integer(portNumber),
				ToPort:     gocf.Integer(portNumber),
			})
	}
	if len(instanceIngressRules) > 0 {
		securityGroup := &gocf.EC2SecurityGroup{
			GroupDescription:     gocf.String(fmt.Sprintf("Security Group for %s", s.GetFullName(e, "-"))),
			SecurityGroupIngress: &instanceIngressRules,
			VpcId:                gocf.String(vpc),
		}
		t.AddResource("SecurityGroup", securityGroup)
	}
}
