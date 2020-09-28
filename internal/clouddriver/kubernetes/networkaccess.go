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

// NetworkAccess defines functionality for supporting Network Access tests.
type NetworkAccess interface {
	ClusterIsDeployed() *bool
	SetupNetworkAccessTestPod() (*apiv1.Pod, error)
	TeardownNetworkAccessTestPod(p *string, e string) error
	AccessURL(pn *string, url *string) (int, error)
}

// NA implements NetworkAccess.
type NA struct {
	k Kubernetes

	testNamespace string
	testImage     string
	testContainer string
	testPodName   string
}

// NewNA creates a new instance of NA with the supplied kubernetes instance.
func NewNA(k Kubernetes) *NA {
	n := &NA{}
	n.k = k

	n.setup()
	return n
}

// NewDefaultNA creates a new instance of NA using the default kubernetes instance.
func NewDefaultNA() *NA {
	n := &NA{}
	n.k = GetKubeInstance()

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
	i := config.Vars.Images.Repository
	if len(i) < 1 {
		i = defaultNAImageRepository
	}
	b := config.Vars.Images.Curl
	if len(b) < 1 {
		b = defaultNATestImage
	}

	n.testImage = i + "/" + b
}

// ClusterIsDeployed verifies if a suitable cluster is deployed.
func (n *NA) ClusterIsDeployed() *bool {
	return n.k.ClusterIsDeployed()
}

// SetupNetworkAccessTestPod creates a pod with characteristics required for testing network access.
func (n *NA) SetupNetworkAccessTestPod() (*apiv1.Pod, error) {
	pname, ns, cname, image := GenerateUniquePodName(n.testPodName), n.testNamespace, n.testContainer, n.testImage
	//let caller handle result:
	return n.k.CreatePod(&pname, &ns, &cname, &image, true, nil)
}

// TeardownNetworkAccessTestPod deletes the test pod with the given name.
func (n *NA) TeardownNetworkAccessTestPod(p *string, e string) error {
	_, exists := os.LookupEnv("DONT_DELETE")
	if !exists {
		ns := n.testNamespace
		err := n.k.DeletePod(p, &ns, false, e) //don't worry about waiting
		return err
	}

	return nil
}

// AccessURL calls the supplied URL and returns the http code
func (n *NA) AccessURL(pn *string, url *string) (int, error) {

	//create a curl command to access the supplied url
	cmd := "curl -s -o /dev/null -I -L -w %{http_code} " + *url
	ns := n.testNamespace
	res := n.k.ExecCommand(&cmd, &ns, pn)
	httpCode := res.Stdout

	log.Printf("[NOTICE] URL: %v HTTP Code: %v Exit Code: %v (error: %v)", *url, httpCode, res.Code, res.Err)

	if res.Err != nil && !res.Internal {
		//error which is not internal (so external!)
		//this means code is from the execution of the command on the cluster

		//Check the exit code.  If it's '6' (Couldn't resolve host.)
		//then we want to nil out the error and return the code as this
		//is an expected condition if access is inhibited
		if res.Code == 6 {
			return res.Code, nil
		}
		//otherwise return both code & error
		return res.Code, res.Err
	}

	//no errors, so just return code
	return strconv.Atoi(httpCode)
}
