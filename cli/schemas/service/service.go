package service

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
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
	readBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var result Service
	err = yaml.Unmarshal(readBytes, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

//MarshalYAML serializes p provided into a YAML document. The return value is a string.
//It will fail if p has missing required values
func (p *Port) MarshalYAML() (interface{}, error) {
	if p.Protocol == "" || p.Port == "" || p.InstanceProtocol == "" || p.InstancePort == "" {
		return "", fmt.Errorf("missing values")
	}

	var buffer bytes.Buffer
	buffer.WriteString(strings.ToLower(p.Protocol))
	buffer.WriteString(":")
	buffer.WriteString(p.Port)
	buffer.WriteString(":")
	buffer.WriteString(strings.ToLower(p.InstanceProtocol))
	buffer.WriteString(":")
	buffer.WriteString(p.InstancePort)
	if p.Certificate != "" {
		buffer.WriteString(":")
		buffer.WriteString(p.Certificate)
	}

	return buffer.String(), nil
}

//UnmarshalYAML parses the yaml element and sets the values of p; it will return an error if the parsing fails, or
//if the format is incorrect
func (p *Port) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var port string
	if err := unmarshal(&port); err != nil {
		return err
	}

	parts := strings.SplitN(port, ":", 5)
	switch len(parts) {
	case 4:
		if strings.ToUpper(parts[0]) != "HTTP" || strings.ToUpper(parts[2]) != "HTTP" {
			return fmt.Errorf("unsupported port protocol")
		}

		p.Protocol = "HTTP"
		p.Port = parts[1]
		p.InstanceProtocol = "HTTP"
		p.InstancePort = parts[3]

	case 5:
		// AWS format: https:443:http:8000:arn:aws:acm:us-west-2:062762192540:certificate/12de3ac5-abcd-461a-1234-9e81250b33d8
		if strings.ToUpper(parts[0]) != "HTTPS" ||
			(strings.ToUpper(parts[2]) != "HTTP" && strings.ToUpper(parts[2]) != "HTTPS") {
			return fmt.Errorf("unsupported port protocol")
		}

		p.Protocol = "HTTPS"
		p.Port = parts[1]
		p.InstanceProtocol = strings.ToUpper(parts[2])
		p.InstancePort = parts[3]
		p.Certificate = parts[4]
	default:
		return fmt.Errorf("invalid port syntax")
	}

	return nil
}

//MarshalYAML serializes e into a YAML document. The return value is a string; It will fail if e has an empty name.
func (e *EnvVar) MarshalYAML() (interface{}, error) {
	if e.Name == "" {
		return "", fmt.Errorf("missing values")
	}

	var buffer bytes.Buffer
	buffer.WriteString(e.Name)
	buffer.WriteString("=")
	if e.Value != "" {
		buffer.WriteString(e.Value)
	}

	return buffer.String(), nil
}

//UnmarshalYAML parses the yaml element and sets the values of e; it will return an error if the parsing fails, or
//if the format is incorrect
func (e *EnvVar) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var envvar string
	if err := unmarshal(&envvar); err != nil {
		return err
	}

	envvar = strings.TrimPrefix(envvar, "=")

	parts := strings.Split(envvar, "=")
	if len(parts) != 2 {
		return fmt.Errorf("Invalid environment variable syntax")
	}

	e.Name = parts[0]
	e.Value = parts[1]
	return nil
}

func (s *Service) String() string {
	yamlBytes, err := yaml.Marshal(s)
	if err != nil {
		return "service-error"
	}

	return string(yamlBytes[:])
}
