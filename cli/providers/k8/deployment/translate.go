package deployment

import (
	"strconv"

	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Translate an i2kit service to a k8 deployment manifest
func Translate(s *service.Service, e *environment.Environment) *appsv1beta1.Deployment {
	deploymentName := s.Name
	// deploymentName := s.GetFullName(e, "-")
	// domain := e.Domain()
	replicas := int32(s.Replicas)
	containers := []apiv1.Container{}
	for name, c := range s.Containers {
		ports := []apiv1.ContainerPort{}
		for _, p := range c.Ports {
			portInt64, _ := strconv.ParseInt(p.InstancePort, 10, 32)
			ports = append(ports, apiv1.ContainerPort{
				Protocol:      apiv1.ProtocolTCP,
				ContainerPort: int32(portInt64),
			})
		}
		envs := []apiv1.EnvVar{}
		for _, e := range c.Environment {
			envs = append(envs, apiv1.EnvVar{
				Name:  e.Name,
				Value: e.Value,
			})
		}
		command := []string{}
		if c.Command != "" {
			command = append(command, c.Command)
		}
		containers = append(containers, apiv1.Container{
			Name:    name,
			Image:   c.Image,
			Ports:   ports,
			Env:     envs,
			Command: command,
		})
	}
	return &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: &replicas,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deploymentName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: containers,
					// DNSConfig: &apiv1.PodDNSConfig{
					// 	Searches: []string{domain},
					// },
				},
			},
		},
	}
}
