package cf

import (
	"io/ioutil"
	"testing"

	"github.com/pchico83/i2kit/cf"
	"k8s.io/client-go/pkg/api"
)

func TestCreateCloudFormationTemplate(t *testing.T) {
	deploymentBytes, err := ioutil.ReadFile("../test/k8s-example")
	if err != nil {
		t.Fatalf("Error read file: %s", err)
	}
	decode := api.Codecs.UniversalDeserializer().Decode
	deployment, _, err := decode(deploymentBytes, nil, nil)
	if err != nil {
		t.Fatalf("Error decode: %s", err)
	}
	res, _ := cf.Translate(deployment, "testing-ami")
	t.Log(res)
}
