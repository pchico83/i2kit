package cf

import (
	"encoding/json"
	"fmt"
	"strconv"

	gocf "github.com/crewjam/go-cloudformation"
	"github.com/pchico83/i2kit/k8"
)

const hostedZone string = "i2kit.com."

// ELB for CF
func loadBalancerSection(deployment *k8.Deployment) *gocf.ElasticLoadBalancingLoadBalancer {
	elb := &gocf.ElasticLoadBalancingLoadBalancer{
		LoadBalancerName: gocf.String("testing-i2kit"),
		Subnets:          gocf.StringList(gocf.String("subnet-7c505724")), // TODO configurable parameter
		SecurityGroups:   []string{"sg-dfce36a5"},                         // TODO configurable parameter
	}
	listeners := gocf.ElasticLoadBalancingListenerList{}
	for _, container := range deployment.Spec.Template.Spec.Containers {
		for _, port := range container.Ports {
			listeners = append(listeners, gocf.ElasticLoadBalancingListener{
				InstancePort:     gocf.String(strconv.Itoa(port.ContainerPort)),
				InstanceProtocol: gocf.String("HTTP"),
				LoadBalancerPort: gocf.String(strconv.Itoa(port.ContainerPort)),
				Protocol:         gocf.String("HTTP"),
			})
		}
	}
	if len(listeners) > 0 {
		elb.Listeners = &listeners
	}
	elb.Listeners = &listeners
	return elb
}

// Auto-scaling group for CF
func asgSection(deployment *k8.Deployment, ami string) (*gocf.AutoScalingAutoScalingGroup, *gocf.AutoScalingLaunchConfiguration) {
	replicas := strconv.Itoa(deployment.Spec.Replicas)
	asg := &gocf.AutoScalingAutoScalingGroup{
		HealthCheckGracePeriod:  gocf.Integer(120),
		LaunchConfigurationName: gocf.Ref("LaunchConfig").String(),
		LoadBalancerNames:       gocf.StringList(gocf.Ref("ELB")),
		MaxSize:                 gocf.String(replicas),
		MinSize:                 gocf.String(replicas),
		VPCZoneIdentifier:       gocf.StringList(gocf.String("subnet-7c505724")), // TODO configurable parameter
	}
	launchConfig := &gocf.AutoScalingLaunchConfiguration{
		ImageId:        gocf.String(ami),
		InstanceType:   gocf.String("t2.small"),
		KeyName:        gocf.String("pablo"),    // TODO configurable parameter
		SecurityGroups: []string{"sg-dfce36a5"}, // TODO configurable parameter
	}
	return asg, launchConfig
}

func elbURLOutput() *gocf.Output {
	return &gocf.Output{
		Description: "The URL of the stack",
		Value:       gocf.Join("", gocf.String("http://"), gocf.GetAtt("ELB", "DNSName")),
	}
}

func route53section(deployment *k8.Deployment) *gocf.Route53RecordSet {
	recordSet := &gocf.Route53RecordSet{
		HostedZoneName:  gocf.String(hostedZone),
		Name:            gocf.String(fmt.Sprintf("%s.%s", deployment.Metadata.Name, hostedZone)),
		Type:            gocf.String("CNAME"),
		TTL:             gocf.String("900"),
		ResourceRecords: gocf.StringList(gocf.GetAtt("ELB", "DNSName")),
	}
	return recordSet
}

// Translate a k8s yaml to a CloudFormation template
func Translate(deployment *k8.Deployment, ami string, pprint bool) ([]byte, error) {
	t := gocf.NewTemplate()
	t.AddResource("ELB", loadBalancerSection(deployment))
	asg, launchConfig := asgSection(deployment, ami)
	t.AddResource("ASG", asg)
	t.AddResource("LaunchConfig", launchConfig)
	t.AddResource("DNSRecord", route53section(deployment))
	t.Outputs["URL"] = elbURLOutput()
	if pprint {
		templateMarshalled, err := json.MarshalIndent(t, "", "    ")
		if err != nil {
			return []byte(""), err
		}
		return templateMarshalled, nil
	}
	templateMarshalled, err := json.Marshal(t)
	if err != nil {
		return []byte(""), err
	}
	return templateMarshalled, nil
}
