package kubernetes_test

import (
	"io/ioutil"
	"testing"

	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"
)

func TestAzureIdentityBindingExists(t *testing.T) {
	i := kubernetes.NewDefaultIAM()

	b, err := i.AzureIdentityBindingExists("default")
	handleResult(&b, err)

	b, err = i.AzureIdentityBindingExists("")
	handleResult(&b, err)

	b, err = i.AzureIdentityBindingExists("blah")
	handleResult(&b, err)
}

func TestCreateAIBinding(t *testing.T) {
	by, _ := ioutil.ReadFile("assets/azure-identity-binding.yaml")
	// b, _ := ioutil.ReadFile("assets/pod-test.yaml")

	i := kubernetes.NewDefaultIAM()

	b, err := i.CreateAIB(by, "demo", "demo-binding", "default")

	handleResult(&b, err)
}
