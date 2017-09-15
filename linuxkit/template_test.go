package linuxkit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadTemplate(t *testing.T) {
	_, err := Read("./aws.yml")
	require.NoError(t, err)
}
