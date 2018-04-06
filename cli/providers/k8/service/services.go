package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	logger "log"

	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//Deploy deploys a i2kit service as k8 service
func Deploy(s *service.Service, e *environment.Environment, c *kubernetes.Clientset, log *logger.Logger) error {
	serviceName := s.GetFullName(e, "-")
	k8Service, err := c.CoreV1().Services("default").Get(serviceName, metav1.GetOptions{})
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return fmt.Errorf("Error getting k8 service: %s", err)
	}
	if k8Service.Name == "" {
		log.Printf("Creating service '%s'...", serviceName)
		k8Service = Translate(s, e)
		_, err = c.CoreV1().Services("default").Create(k8Service)
		if err != nil {
			return fmt.Errorf("Error creating k8 service: %s", err)
		}
		log.Printf("Created service '%s'", serviceName)
	} else {
		log.Printf("Updating service '%s'...", serviceName)
		k8Service.Spec.Ports = GetPorts(s, e)
		_, err = c.CoreV1().Services("default").Update(k8Service)
		if err != nil {
			return fmt.Errorf("Error updating k8 service: %s", err)
		}
		log.Printf("Updated service '%s'", serviceName)
	}

	// TODO: Wait for load balancer to be ready by pinging
	log.Print("Waiting for Pods to register in the external Load Balancer...")
	time.Sleep(45 * time.Second)
	log.Print("Pods registered in the external Load Balancer...")
	return nil
}

//GetEndpoint returns the endpoint of a given k8 service
func GetEndpoint(s *service.Service, e *environment.Environment) (string, error) {
	serviceName := s.GetFullName(e, "-")
	c, tmpfile, err := e.Provider.GetConfigFile()
	if tmpfile != "" {
		defer os.Remove(tmpfile)
	}
	if err != nil {
		return "", err
	}

	tries := 0
	for tries < 10 {
		tries++
		time.Sleep(5 * time.Second)
		k8Service, err := c.CoreV1().Services("default").Get(serviceName, metav1.GetOptions{})
		if err != nil {
			return "", fmt.Errorf("Error getting k8 service: %s", err)
		}
		if len(k8Service.Status.LoadBalancer.Ingress) > 0 {
			return k8Service.Status.LoadBalancer.Ingress[0].Hostname, nil
		}
	}
	return "", fmt.Errorf("External load balancer not created after 50s")
}

//Destroy destroys the k8 service created by a i2kit service
func Destroy(s *service.Service, e *environment.Environment, c *kubernetes.Clientset, log *logger.Logger) error {
	serviceName := s.GetFullName(e, "-")
	log.Printf("Deleting service '%s'...", serviceName)
	err := c.CoreV1().Services("default").Delete(serviceName, &metav1.DeleteOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}
		return fmt.Errorf("Error getting k8 deployment: %s", err)
	}

	log.Printf("Waiting for the service '%s' to be deleted...", serviceName)
	tries := 0
	for tries < 10 {
		tries++
		time.Sleep(5 * time.Second)
		_, err := c.CoreV1().Services("default").Get(serviceName, metav1.GetOptions{})
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return nil
			}
			return fmt.Errorf("Error getting k8 deployment: %s", err)
		}
	}
	return fmt.Errorf("k8 service not deleted after 50s")
}
