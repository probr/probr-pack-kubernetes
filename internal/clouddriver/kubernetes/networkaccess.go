package kubernetes

import (
	"log"
	"os"
	"strconv"

	"gitlab.com/citihub/probr/internal/config"
	apiv1 "k8s.io/api/core/v1"
)

const (
	//default values.  Overrides can be set via the environment.
	defaultNATestNamespace   = "probr-network-access-test-ns" //this needs to be set up as an exculsion in the image registry policy	
	defaultNAImageRepository = "curlimages"
	defaultNATestImage       = "curl"
	defaultNATestContainer   = "na-test"
	defaultNATestPodName     = "na-test-pod"
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
	c config.Config

	testNamespace string
	testImage     string
	testContainer string
	testPodName   string
}

// NewNA ...
func NewNA(k Kubernetes, c config.Config) *NA {
	n := &NA{}
	n.k = k
	n.c = c

	n.setup()
	return n
}

// NewDefaultNA ...
func NewDefaultNA() *NA {
	n := &NA{}
	n.k = GetKubeInstance()
	n.c = config.GetEnvConfigInstance()

	n.setup()
	return n
}

func (n *NA) setup() {

	//just default these for now (not sure we'll ever want to supply):
	n.testNamespace = defaultNATestNamespace
	n.testContainer = defaultNATestContainer
	n.testPodName = defaultNATestPodName

	// image repository + curl from config
	// but default if not supplied
	i := *n.c.GetImageRepository()
	if len(i) < 1 {
		i = defaultNAImageRepository
	}
	b := *n.c.GetCurlImage()
	if len(b) < 1 {
		b = defaultNATestImage
	}

	n.testImage = i + "/" + b
}

// ClusterIsDeployed ...
func (n *NA) ClusterIsDeployed() *bool {
	return n.k.ClusterIsDeployed()
}

//SetupNetworkAccessTestPod creates a pod with characteristics required for testing network access.
func (n *NA) SetupNetworkAccessTestPod() (*apiv1.Pod, error) {
	pname, ns, cname, image := GenerateUniquePodName(n.testPodName), n.testNamespace, n.testContainer, n.testImage
	//let caller handle result:
	return n.k.CreatePod(&pname, &ns, &cname, &image, true, nil)	
}

//TeardownNetworkAccessTestPod ...
func (n *NA) TeardownNetworkAccessTestPod(p *string) error {
	_, exists := os.LookupEnv("DONT_DELETE")
	if !exists {
		ns := n.testNamespace
		err := n.k.DeletePod(p, &ns, false) //don't worry about waiting
		return err
	}

	return nil
}

//AccessURL calls the supplied URL and returns the http code
func (n *NA) AccessURL(pn *string, url *string) (int, error) {

	//create a curl command to access the supplied url
	cmd := "curl -s -o /dev/null -I -L -w %{http_code} " + *url
	ns := n.testNamespace
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
