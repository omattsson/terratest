package azure

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/go-autorest/autorest"
	az "github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

const (
	// AuthFromEnvClient is an env variable supported by the Azure SDK
	AuthFromEnvClient = "AZURE_CLIENT_ID"

	// AuthFromEnvTenant is an env variable supported by the Azure SDK
	AuthFromEnvTenant = "AZURE_TENANT_ID"

	// AuthFromFile is an env variable supported by the Azure SDK
	AuthFromFile = "AZURE_AUTH_LOCATION"
)

// NewAuthorizer creates an Azure authorizer adhering to standard auth mechanisms provided by the Azure Go SDK
// See Azure Go Auth docs here: https://docs.microsoft.com/en-us/go/azure/azure-sdk-go-authorization
// Deprecated: Use NewAzureCredential instead.
func NewAuthorizer() (*autorest.Authorizer, error) {
	// Carry out env var lookups
	_, clientIDExists := os.LookupEnv(AuthFromEnvClient)
	_, tenantIDExists := os.LookupEnv(AuthFromEnvTenant)
	_, fileAuthSet := os.LookupEnv(AuthFromFile)

	// Execute logic to return an authorizer from the correct method
	if clientIDExists && tenantIDExists {
		authorizer, err := auth.NewAuthorizerFromEnvironment()
		return &authorizer, err
	} else if fileAuthSet {
		authorizer, err := auth.NewAuthorizerFromFile(az.PublicCloud.ResourceManagerEndpoint)
		return &authorizer, err
	} else {
		authorizer, err := auth.NewAuthorizerFromCLI()
		return &authorizer, err
	}
}

// NewAzureCredential creates a DefaultAzureCredential configured for the current Azure cloud.
func NewAzureCredential() (*azidentity.DefaultAzureCredential, error) {
	clientCloudConfig, err := getClientCloudConfig()
	if err != nil {
		return nil, err
	}

	return azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
}
