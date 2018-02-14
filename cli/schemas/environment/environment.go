package environment

import (
	"encoding/base64"
	"fmt"
)

//Environment represents a environment.yml file
type Environment struct {
	Provider *Provider `yaml:"provider,omitempty"`
	Docker   *Docker   `yaml:"docker,omitempty"`
}

//Provider represents the info for the cloud provider where the deployment takes place
type Provider struct {
	Size          string `yaml:"size,omitempty"`
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
func (e *Environment) Validate() error {
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

var dockerConfigTemplate = `
{
	"auths": {
		"https://index.docker.io/v1/": {
			"auth": "%s"
		}
	}
}
`

func (e *Environment) B64DockerConfig() string {
	if e.Docker == nil || e.Docker.Username == "" || e.Docker.Password == "" {
		return ""
	}
	auth := fmt.Sprintf("%s:%s", e.Docker.Username, e.Docker.Password)
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))
	config := fmt.Sprintf(dockerConfigTemplate, authEncoded)
	return base64.StdEncoding.EncodeToString([]byte(config))
}
