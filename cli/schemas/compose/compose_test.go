package compose

import (
	"encoding/base64"
	"io/ioutil"
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/stretchr/testify/require"

	"github.com/pchico83/i2kit/cli/schemas/service"
)

func TestE2E(t *testing.T) {
	readBytes, err := ioutil.ReadFile("./examples/service.yml")
	require.NoError(t, err)
	var s service.Service
	err = yaml.Unmarshal(readBytes, &s)
	require.NoError(t, err)
	generatedEncodedCompose, err := Create(&s, "staging.i2kit.com")
	composeBytes, err := ioutil.ReadFile("./examples/compose.yml")
	require.NoError(t, err)
	encoded := base64.StdEncoding.EncodeToString(composeBytes)
	if encoded != generatedEncodedCompose {
		decoded, err := base64.StdEncoding.DecodeString(generatedEncodedCompose)
		require.NoError(t, err)
		t.Fatal(string(decoded))
	}
}
