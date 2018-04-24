package compose

import (
	"encoding/base64"
	"io/ioutil"
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/stretchr/testify/require"

	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

func TestEncodedCompose(t *testing.T) {
	readBytes, err := ioutil.ReadFile("./examples/service.yml")
	require.NoError(t, err)
	var s service.Service
	err = yaml.Unmarshal(readBytes, &s)
	require.NoError(t, err)
	e := &environment.Environment{
		Name: "staging",
		Provider: &environment.Provider{
			Region:     "us-west-2",
			HostedZone: "i2kit.com.",
		},
	}
	generatedEncodedCompose, err := Create(&s, e)
	require.NoError(t, err)
	composeBytes, err := ioutil.ReadFile("./examples/compose.yml")
	require.NoError(t, err)
	encoded := base64.StdEncoding.EncodeToString(composeBytes)
	if encoded != generatedEncodedCompose {
		decoded, err := base64.StdEncoding.DecodeString(generatedEncodedCompose)
		require.NoError(t, err)
		t.Fatal(string(decoded))
	}
}

func Test_parsePorts(t *testing.T) {
	tests := []struct {
		name     string
		stateful bool
		ports    []*service.Port
		expected []string
	}{
		{name: "single-http", stateful: false, ports: []*service.Port{&service.Port{Port: "80", InstancePort: "80"}}, expected: []string{"80:80"}},
		{name: "single-https", stateful: false, ports: []*service.Port{&service.Port{Port: "443", InstancePort: "80"}}, expected: []string{"80:80"}},
		{name: "single-stateful", stateful: true, ports: []*service.Port{&service.Port{Port: "54320", InstancePort: "5432"}}, expected: []string{"54320:5432"}},
		{name: "duplicated", stateful: false, ports: []*service.Port{&service.Port{Port: "80", InstancePort: "80"}, &service.Port{Port: "80", InstancePort: "80"}, &service.Port{Port: "80", InstancePort: "80"}}, expected: []string{"80:80"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parsePorts(tt.stateful, tt.ports)
			if len(result) != len(tt.expected) {
				t.Errorf("didn't got the expected results: %+v", result)
				return
			}
			for i := range result {
				if *result[i] != tt.expected[i] {
					t.Errorf("parsePorts(): got %s; expected: %s", *result[i], tt.expected[i])
				}
			}
		})
	}
}
