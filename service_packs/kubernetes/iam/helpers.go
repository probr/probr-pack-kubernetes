package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/coreengine"
	"github.com/citihub/probr/internal/summary"
	"github.com/citihub/probr/service_packs/kubernetes"
	"github.com/cucumber/godog"

	aibv1 "github.com/Azure/aad-pod-identity/pkg/apis/aadpodidentity"

	apiv1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type scenarioState struct {
	name           string
	audit          *summary.ScenarioAudit
	probe          *summary.Probe
	httpStatusCode int
	podName        string
	podState       kubernetes.PodState
	useDefaultNS   bool
	wildcardRoles  interface{}
}

const (
	//default values.  Overrides can be supplied via the environment.
	defaultIAMProbeNamespace = "probr-rbac-test-ns"
	//NOTE: either the above namespace needs to be added to the exclusion list on the
	//container registry rule or image need to be available in the allowed (probably internal) registry
	defaultIAMProbeContainer = "iam-test"
	defaultIAMProbePodName   = "iam-test-pod"
)

// IAMProbeCommand defines commands for use in testing IAM
type IAMProbeCommand int

// enum supporting IAMProbeCommand
const (
	CatAzJSON IAMProbeCommand = iota
	CurlAuthToken
)

func (c IAMProbeCommand) String() string {
	return [...]string{"cat /etc/kubernetes/azure.json",
		"curl http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https%3A%2F%2Fmanagement.azure.com%2F -H Metadata:true -s"}[c]
}

// IdentityAccessManagement encapsulates functionality for querying and probing Identity and Access Management setup.
type IdentityAccessManagement interface {
	AzureIdentityExists(useDefaultNS bool) (bool, error)
	AzureIdentityBindingExists(useDefaultNS bool) (bool, error)
	CreateAIB() error
	CreateIAMProbePod(y []byte, useDefaultNS bool, probe *summary.Probe) (*apiv1.Pod, error)
	DeleteIAMProbePod(n string, useDefaultNS bool, e string) error
	ExecuteVerificationCmd(pn string, cmd IAMProbeCommand, useDefaultNS bool) (*kubernetes.CmdExecutionResult, error)
	GetAccessToken(pn string, useDefaultNS bool) (*string, error)
}

// IAM implements the IdentityAccessManagement interface.
type IAM struct {
	k kubernetes.Kubernetes

	probeNamespace string
	probeImage     string
	probeContainer string
	probePodName   string

	azureIdentityBinding  string
	azureIdentityName     string
	azureIdentitySelector string
}

// NewDefaultIAM creates a new IAM instance using the default kubernetes provider.
func NewDefaultIAM() *IAM {
	i := &IAM{}
	i.k = kubernetes.GetKubeInstance()

	i.setenv()
	return i
}

func (i *IAM) setenv() {

	//just default these for now (not sure we'll ever want to supply):
	i.probeNamespace = defaultIAMProbeNamespace
	i.probeContainer = defaultIAMProbeContainer
	i.probePodName = defaultIAMProbePodName

	// Extract registry and image info from config
	i.probeImage = config.Vars.AuthorisedContainerRegistry + "/" + config.Vars.ProbeImage

	// Set the Azure Identity vars
	// azureIdentityBinding - name of AzureIdentityBinding
	// azureIdentityName - name of the AzureIdentity
	// azureIdentitySelector - to allow selection of the binding on pod creation
	i.azureIdentityBinding = "probr-aib"
	i.azureIdentityName = "probr-probe"
	i.azureIdentitySelector = "aadpodidbinding"
}

// creates an AzureIdentityBinding object to a specified AzureIdentity in a specified non-default namespace
func (i *IAM) createAIBObject() runtime.Object {
	// Create an AIB object and assign attributes using input parameters
	aib := aibv1.AzureIdentityBinding{}

	aib.TypeMeta.Kind = "AzureIdentityBinding"
	aib.TypeMeta.APIVersion = "aadpodidentity.k8s.io/v1"
	aib.ObjectMeta.Name = i.azureIdentityBinding
	aib.ObjectMeta.Namespace = i.probeNamespace
	aib.Spec.AzureIdentity = i.azureIdentityName
	aib.Spec.Selector = i.azureIdentitySelector

	// Copy into a runtime.Object which is required for the api request
	runtimeAib := aib.DeepCopyObject()

	return runtimeAib
}

// CreateAIB creates an AzureIdentityBinding in the cluster
func (i *IAM) CreateAIB() error {

	// Obtain the kubernetes cluster client connection
	c, _ := i.k.GetClient()

	// Create an runtime AIB object and assign attributes using input parameters
	runtimeAib := i.createAIBObject()

	// Create the namespace
	apiNS := apiv1.Namespace{}
	apiNS.ObjectMeta.Name = i.probeNamespace
	c.CoreV1().Namespaces().Create(context.TODO(), &apiNS, metav1.CreateOptions{})

	// set the api path for the aadpodidentity package which include the azureidentitybindings custom resource definition
	apiPath := "apis/aadpodidentity.k8s.io/v1"

	//	Create a rest api client Post request object
	request := c.CoreV1().RESTClient().Post().AbsPath(apiPath).Namespace(i.probeNamespace).Resource("azureidentitybindings").Body(runtimeAib)

	//	Call the api to execute the Post request and create the AIB in the cluster
	response := request.Do(context.TODO())
	log.Printf("[DEBUG] RAW Result: %+v", response)

	//	Handle response error - ignore the error if the AIB already exists
	if (response.Error() != nil) && (!errors.IsAlreadyExists(response.Error())) {
		log.Printf("[DEBUG] AIB creation Error: %+v", response.Error())
		return response.Error()
	}

	b, _ := response.Raw()
	bs := string(b)
	log.Printf("[DEBUG] STRING result: %+v", bs)

	j := kubernetes.K8SJSON{}
	json.Unmarshal(b, &j)

	log.Printf("[DEBUG] JSON result: %+v", j)

	return nil
}

