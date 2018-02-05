package compose

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pchico83/i2kit/cli/service"
)

func TestE2E(t *testing.T) {
	reader, err := os.Open("./examples/i2kit.yml")
	if err != nil {
		t.Fatalf(err.Error())
	}

	service, err := service.Read(reader)
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
