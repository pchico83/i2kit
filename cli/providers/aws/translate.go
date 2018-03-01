package aws

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	gocf "github.com/crewjam/go-cloudformation"
	"github.com/google/uuid"
	"github.com/pchico83/i2kit/cli/schemas/compose"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

var amisPerRegion = map[string]string{
	"us-west-2": "ami-7293320a",
}

// Translate an i2kit service to a AWS CloudFormation template
func Translate(s *service.Service, e *environment.Environment, signal string) (string, error) {
	ami, ok := amisPerRegion[e.Provider.Region]
	if !ok {
		return "", fmt.Errorf("Region'%s' is not supported", e.Provider.Region)
	}
	t := gocf.NewTemplate()
	err := loadASG(t, s, e, ami, signal)
	if err != nil {
		return "", err
	}
	err = loadELB(t, s, e)
	if err != nil {
		return "", err
	}
	loadIAM(t, s, e)
	loadLogGroup(t, s, e)
	loadRoute53(t, s, e)
	marshalledTemplate := []byte("")
	if marshalledTemplate, err = json.Marshal(t); err != nil {
		return "", err
	}
	return string(marshalledTemplate), nil
}

func loadASG(t *gocf.Template, s *service.Service, e *environment.Environment, ami, signal string) error {
	domain := e.Provider.HostedZone[:len(e.Provider.HostedZone)-1]
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
		VPCZoneIdentifier:       gocf.StringList(gocf.String(e.Provider.Subnet)),
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
	launchConfig := &gocf.AutoScalingLaunchConfiguration{
		ImageId:            gocf.String(ami),
		InstanceType:       gocf.String(s.GetSize(e)),
		KeyName:            gocf.String(e.Provider.Keypair),
		SecurityGroups:     []string{e.Provider.SecurityGroup},
		IamInstanceProfile: gocf.Ref("InstanceProfile").String(),
		UserData:           gocf.String(userData(s.Name, encodedCompose, e, signal)),
	}
	t.AddResource("LaunchConfig", launchConfig)
	return nil
}

func userData(containerName, encodedCompose string, e *environment.Environment, signal string) string {
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
		e.B64DockerConfig(),
		uniqueOperationID,
		containerName,
		e.Provider.Region,
		signal,
		containerName,
	)
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func loadELB(t *gocf.Template, s *service.Service, e *environment.Environment) error {
	healthCheckPort := ""
	listeners := gocf.ElasticLoadBalancingListenerList{}
	for _, container := range s.Containers {
		for _, port := range container.Ports {
			if healthCheckPort == "" {
				healthCheckPort = port.InstancePort
			}
			certificate := port.Certificate
			if certificate == "" {
				certificate = e.Provider.Certificate
			}
			if certificate == "" && (port.Protocol == "HTTPS" || port.Protocol == "SSL") {
				return fmt.Errorf("Port '%s:%s' requires a certificate", port.Protocol, port.Port)
			}
			listeners = append(listeners, gocf.ElasticLoadBalancingListener{
				InstancePort:     gocf.String(port.InstancePort),
				InstanceProtocol: gocf.String(port.InstanceProtocol),
				LoadBalancerPort: gocf.String(port.Port),
				Protocol:         gocf.String(port.Protocol),
				SSLCertificateId: gocf.String(certificate),
			})
		}
	}
	if len(listeners) == 0 {
		return nil
	}
	elb := &gocf.ElasticLoadBalancingLoadBalancer{
		LoadBalancerName: gocf.String(s.Name),
		Subnets:          gocf.StringList(gocf.String(e.Provider.Subnet)),
		SecurityGroups:   []string{e.Provider.SecurityGroup},
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
		PolicyName: gocf.String(s.Name),
		PolicyDocument: &map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": map[string]interface{}{
				"Effect":   "Allow",
				"Action":   []string{"logs:CreateLogStream", "logs:PutLogEvents"},
				"Resource": fmt.Sprintf("arn:aws:logs:us-west-2:*:log-group:i2kit-%s:log-stream:%s-*", s.Name, s.Name),
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
	logGroup := &gocf.LogsLogGroup{
		LogGroupName:    gocf.String(fmt.Sprintf("i2kit-%s", s.Name)),
		RetentionInDays: gocf.Integer(30),
	}
	t.AddResource("LogGroup", logGroup)
}

func loadRoute53(t *gocf.Template, s *service.Service, e *environment.Environment) {
	recordName := fmt.Sprintf("%s.%s", s.Name, e.Provider.HostedZone)
	recordSet := &gocf.Route53RecordSet{
		HostedZoneName:  gocf.String(e.Provider.HostedZone),
		Name:            gocf.String(recordName),
		Type:            gocf.String("CNAME"),
		TTL:             gocf.String("900"),
		ResourceRecords: gocf.StringList(gocf.GetAtt("ELB", "DNSName")),
	}
	t.AddResource("DNSRecord", recordSet)
}
