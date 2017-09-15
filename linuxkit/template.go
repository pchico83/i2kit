package linuxkit

//Template represents a linuxkit template
type Template struct {
	Kernel   *KernelType
	Init     []*string
	Onboot   []*ContainerType
	Services []*ContainerType
	Trust    *TrustType
}

//KernelType represents the kernel field
type KernelType struct {
	Image   string
	Cmdline string
}

//ContainerType represents the kernel field
type ContainerType struct {
	Name    string
	Image   string
	Command []*string
}

//TrustType represents the trus field
type TrustType struct {
	Org []*string
}

//Read returns a linuxkit template given a path to a linuxkit.yml file
func Read(path string) (*Template, error) {
	tYml, err := readYml(path)
	if err != nil {
		return nil, err
	}
	t, err := createTemplate(tYml)
	if err != nil {
		return nil, err
	}
	return t, nil
}
