package linuxkit

import (
	"testing"

	"github.com/pchico83/i2kit/k8"
	"github.com/stretchr/testify/require"
)

func TestTranslateDeployment(t *testing.T) {
	deployment, err := k8.Read("../k8/templates/test.yml")
	require.NoError(t, err)
	_, err = GetTemplate(deployment)
	require.NoError(t, err)
}
