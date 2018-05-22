package compose

import (
	"encoding/base64"
	"fmt"

	"github.com/pchico83/i2kit/cli/schemas/environment"
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
	Command     string         `yaml:"command,omitempty"`
	Ports       []*string      `yaml:"ports,omitempty"`
	Environment []*string      `yaml:"environment,omitempty"`
	Logging     *LoggingDriver `yaml:"logging,omitempty"`
	Restart     string         `yaml:"restart,omitempty"`
	DNSSearch   []*string      `yaml:"dns_search,omitempty"`
}

//LoggingDriver represents a docker logging driver
type LoggingDriver struct {
	Driver  string            `yaml:"driver,omitempty"`
	Options map[string]string `yaml:"options,omitempty"`
}

//Create returns a compose base64 encoded given a service object
func Create(s *service.Service, e *environment.Environment) (string, error) {
	domain := e.Domain()
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
			Restart:     "on-failure",
			DNSSearch:   []*string{&domain},
		}

		compose.Services[cName].Ports = parsePorts(s.Stateful, c.Ports)

		for _, env := range c.Environment {
			if env.Value == "" {
				for _, secret := range e.Secrets {
					if env.Name == secret.Name {
						env.Value = secret.Value
					}
				}
			}
			composeEnvVar := fmt.Sprintf("%s=%s", env.Name, env.Value)
			compose.Services[cName].Environment = append(
				compose.Services[cName].Environment,
				&composeEnvVar,
			)
		}
		compose.Services[cName].Logging = &LoggingDriver{
			Driver: "awslogs",
			Options: map[string]string{
				"awslogs-region": e.Provider.Region,
				"awslogs-group":  fmt.Sprintf("i2kit-%s", s.GetFullName(e, "-")),
				"tag":            fmt.Sprintf("%s-${INSTANCE_ID}", cName),
			},
		}
	}
	composeBytes, err := yaml.Marshal(compose)
	if err != nil {
		return "", err
	}
	composeEncoded := base64.StdEncoding.EncodeToString(composeBytes)
	return composeEncoded, nil
}

func parsePorts(stateful bool, ports []*service.Port) []*string {
	parsedPorts := make([]*string, 0)

	for _, p := range ports {
		var composePort string
		if stateful {
			composePort = fmt.Sprintf("%s:%s", p.Port, p.InstancePort)
		} else {
			composePort = fmt.Sprintf("%s:%s", p.InstancePort, p.InstancePort)
		}
		parsedPorts = append(parsedPorts, &composePort)
	}

	removeDuplicates(&parsedPorts)
	return parsedPorts
}

func removeDuplicates(xs *[]*string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[*x] {
			found[*x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}
