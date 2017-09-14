package stack

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type deploymentYml struct {
	Services map[string]*serviceYml
}

type serviceYml struct {
	Size       string
	Min        int
	Max        int
	Ports      []string
	Links      []string
	Containers map[string]*containerYml
}

type containerYml struct {
	Image string
}

func readYml(path string) (*stackYml, error) {
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

func createDeployment(name string, s *deploymentYml) (*Deployment, error) {
	result := Deployment{
		Name:     name,
		Services: make(map[string]*Service),
	}
	for sName, sYml := range s.Services {
		result.Services[sName] = &Service{
			Name:       sName,
			Size:       sYml.Size,
			Min:        sYml.Min,
			Max:        sYml.Max,
			Links:      []*Link{},
			Ports:      []*Port{},
			Containers: make(map[string]*Container),
		}
		err := createLinks(result.Services[sName], sYml)
		if err != nil {
			return nil, err
		}
		err = createPorts(result.Services[sName], sYml)
		if err != nil {
			return nil, err
		}
		createContainers(result.Services[sName], sYml)
	}
	return &result, nil
}

func createLinks(s *Service, sYml *serviceYml) error {
	for _, linkYml := range sYml.Links {
		parts := strings.Split(linkYml, ":")
		link := &Link{}
		switch len(parts) {
		case 1:
			link = &Link{
				Alias:   parts[0],
				Service: parts[0],
			}
		case 2:
			link = &Link{
				Alias:   parts[0],
				Service: parts[1],
			}
		default:
			return fmt.Errorf("Links in service %s must be of the form (ALIAS:)SERVICE", s.Name)
		}
		s.Links = append(s.Links, link)
	}
	return nil
}

func createPorts(s *Service, sYml *serviceYml) error {
	for _, portYml := range sYml.Ports {
		parts := strings.Split(portYml, ":")
		port := &Port{}
		switch len(parts) {
		case 3:
			number, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("Port numbers in service %s must be integers", s.Name)
			}
			port = &Port{
				Container: parts[0],
				Number:    number,
				Protocol:  parts[2],
			}
		default:
			return fmt.Errorf("Ports in service %s must be of the form CONTAINER:NUMBER:PROTOCOL", s.Name)
		}
		s.Ports = append(s.Ports, port)
	}
	return nil
}

func createContainers(s *Service, sYml *serviceYml) {
	for cName, cYml := range sYml.Containers {
		s.Containers[cName] = &Container{
			Name:  cName,
			Image: cYml.Image,
		}
	}
}
