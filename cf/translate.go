package cf

import (
	"encoding/json"

	gocf "github.com/crewjam/go-cloudformation"
	"github.com/pchico83/i2kit/k8"
)

// ELB for CF
func loadBalancerSection(deployment k8.Deployment) gocf.ElasticLoadBalancingLoadBalancer {
	elb := gocf.ElasticLoadBalancingLoadBalancer{
		LoadBalancerName: gocf.String("testing-i2kit"),
		Subnets:          gocf.StringList(gocf.String("hello")),
	}
	listeners := gocf.ElasticLoadBalancingListenerList{}
	for _, container := range deployment.Spec.Template.Spec.Containers {
		for _, port := range container.Ports {
			listeners = append(listeners, gocf.ElasticLoadBalancingListener{
				InstancePort:     port.ContainerPort,
				InstanceProtocol: gocf.String("HTTP"),
				LoadBalancerPort: port.ContainerPort,
				Protocol:         gocf.String("HTTP"),
			})
		}
	}
	if len(listeners) > 0 {
		elb.Listeners = &listeners
	}
	return elb
}

// Auto-scaling group for CF
func asgSection(deployment k8.Deployment, ami string) (*gocf.AutoScalingAutoScalingGroup, *gocf.AutoScalingLaunchConfiguration) {
	// TODO parse number of instances from k8sFile
	instances := "1"
	asg := &gocf.AutoScalingAutoScalingGroup{
		HealthCheckGracePeriod:  gocf.Integer(120),
		LaunchConfigurationName: gocf.String(`{ "Ref" : "LaunchConfig" }`),
		LoadBalancerNames:       gocf.StringList(gocf.String(`{ "Ref" : "ELB" }`)),
		MaxSize:                 gocf.String(instances),
		MinSize:                 gocf.String(instances),
		VPCZoneIdentifier:       gocf.StringList(gocf.String(`["subnet-3f087e57"]`)),
	}
	launchConfig := &gocf.AutoScalingLaunchConfiguration{
		ImageId:        gocf.String(ami),
		InstanceType:   gocf.String("t2.micro"),
		KeyName:        gocf.String("pablo"),
		SecurityGroups: []string{"sg-5d42b836"},
	}
	return asg, launchConfig
}

// Translate a k8s yaml to a CloudFormation template
func Translate(deployment k8.Deployment, ami string) ([]byte, error) {
	t := gocf.NewTemplate()
	t.AddResource("ELB", loadBalancerSection(deployment))
	asg, launchConfig := asgSection(deployment, ami)
	t.AddResource("ASG", asg)
	t.AddResource("LaunchConfig", launchConfig)
	t.Outputs["URL"] = &gocf.Output{
		Description: "The URL of the stack",
		Value:       `{ "Fn::Join" : [ "", [ "http://", { "Fn::GetAtt" : [ "ELB", "DNSName" ]}]]}`,
	}
	templateMarshalled, error := json.Marshal(t)
	if error != nil {
		return []byte(""), error
	}
	return templateMarshalled, nil
}
