package connection

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/citihub/probr-sdk/utils"
)

// AzureResourceGroup ...
type AzureResourceGroup struct {
	ctx                   context.Context
	credentials           AzureCredentials
	azResourceGroupClient resources.GroupsClient
}

// NewResourceGroup provides a new instance of AzureResourceGroup
func NewResourceGroup(c context.Context, creds AzureCredentials) (rg *AzureResourceGroup, err error) {

	// Guard clause - context
	if c == nil {
		err = utils.ReformatError("Context instance cannot be nil")
		return
	}

	// Guard clause - authorizer
	if creds.Authorizer == nil {
		err = utils.ReformatError("Authorizer instance cannot be nil")
		return
	}

	rg = &AzureResourceGroup{
		ctx:         c,
		credentials: creds,
	}

	// Create an azure resource group client object via the connection config vars
	var grpErr error
	rg.azResourceGroupClient, grpErr = rg.getResourceGroupClient(creds)
	if grpErr != nil {
		err = utils.ReformatError("Failed to initialize Azure Group client: %v", grpErr)
		return
	}

	return
}

// Get an existing Resource Group by name
func (rg *AzureResourceGroup) Get(name string) (resources.Group, error) {
	log.Printf("[DEBUG] getting a Resource Group '%s'", name)
	return rg.azResourceGroupClient.Get(rg.ctx, name)
}

func (rg *AzureResourceGroup) getResourceGroupClient(creds AzureCredentials) (rgClient resources.GroupsClient, err error) {

	// Create an azure resource group client object via the connection config vars
	rgClient = resources.NewGroupsClient(creds.SubscriptionID)

	rgClient.Authorizer = creds.Authorizer

	return
}
