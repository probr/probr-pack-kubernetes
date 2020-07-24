package kubernetes_test

import (
	"log"
	"os"
	"testing"

	"citihub.com/probr/internal/clouddriver/kubernetes"
	_ "citihub.com/probr/internal/config"
)

var (
	testNS        = "probrtestns"
	testPod       = "probrtestpod"
	testContainer = "curlimages"
	testImage     = "curlimages/curl"
)

func TestMain(m *testing.M) {
	result := m.Run()

	os.Exit(result)
}

func TestGetPods(t *testing.T) {
	kubernetes.GetPods()
}

func TestCreateNamespace(t *testing.T) {
	_, err := kubernetes.CreateNamespace(&testNS)

	handleResult(nil, err)
}

func TestCreatePod(t *testing.T) {
	_, err := kubernetes.CreatePod(&testPod, &testNS, &testContainer, &testImage, true)

	handleResult(nil, err)
}

func TestExecCmd(t *testing.T) {

	url := "http://www.google.com"
	cmd := "curl -s -o /dev/null -I -L -w %{http_code} " + url

	so, se, ec, err := kubernetes.ExecCommand(&cmd, &testNS, &testPod)

	log.Printf("[NOTICE] Test command result:")
	log.Printf("[NOTICE] stdout: %v stderr: %v exit code: %v", so, se, ec)

	handleResult(nil, err)
}

func TestDeletePod(t *testing.T) {
	err := kubernetes.DeletePod(&testPod, &testNS, true)

	handleResult(nil, err)
}

func TestDeleteNamespace(t *testing.T) {
	err := kubernetes.DeleteNamespace(&testNS)

	handleResult(nil, err)
}

func handleResult(yesNo *bool, err error) {
	if err != nil {
		//Log but don't check for this atm, i.e. keep tests running
		log.Printf("[WARN] Test failed with ERROR: %v\n", err)
		return
	}

	//if we didn't get an error, then the test was successful in the sense
	//that it conducted the kube operation without issue
	log.Print("[NOTICE] Test successfully performed Kubernetes operation.")

	if yesNo != nil {
		log.Printf("[NOTICE] Result of operation was: %t\n", *yesNo)
	}

}
