// Package general provides the implementation required to execute the feature-based test cases
// described in the the 'events' directory.
package general

import (
	"fmt"
	"strings"

	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/api/rbac/v1"

	"github.com/cucumber/godog"

	"github.com/citihub/probr/config"
	"github.com/citihub/probr/internal/utils"
	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/citihub/probr/service_packs/kubernetes"
)

type ProbeStruct struct{}

var Probe ProbeStruct

// General
func (s *scenarioState) aKubernetesClusterIsDeployed() error {
	description, payload, error := kubernetes.ClusterIsDeployed()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, error)
	}()
	return error //  ClusterIsDeployed will create a fatal error if kubeconfig doesn't validate
}

//@CIS-5.1.3
func (s *scenarioState) iInspectTheThatAreConfigured(roleLevel string) error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	if roleLevel == "Cluster Roles" {
		l, e := kubernetes.GetKubeInstance().GetClusterRolesByResource("*")
		err = e
		s.wildcardRoles = l
	} else if roleLevel == "Roles" {
		l, e := kubernetes.GetKubeInstance().GetRolesByResource("*")
		err = e
		s.wildcardRoles = l
	}
	if err != nil {
		err = utils.ReformatError("error raised when retrieving '%v': %v", roleLevel, err)
	}

	description = fmt.Sprintf("Ensures that %s are configured. Retains wildcard roles in state for following steps. Passes if retrieval command does not have error.", roleLevel)
	payload = struct {
		PodState kubernetes.PodState
		PodName  string
	}{s.podState, s.podState.PodName}
	return err
}

func (s *scenarioState) iShouldOnlyFindWildcardsInKnownAndAuthorisedConfigurations() error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	//we strip out system/known entries in the cluster roles & roles call
	var wildcardCount int
	//	wildcardCount := len(s.wildcardRoles.([]interface{}))
	switch s.wildcardRoles.(type) {
	case *[]v1.Role:
		wildCardRoles := s.wildcardRoles.(*[]rbacv1.Role)
		wildcardCount = len(*wildCardRoles)
	case *[]v1.ClusterRole:
		wildCardRoles := s.wildcardRoles.(*[]rbacv1.ClusterRole)
		wildcardCount = len(*wildCardRoles)
	default:
	}

	if wildcardCount > 0 {
		err = utils.ReformatError("roles exist with wildcarded resources")
	}

	description = "Examines scenario state's wildcard roles. Passes if no wildcard roles are found."
	payload = struct {
		PodState kubernetes.PodState
		PodName  string
	}{s.podState, s.podState.PodName}

	return err
}

//@CIS-5.6.3
func (s *scenarioState) iAttemptToCreateADeploymentWhichDoesNotHaveASecurityContext() error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	cname := "probr-general"
	podName := kubernetes.GenerateUniquePodName(cname)
	image := config.Vars.ServicePacks.Kubernetes.AuthorisedContainerRegistry + "/" + config.Vars.ServicePacks.Kubernetes.ProbeImage

	//create pod with nil security context
	pod, podAudit, err := kubernetes.GetKubeInstance().CreatePod(podName, "probr-general-test-ns", cname, image, true, nil, s.probe)

	err = kubernetes.ProcessPodCreationResult(&s.podState, pod, kubernetes.UndefinedPodCreationErrorReason, err)

	description = "Attempts to create a deployment without a security context. Retains the status of the deployment in scenario state for following steps. Passes if created, or if an expected error is encountered."
	payload = kubernetes.PodPayload{Pod: pod, PodAudit: podAudit}
	return err
}

func (s *scenarioState) theDeploymentIsRejected() error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	//looking for a non-nil creation error
	if s.podState.CreationError == nil {
		err = utils.ReformatError("pod %v was created successfully. Test fail.", s.podState.PodName)
	}

	description = "Looks for a creation error on the current scenario state. Passes if error is found, because it should have been rejected."
	payload = struct {
		PodState kubernetes.PodState
		PodName  string
	}{s.podState, s.podState.PodName}

	return err
}

