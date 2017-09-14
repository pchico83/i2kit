package k8

//Deployment represents a k8 deployment file
type Deployment struct {
	Metadata *MetadataType
	Spec     *DeploymentSpecType
}

//MetadataType represents k8 deployment metadata
type MetadataType struct {
	Name string
}

//DeploymentSpecType represents a k8 deployment spec
type DeploymentSpecType struct {
	Replicas int
	Template *TemplateType
}

//TemplateType represents a k8 deployment template
type TemplateType struct {
	Spec *ContainerSpecType
}

//ContainerSpecType represents a k8 container spec
type ContainerSpecType struct {
	Containers map[string]*Container
}

//Container represents a container in a k8 deployment
type Container struct {
	Name  string
	Image string
	Ports []*Port
}

//Port represents a container port
type Port struct {
	ContainerPort int
}

//Read returns a k8 deployment structure given a path to deployment.yml file
func Read(path string) (*Deployment, error) {
	dYml, err := readYml(path)
	if err != nil {
		return nil, err
	}
	d, err := createDeployment(dYml)
	if err != nil {
		return nil, err
	}
	return d, nil
}
