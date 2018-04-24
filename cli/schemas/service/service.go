package service

import (
	"fmt"
	"regexp"

	"github.com/pchico83/i2kit/cli/schemas/environment"
)

var isAlphaNumeric = regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString

//Service represents a service.yml file
type Service struct {
	Name         string                `yaml:"name,omitempty"`
	Replicas     int                   `yaml:"replicas,omitempty"`
	Stateful     bool                  `yaml:"stateful,omitempty"`
	Public       bool                  `yaml:"public,omitempty"`
	InstanceType string                `yaml:"instance_type,omitempty"`
	Containers   map[string]*Container `yaml:"containers,omitempty"`
}

//Container represents a container in a service.yml file
type Container struct {
	Image       string    `yaml:"image,omitempty"`
	Command     string    `yaml:"command,omitempty"`
	Ports       []*Port   `yaml:"ports,omitempty"`
	Environment []*EnvVar `yaml:"environment,omitempty"`
}

//Port represents a container port
type Port struct {
	Certificate      string
	InstanceProtocol string
	InstancePort     string
	Protocol         string
	Port             string
}

//EnvVar represents a container envvar
type EnvVar struct {
	Name  string
	Value string
}

//Validate returns an error for invalid service.yml files
func (s *Service) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("'service.name' is mandatory")
	}
	if !isAlphaNumeric(s.Name) {
		return fmt.Errorf("'service.name' only allows alphanumeric characters")
	}
	if s.Stateful && s.Replicas != 1 {
		return fmt.Errorf("Stateful services can only have one replica")
	}
	return nil
}

//GetFullName returns the service full name taking into account the environment name
func (s *Service) GetFullName(e *environment.Environment, sep string) string {
	return fmt.Sprintf("%s%s%s", s.Name, sep, e.Name)
}

//GetInstanceType returns the service size taking into account default values
func (s *Service) GetInstanceType(e *environment.Environment) string {
	if s.InstanceType != "" {
		return s.InstanceType
	}
	if e.Provider != nil && e.Provider.InstanceType != "" {
		return e.Provider.InstanceType
	}
	return "t2.small"
}

//GetPorts returns the list of ports of a service
func (s *Service) GetPorts() []*Port {
	result := []*Port{}
	for _, container := range s.Containers {
		for _, port := range container.Ports {
			result = append(result, port)
		}
	}
	return result
}
