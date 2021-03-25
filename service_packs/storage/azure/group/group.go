package group

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/citihub/probr/service_packs/storage/azure"
)

// Create creates a new Resource Group in the default location (configured using the AZURE_LOCATION environment variable).
func Create(ctx context.Context, name string) (resources.Group, error) {
	log.Printf("[INFO] creating Resource Group '%s' in location: %v", name, azure.ResourceLocation())
	return client().CreateOrUpdate(
		ctx,
		name,
		resources.Group{
			Location: to.StringPtr(azure.ResourceLocation()),
		})
}

// Get an existing Resource Group by name
func Get(ctx context.Context, name string) (resources.Group, error) {
	log.Printf("[DEBUG] getting a Resource Group '%s'", name)
	return client().Get(ctx, name)
}

// CreateWithTags creates a new Resource Group in the default location (configured using the AZURE_LOCATION environment variable) and sets the supplied tags.
func CreateWithTags(ctx context.Context, name string, tags map[string]*string) (resources.Group, error) {
	log.Printf("[INFO] creating Resource Group '%s' on location: '%v'", name, azure.ResourceLocation())
	return client().CreateOrUpdate(
		ctx,
		name,
		resources.Group{
			Location: to.StringPtr(azure.ResourceLocation()),
			Tags:     tags,
		})
}

func client() resources.GroupsClient {

	// Create an azure resource group client object via the connection config vars
	c := resources.NewGroupsClient(azure.SubscriptionID())

	// Create an authorization object via the connection config vars
	authorizer := auth.NewClientCredentialsConfig(azure.ClientID(), azure.ClientSecret(), azure.TenantID())

	authorizerToken, err := authorizer.Authorizer()
	if err == nil {
		c.Authorizer = authorizerToken
	} else {
		log.Printf("[ERROR] Unable to authorise Resource Group client: %v", err)
	}
	return c
}
