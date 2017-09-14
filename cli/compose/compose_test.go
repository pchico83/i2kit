package compose

import (
	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"bitbucket.org/riberaproject/cli/service"
)

func TestE2E(t *testing.T) {
	service, err := service.Read("./examples/i2kit.yml")
	require.NoError(t, err)
	generatedEncodedCompose, err := Create(service, "staging.i2kit.com")
	composeBytes, err := ioutil.ReadFile("./examples/compose.yml")
	require.NoError(t, err)
	encoded := base64.StdEncoding.EncodeToString(composeBytes)
	if encoded != generatedEncodedCompose {
		decoded, err := base64.StdEncoding.DecodeString(generatedEncodedCompose)
		require.NoError(t, err)
		t.Fatal(string(decoded))
	}
}
