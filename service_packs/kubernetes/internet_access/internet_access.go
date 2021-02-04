package internet_access

import (
	"fmt"
	"log"

	"github.com/cucumber/godog"
	apiv1 "k8s.io/api/core/v1"

	"github.com/citihub/probr/internal/utils"
	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/citihub/probr/service_packs/kubernetes"
)

type ProbeStruct struct{}

var Probe ProbeStruct

var ia_ps scenarioState

func init() {
	ia_ps = scenarioState{}
}

// NetworkAccess is the section of the kubernetes package which provides the kubernetes interactions required to support
// network access scenarios.
var na NetworkAccess

// SetNetworkAccess allows injection of a specific NetworkAccess helper.
func SetNetworkAccess(n NetworkAccess) {
	na = n
}

// General
func (s *scenarioState) aKubernetesClusterIsDeployed() error {
	description, payload, err := kubernetes.ClusterIsDeployed()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()
	return err //  ClusterIsDeployed will create a fatal error if kubeconfig doesn't validate
}

func (s *scenarioState) aPodIsDeployedInTheCluster() error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	var podAudit *kubernetes.PodAudit
	var pod *apiv1.Pod
	if s.podName != "" {
		//only one pod is needed for all scenarios in this probe
		log.Printf("[DEBUG] Pod %v has already been created - reusing the pod", s.podName)
	} else {
		pd, pa, e := na.SetupNetworkAccessProbePod(s.probe)
		podAudit = pa
		pod = pd
		if e != nil {
			err = e
		} else if pod == nil {
			err = utils.ReformatError("Failed to setup network access test pod")
			log.Print(err)
		} else {
			s.podName = pod.GetObjectMeta().GetName()
		}
	}

	description = fmt.Sprintf("Verifying the Pod %s deployed in the cluster", s.podState.PodName)
	payload = kubernetes.PodPayload{Pod: pod, PodAudit: podAudit}

	return err
}

func (s *scenarioState) aProcessInsideThePodEstablishesADirectHTTPSConnectionTo(url string) error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	code, err := na.AccessURL(&s.podName, &url)

	if err != nil {
		err = utils.ReformatError("[ERROR] Error raised when attempting to access URL: %v", err)
		log.Print(err)
	}

	//hold on to the code
	s.httpStatusCode = code

	description = fmt.Sprintf("Proces inside the pod established http connection with url '%s',", url)
	payload = struct {
		PodState kubernetes.PodState
		PodName  string
	}{s.podState, s.podState.PodName}

	return err
}

func (s *scenarioState) accessIs(accessResult string) error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	if accessResult == "blocked" {
		//then the result should be anything other than 200
		if s.httpStatusCode == 200 {
			//it's a fail:
			err = utils.ReformatError("got HTTP Status Code %v - failed", s.httpStatusCode)
		}
	}

	description = fmt.Sprintf("The access result is %s,", accessResult)
	payload = struct {
		PodState kubernetes.PodState
		PodName  string
	}{s.podState, s.podState.PodName}

	return err
}

func (p ProbeStruct) Name() string {
	return "internet_access"
}

func (p ProbeStruct) Path() string {
	return coreengine.GetFeaturePath("service_packs", "kubernetes", p.Name())
}

// iaProbeInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func (p ProbeStruct) ProbeInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {}) //nothing for now

	ctx.AfterSuite(func() {
		na.TeardownNetworkAccessProbePod(ia_ps.podName, p.Name())
	})

	//check dependencies ...
	if na == nil {
		// not been given one so set default
		na = NewDefaultNA()
	}
}

// iaScenarioInitialize initialises the specific test steps.  This is essentially the creation of the test
// which reflects the tests described in the events directory.  There must be a test step registered for
// each line in the feature files. Note: Godog will output stub steps and implementations if it doesn't find
// a step / function defined.  See: https://github.com/cucumber/godog#example.
func (p ProbeStruct) ScenarioInitialize(ctx *godog.ScenarioContext) {

	ctx.BeforeScenario(func(s *godog.Scenario) {
		beforeScenario(&ia_ps, p.Name(), s)
	})

	ctx.Step(`^a Kubernetes cluster is deployed$`, ia_ps.aKubernetesClusterIsDeployed)
	ctx.Step(`^a pod is deployed in the cluster$`, ia_ps.aPodIsDeployedInTheCluster)
	ctx.Step(`^a process inside the pod establishes a direct http\(s\) connection to "([^"]*)"$`, ia_ps.aProcessInsideThePodEstablishesADirectHTTPSConnectionTo)
	ctx.Step(`^access is "([^"]*)"$`, ia_ps.accessIs)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		ia_ps.httpStatusCode = 0
		coreengine.LogScenarioEnd(s)
	})
}
