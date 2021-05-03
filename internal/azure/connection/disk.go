package connection

import (
	"context"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-03-01/compute"
	"github.com/citihub/probr-sdk/utils"
)

// AzureDisk contains the state of the Azure Disk context
type AzureDisk struct {
	ctx          context.Context
	credentials  AzureCredentials
	azDiskClient compute.DisksClient
}

// NewDisk validates the context and credentials, and retrieves the corresponding Disks Client
func NewDisk(ctx context.Context, creds AzureCredentials) (dsk *AzureDisk, err error) {
	if ctx == nil {
		err = utils.ReformatError("Context instance cannot be nil")
		return
	}
	if creds.Authorizer == nil {
		err = utils.ReformatError("Authorizer instance cannot be nil")
		return
	}

	dsk = &AzureDisk{
		ctx:         ctx,
		credentials: creds,
	}

	// Create an azure storage account client object via the connection config vars
	var dskErr error
	dsk.azDiskClient, dskErr = dsk.getDisksClient(creds)
	if dskErr != nil {
		err = utils.ReformatError("Failed to initialize Azure Disks client: %v", dskErr)
		return
	}
	return
}

// Create an azure container services client object via the connection config vars
func (dsk *AzureDisk) getDisksClient(creds AzureCredentials) (dskClient compute.DisksClient, err error) {
	dskClient = compute.NewDisksClient(creds.SubscriptionID)
	dskClient.Authorizer = creds.Authorizer

	return
}

// GetDisk Retrieves the specified disk from the specified resource group
func (dsk *AzureDisk) GetDisk(resourceGroup string, diskName string) (d compute.Disk, err error) {
	d, err = dsk.azDiskClient.Get(dsk.ctx, resourceGroup, diskName)
	log.Printf("[DEBUG] GetDisk.d: %v", d)
	return
}

// ParseDiskDetails returns the resource group name and disk name of the Managed Disk
func (dsk *AzureDisk) ParseDiskDetails(diskURI string) (resourceGroupName, diskName string) {
	s := strings.Split(diskURI, "/")
	resourceGroupName = s[4]
	diskName = s[8]
	return
}

// GetJSONRepresentation returns the JSON representation of an AKS cluster, similar to az aks show. NOTE that the output from this function has differences to the az cli that needs to be accomodated if you are using the JSON created by this function.
func (dsk *AzureDisk) GetJSONRepresentation(resourceGroupName string, diskName string) (dskJSON []byte, err error) {
	var d compute.Disk
	d, err = dsk.azDiskClient.Get(dsk.ctx, resourceGroupName, diskName)
	if err != nil {
		log.Printf("Error getting ContainerServiceClient: %v", err)
		return
	}
	dskJSON, err = d.MarshalJSON()
	return
}
