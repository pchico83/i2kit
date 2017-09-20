package linuxkit

import (
	"fmt"
	"io/ioutil"
	"os"

	v1beta1 "k8s.io/api/extensions/v1beta1"
)

//GetTemplate generates a linuxkit template from a k8 deployment object
func GetTemplate(deployment *v1beta1.Deployment) (string, error) {
	t, err := read("./aws.yml")
	if err != nil {
		return "", err
	}
	for _, container := range deployment.Spec.Template.Spec.Containers {
		t.Services = append(
			t.Services,
			&containerYml{
				Name:  container.Name,
				Image: container.Image,
			},
		)
	}
	file, err := ioutil.TempFile(
		os.TempDir(),
		fmt.Sprintf("%s-i2kit-", deployment.GetObjectMeta().GetName()),
	)
	if err != nil {
		return "", err
	}
	write(t, file.Name())
	return file.Name(), nil
}
