package k8

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDeployment(t *testing.T) {
	_, err := Read("./templates/test.yml")
	require.NoError(t, err)
}
