package aws

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	gocf "github.com/crewjam/go-cloudformation"
	"github.com/google/uuid"
	"github.com/pchico83/i2kit/cli/providers/aws/ec2"
	"github.com/pchico83/i2kit/cli/schemas/compose"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

var amisPerRegion = map[string]string{
	"us-east-2":      "ami-1b90a67e",
	"us-east-1":      "ami-cb17d8b6",
	"us-west-2":      "ami-05b5277d",
	"us-west-1":      "ami-9cbbaffc",
	"eu-west-3":      "ami-914afcec",
	"eu-west-2":      "ami-a48d6bc3",
	"eu-west-1":      "ami-bfb5fec6",
	"eu-central-1":   "ami-ac055447",
	"ap-northeast-2": "ami-ba74d8d4",
	"ap-northeast-1": "ami-5add893c",
	"ap-southeast-2": "ami-4cc5072e",
	"ap-southeast-1": "ami-acbcefd0",
	"ca-central-1":   "ami-a535b2c1",
	"ap-south-1":     "ami-2149114e",
	"sa-east-1":      "ami-d3bce9bf",
}

// Translate an i2kit service to a AWS CloudFormation template
func Translate(s *service.Service, e *environment.Environment, config *aws.Config) (string, error) {
	ami, ok := amisPerRegion[e.Provider.Region]
	if !ok {
		return "", fmt.Errorf("Region'%s' is not supported", e.Provider.Region)
	}
	t := gocf.NewTemplate()
	vpc, err := ec2.GetVPC(e, config)
	if err != nil {
		return "", err
	}
	encodedCompose, err := compose.Create(s, e)
	if err != nil {
		return "", err
	}
	if s.Stateful {
		if err = loadStateful(t, s, e, ami, vpc, encodedCompose); err != nil {
			return "", err
		}
	} else {
		if err = loadStateless(t, s, e, ami, vpc, encodedCompose); err != nil {
			return "", err
		}
	}
	loadIAM(t, s, e)
	loadLogGroup(t, s, e)
	ports := s.GetPorts()
	if len(ports) > 0 && e.Provider.HostedZone != "" {
		loadRoute53(t, s, e)
	}
	marshalledTemplate, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(marshalledTemplate), nil
}

