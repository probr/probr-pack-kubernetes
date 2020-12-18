package iam

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
)

// CreateAIB creates an AzureIdentityBinding to a specified AzureIdentity in a specified non-default namespace
func TestSetEnv(t *testing.T) {

	// Initialise IAM
	i := &IAM{}
	i.setenv()

	if i.probePodName == "" {
		t.Logf("probePodName not set")
		t.Fail()
	}
	if i.probeImage == "" {
		t.Logf("probeImage not set")
		t.Fail()
	}
	if i.azureIdentitySelector == "" {
		t.Logf("azureIdentitySelector not set")
		t.Fail()
	}
}

func TestCreateAIBObject(t *testing.T) {

	// Initialise IAM
	i := &IAM{}
	i.setenv()

	namespace := "default"
	aibName := "testAib"
	aiName := "testAi"
	runtimeAib := i.createAIBObject(namespace, aibName, aiName)

	// Check returned type
	_, typeCastOK := runtimeAib.(runtime.Object)
	if !typeCastOK {
		t.Logf("wrong type returned")
		t.Fail()
	}
}
