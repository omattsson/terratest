package azure

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	autorestAzure "github.com/Azure/go-autorest/autorest/azure"
)

func getClientCloudConfig() (cloud.Configuration, error) {
	envName := getDefaultEnvironmentName()
	switch strings.ToUpper(envName) {
	case "AZURECHINACLOUD":
		return cloud.AzureChina, nil
	case "AZUREUSGOVERNMENTCLOUD":
		return cloud.AzureGovernment, nil
	case "AZUREPUBLICCLOUD":
		return cloud.AzurePublic, nil
	case "AZURESTACKCLOUD":
		env, err := autorestAzure.EnvironmentFromName(envName)
		if err != nil {
			return cloud.Configuration{}, err
		}
		c := cloud.Configuration{
			ActiveDirectoryAuthorityHost: env.ActiveDirectoryEndpoint,
			Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
				cloud.ResourceManager: {
					Audience: env.TokenAudience,
					Endpoint: env.ResourceManagerEndpoint,
				},
			},
		}
		return c, nil
	default:
		return cloud.Configuration{},
			fmt.Errorf("no cloud environment matching the name: %s. "+
				"Available values are: "+
				"AzurePublicCloud (default), "+
				"AzureUSGovernmentCloud, "+
				"AzureChinaCloud or "+
				"AzureStackCloud",
				envName)
	}
}

// getArmClientOptions returns *arm.ClientOptions configured for the current Azure cloud.
// Use this when creating ARM client factories to ensure consistent cloud configuration.
func getArmClientOptions() (*arm.ClientOptions, error) {
	clientCloudConfig, err := getClientCloudConfig()
	if err != nil {
		return nil, err
	}
	return &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: clientCloudConfig,
		},
	}, nil
}

// getArmEndpoint returns the Azure Resource Manager endpoint for the current Azure cloud.
func getArmEndpoint() (string, error) {
	clientCloudConfig, err := getClientCloudConfig()
	if err != nil {
		return "", err
	}
	rmConfig, ok := clientCloudConfig.Services[cloud.ResourceManager]
	if !ok {
		return "", fmt.Errorf("no Resource Manager service configuration found for the current cloud")
	}
	return rmConfig.Endpoint, nil
}
