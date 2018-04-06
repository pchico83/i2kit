package deployment

import (
	"fmt"
	"strings"
	"time"

	logger "log"

	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//Deploy deploys a i2kit service as k8 deployment
func Deploy(s *service.Service, e *environment.Environment, c *kubernetes.Clientset, log *logger.Logger) error {
	deploymentName := s.Name
	// deploymentName := s.GetFullName(e, "-")
	dClient := c.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

	d, err := dClient.Get(deploymentName, metav1.GetOptions{})
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return fmt.Errorf("Error getting k8 deployment: %s", err)
	}

	if d.Name == "" {
		log.Printf("Creating deployment '%s'...", deploymentName)
		d = Translate(s, e)
		_, err = dClient.Create(d)
		if err != nil {
			return fmt.Errorf("Error creating k8 deployment: %s", err)
		}
		log.Printf("Created deployment %s", deploymentName)
	} else {
		log.Printf("Updating deployment '%s'...", deploymentName)
		d = Translate(s, e)
		_, err = dClient.Update(d)
		if err != nil {
			return fmt.Errorf("Error updating k8 deployment: %s", err)
		}
	}

	log.Printf("Waiting for the deployment '%s' to be ready...", deploymentName)
	tries := 0
	for tries < 60 {
		tries++
		time.Sleep(5 * time.Second)
		d, err = dClient.Get(deploymentName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("Error getting k8 deployment: %s", err)
		}
		if d.Status.ReadyReplicas == int32(s.Replicas) && d.Status.UpdatedReplicas == int32(s.Replicas) {
			log.Printf("k8 deployment '%s' is ready", deploymentName)
			return nil
		}
	}
	return fmt.Errorf("k8 deployment not ready after 300s")
}

//Destroy destroys the k8 deployment created by a i2kit service
func Destroy(s *service.Service, e *environment.Environment, c *kubernetes.Clientset, log *logger.Logger) error {
	deploymentName := s.Name
	// deploymentName := s.GetFullName(e, "-")
	log.Printf("Deleting deployment '%s'...", deploymentName)
	dClient := c.AppsV1beta1().Deployments(apiv1.NamespaceDefault)
	deletePolicy := metav1.DeletePropagationForeground
	err := dClient.Delete(deploymentName, &metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}
		return fmt.Errorf("Error getting k8 deployment: %s", err)
	}

	log.Printf("Waiting for the deployment '%s' to be deleted...", deploymentName)
	tries := 0
	for tries < 10 {
		tries++
		time.Sleep(5 * time.Second)
		_, err := dClient.Get(deploymentName, metav1.GetOptions{})
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return nil
			}
			return fmt.Errorf("Error getting k8 deployment: %s", err)
		}
	}
	return fmt.Errorf("k8 service not deleted after 50s")
}
