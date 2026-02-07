//go:build azure
// +build azure

package azure

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetClientCloudConfigDefaultsToPublic(t *testing.T) {
	orig, existed := os.LookupEnv(AzureEnvironmentEnvName)
	os.Unsetenv(AzureEnvironmentEnvName)
	defer func() {
		if existed {
			os.Setenv(AzureEnvironmentEnvName, orig)
		}
	}()

	config, err := getClientCloudConfig()
	require.NoError(t, err)
	assert.Equal(t, cloud.AzurePublic, config)
}

func TestGetClientCloudConfigPublic(t *testing.T) {
	orig := os.Getenv(AzureEnvironmentEnvName)
	os.Setenv(AzureEnvironmentEnvName, "AzurePublicCloud")
	defer os.Setenv(AzureEnvironmentEnvName, orig)

	config, err := getClientCloudConfig()
	require.NoError(t, err)
	assert.Equal(t, cloud.AzurePublic, config)
}

func TestGetClientCloudConfigChina(t *testing.T) {
	orig := os.Getenv(AzureEnvironmentEnvName)
	os.Setenv(AzureEnvironmentEnvName, "AzureChinaCloud")
	defer os.Setenv(AzureEnvironmentEnvName, orig)

	config, err := getClientCloudConfig()
	require.NoError(t, err)
	assert.Equal(t, cloud.AzureChina, config)
}

func TestGetClientCloudConfigGovernment(t *testing.T) {
	orig := os.Getenv(AzureEnvironmentEnvName)
	os.Setenv(AzureEnvironmentEnvName, "AzureUSGovernmentCloud")
	defer os.Setenv(AzureEnvironmentEnvName, orig)

	config, err := getClientCloudConfig()
	require.NoError(t, err)
	assert.Equal(t, cloud.AzureGovernment, config)
}

func TestGetClientCloudConfigInvalidName(t *testing.T) {
	orig := os.Getenv(AzureEnvironmentEnvName)
	os.Setenv(AzureEnvironmentEnvName, "InvalidCloud")
	defer os.Setenv(AzureEnvironmentEnvName, orig)

	_, err := getClientCloudConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no cloud environment matching the name")
}

func TestGetArmClientOptionsReturnsCorrectCloudConfig(t *testing.T) {
	orig := os.Getenv(AzureEnvironmentEnvName)
	os.Setenv(AzureEnvironmentEnvName, "AzurePublicCloud")
	defer os.Setenv(AzureEnvironmentEnvName, orig)

	opts, err := getArmClientOptions()
	require.NoError(t, err)
	require.NotNil(t, opts)
	assert.Equal(t, cloud.AzurePublic, opts.Cloud)
}

func TestGetArmEndpointReturnsPublicEndpoint(t *testing.T) {
	orig := os.Getenv(AzureEnvironmentEnvName)
	os.Setenv(AzureEnvironmentEnvName, "AzurePublicCloud")
	defer os.Setenv(AzureEnvironmentEnvName, orig)

	endpoint, err := getArmEndpoint()
	require.NoError(t, err)
	assert.Equal(t, "https://management.azure.com", endpoint)
}

func TestGetArmEndpointReturnsChinaEndpoint(t *testing.T) {
	orig := os.Getenv(AzureEnvironmentEnvName)
	os.Setenv(AzureEnvironmentEnvName, "AzureChinaCloud")
	defer os.Setenv(AzureEnvironmentEnvName, orig)

	endpoint, err := getArmEndpoint()
	require.NoError(t, err)
	assert.Equal(t, "https://management.chinacloudapi.cn", endpoint)
}

func TestGetArmEndpointReturnsGovernmentEndpoint(t *testing.T) {
	orig := os.Getenv(AzureEnvironmentEnvName)
	os.Setenv(AzureEnvironmentEnvName, "AzureUSGovernmentCloud")
	defer os.Setenv(AzureEnvironmentEnvName, orig)

	endpoint, err := getArmEndpoint()
	require.NoError(t, err)
	assert.Equal(t, "https://management.usgovcloudapi.net", endpoint)
}
