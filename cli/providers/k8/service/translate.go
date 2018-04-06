package service

import (
	"strconv"

	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

//Translate an i2kit service to a k8 service manifest
func Translate(s *service.Service, e *environment.Environment) *apiv1.Service {
	serviceName := s.GetFullName(e, "-")
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{"app": serviceName},
			Type:     apiv1.ServiceTypeLoadBalancer,
			Ports:    GetPorts(s, e),
		},
	}
}

//GetPorts returns rhe k8 ports of an i2kit service
func GetPorts(s *service.Service, e *environment.Environment) []apiv1.ServicePort {
	ports := []apiv1.ServicePort{}
	for _, port := range s.GetPorts() {
		portInt64, _ := strconv.ParseInt(port.Port, 10, 32)
		ports = append(ports, apiv1.ServicePort{
			Port:       int32(portInt64),
			TargetPort: intstr.IntOrString{StrVal: port.InstancePort},
		})
	}
	return ports

}
