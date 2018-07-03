package templates

import (
	"fmt"
	"strconv"

	gocf "github.com/crewjam/go-cloudformation"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

func stateless(t *gocf.Template, s *service.Service, e *environment.Environment, encodedCompose string) error {
	if err := asg(t, s, e, encodedCompose); err != nil {
		return err
	}
	if len(s.GetPorts()) == 0 {
		return nil
	}
	return elb(t, s, e)
}

func asg(t *gocf.Template, s *service.Service, e *environment.Environment, encodedCompose string) error {
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
				VpcId:                gocf.String(e.Provider.VPC),
			}
			t.AddResource("InstanceSecurityGroup", securityGroup)
			securityGroups = append(securityGroups, gocf.Ref("InstanceSecurityGroup").String())
			securityGroup = &gocf.EC2SecurityGroup{
				GroupDescription:     gocf.String(fmt.Sprintf("ELB Security Group for %s", s.GetFullName(e, "-"))),
				SecurityGroupIngress: &loadbalancerIngressRules,
				VpcId:                gocf.String(e.Provider.VPC),
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
		ImageId:            gocf.String(e.Provider.Ami),
		InstanceType:       gocf.String(s.GetInstanceType(e)),
		KeyName:            gocf.String(e.Provider.Keypair),
		SecurityGroups:     securityGroups,
		IamInstanceProfile: gocf.Ref("InstanceProfile").String(),
		UserData:           gocf.String(userData(s.GetFullName(e, "-"), encodedCompose, e)),
	}
	t.AddResource("LaunchConfig", launchConfig)
	return nil
}

func elb(t *gocf.Template, s *service.Service, e *environment.Environment) error {
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
