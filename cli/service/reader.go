package service

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type serviceYml struct {
	Name       string
	Replicas   int
	Containers map[string]*containerYml
}

type containerYml struct {
	Image       string
	Command     string
	Ports       []*string
	Environment []*string
}

func readYml(path string) (*serviceYml, error) {
	ymlBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var result serviceYml
	err = yaml.Unmarshal(ymlBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func createService(s *serviceYml) (*Service, error) {
	result := Service{
		Name:       s.Name,
		Replicas:   s.Replicas,
		Containers: make(map[string]*Container),
	}
	for name, cYml := range s.Containers {
		c := &Container{
			Image:       cYml.Image,
			Command:     cYml.Command,
			Ports:       []*Port{},
			Environment: []*EnvVar{},
		}
		result.Containers[name] = c
		if cYml.Ports != nil {
			extractPorts, err := extractPorts(name, cYml.Ports)
			if err != nil {
				return nil, err
			}

			c.Ports = extractPorts
		}

		if cYml.Environment != nil {
			for _, e := range cYml.Environment {
				parts := strings.Split(*e, "=")
				if len(parts) == 2 {
					c.Environment = append(
						c.Environment,
						&EnvVar{Name: parts[0], Value: parts[1]},
					)
				}
			}
		}
	}
	return &result, nil
}

func extractPorts(name string, ports []*string) ([]*Port, error) {
	var parsedPorts []*Port
	for _, p := range ports {
		parts := strings.SplitN(*p, ":", 5)
		switch len(parts) {
		case 4:
			if strings.ToUpper(parts[0]) != "HTTP" || strings.ToUpper(parts[2]) != "HTTP" {
				return nil, fmt.Errorf("Unsupported port protocol in container '%s'", name)
			}
			parsedPorts = append(
				parsedPorts,
				&Port{
					Protocol:         "HTTP",
					Port:             parts[1],
					InstanceProtocol: "HTTP",
					InstancePort:     parts[3],
				},
			)
		case 5:
			// AWS format: https:443:http:8000:arn:aws:acm:us-west-2:062762192540:certificate/12de3ac5-abcd-461a-1234-9e81250b33d8
			if strings.ToUpper(parts[0]) != "HTTPS" ||
				(strings.ToUpper(parts[2]) != "HTTP" && strings.ToUpper(parts[2]) != "HTTPS") {
				return nil, fmt.Errorf("Unsupported port protocol in container '%s'", name)
			}
			parsedPorts = append(
				parsedPorts,
				&Port{
					Protocol:         "HTTPS",
					Port:             parts[1],
					InstanceProtocol: strings.ToUpper(parts[2]),
					InstancePort:     parts[3],
					Certificate:      parts[4],
				},
			)
		default:
			return nil, fmt.Errorf("Wrong port syntax in container '%s'", name)
		}
	}

	return parsedPorts, nil
}
