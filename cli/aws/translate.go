package aws

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"bitbucket.org/riberaproject/cli/compose"
	"bitbucket.org/riberaproject/cli/service"
	gocf "github.com/crewjam/go-cloudformation"
	"github.com/google/uuid"
)

var amisPerRegion = map[string]string{
	"us-west-2": "ami-7293320a",
}

// Translate an i2kit service to a AWS CloudFormation template
func Translate(s *service.Service, space string, signal string) (string, error) {
	region := os.Getenv("I2KIT_REGION")
	ami, ok := amisPerRegion[region]
	if !ok {
		return "", fmt.Errorf("Region'%s' is not supported", region)
	}
	t := gocf.NewTemplate()
	hostedZone := os.Getenv("I2KIT_HOSTED_ZONE")
	err := loadASG(t, s, space, ami, hostedZone, signal)
	if err != nil {
		return "", err
	}
	resourceName := s.Name
	if space != "" {
		resourceName = fmt.Sprintf("%s-%s", s.Name, space)
	}
	loadELB(t, s, resourceName)
	loadIAM(t, resourceName)
	loadLogGroup(t, resourceName)
	loadRoute53(t, s.Name, space, hostedZone)
	marshalledTemplate := []byte("")
	if marshalledTemplate, err = json.Marshal(t); err != nil {
		return "", err
	}
	return string(marshalledTemplate), nil
}

