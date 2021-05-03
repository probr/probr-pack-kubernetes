package connection

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/authorization/mgmt/2015-07-01/authorization"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2018-03-31/containerservice"

	"github.com/citihub/probr-sdk/utils"
)

// AzureManagedCluster holds the state of the Azure context
type AzureManagedCluster struct {
	ctx                     context.Context
	credentials             AzureCredentials
	azManagedClustersClient containerservice.ManagedClustersClient
}

//var azConnection connection.Azure // Provides functionality to interact with Azure

// NewContainerService provides a new instance of AzureContainerService
func NewContainerService(c context.Context, creds AzureCredentials) (cs *AzureManagedCluster, err error) {

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

	cs = &AzureManagedCluster{
		ctx:         c,
		credentials: creds,
	}

	// Create an azure storage account client object via the connection config vars
	var csErr error
	cs.azManagedClustersClient, csErr = cs.getManagedClusterClient(creds)
	if csErr != nil {
		err = utils.ReformatError("Failed to initialize Azure Kubernetes Service client: %v", csErr)
		return
	}

	return
}

// GetJSONRepresentation returns the JSON representation of an AKS cluster, similar to az aks show. NOTE that the output from this function has differences to the az cli that needs to be accomodated if you are using the JSON created by this function.
func (amc *AzureManagedCluster) GetJSONRepresentation(resourceGroupName string, clusterName string) (aksJSON []byte, err error) {
	var cs containerservice.ManagedCluster
	cs, err = amc.azManagedClustersClient.Get(amc.ctx, resourceGroupName, clusterName)
	if err != nil {
		log.Printf("Error getting ContainerServiceClient: %v", err)
		return
	}
	aksJSON, err = cs.MarshalJSON()
	return
}

// GetClusterAdminCredentials returns a base64 encoded kubeconfig file for the cluster admin (equivalent to az get-credentials --admin)
func (amc *AzureManagedCluster) GetClusterAdminCredentials(resourceGroupName, clusterName string) (kubeconfigBase64 string, err error) {

	credResults, err := amc.azManagedClustersClient.ListClusterAdminCredentials(amc.ctx, resourceGroupName, clusterName)

	if err != nil {
		log.Printf("Error getting Cluster Admin credentials: %v", err)
	}

	kc := (*credResults.Kubeconfigs)[0]
	kubeconfigBase64 = string(*kc.Value)

	return
}

// ClusterHasRoleAssignment looks through the Azure role assignments on the cluster and returns true if it find the role assigned.  Note that the roleDefName is the UUID of the role not the friendly name of the role.
func (amc *AzureManagedCluster) ClusterHasRoleAssignment(resourceGroupName, clusterName, roleDefName string) (present bool, err error) {
	roleAssignmentsClient, _ := amc.getRoleAssignmentsClient(amc.credentials)

	roleDefinitionID := fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Authorization/roleDefinitions/%s", amc.credentials.SubscriptionID, roleDefName)

	log.Printf("[DEBUG] Checking Managed Cluster for role name")
	filter := fmt.Sprintf("atScope()")

	resourceProviderNamespace := "Microsoft.ContainerService"
	parentResourcePath := ""
	resourceType := "managedClusters"

	res, err := roleAssignmentsClient.ListForResource(amc.ctx, resourceGroupName, resourceProviderNamespace, parentResourcePath, resourceType, clusterName, filter)
	if err != nil {
		log.Printf("Error listing role assignments. %v", err)
		return false, err
	}

	for _, v := range res.Values() {
		log.Printf("[DEBUG] Found role. ID: %s, Name: %s; Type: %s", *v.ID, *v.Name, *v.Type)
		log.Printf("[DEBUG] Role Definition ID: %s", *v.Properties.RoleDefinitionID)
		// TODO:
		if *v.Properties.RoleDefinitionID == roleDefinitionID {
			log.Printf("Found role assigned to cluster. Returning.")
			return true, utils.ReformatError("Role assignment found for role %s", roleDefName)
		}
	}

	return false, nil

}

func (amc *AzureManagedCluster) getManagedClusterClient(creds AzureCredentials) (csClient containerservice.ManagedClustersClient, err error) {

	log.Printf("Credentials: Subscription: %s", creds.SubscriptionID)
	csClient = containerservice.NewManagedClustersClient(creds.SubscriptionID)
	csClient.Authorizer = creds.Authorizer

	return
}

//TODO: put this in an RBAC .go file
func (amc *AzureManagedCluster) getRoleAssignmentsClient(creds AzureCredentials) (authorization.RoleAssignmentsClient, error) {
	roleClient := authorization.NewRoleAssignmentsClient(creds.SubscriptionID)
	roleClient.Authorizer = creds.Authorizer
	return roleClient, nil
}
