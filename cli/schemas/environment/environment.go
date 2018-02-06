package environment

import (
	"fmt"
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

//Environment represents a environment.yml file
type Environment struct {
	Provider *Provider `yaml:"provider,omitempty"`
	Docker   *Docker   `yaml:"docker,omitempty"`
}

//Provider represents the info for the cloud provider where the deployment takes place
type Provider struct {
	AccessKey     string `yaml:"access_key,omitempty"`
	SecretKey     string `yaml:"secret_key,omitempty"`
	Region        string `yaml:"region,omitempty"`
	Subnet        string `yaml:"subnet,omitempty"`
	SecurityGroup string `yaml:"security_group,omitempty"`
	Keypair       string `yaml:"keypair,omitempty"`
	HostedZone    string `yaml:"hosted_zone,omitempty"`
}

//Docker represents Docker Hub credentials
type Docker struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

//Validate returns an error for invalid environment.yml files
func Validate(reader io.Reader) error {
	e, err := Read(reader)
	if err != nil {
		return err
	}
	if e.Provider == nil {
		return nil
	}

	if e.Provider.AccessKey == "" {
		return fmt.Errorf("'provider.access_key' cannot be empty")
	}
	if e.Provider.SecretKey == "" {
		return fmt.Errorf("'provider.secret_key' cannot be empty")
	}
	if e.Provider.Region == "" {
		return fmt.Errorf("'provider.region' cannot be empty")
	}
	if e.Provider.Subnet == "" {
		return fmt.Errorf("'provider.subnet' cannot be empty")
	}
	if e.Provider.SecurityGroup == "" {
		return fmt.Errorf("'provider.security_group' cannot be empty")
	}
	if e.Provider.Keypair == "" {
		return fmt.Errorf("'provider.keypair' cannot be empty")
	}
	if e.Provider.HostedZone == "" {
		return fmt.Errorf("'provider.hosted_zone' cannot be empty")
	}
	return nil
}

//Read returns a Environment structure given a reader to a env.yml file
func Read(reader io.Reader) (*Environment, error) {
	readBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var result Environment
	err = yaml.Unmarshal(readBytes, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
