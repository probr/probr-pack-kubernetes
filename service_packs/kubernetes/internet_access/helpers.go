package internet_access

import (
	"log"
	"os"
	"strconv"

	"github.com/citihub/probr/audit"
	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/citihub/probr/service_packs/kubernetes"
	"github.com/cucumber/godog"
	apiv1 "k8s.io/api/core/v1"
)

type scenarioState struct {
	name           string
	audit          *audit.ScenarioAudit
	probe          *audit.Probe
	httpStatusCode int
	podName        string
	podState       kubernetes.PodState
}

func beforeScenario(s *scenarioState, probeName string, gs *godog.Scenario) {
	s.name = gs.Name
	s.probe = audit.State.GetProbeLog(probeName)
	s.audit = audit.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	coreengine.LogScenarioStart(gs)
}

const (
	defaultNAProbeContainer = "na-test"
	defaultNAProbePodName   = "na-test-pod"
)

// NetworkAccess defines functionality for supporting Network Access tests.
type NetworkAccess interface {
	ClusterIsDeployed() *bool
	SetupNetworkAccessProbePod(probe *audit.Probe) (*apiv1.Pod, *kubernetes.PodAudit, error)
	TeardownNetworkAccessProbePod(p string, e string) error
	AccessURL(pn *string, url *string) (int, error)
}

// NA implements NetworkAccess.
type NA struct {
	k kubernetes.Kubernetes

	probeImage     string
	probeContainer string
	probePodName   string
}

// NewNA creates a new instance of NA with the supplied kubernetes instance.
func NewNA(k kubernetes.Kubernetes) *NA {
	n := &NA{}
	n.k = k

	n.setup()
	return n
}

// NewDefaultNA creates a new instance of NA using the default kubernetes instance.
func NewDefaultNA() *NA {
	n := &NA{}
	n.k = kubernetes.GetKubeInstance()

	n.setup()
	return n
}

func (n *NA) setup() {

	//just default these for now (not sure we'll ever want to supply):
	n.probeContainer = defaultNAProbeContainer
	n.probePodName = defaultNAProbePodName

	// Extract registry and image info from config
	n.probeImage = config.Vars.ServicePacks.Kubernetes.AuthorisedContainerRegistry + "/" + config.Vars.ServicePacks.Kubernetes.ProbeImage
}

// ClusterIsDeployed verifies if a suitable cluster is deployed.
func (n *NA) ClusterIsDeployed() *bool {
	return n.k.ClusterIsDeployed()
}

// SetupNetworkAccessProbePod creates a pod with characteristics required for testing network access.
func (n *NA) SetupNetworkAccessProbePod(probe *audit.Probe) (*apiv1.Pod, *kubernetes.PodAudit, error) {
	pname, ns, cname, image := kubernetes.GenerateUniquePodName(n.probePodName), kubernetes.Namespace, n.probeContainer, n.probeImage
	//let caller handle result:
	return n.k.CreatePod(pname, ns, cname, image, true, nil, probe)
}

// TeardownNetworkAccessProbePod deletes the test pod with the given name.
func (n *NA) TeardownNetworkAccessProbePod(p string, e string) error {
	_, exists := os.LookupEnv("DONT_DELETE")
	if !exists {
		err := n.k.DeletePod(p, kubernetes.Namespace, e) //don't worry about waiting
		return err
	}

	return nil
}

// AccessURL calls the supplied URL and returns the http code
func (n *NA) AccessURL(pn *string, url *string) (int, error) {

	//create a curl command to access the supplied url
	cmd := "curl -s -o /dev/null -I -L -w %{http_code} " + *url
	res := n.k.ExecCommand(cmd, kubernetes.Namespace, pn)
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
