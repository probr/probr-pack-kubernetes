package kubernetes

import (
	// "fmt"
	"log"
	"strconv"

	apiv1 "k8s.io/api/core/v1"
)

const (
	//TODO: default to these values for MVP - need to expose in future
	testNamespace = "network-access-test-ns"
	testImage     = "curlimages/curl"
	testContainer = "curlimages"
	testPodName   = "na-test-pod"
)

//SetupNetworkAccessTestPod creates a pod with characteristics required for testing network access.
func SetupNetworkAccessTestPod() (*apiv1.Pod, error) {
	pname, ns, cname, image := testPodName, testNamespace, testContainer, testImage
	p, err := CreatePod(&pname, &ns, &cname, &image)

	if err != nil {
		return nil, err
	}

	return p, nil
}

//TeardownNetworkAccessTestPod ...
func TeardownNetworkAccessTestPod() error {	
	ns := testNamespace
	err := DeleteNamespace(&ns)

	return err
}

//AccessURL calls the supplied URL and returns the http code
func AccessURL(url *string) (int, error) {

	//create a curl command to access the supplied url		
	cmd := "curl -s -o /dev/null -I -L -w %{http_code} " + *url
	ns, pn := testNamespace, testPodName
	httpCode, _, err := ExecCommand(&cmd, &ns, &pn)

	if err != nil {
		return -1, err
	}

	log.Printf("[NOTICE] URL: %v HTTP Code: %v", *url, httpCode)

	httpStatusCode, err := strconv.Atoi(httpCode)
	if err != nil {
		return -1, err
	}
	
	return httpStatusCode, nil
}