//@CIS-6.10.1
// PENDING IMPLEMENTATION
func (s *scenarioState) iShouldNotBeAbleToAccessTheKubernetesWebUI() error {
	//TODO: will be difficult to test this.  To access it, a proxy needs to be created:
	//az aks browse --resource-group rg-probr-all-policies --name ProbrAllPolicies
	//which will then open a browser at:
	//http://127.0.0.1:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#/login
	//I don't think this is going to be easy to do from here
	//Is there another test?  Or is it sufficient to verify that no kube-dashboard is running?

	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	return godog.ErrPending
}

func (s *scenarioState) theKubernetesWebUIIsDisabled() error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	//look for the dashboard pod in the kube-system ns
	pl, err := kubernetes.GetKubeInstance().GetPods("kube-system")
	var name string

	if err != nil {
		err = utils.ReformatError("Probe step not run. Error raised when trying to retrieve pods: %v", err)
	} else {
		//a "pass" is the absence of a "kubernetes-dashboard" pod
		for _, v := range pl.Items {
			if strings.HasPrefix(v.Name, "kubernetes-dashboard") {
				err = utils.ReformatError("kubernetes-dashboard pod found (%v) - test fail", v.Name)
				name = v.Name
			}
		}
	}

	description = "Attempts to find a pod in the 'kube-system' namespace with the prefix 'kubernetes-dashboard'. Passes if no pod is returned."
	payload = struct {
		PodState         kubernetes.PodState
		PodName          string
		PodDashBoardName string
	}{s.podState, s.podState.PodName, name}

	return err
}

func (p ProbeStruct) Name() string {
	return "general"
}

func (p ProbeStruct) Path() string {
	return coreengine.GetFeaturePath("service_packs", "kubernetes", p.Name())
}

// genProbeInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func (p ProbeStruct) ProbeInitialize(ctx *godog.TestSuiteContext) {

	ctx.BeforeSuite(func() {}) //nothing for now

	ctx.AfterSuite(func() {})

}

// genScenarioInitialize initialises the specific test steps.  This is essentially the creation of the test
// which reflects the tests described in the events directory.  There must be a test step registered for
// each line in the feature files. Note: Godog will output stub steps and implementations if it doesn't find
// a step / function defined.  See: https://github.com/cucumber/godog#example.
func (p ProbeStruct) ScenarioInitialize(ctx *godog.ScenarioContext) {
	ps := scenarioState{}

	ctx.BeforeScenario(func(s *godog.Scenario) {
		beforeScenario(&ps, p.Name(), s)
	})

	//general
	ctx.Step(`^a Kubernetes cluster is deployed$`, ps.aKubernetesClusterIsDeployed)

	//@CIS-5.1.3
	ctx.Step(`^I inspect the "([^"]*)" that are configured$`, ps.iInspectTheThatAreConfigured)
	ctx.Step(`^I should only find wildcards in known and authorised configurations$`, ps.iShouldOnlyFindWildcardsInKnownAndAuthorisedConfigurations)

	//@CIS-5.6.3
	ctx.Step(`^I attempt to create a deployment which does not have a Security Context$`, ps.iAttemptToCreateADeploymentWhichDoesNotHaveASecurityContext)
	ctx.Step(`^the deployment is rejected$`, ps.theDeploymentIsRejected)

	ctx.Step(`^I should not be able to access the Kubernetes Web UI$`, ps.iShouldNotBeAbleToAccessTheKubernetesWebUI)
	ctx.Step(`^the Kubernetes Web UI is disabled$`, ps.theKubernetesWebUIIsDisabled)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		kubernetes.GetKubeInstance().DeletePod(ps.podState.PodName, "probr-general-test-ns", p.Name())
		coreengine.LogScenarioEnd(s)
	})
}