// AzureIdentityExists gets the AzureIdentityBindings and filter for namespace (if supplied)
func (i *IAM) AzureIdentityExists(useDefaultNS bool) (bool, error) {
	//need to make a 'raw' call to get the AIBs
	//the AIB's are in the API group: "apis/aadpodidentity.k8s.io/v1/azureidentity"

	return i.filteredRawResourceGrp("apis/aadpodidentity.k8s.io/v1/azureidentities", "namespace", *i.getNamespace(useDefaultNS))
}

// AzureIdentityBindingExists gets the AzureIdentityBindings and filter for namespace (if supplied)
func (i *IAM) AzureIdentityBindingExists(useDefaultNS bool) (bool, error) {
	//need to make a 'raw' call to get the AIBs
	//the AIB's are in the API group: "apis/aadpodidentity.k8s.io/v1/azureidentitybindings"

	return i.filteredRawResourceGrp("apis/aadpodidentity.k8s.io/v1/azureidentitybindings", "namespace", *i.getNamespace(useDefaultNS))
}

func (i *IAM) filteredRawResourceGrp(g string, k string, f string) (bool, error) {
	j, err := i.k.GetRawResourcesByGrp(g)

	if err != nil {
		return false, err
	}

	for _, r := range j.Items {
		n := r.Metadata[k]
		if strings.HasPrefix(n, f) {
			return true, nil
		}
	}

	//false if none found in group g with key k and prefix f
	return false, nil
}

// CreateIAMProbePod creates a pod configured for IAM test cases.
func (i *IAM) CreateIAMProbePod(y []byte, useDefaultNS bool, probe *summary.Probe) (*apiv1.Pod, error) {
	n := kubernetes.GenerateUniquePodName(i.probePodName)

	pod, err := i.k.CreatePodFromYaml(y, n, *i.getNamespace(useDefaultNS),
		i.probeImage, *i.getAadPodIDBinding(useDefaultNS), true, probe)
	return pod, err
}

// DeleteIAMProbePod deletes the IAM test pod with the supplied name.
func (i *IAM) DeleteIAMProbePod(n string, useDefaultNS bool, e string) error {
	return i.k.DeletePod(&n, i.getNamespace(useDefaultNS), false, e) //don't worry about waiting
}

// ExecuteVerificationCmd executes a verification command against the supplied pod name.
func (i *IAM) ExecuteVerificationCmd(pn string, cmd IAMProbeCommand, useDefaultNS bool) (*kubernetes.CmdExecutionResult, error) {
	c := cmd.String()
	res := i.k.ExecCommand(&c, i.getNamespace(useDefaultNS), &pn)

	log.Printf("[NOTICE] ExecuteVerificationCmd: %v stdout: %v exit code: %v (error: %v)", cmd, res.Stdout, res.Code, res.Err)

	return res, nil

}

// GetAccessToken attempts to retrieve an access token by executing a curl command requesting a token for the Azure Resource Manager.
func (i *IAM) GetAccessToken(pn string, useDefaultNS bool) (*string, error) {
	//curl for the auth token ... need to supply appropiate ns
	res, err := i.ExecuteVerificationCmd(pn, CurlAuthToken, useDefaultNS)
	log.Printf("[NOTICE] curl result: %v", res)

	if err != nil {
		//this is an error from trying to execute the command as opposed to
		//the command itself returning an error
		return nil, fmt.Errorf("error raised trying to execute auth token command - %v", err)
	}

	//try and extract token
	var a struct {
		AccessToken string `json:"access_token,omitempty"`
	}
	json.Unmarshal([]byte(res.Stdout), &a)

	log.Printf("[DEBUG] Access Token JSON result: %+v", a)

	return &a.AccessToken, nil
}

func (i *IAM) getNamespace(useDefaultNS bool) *string {
	var ns string
	if useDefaultNS {
		ns = "default"
	} else {
		ns = i.probeNamespace
	}

	return &ns
}

func (i *IAM) getAadPodIDBinding(useDefaultNS bool) *string {
	//return the value for the following pod label
	// labels:
	// 	aadpodidbinding:
	//the value in this label should match the selector value in
	//the AzureIdentityBinding

	var b string
	if useDefaultNS {
		//if the default namespace, we can get the value from the config
		//this can be specified via config file or env and could vary
		//between deployment situations.  If not supplied the default
		//will be returned.
		b = config.Vars.CloudProviders.Azure.Identity.DefaultNamespaceAIB
	} else {
		//if not the default namespace, then we are testing a specific
		//identity binding set up as part of the probr run.
		b = i.azureIdentityBinding
	}

	return &b
}

func beforeScenario(s *scenarioState, probeName string, gs *godog.Scenario) {
	s.name = gs.Name
	s.probe = summary.State.GetProbeLog(probeName)
	s.audit = summary.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	coreengine.LogScenarioStart(gs)
}
