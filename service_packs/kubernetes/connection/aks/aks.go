// Package aks is a wrapper for the connection to the Azure Kubernetes API
package aks

import (
	"log"

	aibv1 "github.com/Azure/aad-pod-identity/pkg/apis/aadpodidentity"
	"github.com/citihub/probr/service_packs/kubernetes/connection"
)

// AKS implements the Azure Kubernetes Service wrapper
type AKS struct {
	conn connection.Connection
}

// NewAKS creates a new AKS instance taking a Connection instance as argument.
func NewAKS(connection connection.Connection) *AKS {
	aks := &AKS{}
	aks.conn = connection

	// Guard clause: Check valid connection instance
	if connection == nil {
		log.Fatal("Connection instance cannot be nil")
	}

	return aks
}

// CreateAIB creates an AzureIdentityBinding in the cluster, 409 error if it already exists
func (aks *AKS) CreateAIB(namespace, aibName, aiName string) (resource connection.APIResource, err error) {

	aib := aibv1.AzureIdentityBinding{}

	aib.TypeMeta.Kind = "AzureIdentityBinding"
	aib.TypeMeta.APIVersion = "aadpodidentity.k8s.io/v1"
	aib.ObjectMeta.Namespace = namespace
	aib.ObjectMeta.Name = aibName
	aib.Spec.AzureIdentity = aiName
	aib.Spec.Selector = "aadpodidbinding"
	// Copy into a runtime.Object which is required for the api request
	runtimeAib := aib.DeepCopyObject()

	// set the api path for the aadpodidentity package which include the azureidentitybindings custom resource definition
	apiPath := "apis/aadpodidentity.k8s.io/v1"

	resource, err = aks.conn.PostRawResource(apiPath, namespace, "azureidentitybindings", runtimeAib)
	log.Printf("Resource %v", resource)

	return
}

// GetIdentityByNameAndNamespace queries cluster and returns resource, 404 error if not found
func (aks *AKS) GetIdentityByNameAndNamespace(azureIdentityName, namespace string) (resource connection.APIResource, err error) {
	// Azure Identities are implemented as K8s Custom Resource Definition
	// Need to make a 'raw' call to the corresponding K8s endpoint
	// The K8s api endpoint for Azure Indentity is: 		"apis/aadpodidentity.k8s.io/v1/azureidentities"
	apiEndPoint := "apis/aadpodidentity.k8s.io/v1"
	resourceType := "azureidentities"

	resource, err = aks.conn.GetRawResourceByName(apiEndPoint, namespace, resourceType, azureIdentityName)

	return
}

// GetIdentityBindingByNameAndNamespace queries cluster and returns resource, 404 eror if not found
func (aks *AKS) GetIdentityBindingByNameAndNamespace(azureIdentityBindingName, namespace string) (resource connection.APIResource, err error) {
	// Azure Identity Bindings are implemented as K8s Custom Resource Definition
	// Need to make a 'raw' call to the corresponding K8s endpoint
	// The K8s api endpoint for Azure Indentity Binding is:	"apis/aadpodidentity.k8s.io/v1/azureidentitybindings"
	apiEndPoint := "apis/aadpodidentity.k8s.io/v1"
	resourceType := "azureidentitybindings"

	resource, err = aks.conn.GetRawResourceByName(apiEndPoint, namespace, resourceType, azureIdentityBindingName)

	return
}
