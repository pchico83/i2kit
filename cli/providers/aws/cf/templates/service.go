package templates

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	gocf "github.com/crewjam/go-cloudformation"
	"github.com/pchico83/i2kit/cli/schemas/compose"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

var amisPerRegion = map[string]string{
	"us-east-2":      "ami-1b90a67e",
	"us-east-1":      "ami-cb17d8b6",
	"us-west-2":      "ami-05b5277d",
	"us-west-1":      "ami-9cbbaffc",
	"eu-west-3":      "ami-914afcec",
	"eu-west-2":      "ami-a48d6bc3",
	"eu-west-1":      "ami-bfb5fec6",
	"eu-central-1":   "ami-ac055447",
	"ap-northeast-2": "ami-ba74d8d4",
	"ap-northeast-1": "ami-5add893c",
	"ap-southeast-2": "ami-4cc5072e",
	"ap-southeast-1": "ami-acbcefd0",
	"ca-central-1":   "ami-a535b2c1",
	"ap-south-1":     "ami-2149114e",
	"sa-east-1":      "ami-d3bce9bf",
}

// Service translates an i2kit service to a AWS CloudFormation template
func Service(s *service.Service, e *environment.Environment, config *aws.Config) (string, error) {
	t := gocf.NewTemplate()
	ami, ok := amisPerRegion[e.Provider.Region]
	if !ok {
		return "", fmt.Errorf("Region'%s' is not supported", e.Provider.Region)
	}
	e.Provider.Ami = ami
	encodedCompose, err := compose.Create(s, e)
	if err != nil {
		return "", err
	}
	if s.Stateful {
		if err = stateful(t, s, e, encodedCompose); err != nil {
			return "", err
		}
	} else {
		if err = stateless(t, s, e, encodedCompose); err != nil {
			return "", err
		}
	}
	logGroup(t, s, e)
	ports := s.GetPorts()
	if len(ports) > 0 && e.Provider.HostedZone != "" {
		route53(t, s, e)
	}
	marshalledTemplate, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(marshalledTemplate), nil
}
