package connection

import (
	"context"
	"log"
	"sync"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-03-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/citihub/probr-sdk/utils"
)

// AzureCredentials ...
type AzureCredentials struct {
	SubscriptionID, ClientID, TenantID, ClientSecret string
	Authorizer                                       autorest.Authorizer
}

// AzureConnection simplifies the connection with cloud provider
type AzureConnection struct {
	isCloudAvailable error
	ctx              context.Context
	credentials      AzureCredentials
	ResourceGroup    *AzureResourceGroup  // Client obj to interact with Azure Resource Groups
	StorageAccount   *AzureStorageAccount // Client obj to interact with Azure Storage Accounts
	ManagedCluster   *AzureManagedCluster // Client obj to interact with Azure Kubernetes Service
	Disk             *AzureDisk           // Client obj to interact with Azure Disks
}

// Azure interface defining all azure methods
type Azure interface {
	IsCloudAvailable() error
	GetResourceGroupByName(name string) (resources.Group, error)
	// Storage Account Functions
	CreateStorageAccount(accountName, accountGroupName string, tags map[string]*string, httpsOnly bool, networkRuleSet *storage.NetworkRuleSet) (storage.Account, error)
	DeleteStorageAccount(resourceGroupName, accountName string) error
	//AKS Functions
	GetManagedClusterJSON(resourceGroupName, clusterName string) ([]byte, error)
	GetManagedClusterAdminCredentials(resourceGroupName, clusterName string) (string, error)
	ClusterHasRoleAssignment(resourceGroupName, clusterName, roleDefName string) (bool, error)
	//Azure Disk functions
	GetDisk(resourceGroupName string, diskName string) (d compute.Disk, err error)
	ParseDiskDetails(diskURI string) (resourceGroupName, diskName string)
	GetJSONRepresentation(resourceGroupName string, diskName string) (dskJSON []byte, err error)
}

var instance *AzureConnection
var once sync.Once

// NewAzureConnection provides a singleton instance of AzureConnection. Initializes all internal clients to interact with Azure.
func NewAzureConnection(c context.Context, subscriptionID, tenantID, clientID, clientSecret string) (azConn *AzureConnection) {
	once.Do(func() {
		// Guard clause
		if c == nil {
			instance.isCloudAvailable = utils.ReformatError("Context instance cannot be nil")
			return
		}

		instance = &AzureConnection{
			ctx: c,
			credentials: AzureCredentials{
				SubscriptionID: subscriptionID,
				TenantID:       tenantID,
				ClientID:       clientID,
				ClientSecret:   clientSecret,
			},
		}

		// Create an authorization object via the connection config vars
		clientCredentialsConfig := auth.NewClientCredentialsConfig(clientID, clientSecret, tenantID)
		authorizer, authErr := clientCredentialsConfig.Authorizer()
		if authErr == nil {
			instance.credentials.Authorizer = authorizer
		} else {
			instance.isCloudAvailable = utils.ReformatError("Failed to initialize Azure Authorizer: %v", authErr)
			return
		}

		// Create an azure resource group client object via the connection config vars
		var grpErr error
		instance.ResourceGroup, grpErr = NewResourceGroup(c, instance.credentials)
		if grpErr != nil {
			instance.isCloudAvailable = utils.ReformatError("Failed to initialize Azure Resource Group: %v", grpErr)
			return
		}

		// Create an azure resource group client object via the connection config vars
		var saErr error
		instance.StorageAccount, grpErr = NewStorageAccount(c, instance.credentials)
		if saErr != nil {
			instance.isCloudAvailable = utils.ReformatError("Failed to initialize Azure Storage Account: %v", grpErr)
			return
		}

		var csErr error
		instance.ManagedCluster, csErr = NewContainerService(c, instance.credentials)
		if csErr != nil {
			instance.isCloudAvailable = utils.ReformatError("Failed to initialize Azure Kubernetes Service: %v", grpErr)
		}

		var dskErr error
		instance.Disk, dskErr = NewDisk(c, instance.credentials)
		if dskErr != nil {
			instance.isCloudAvailable = utils.ReformatError("Failed to initialize Azure Disk: %v", grpErr)
		}
	})
	return instance
}

