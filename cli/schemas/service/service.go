package service

import (
	"io"
)

//Service represents a service.yml file
type Service struct {
	Name       string
	Replicas   int
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
func Validate(reader io.Reader) error {
	return nil
}

//Read returns a Service structure given a reader to a service.yml file
func Read(reader io.Reader) (*Service, error) {
	sYml, err := readYml(reader)
	if err != nil {
		return nil, err
	}
	return create(sYml)
}
