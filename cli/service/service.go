package service

//Service represents a i2kit.yml file
type Service struct {
	Name       string
	Replicas   int
	Containers map[string]*Container
}

//Container represents a container in a i2kit.yml file
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

//Read returns a Service structure given a path to a i2kit.yml file
func Read(path string) (*Service, error) {
	sYml, err := readYml(path)
	if err != nil {
		return nil, err
	}
	s, err := createService(sYml)
	if err != nil {
		return nil, err
	}
	return s, nil
}
