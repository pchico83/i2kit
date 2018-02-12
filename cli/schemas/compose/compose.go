package compose

import (
	"encoding/base64"
	"fmt"

	"github.com/pchico83/i2kit/cli/schemas/service"

	yaml "gopkg.in/yaml.v2"
)

//Compose represents a docker-compose.yml file
type Compose struct {
	Version  string
	Services map[string]*Service
}

//Service represents a service in a docker-compose.yml file
type Service struct {
	Image       string
	Command     string    `yaml:"command,omitempty"`
	Ports       []*string `yaml:"ports,omitempty"`
	Environment []*string `yaml:"environment,omitempty"`
	NetworkMode string    `yaml:"network_mode,omitempty"`
	DNSSearch   []*string `yaml:"dns_search,omitempty"`
}

//Create returns a compose base64 encoded given a service object
func Create(s *service.Service, domain string) (string, error) {
	compose := &Compose{
		Version:  "3.4",
		Services: make(map[string]*Service),
	}
	for cName, c := range s.Containers {
		compose.Services[cName] = &Service{
			Image:       c.Image,
			Command:     c.Command,
			Ports:       []*string{},
			Environment: []*string{},
			NetworkMode: "bridge",
			DNSSearch:   []*string{&domain},
		}
		for _, p := range c.Ports {
			composePort := fmt.Sprintf("%s:%s", p.InstancePort, p.InstancePort)
			compose.Services[cName].Ports = append(
				compose.Services[cName].Ports,
				&composePort,
			)
		}
		for _, e := range c.Environment {
			composeEnvVar := fmt.Sprintf("%s=%s", e.Name, e.Value)
			compose.Services[cName].Environment = append(
				compose.Services[cName].Environment,
				&composeEnvVar,
			)
		}
	}
	composeBytes, err := yaml.Marshal(compose)
	if err != nil {
		return "", err
	}
	composeEncoded := base64.StdEncoding.EncodeToString(composeBytes)
	return composeEncoded, nil
}