func loadASG(t *gocf.Template, s *service.Service, space string, ami, hostedZone, signal string) error {
	var domain string
	if hostedZone != "" {
		if space != "" {
			domain = fmt.Sprintf("%s.%s", space, hostedZone)
		} else {
			domain = hostedZone
		}
		domain = domain[:len(domain)-1]
	}
	encodedCompose, err := compose.Create(s, domain)
	if err != nil {
		return err
	}
	replicas := strconv.Itoa(s.Replicas)
	asg := &gocf.AutoScalingAutoScalingGroup{
		HealthCheckGracePeriod:  gocf.Integer(15),
		LaunchConfigurationName: gocf.Ref("LaunchConfig").String(),
		LoadBalancerNames:       gocf.StringList(gocf.Ref("ELB")),
		MaxSize:                 gocf.String(replicas),
		MinSize:                 gocf.String(replicas),
		VPCZoneIdentifier:       gocf.StringList(gocf.String(os.Getenv("I2KIT_SUBNET"))),
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
	asgResource := &gocf.Resource{Properties: asg, UpdatePolicy: updatePolicy}
	t.Resources["ASG"] = asgResource
	instanceType := os.Getenv("I2KIT_INSTANCE_TYPE")
	if instanceType == "" {
		instanceType = "t2.micro"
	}
	containerName := s.Name
	if space != "" {
		containerName = fmt.Sprintf("%s-%s", s.Name, space)
	}
	launchConfig := &gocf.AutoScalingLaunchConfiguration{
		ImageId:            gocf.String(ami),
		InstanceType:       gocf.String("t2.small"),
		KeyName:            gocf.String(os.Getenv("I2KIT_KEYPAIR")),
		SecurityGroups:     []string{os.Getenv("I2KIT_SECURITY_GROUP")},
		IamInstanceProfile: gocf.Ref("InstanceProfile").String(),
		UserData:           gocf.String(userData(containerName, encodedCompose, signal)),
	}
	t.AddResource("LaunchConfig", launchConfig)
	return nil
}

func userData(containerName, encodedCompose, signal string) string {
	region := os.Getenv("I2KIT_REGION")
	uniqueOperationID := uuid.New().String()
	value := fmt.Sprintf(
		`#!/bin/bash

set -e
sudo docker run \
	--name %s \
	-e COMPOSE=%s \
	-e CONFIG=%s \
	-e UNIQUE_OPERATION_ID=%s \
	-e STACK=%s \
	-e REGION=%s \
	-e SIGNAL=%s \
	-v /var/run/docker.sock:/var/run/docker.sock \
	--log-driver=awslogs \
	--log-opt awslogs-region=us-west-2 \
	--log-opt awslogs-group=i2kit-%s \
	--log-opt tag='{{ with split .Name ":" }}{{join . "_"}}{{end}}-{{.ID}}' \
	riberaproject/agent`,
		containerName,
		encodedCompose,
		os.Getenv("I2KIT_DOCKER_CONFIG"),
		uniqueOperationID,
		containerName,
		region,
		signal,
		containerName,
	)
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func loadELB(t *gocf.Template, s *service.Service, elbName string) {
	instancePort := ""
	listeners := gocf.ElasticLoadBalancingListenerList{}
	for _, container := range s.Containers {
		for _, port := range container.Ports {
			if instancePort == "" {
				instancePort = port.InstancePort
			}
			listeners = append(listeners, gocf.ElasticLoadBalancingListener{
				InstancePort:     gocf.String(port.InstancePort),
				InstanceProtocol: gocf.String(port.InstanceProtocol),
				LoadBalancerPort: gocf.String(port.Port),
				Protocol:         gocf.String(port.Protocol),
				SSLCertificateId: gocf.String(port.Certificate),
			})
		}
	}
	if len(listeners) == 0 {
		return
	}
	elb := &gocf.ElasticLoadBalancingLoadBalancer{
		LoadBalancerName: gocf.String(elbName),
		Subnets:          gocf.StringList(gocf.String(os.Getenv("I2KIT_SUBNET"))),
		SecurityGroups:   []string{os.Getenv("I2KIT_SECURITY_GROUP")},
		HealthCheck: &gocf.ElasticLoadBalancingHealthCheck{
			HealthyThreshold:   gocf.String("2"),
			Interval:           gocf.String("15"),
			Target:             gocf.String(fmt.Sprintf("TCP:%s", instancePort)),
			Timeout:            gocf.String("10"),
			UnhealthyThreshold: gocf.String("2"),
		},
	}
	elb.Listeners = &listeners
	t.Outputs["elbURL"] = &gocf.Output{
		Description: "The URL of the stack",
		Value:       gocf.Join("", gocf.String("http://"), gocf.GetAtt("ELB", "DNSName")),
	}
	t.Outputs["elbName"] = &gocf.Output{
		Description: "Load balancer name",
		Value:       gocf.Ref("ELB"),
	}
	t.AddResource("ELB", elb)
}

func loadIAM(t *gocf.Template, policyName string) {
	policy := gocf.IAMPolicies{
		PolicyName: gocf.String(policyName),
		PolicyDocument: &map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": map[string]interface{}{
				"Effect":   "Allow",
				"Action":   []string{"logs:CreateLogStream", "logs:PutLogEvents"},
				"Resource": fmt.Sprintf("arn:aws:logs:us-west-2:*:log-group:i2kit-%s:log-stream:%s-*", policyName, policyName),
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

func loadLogGroup(t *gocf.Template, logGroupName string) {
	logGroup := &gocf.LogsLogGroup{
		LogGroupName:    gocf.String(fmt.Sprintf("i2kit-%s", logGroupName)),
		RetentionInDays: gocf.Integer(30),
	}
	t.AddResource("LogGroup", logGroup)
}

func loadRoute53(t *gocf.Template, name, space, hostedZone string) {
	recordName := ""
	if space == "" {
		recordName = fmt.Sprintf("%s.%s", name, hostedZone)
	} else {
		recordName = fmt.Sprintf("%s.%s.%s", name, space, hostedZone)
	}
	recordSet := &gocf.Route53RecordSet{
		HostedZoneName:  gocf.String(hostedZone),
		Name:            gocf.String(recordName),
		Type:            gocf.String("CNAME"),
		TTL:             gocf.String("900"),
		ResourceRecords: gocf.StringList(gocf.GetAtt("ELB", "DNSName")),
	}
	t.AddResource("DNSRecord", recordSet)
}
