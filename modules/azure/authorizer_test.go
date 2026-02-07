//go:build azure
// +build azure

package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAzureCredential(t *testing.T) {
	cred, err := NewAzureCredential()
	require.NoError(t, err)
	require.NotNil(t, cred)
}
