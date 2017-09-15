package linuxkit

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type templateYml struct {
	Kernel   *kernelYml
	Init     []*string
	Onboot   []*containerYml
	Services []*containerYml
	Trust    *trustYml
}

type kernelYml struct {
	Image   string
	Cmdline string
}

type containerYml struct {
	Name    string
	Image   string
	Command []*string
}

type trustYml struct {
	Org []*string
}

func readYml(path string) (*templateYml, error) {
	ymlBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var result templateYml
	err = yaml.Unmarshal(ymlBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func createTemplate(t *templateYml) (*Template, error) {
	result := Template{
		Kernel:   &KernelType{Image: t.Kernel.Image, Cmdline: t.Kernel.Cmdline},
		Init:     []*string{},
		Onboot:   []*ContainerType{},
		Services: []*ContainerType{},
		Trust:    &TrustType{Org: []*string{}},
	}
	for _, init := range t.Init {
		result.Init = append(result.Init, init)
	}
	for _, onboot := range t.Onboot {
		result.Onboot = append(
			result.Onboot,
			&ContainerType{
				Name:    onboot.Name,
				Image:   onboot.Image,
				Command: onboot.Command,
			},
		)
	}
	for _, service := range t.Services {
		result.Services = append(
			result.Services,
			&ContainerType{
				Name:    service.Name,
				Image:   service.Image,
				Command: service.Command,
			},
		)
	}
	for _, org := range t.Trust.Org {
		result.Trust.Org = append(result.Trust.Org, org)
	}
	return &result, nil
}
