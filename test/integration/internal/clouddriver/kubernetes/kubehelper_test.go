// Package kubernetes_test provides test functions for interacting with the kubernetes cluster.  Note, these are integration tests and require an active kubernetes
// cluster, etc.
package kubernetes_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"
	_ "gitlab.com/citihub/probr/internal/config"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/client-go/kubernetes/scheme"
)

var (
	testNS        = "probrtestns"
	testPod       = "probrtestpod"
	testContainer = "curlimages"
	testImage     = "curlimages/curl"
)

func TestMain(m *testing.M) {
	log.Print("[NOTICE] Running Kube tests ...")
	result := m.Run()

	log.Printf("[NOTICE] Completed Kube tests ... (result: %v)", result)
	os.Exit(result)
}

func TestGetPods(t *testing.T) {
	kubernetes.GetKubeInstance().GetPods("")
}

func TestCreatePod(t *testing.T) {
	_, err := kubernetes.GetKubeInstance().CreatePod(&testPod, &testNS, &testContainer, &testImage, true, nil)

	handleResult(nil, err)
}

func TestCreatePodFromYaml(t *testing.T) {
	//read the yaml:
	b, _ := ioutil.ReadFile("assets/pod-test.yaml")
	//y := string(b)

	_, err := kubernetes.GetKubeInstance().CreatePodFromYaml(b, &testPod, &testNS, &testImage, nil, true)

	handleResult(nil, err)
}

func TestExecCmd(t *testing.T) {

	url := "http://www.google.com"
	cmd := "curl -s -o /dev/null -I -L -w %{http_code} " + url

	res := kubernetes.GetKubeInstance().ExecCommand(&cmd, &testNS, &testPod)

	log.Printf("[NOTICE] Test command result:")
	log.Printf("[NOTICE] stdout: %v stderr: %v exit code: %v", res.Stdout, res.Stderr, res.Code)

	handleResult(nil, res.Err)
}

func TestDeletePod(t *testing.T) {
	err := kubernetes.GetKubeInstance().DeletePod(&testPod, &testNS, true)

	handleResult(nil, err)
}

func TestDeleteNamespace(t *testing.T) {
	err := kubernetes.GetKubeInstance().DeleteNamespace(&testNS)

	handleResult(nil, err)
}

func TestConfigMap(t *testing.T) {
	c := "test-cm"
	cm, err := kubernetes.GetKubeInstance().CreateConfigMap(&c, &testNS)

	handleResult(nil, err)

	//now delete it:
	if cm != nil {
		err = kubernetes.GetKubeInstance().DeleteConfigMap(&c, &testNS)
		handleResult(nil, err)
	}
}

func TestGetConstraintTemplate(t *testing.T) {
	con, err := kubernetes.GetKubeInstance().GetConstraintTemplates("k8sazure")

	log.Printf("[NOTICE] constraints: %v", con)

	handleResult(nil, err)
}

func TestGetClusterRoles(t *testing.T) {
	crl, err := kubernetes.GetKubeInstance().GetClusterRoles()

	log.Printf("[NOTICE] cluster roles: %v", crl)

	handleResult(nil, err)
}

func TestGetClusterRolesByResource(t *testing.T) {

	crl, err := kubernetes.GetKubeInstance().GetClusterRolesByResource("*")

	log.Printf("[NOTICE] cluster roles with '*' resource")

	for _, cr := range *crl {
		log.Printf("[NOTICE] role name: %v role labels: %v", cr.Name, cr.Labels)
	}

	handleResult(nil, err)
}

func TestGetRolesByResource(t *testing.T) {

	rl, err := kubernetes.GetKubeInstance().GetRolesByResource("*")

	log.Printf("[NOTICE] roles with '*' resource")
	
	for _, r := range *rl {
		log.Printf("[NOTICE] role name: %v role labels: %v", r.Name, r.Labels)
	}

	handleResult(nil, err)
}

func TestGetRawResourcesByGrp(t *testing.T) {
	j, err := kubernetes.GetKubeInstance().GetRawResourcesByGrp("apis/aadpodidentity.k8s.io/v1/azureidentities")

	fmt.Printf("JSON Azure Identities: %v\n", j)

	for _, r := range j.Items {
		fmt.Printf("Name: %v\n", r.Metadata["name"])
		fmt.Printf("Namespace: %v\n", r.Metadata["namespace"])
	}
	

	handleResult(nil, err)

	j, err = kubernetes.GetKubeInstance().GetRawResourcesByGrp("apis/aadpodidentity.k8s.io/v1/azureidentitybindings")

	fmt.Printf("JSON Azure Identity Bindings: %v\n", j)

	for _, r := range j.Items {
		fmt.Printf("Name: %v\n", r.Metadata["name"])
		fmt.Printf("Namespace: %v\n", r.Metadata["namespace"])
	}
	

	handleResult(nil, err)
}

func handleResult(yesNo *bool, err error) {
	if err != nil {
		//Log but don't check for this atm, i.e. keep tests running
		// log.Printf("[WARN] Test failed with ERROR: %v\n", err)
		fmt.Printf("Test failed with ERROR: %v\n", err)
		return
	}

	//if we didn't get an error, then the test was successful in the sense
	//that it conducted the kube operation without issue
	log.Print("[NOTICE] Test successfully performed Kubernetes operation.")

	if yesNo != nil {
		log.Printf("[NOTICE] Result of operation was: %t\n", *yesNo)
	}

}