// IsCloudAvailable verifies that the connection instantiation did not report a failure
func (az *AzureConnection) IsCloudAvailable() error {
	return az.isCloudAvailable
}

// GetResourceGroupByName returns an existing Resource Group by name
func (az *AzureConnection) GetResourceGroupByName(name string) (resources.Group, error) {
	log.Printf("[DEBUG] getting Resource Group '%s'", name)
	return az.ResourceGroup.Get(name)
}

// CreateStorageAccount creates a storage account
func (az *AzureConnection) CreateStorageAccount(accountName, accountGroupName string, tags map[string]*string, httpsOnly bool, networkRuleSet *storage.NetworkRuleSet) (storage.Account, error) {
	log.Printf("[DEBUG] creating Storage Account '%s'", accountName)
	return az.StorageAccount.Create(accountName, accountGroupName, tags, httpsOnly, networkRuleSet)
}

// DeleteStorageAccount deletes a storage account
func (az *AzureConnection) DeleteStorageAccount(resourceGroupName, accountName string) error {
	log.Printf("[DEBUG] deleting Storage Account '%s'", accountName)
	return az.StorageAccount.Delete(resourceGroupName, accountName)
}

// GetManagedClusterJSON returns the JSON representation of an AKS cluster, similar to az aks show. NOTE that the output from this function has differences to the az cli that needs to be accomodated if you are using the JSON created by this function.
func (az *AzureConnection) GetManagedClusterJSON(resourceGroupName, clusterName string) ([]byte, error) {
	log.Printf("[DEBUG] getting JSON for AKS Cluster '%s'", clusterName)
	return az.ManagedCluster.GetJSONRepresentation(resourceGroupName, clusterName)
}

// GetManagedClusterAdminCredentials returns a base64 encoded kubeconfig file for the cluster admin (equivalent to az get-credentials --admin)
func (az *AzureConnection) GetManagedClusterAdminCredentials(resourceGroupName, clusterName string) (string, error) {
	log.Printf("[DEBUG] getting Cluster Admin credentials for AKS Cluster '%s'", clusterName)
	return az.ManagedCluster.GetClusterAdminCredentials(resourceGroupName, clusterName)
}

// ClusterHasRoleAssignment looks through the Azure role assignments on the cluster and returns true if it find the role assigned.  Note that the roleDefName is the UUID of the role not the friendly name of the role.
func (az *AzureConnection) ClusterHasRoleAssignment(resourceGroupName, clusterName, roleDefName string) (bool, error) {
	log.Printf("[DEBUG] Checking if cluster has Kube Cluster Admin role assignments")
	return az.ManagedCluster.ClusterHasRoleAssignment(resourceGroupName, clusterName, roleDefName)
}

// GetDisk returns the disk client
func (az *AzureConnection) GetDisk(resourceGroupName string, diskName string) (compute.Disk, error) {
	return az.Disk.GetDisk(resourceGroupName, diskName)
}

// ParseDiskDetails parses the resource group name and disk name from an Azure Disk URI
func (az *AzureConnection) ParseDiskDetails(diskURI string) (string, string) {
	return az.Disk.ParseDiskDetails(diskURI)
}

// GetJSONRepresentation returns the JSON representation of an AKS cluster, similar to az aks show. NOTE that the output from this function has differences to the az cli that needs to be accomodated if you are using the JSON created by this function.
func (az *AzureConnection) GetJSONRepresentation(resourceGroupName string, diskName string) (dskJSON []byte, err error) {
	return az.Disk.GetJSONRepresentation(resourceGroupName, diskName)
}
