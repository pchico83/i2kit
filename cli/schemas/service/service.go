package service

import (
	"fmt"

	"github.com/pchico83/i2kit/cli/schemas/environment"
)

//Service represents a service.yml file
type Service struct {
	Name         string
	Replicas     int
	Stateful     bool
	InstanceType string
	Containers   map[string]*Container
}

//Container represents a container in a service.yml file
type Container struct {
	Image       string
	Command     string
	Ports       []*Port
	Environment []*EnvVar
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
	if s.Stateful && s.Replicas != 1 {
		return fmt.Errorf("Stateful services can only have one replica")
	}
	return nil
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
