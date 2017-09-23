package linuxkit

import (
	"io/ioutil"

	"github.com/moby/tool/src/moby"
	v1beta1 "k8s.io/api/extensions/v1beta1"
)

//GetTemplate generates a linuxkit template from a k8 deployment object
func GetTemplate(deployment *v1beta1.Deployment) (*moby.Moby, error) {
	configBytes, err := ioutil.ReadFile("./linuxkit/aws.yml")
	if err != nil {
		return nil, err
	}
	mobyConfig, err := moby.NewConfig(configBytes)
	if err != nil {
		return nil, err
	}
	for _, container := range deployment.Spec.Template.Spec.Containers {
		mobyConfig.Services = append(
			mobyConfig.Services,
			moby.Image{
				Name:  container.Name,
				Image: container.Image,
			},
		)
	}
	return &mobyConfig, nil
}
