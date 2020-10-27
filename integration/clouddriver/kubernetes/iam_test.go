package kubernetes_test

import (
	"io/ioutil"
	"testing"

	"github.com/citihub/probr/internal/clouddriver/kubernetes"
)

func TestAzureIdentityBindingExists(t *testing.T) {
	i := kubernetes.NewDefaultIAM()

	b, err := i.AzureIdentityBindingExists(true)
	handleResult(&b, err)

	b, err = i.AzureIdentityBindingExists(false)
	handleResult(&b, err)
}

func TestCreateAIBinding(t *testing.T) {
	by, _ := ioutil.ReadFile("assets/azure-identity-binding.yaml")	

	i := kubernetes.NewDefaultIAM()

	b, err := i.CreateAIB(by, "demo", "demo-binding", "default")

	handleResult(&b, err)
}

func TestExecVerificationCmd(t *testing.T) {
	i := kubernetes.NewDefaultIAM()

	b, _ := ioutil.ReadFile("assets/pod-test.yaml")

	pd, _ := i.CreateIAMProbePod(b, true)

	i.ExecuteVerificationCmd(pd.Name, kubernetes.CurlAuthToken, true)
}
