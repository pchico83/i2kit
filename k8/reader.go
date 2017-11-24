package k8

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type deploymentYml struct {
	Metadata *metadataYml
	Spec     *deploymentSpecYml
}

type metadataYml struct {
	Name string
}

type deploymentSpecYml struct {
	Replicas int
	Template *templateYml
}

type templateYml struct {
	Spec *containerSpecYml
}

type containerSpecYml struct {
	Containers []*containerYml
}

type containerYml struct {
	Name  string
	Image string
	Ports []*portYml
}

type portYml struct {
	ContainerPort int `yaml:"containerPort"`
}

func readYml(path string) (*deploymentYml, error) {
	ymlBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var result deploymentYml
	err = yaml.Unmarshal(ymlBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func createDeployment(d *deploymentYml) (*Deployment, error) {
	result := Deployment{
		Metadata: &MetadataType{Name: d.Metadata.Name},
		Spec: &DeploymentSpecType{
			Replicas: d.Spec.Replicas,
			Template: &TemplateType{
				Spec: &ContainerSpecType{Containers: make(map[string]*Container)},
			},
		},
	}
	if d.Spec.Template.Spec.Containers == nil {
		return &result, nil
	}
	for _, cYml := range d.Spec.Template.Spec.Containers {
		c := &Container{
			Name:  cYml.Name,
			Image: cYml.Image,
			Ports: []*Port{},
		}
		result.Spec.Template.Spec.Containers[cYml.Name] = c
		if cYml.Ports == nil {
			continue
		}
		for _, pYml := range cYml.Ports {
			c.Ports = append(c.Ports, &Port{ContainerPort: pYml.ContainerPort})
		}
	}
	return &result, nil
}
