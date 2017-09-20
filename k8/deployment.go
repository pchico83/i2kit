package k8

import (
	"bytes"
	"io/ioutil"

	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//Read returns a k8 deployment structure given a path to deployment.yml file
func Read(path string) (*v1beta1.Deployment, error) {
	ymlBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var deployment v1beta1.Deployment
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(ymlBytes)), 1024)
	err = decoder.Decode(&deployment)
	if err != nil {
		return nil, err
	}
	return &deployment, nil
}
