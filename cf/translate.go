package cf

import (
	"encoding/json"
	"fmt"

	gocf "github.com/crewjam/go-cloudformation"
	"k8s.io/client-go/pkg/runtime"
)

// ELB for CF
func loadBalancerSection(k8sDeployment runtime.Object) gocf.ElasticLoadBalancingLoadBalancer {
	fmt.Printf("Ojbect %v", k8sDeployment)
	elb := gocf.ElasticLoadBalancingLoadBalancer{
		LoadBalancerName: gocf.String("testing-i2kit"),
		Subnets:          gocf.StringList(gocf.String("hello")),
	}
	/*
		listeners := gocf.ElasticLoadBalancingListenerList{}
		for _, container := range k8sDeployment["containers"] {
			if _, exists := container["ports"]; exists {
				for _, port := range container["ports"] {
					listeners = append(listeners, gocf.ElasticLoadBalancingListener{
						InstancePort:     port["containerPort"],
						InstanceProtocol: gocf.String("HTTP"),
						LoadBalancerPort: port["containerPort"],
						Protocol:         gocf.String("HTTP"),
					})
				}
			}
		}
		if len(listeners) > 0 {
			elb.Listeners = &listeners
		}
	*/
	return elb
}

// Auto-scaling group for CF
func asgSection(k8sDeployment runtime.Object, ami string) (*gocf.AutoScalingAutoScalingGroup, *gocf.AutoScalingLaunchConfiguration) {
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

// Translates a k8s yaml to a CloudFormation template
func Translate(k8sDeployment runtime.Object, ami string) ([]byte, error) {
	k8sDeployment.GetObjectKind()
	t := gocf.NewTemplate()
	t.AddResource("ELB", loadBalancerSection(k8sDeployment))
	asg, launchConfig := asgSection(k8sDeployment, ami)
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
