package service

import (
	"github.com/pchico83/i2kit/cli/schemas/environment"
)

//Service represents a service.yml file
type Service struct {
	Name       string
	Replicas   int
	Size       string
	Containers map[string]*Container
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
	return nil
}

//GetSize returns the service size taking into account default values
func (s *Service) GetSize(e *environment.Environment) string {
	if s.Size != "" {
		return s.Size
	}
	if e.Provider != nil && e.Provider.Size != "" {
		return e.Provider.Size
	}
	return "t2.nano"
}
