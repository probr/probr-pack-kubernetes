package kubernetes

import (
	"log"
	"os"
	"strconv"

	apiv1 "k8s.io/api/core/v1"
)

const (
	//TODO: default to these values for MVP - need to expose in future
	testNamespace = "probr-network-access-test-ns" //this needs to be set up as an exculsion in the image registry policy
	testImage     = "curlimages/curl"
	testContainer = "curlimages"
	testPodName   = "na-test-pod"
)

// NetworkAccess ...
type NetworkAccess interface {
	ClusterIsDeployed() *bool
	SetupNetworkAccessTestPod() (*apiv1.Pod, error)
	TeardownNetworkAccessTestPod(p *string) error
	AccessURL(pn *string, url *string) (int, error)
}

// NA ...
type NA struct {
	k Kubernetes
}

// NewNA ...
func NewNA(k Kubernetes) *NA {
	n := &NA{}
	n.k = k

	return n
}

// NewDefaultNA ...
func NewDefaultNA() *NA {
	n := &NA{}
	n.k = GetKubeInstance()

	return n
}

// ClusterIsDeployed ...
func (n *NA) ClusterIsDeployed() *bool {
	return n.k.ClusterIsDeployed()
}

//SetupNetworkAccessTestPod creates a pod with characteristics required for testing network access.
func (n *NA) SetupNetworkAccessTestPod() (*apiv1.Pod, error) {
	pname, ns, cname, image := GenerateUniquePodName(testPodName), testNamespace, testContainer, testImage
	p, err := n.k.CreatePod(&pname, &ns, &cname, &image, true, nil)

	if err != nil {
		return nil, err
	}

	return p, nil
}

//TeardownNetworkAccessTestPod ...
func (n *NA) TeardownNetworkAccessTestPod(p *string) error {
	_, exists := os.LookupEnv("DONT_DELETE")
	if !exists {
		ns := testNamespace
		err := n.k.DeletePod(p, &ns, false) //don't worry about waiting
		return err
	}

	return nil
}

//AccessURL calls the supplied URL and returns the http code
func (n *NA) AccessURL(pn *string, url *string) (int, error) {

	//create a curl command to access the supplied url
	cmd := "curl -s -o /dev/null -I -L -w %{http_code} " + *url
	ns := testNamespace
	httpCode, _, ex, err := n.k.ExecCommand(&cmd, &ns, pn)

	if err != nil {
		//check the exit code.  If it's '6' (Couldn't resolve host.)
		//then we want to nil out the error and return the code as this
		//is an expected condition if access is inhibited
		if ex == 6 {
			return ex, nil
		}
		return -1, err
	}

	log.Printf("[NOTICE] URL: %v HTTP Code: %v", *url, httpCode)

	httpStatusCode, err := strconv.Atoi(httpCode)
	if err != nil {
		return -1, err
	}

	return httpStatusCode, nil
}
