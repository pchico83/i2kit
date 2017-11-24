package linuxkit

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type templateYml struct {
	Kernel   *kernelYml
	Init     []*string
	Onboot   []*containerYml
	Services []*containerYml
	Trust    *trustYml
}

type kernelYml struct {
	Image   string
	Cmdline string
}

type containerYml struct {
	Name    string
	Image   string
	Command []*string
}

type trustYml struct {
	Org []*string
}

func read(path string) (*templateYml, error) {
	ymlBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var result templateYml
	err = yaml.Unmarshal(ymlBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func write(t *templateYml, path string) error {
	templateBytes, err := yaml.Marshal(t)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, templateBytes, 0644)
}
