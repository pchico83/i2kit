package cf

import (
	"fmt"
	"testing"

	"github.com/pchico83/i2kit/k8"

	"github.com/stretchr/testify/require"
)

func TestTranslateK8sDeploymentToCF(t *testing.T) {
	k8sInfo, err := k8.Read("../k8/templates/test.yml")
	require.NoError(t, err)
	template, err := Translate(k8sInfo, "test-ami", true)
	require.NoError(t, err)
	fmt.Print(string(template))
}
