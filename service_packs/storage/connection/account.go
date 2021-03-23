package connection

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/citihub/probr/service_packs/storage/azure"
)

// DeleteAccount - deletes a storage account given the azure contect, resource group and account name
func DeleteAccount(ctx context.Context, resourceGroupName, accountName string) error {

	c := accountClient()

	_, err := c.Delete(ctx, resourceGroupName, accountName)

	return err
}

// CreateWithNetworkRuleSet starts creation of a new Storage Account and waits for the account to be created.
func CreateWithNetworkRuleSet(ctx context.Context, accountName, accountGroupName string, tags map[string]*string, httpsOnly bool, networkRuleSet *storage.NetworkRuleSet) (storage.Account, error) {

	var sa storage.Account
	c := accountClient()

	r, err := c.CheckNameAvailability(
		ctx,
		storage.AccountCheckNameAvailabilityParameters{
			Name: to.StringPtr(accountName),
			Type: to.StringPtr("Microsoft.Storage/storageAccounts"),
		})
	if err != nil {
		return sa, err
	}

	if *r.NameAvailable != true {
		return sa, fmt.Errorf(
			"storage account name [%sa] not available: %v\nserver message: %v",
			accountName, err, *r.Message)
	}

	networkRuleSetParam := &storage.AccountPropertiesCreateParameters{
		EnableHTTPSTrafficOnly: to.BoolPtr(httpsOnly),
		NetworkRuleSet:         networkRuleSet,
	}

	future, err := c.Create(
		ctx,
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

	if err != nil {
		return sa, err
	}

	err = future.WaitForCompletionRef(ctx, c.Client)
	if err != nil {
		return sa, err
	}

	return future.Result(c)
}

// AccountProperties returns the properties for the specified storage account including but not limited to name, SKU name, location, and account status
func AccountProperties(ctx context.Context, rgName, accountName string) (storage.Account, error) {
	return accountClient().GetProperties(ctx, rgName, accountName, "")
}

// AccountPrimaryKey return the primary key
func AccountPrimaryKey(ctx context.Context, accountName, accountGroupName string) string {
	response, err := getAccountKeys(ctx, accountName, accountGroupName)
	if err != nil {
		//log.Fatalf("failed to list keys: %v", err)
	}
	return *(((*response.Keys)[0]).Value)
}

func getAccountKeys(ctx context.Context, accountName, accountGroupName string) (storage.AccountListKeysResult, error) {
	return accountClient().ListKeys(ctx, accountGroupName, accountName, "")
}

func accountClient() storage.AccountsClient {

	// Create an azure storage account client object via the connection config vars
	c := storage.NewAccountsClient(azure.SubscriptionID())

	// Create an authorization object via the connection config vars
	authorizer := auth.NewClientCredentialsConfig(azure.ClientID(), azure.ClientSecret(), azure.TenantID())

	authorizerToken, err := authorizer.Authorizer()
	if err == nil {
		c.Authorizer = authorizerToken
	} else {
		//log.Printf("[ERROR] Unable to authorise Storage Account accountClient: %v", err)
	}
	return c
}
