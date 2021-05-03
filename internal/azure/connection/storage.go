package connection

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/citihub/probr-pack-kubernetes/internal/azure"
	"github.com/citihub/probr-sdk/utils"
)

// AzureStorageAccount ...
type AzureStorageAccount struct {
	ctx                    context.Context
	credentials            AzureCredentials
	azStorageAccountClient storage.AccountsClient
}

// NewStorageAccount provides a new instance of AzureStorageAccount
func NewStorageAccount(c context.Context, creds AzureCredentials) (sa *AzureStorageAccount, err error) {

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

	sa = &AzureStorageAccount{
		ctx:         c,
		credentials: creds,
	}

	// Create an azure storage account client object via the connection config vars
	var saErr error
	sa.azStorageAccountClient, saErr = sa.getStorageAccountClient(creds)
	if saErr != nil {
		err = utils.ReformatError("Failed to initialize Azure Storage Account client: %v", saErr)
		return
	}

	return
}

func (sa *AzureStorageAccount) getStorageAccountClient(creds AzureCredentials) (saClient storage.AccountsClient, err error) {

	// Create an azure storage account client object via the connection config vars
	saClient = storage.NewAccountsClient(creds.SubscriptionID)

	saClient.Authorizer = creds.Authorizer

	return
}

// Create starts creation of a new Storage Account and waits for the account to be created.
func (sa *AzureStorageAccount) Create(accountName, accountGroupName string, tags map[string]*string, httpsOnly bool, networkRuleSet *storage.NetworkRuleSet) (storage.Account, error) {

	log.Printf("[DEBUG] creating Storage Account '%s'", accountName)

	var storageAccount storage.Account

	checkNameResult, checkNameErr := sa.azStorageAccountClient.CheckNameAvailability(
		sa.ctx,
		storage.AccountCheckNameAvailabilityParameters{
			Name: to.StringPtr(accountName),
			Type: to.StringPtr("Microsoft.Storage/storageAccounts"),
		})
	if checkNameErr != nil {
		err := utils.ReformatError("Failed checking name availability: %v", checkNameErr)
		return storageAccount, err
	}
	if *checkNameResult.NameAvailable != true {
		err := utils.ReformatError("Provided name for storage account '%s' is not available: %v", accountName, *checkNameResult.Message)
		return storageAccount, err
	}

	networkRuleSetParam := &storage.AccountPropertiesCreateParameters{
		EnableHTTPSTrafficOnly: to.BoolPtr(httpsOnly),
		NetworkRuleSet:         networkRuleSet,
	}

	future, createErr := sa.azStorageAccountClient.Create(
		sa.ctx,
		accountGroupName,
		accountName,
		storage.AccountCreateParameters{
			Sku: &storage.Sku{
				Name: storage.StandardLRS},
			Kind:                              storage.Storage,
			Location:                          to.StringPtr(azure.ResourceLocation()),
			AccountPropertiesCreateParameters: networkRuleSetParam,
			Tags:                              tags,
		})
	if createErr != nil {
		return storageAccount, createErr
	}

	waitErr := future.WaitForCompletionRef(sa.ctx, sa.azStorageAccountClient.Client)
	if waitErr != nil {
		return storageAccount, waitErr
	}

	return future.Result(sa.azStorageAccountClient)

}

// Delete deletes a storage account given the resource group and account name
func (sa *AzureStorageAccount) Delete(resourceGroupName, accountName string) error {

	log.Printf("[DEBUG] deleting Storage Account '%s' from Resource Group '%s'", accountName, resourceGroupName)

	_, err := sa.azStorageAccountClient.Delete(sa.ctx, resourceGroupName, accountName)

	return err
}
