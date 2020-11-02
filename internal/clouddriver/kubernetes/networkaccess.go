package kubernetes

import (
	"log"
	"os"
	"strconv"

	"github.com/citihub/probr/internal/config"
	apiv1 "k8s.io/api/core/v1"
)

const (
	//default values.  Overrides can be set via the environment.
	defaultNAProbeNamespace  = "probr-network-access-test-ns" //this needs to be set up as an exculsion in the image registry policy
	defaultNAImageRepository = "curlimages"
	defaultNAProbeImage      = "curl"
	defaultNAProbeContainer  = "na-test"
	defaultNAProbePodName    = "na-test-pod"
)

// NetworkAccess defines functionality for supporting Network Access tests.
type NetworkAccess interface {
	ClusterIsDeployed() *bool
	SetupNetworkAccessProbePod() (*apiv1.Pod, *PodAudit, error)
	TeardownNetworkAccessProbePod(p *string, e string) error
	AccessURL(pn *string, url *string) (int, error)
}

// NA implements NetworkAccess.
type NA struct {
	k Kubernetes

	probeNamespace string
	probeImage     string
	probeContainer string
	probePodName   string
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
	n.probeNamespace = defaultNAProbeNamespace
	n.probeContainer = defaultNAProbeContainer
	n.probePodName = defaultNAProbePodName

	// image repository + curl from config
	// but default if not supplied
	i := config.Vars.ImagesRepository
	//need to fudge for 'curl' as it's registered as curlimages/curl
	//on docker, so if we've been given a repository from the config
	//and it's 'docker.io' then ignore it and set default (curlimages)
	if len(i) < 1 || i == "docker.io" {
		i = defaultNAImageRepository
	}

	n.probeImage = i + "/" + defaultNAProbeImage
}

// ClusterIsDeployed verifies if a suitable cluster is deployed.
func (n *NA) ClusterIsDeployed() *bool {
	return n.k.ClusterIsDeployed()
}

// SetupNetworkAccessProbePod creates a pod with characteristics required for testing network access.
func (n *NA) SetupNetworkAccessProbePod() (*apiv1.Pod, *PodAudit, error) {
	pname, ns, cname, image := GenerateUniquePodName(n.probePodName), n.probeNamespace, n.probeContainer, n.probeImage
	//let caller handle result:
	return n.k.CreatePod(pname, ns, cname, image, true, nil)
}

// TeardownNetworkAccessProbePod deletes the test pod with the given name.
func (n *NA) TeardownNetworkAccessProbePod(p *string, e string) error {
	_, exists := os.LookupEnv("DONT_DELETE")
	if !exists {
		ns := n.probeNamespace
		err := n.k.DeletePod(p, &ns, false, e) //don't worry about waiting
		return err
	}

	return nil
}

// AccessURL calls the supplied URL and returns the http code
func (n *NA) AccessURL(pn *string, url *string) (int, error) {

	//create a curl command to access the supplied url
	cmd := "curl -s -o /dev/null -I -L -w %{http_code} " + *url
	ns := n.probeNamespace
	res := n.k.ExecCommand(&cmd, &ns, pn)
	httpCode := res.Stdout

	log.Printf("[INFO] URL: %v HTTP Code: %v Exit Code: %v (error: %v)", *url, httpCode, res.Code, res.Err)

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
