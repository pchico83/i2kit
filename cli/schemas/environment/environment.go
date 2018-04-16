package environment

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	//AWS is the amazon web services provider
	AWS = "aws"
	//K8 is the kubernetes provider
	K8 = "k8"
)

var isAlphaNumeric = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9]+$`).MatchString

//Environment represents a environment.yml file
type Environment struct {
	Name        string       `yaml:"name,omitempty"`
	DNSProvider *DNSProvider `yaml:"dns,omitempty"`
	Provider    *Provider    `yaml:"provider,omitempty"`
	Docker      *Docker      `yaml:"docker,omitempty"`
}

//DNSProvider represents the info for the cloud provider where the DNS is created
type DNSProvider struct {
	AccessKey    string `yaml:"access_key,omitempty"`
	SecretKey    string `yaml:"secret_key,omitempty"`
	HostedZone   string `yaml:"hosted_zone,omitempty"`
	HostedZoneID string `yaml:"hosted_zone_id,omitempty"`
}

//Provider represents the info for the cloud provider where the deployment takes place
type Provider struct {
	Type          string    `yaml:"type,omitempty"`
	Config        string    `yaml:"config,omitempty"`
	InstanceType  string    `yaml:"instance_type,omitempty"`
	Certificate   string    `yaml:"certificate,omitempty"`
	AccessKey     string    `yaml:"access_key,omitempty"`
	SecretKey     string    `yaml:"secret_key,omitempty"`
	Region        string    `yaml:"region,omitempty"`
	Subnets       []*string `yaml:"subnets,omitempty"`
	SecurityGroup string    `yaml:"security_group,omitempty"`
	Keypair       string    `yaml:"keypair,omitempty"`
	HostedZone    string    `yaml:"hosted_zone,omitempty"`
}

//Docker represents Docker Hub credentials
type Docker struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

//Validate returns an error for invalid environment.yml files
func (e *Environment) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("'environment.name' is mandatory")
	}
	if !isAlphaNumeric(e.Name) {
		return fmt.Errorf("'environment.name' only allows alphanumeric characters")
	}
	if e.Provider == nil {
		return nil
	}
	if err := e.Provider.Validate(); err != nil {
		return err
	}
	if e.DNSProvider == nil {
		if e.Provider.HostedZone == "" {
			return fmt.Errorf("'environment.provider.hosted_zone' must be defined if no dns provider is defined")
		}
		return nil
	}
	return e.DNSProvider.Validate()
}

//Domain returns the seacrh domain for a given environment
func (e *Environment) Domain() string {
	var domain string
	if e.Provider.HostedZone == "" {
		domain = strings.TrimSuffix(e.DNSProvider.HostedZone, ".")
	} else {
		domain = strings.TrimSuffix(e.Provider.HostedZone, ".")
	}
	return fmt.Sprintf("%s.%s", e.Name, domain)
}

//GetConfig returns a config aws object
func (p *Provider) GetConfig() *aws.Config {
	awsConfig := &aws.Config{
		Region:      aws.String(p.Region),
		Credentials: credentials.NewStaticCredentials(p.AccessKey, p.SecretKey, ""),
	}
	return awsConfig
}

//GetConfigFile returns a config k8 file
func (p *Provider) GetConfigFile() (*kubernetes.Clientset, string, error) {
	sDec, err := base64.StdEncoding.DecodeString(p.Config)
	if err != nil {
		return nil, "", fmt.Errorf("Error decoding k8 config: %s", err)
	}
	configFile, err := ioutil.TempFile("", "k8-config")
	if err != nil {
		return nil, "", fmt.Errorf("Error creating tmp file: %s", err)
	}
	err = ioutil.WriteFile(configFile.Name(), sDec, 0400)
	if err != nil {
		return nil, configFile.Name(), fmt.Errorf("Error writing to tmp file: %s", err)
	}
	config, err := clientcmd.BuildConfigFromFlags("", configFile.Name())
	if err != nil {
		return nil, configFile.Name(), fmt.Errorf("Error reading k8 config: %s", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, configFile.Name(), fmt.Errorf("Error creating k8 client: %s", err)
	}
	return clientset, configFile.Name(), nil
}

//GetType returns the type of a provider
func (p *Provider) GetType() string {
	return strings.ToLower(p.Type)
}

//Validate returns an error for invalid providers
func (p *Provider) Validate() error {
	switch p.GetType() {
	case "":
		return fmt.Errorf("'provider.type' cannot be empty")
	case AWS:
		if p.AccessKey == "" {
			return fmt.Errorf("'provider.access_key' cannot be empty")
		}
		if p.SecretKey == "" {
			return fmt.Errorf("'provider.secret_key' cannot be empty")
		}
		if p.Region == "" {
			return fmt.Errorf("'provider.region' cannot be empty")
		}
		if p.Subnets == nil || len(p.Subnets) == 0 {
			return fmt.Errorf("'provider.subnets' cannot be empty")
		}
		if p.SecurityGroup == "" {
			return fmt.Errorf("'provider.security_group' cannot be empty")
		}
		if p.Keypair == "" {
			return fmt.Errorf("'provider.keypair' cannot be empty")
		}
		return nil
	case K8:
		if p.Config == "" {
			return fmt.Errorf("'provider.config' cannot be empty")
		}
		_, err := base64.StdEncoding.DecodeString(p.Config)
		if err != nil {
			return fmt.Errorf("'provider.config' is not a valid base64 encoded string")
		}
		return nil
	default:
		return fmt.Errorf("'provider.type' '%s' is not supported", p.GetType())
	}
}

//GetConfig returns a config aws object
func (p *DNSProvider) GetConfig() *aws.Config {
	awsConfig := &aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials(p.AccessKey, p.SecretKey, ""),
	}
	return awsConfig
}

//Validate returns an error for invalid providers
func (p *DNSProvider) Validate() error {
	if p.AccessKey == "" {
		return fmt.Errorf("'provider.access_key' cannot be empty")
	}
	if p.SecretKey == "" {
		return fmt.Errorf("'provider.secret_key' cannot be empty")
	}
	if p.HostedZone == "" {
		return fmt.Errorf("'provider.hosted_zone' cannot be empty")
	}
	svc := route53.New(session.New(), p.GetConfig())
	hostedZonesInput := &route53.ListHostedZonesByNameInput{
		DNSName:  aws.String(p.HostedZone),
		MaxItems: aws.String("1"),
	}
	resp, err := svc.ListHostedZonesByName(hostedZonesInput)
	if err != nil {
		return err
	}
	if len(resp.HostedZones) != 1 {
		return fmt.Errorf("Hosted zone '%s' not found", p.HostedZone)
	}
	p.HostedZoneID = *resp.HostedZones[0].Id
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

//B64DockerConfig conputes the base64 format of docker credentials
func (e *Environment) B64DockerConfig() string {
	if e.Docker == nil || e.Docker.Username == "" || e.Docker.Password == "" {
		return ""
	}
	auth := fmt.Sprintf("%s:%s", e.Docker.Username, e.Docker.Password)
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))
	config := fmt.Sprintf(dockerConfigTemplate, authEncoded)
	return base64.StdEncoding.EncodeToString([]byte(config))
}