func loadStateful(t *gocf.Template, s *service.Service, e *environment.Environment, ami, vpc, encodedCompose string) error {
	rand.Seed(time.Now().Unix())
	subnet := *e.Provider.Subnets[rand.Intn(len(e.Provider.Subnets))]
	securityGroups := &gocf.StringListExpr{
		Literal: []*gocf.StringExpr{gocf.String(e.Provider.SecurityGroup)}}
	if s.Public {
		loadSG(t, s, e, vpc)
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

func loadSG(t *gocf.Template, s *service.Service, e *environment.Environment, vpc string) {
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

func loadStateless(t *gocf.Template, s *service.Service, e *environment.Environment, ami, vpc, encodedCompose string) error {
	if err := loadASG(t, s, e, ami, vpc, encodedCompose); err != nil {
		return err
	}
	if len(s.GetPorts()) == 0 {
		return nil
	}
	return loadELB(t, s, e)
}

func loadASG(t *gocf.Template, s *service.Service, e *environment.Environment, ami, vpc, encodedCompose string) error {
	securityGroups := []gocf.Stringable{gocf.String(e.Provider.SecurityGroup)}
	loadBalancerNames := gocf.StringList()
	instanceIngressRules := gocf.EC2SecurityGroupRuleList{}
	loadbalancerIngressRules := gocf.EC2SecurityGroupRuleList{}
	for _, port := range s.GetPorts() {
		instancePortNumber, _ := strconv.ParseInt(port.InstancePort, 10, 64)
		instanceIngressRules = append(
			instanceIngressRules,
			gocf.EC2SecurityGroupRule{
				SourceSecurityGroupIdXXSecurityGroupIngressXOnlyX: gocf.Ref("ELBSecurityGroup").String(),
				IpProtocol: gocf.String("tcp"),
				FromPort:   gocf.Integer(instancePortNumber),
				ToPort:     gocf.Integer(instancePortNumber),
			})
		portNumber, _ := strconv.ParseInt(port.Port, 10, 64)
		loadbalancerIngressRules = append(
			loadbalancerIngressRules,
			gocf.EC2SecurityGroupRule{
				CidrIp:     gocf.String("0.0.0.0/0"),
				IpProtocol: gocf.String("tcp"),
				FromPort:   gocf.Integer(portNumber),
				ToPort:     gocf.Integer(portNumber),
			})
	}
	if len(instanceIngressRules) > 0 {
		if s.Public {
			securityGroup := &gocf.EC2SecurityGroup{
				GroupDescription:     gocf.String(fmt.Sprintf("Instance Security Group for %s", s.GetFullName(e, "-"))),
				SecurityGroupIngress: &instanceIngressRules,
				VpcId:                gocf.String(vpc),
			}
			t.AddResource("InstanceSecurityGroup", securityGroup)
			securityGroups = append(securityGroups, gocf.Ref("InstanceSecurityGroup").String())
			securityGroup = &gocf.EC2SecurityGroup{
				GroupDescription:     gocf.String(fmt.Sprintf("ELB Security Group for %s", s.GetFullName(e, "-"))),
				SecurityGroupIngress: &loadbalancerIngressRules,
				VpcId:                gocf.String(vpc),
			}
			t.AddResource("ELBSecurityGroup", securityGroup)
		}
		loadBalancerNames = gocf.StringList(gocf.Ref("ELB"))
	}

	subnets := &gocf.StringListExpr{Literal: []*gocf.StringExpr{}}
	for _, item := range e.Provider.Subnets {
		subnets.Literal = append(subnets.Literal, gocf.String(*item))
	}
	replicas := strconv.Itoa(s.Replicas)
	asg := &gocf.AutoScalingAutoScalingGroup{
		HealthCheckGracePeriod:  gocf.Integer(15),
		LaunchConfigurationName: gocf.Ref("LaunchConfig").String(),
		LoadBalancerNames:       loadBalancerNames,
		MaxSize:                 gocf.String(replicas),
		MinSize:                 gocf.String(replicas),
		VPCZoneIdentifier:       subnets,
	}
	creationPolicy := &gocf.CreationPolicy{
		ResourceSignal: &gocf.CreationPolicyResourceSignal{
			Count:   gocf.Integer(int64(s.Replicas)),
			Timeout: gocf.String("PT15M"),
		},
	}
	updatePolicy := &gocf.UpdatePolicy{
		AutoScalingRollingUpdate: &gocf.UpdatePolicyAutoScalingRollingUpdate{
			MaxBatchSize:          gocf.Integer(1),
			PauseTime:             gocf.String("PT15M"),
			WaitOnResourceSignals: gocf.Bool(true),
			SuspendProcesses: gocf.StringList(
				gocf.String("HealthCheck"),
				gocf.String("ReplaceUnhealthy"),
				gocf.String("AZRebalance"),
				gocf.String("AlarmNotification"),
				gocf.String("ScheduledActions"),
			),
		},
	}
	asgResource := &gocf.Resource{
		Properties:     asg,
		CreationPolicy: creationPolicy,
		UpdatePolicy:   updatePolicy,
	}
	t.Resources["ASG"] = asgResource

	launchConfig := &gocf.AutoScalingLaunchConfiguration{
		ImageId:            gocf.String(ami),
		InstanceType:       gocf.String(s.GetInstanceType(e)),
		KeyName:            gocf.String(e.Provider.Keypair),
		SecurityGroups:     securityGroups,
		IamInstanceProfile: gocf.Ref("InstanceProfile").String(),
		UserData:           gocf.String(userData(s.GetFullName(e, "-"), encodedCompose, e)),
	}
	t.AddResource("LaunchConfig", launchConfig)
	return nil
}

func userData(containerName, encodedCompose string, e *environment.Environment) string {
	value := fmt.Sprintf(
		`#!/bin/bash

set -e
INSTANCE_ID=$(curl http://169.254.169.254/latest/meta-data/instance-id)
sudo docker run \
	--name %s \
	-e COMPOSE=%s \
	-e CONFIG=%s \
	-e INSTANCE_ID=$INSTANCE_ID \
	-e STACK=%s \
	-e REGION=%s \
	-v /var/run/docker.sock:/var/run/docker.sock \
	--log-driver=awslogs \
	--log-opt awslogs-region=%s \
	--log-opt awslogs-group=i2kit-%s \
	--log-opt tag=i2kit-$INSTANCE_ID \
	okteto/agent:1.0`,
		containerName,
		encodedCompose,
		e.B64DockerConfig(),
		containerName,
		e.Provider.Region,
		e.Provider.Region,
		containerName,
	)
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func loadELB(t *gocf.Template, s *service.Service, e *environment.Environment) error {
	healthCheckPort := ""
	listeners := gocf.ElasticLoadBalancingListenerList{}
	for _, port := range s.GetPorts() {
		if healthCheckPort == "" {
			healthCheckPort = port.InstancePort
		}
		certificate := port.Certificate
		if certificate == "" && (port.Protocol == service.HTTPS || port.Protocol == service.SSL) {
			if e.Provider.Certificate == "" {
				return fmt.Errorf("Port '%s:%s' requires a certificate", port.Protocol, port.Port)
			}
			certificate = e.Provider.Certificate
		}
		listeners = append(listeners, gocf.ElasticLoadBalancingListener{
			InstancePort:     gocf.String(port.InstancePort),
			InstanceProtocol: gocf.String(port.InstanceProtocol),
			LoadBalancerPort: gocf.String(port.Port),
			Protocol:         gocf.String(port.Protocol),
			SSLCertificateId: gocf.String(certificate),
		})
	}
	subnets := &gocf.StringListExpr{Literal: []*gocf.StringExpr{}}
	for _, item := range e.Provider.Subnets {
		subnets.Literal = append(subnets.Literal, gocf.String(*item))
	}
	securityGroups := []gocf.Stringable{gocf.String(e.Provider.SecurityGroup)}
	if s.Public {
		securityGroups = append(securityGroups, gocf.Ref("ELBSecurityGroup").String())
	}
	schema := "internal"
	if s.Public {
		schema = "internet-facing"
	}
	crossZone := false
	if len(subnets.Literal) > 0 {
		crossZone = true
	}
	elb := &gocf.ElasticLoadBalancingLoadBalancer{
		Subnets:        subnets,
		Scheme:         gocf.String(schema),
		CrossZone:      gocf.Bool(crossZone),
		SecurityGroups: securityGroups,
		HealthCheck: &gocf.ElasticLoadBalancingHealthCheck{
			HealthyThreshold:   gocf.String("2"),
			Interval:           gocf.String("15"),
			Target:             gocf.String(fmt.Sprintf("TCP:%s", healthCheckPort)),
			Timeout:            gocf.String("10"),
			UnhealthyThreshold: gocf.String("2"),
		},
	}
	elb.Listeners = &listeners
	t.Outputs["elbURL"] = &gocf.Output{
		Description: "The URL of the stack",
		Value:       gocf.GetAtt("ELB", "DNSName"),
	}
	t.Outputs["elbName"] = &gocf.Output{
		Description: "Load balancer name",
		Value:       gocf.Ref("ELB"),
	}
	t.AddResource("ELB", elb)
	return nil
}

func loadIAM(t *gocf.Template, s *service.Service, e *environment.Environment) {
	policy := gocf.IAMPolicies{
		PolicyName: gocf.String(s.GetFullName(e, "-")),
		PolicyDocument: &map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": map[string]interface{}{
				"Effect":   "Allow",
				"Action":   []string{"logs:CreateLogStream", "logs:PutLogEvents"},
				"Resource": fmt.Sprintf("arn:aws:logs:%s:*:log-group:i2kit-%s:log-stream:*", e.Provider.Region, s.GetFullName(e, "-")),
			},
		},
	}
	role := &gocf.IAMRole{
		AssumeRolePolicyDocument: &map[string]interface{}{
			"Statement": map[string]interface{}{
				"Effect":    "Allow",
				"Principal": map[string]interface{}{"Service": []string{"ec2.amazonaws.com"}},
				"Action":    []string{"sts:AssumeRole"},
			},
		},
		Path:     gocf.String("/"),
		Policies: &gocf.IAMPoliciesList{policy},
	}
	t.AddResource("Role", role)
	instanceProfile := &gocf.IAMInstanceProfile{
		Path:  gocf.String("/"),
		Roles: gocf.StringList(gocf.Ref("Role")),
	}
	t.AddResource("InstanceProfile", instanceProfile)
}

func loadLogGroup(t *gocf.Template, s *service.Service, e *environment.Environment) {
	randomID := uuid.New()
	logGroupResource := &gocf.Resource{
		Properties: &gocf.LogsLogGroup{
			LogGroupName:    gocf.String(fmt.Sprintf("i2kit-%s-%d", s.GetFullName(e, "-"), randomID.ID())),
			RetentionInDays: gocf.Integer(30),
		},
		DeletionPolicy: "Retain",
	}
	t.Resources["LogGroup"] = logGroupResource
}

func loadRoute53(t *gocf.Template, s *service.Service, e *environment.Environment) {
	var resourceRecords *gocf.StringListExpr
	var dependsOn []string
	if s.Stateful {
		resourceRecords = gocf.StringList(gocf.GetAtt("EC2Instance", "PublicDnsName"))
		dependsOn = append(dependsOn, "EIP")
	} else {
		resourceRecords = gocf.StringList(gocf.GetAtt("ELB", "DNSName"))
	}
	recordName := fmt.Sprintf("%s.%s", s.GetFullName(e, "."), e.Provider.HostedZone)
	recordSetProperties := &gocf.Route53RecordSet{
		HostedZoneName:  gocf.String(e.Provider.HostedZone),
		Name:            gocf.String(recordName),
		Type:            gocf.String("CNAME"),
		TTL:             gocf.String("60"),
		ResourceRecords: resourceRecords,
	}
	resourceSetResource := &gocf.Resource{
		Properties: recordSetProperties,
		DependsOn:  dependsOn,
	}
	t.Resources["DNSRecord"] = resourceSetResource
}
