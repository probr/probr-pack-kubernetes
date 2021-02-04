package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/summary"
	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/citihub/probr/service_packs/kubernetes"
	"github.com/cucumber/godog"

	aibv1 "github.com/Azure/aad-pod-identity/pkg/apis/aadpodidentity"

	apiv1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type scenarioState struct {
	name         string
	audit        *summary.ScenarioAudit
	probe        *summary.Probe
	podState     kubernetes.PodState
	useDefaultNS bool
}

const (
	//default values.  Overrides can be supplied via the environment.
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
	AzureIdentityExists(namespace, aiName string) (bool, error)
	AzureIdentityBindingExists(namespace, aibName string) (bool, error)
	CreateAIB(useDefaultNS bool, aibName, aiName string) error
	CreateIAMProbePod(y []byte, useDefaultNS bool, aibName string, probe *summary.Probe) (*apiv1.Pod, error)
	DeleteIAMProbePod(n string, useDefaultNS bool, e string) error
	ExecuteVerificationCmd(pn string, cmd IAMProbeCommand, useDefaultNS bool) (*kubernetes.CmdExecutionResult, error)
	GetAccessToken(pn string, useDefaultNS bool) (*string, error)
}

// IAM implements the IdentityAccessManagement interface.
type IAM struct {
	k kubernetes.Kubernetes

	probeImage     string
	probeContainer string
	probePodName   string

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
	i.probeContainer = defaultIAMProbeContainer
	i.probePodName = defaultIAMProbePodName

	// Extract registry and image info from config
	i.probeImage = config.Vars.ServicePacks.Kubernetes.AuthorisedContainerRegistry + "/" + config.Vars.ServicePacks.Kubernetes.ProbeImage

	// Set the Azure Identity vars
	// azureIdentitySelector - to allow selection of the binding on pod creation
	i.azureIdentitySelector = "aadpodidbinding"
}

// creates an AzureIdentityBinding object to a specified AzureIdentity in a specified non-default namespace
func (i *IAM) createAIBObject(namespace, aibName, aiName string) runtime.Object {
	// Create an AIB object and assign attributes using input parameters
	aib := aibv1.AzureIdentityBinding{}

	aib.TypeMeta.Kind = "AzureIdentityBinding"
	aib.TypeMeta.APIVersion = "aadpodidentity.k8s.io/v1"
	aib.ObjectMeta.Namespace = namespace
	aib.ObjectMeta.Name = aibName
	aib.Spec.AzureIdentity = aiName
	aib.Spec.Selector = i.azureIdentitySelector

	// Copy into a runtime.Object which is required for the api request
	runtimeAib := aib.DeepCopyObject()

	return runtimeAib
}

// CreateAIB creates an AzureIdentityBinding in the cluster
func (i *IAM) CreateAIB(useDefaultNS bool, aibName, aiName string) error {

	// Obtain the kubernetes cluster client connection
	c, _ := i.k.GetClient()
	namespace := i.getNamespace(useDefaultNS)

	if !useDefaultNS {
		// Create the non-default namespace (will succeed if the ns already exists)
		apiNS := apiv1.Namespace{}
		apiNS.ObjectMeta.Name = namespace
		c.CoreV1().Namespaces().Create(context.TODO(), &apiNS, metav1.CreateOptions{})
	}

	// Create an runtime AIB object and assign attributes using input parameters
	runtimeAib := i.createAIBObject(namespace, aibName, aiName)

	// set the api path for the aadpodidentity package which include the azureidentitybindings custom resource definition
	apiPath := "apis/aadpodidentity.k8s.io/v1"

	//	Create a rest api client Post request object
	request := c.CoreV1().RESTClient().Post().AbsPath(apiPath).Namespace(namespace).Resource("azureidentitybindings").Body(runtimeAib)

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
func (i *IAM) AzureIdentityExists(namespace, aiName string) (bool, error) {
	//need to make a 'raw' call to get the AIBs
	//the AIB's are in the API group: "apis/aadpodidentity.k8s.io/v1/azureidentity"

	return i.filteredRawResourceGrp("apis/aadpodidentity.k8s.io/v1/azureidentities", namespace, aiName)
}

// AzureIdentityBindingExists gets the AzureIdentityBindings and filter for namespace (if supplied)
func (i *IAM) AzureIdentityBindingExists(namespace, aibName string) (bool, error) {
	//need to make a 'raw' call to get the AIBs
	//the AIB's are in the API group: "apis/aadpodidentity.k8s.io/v1/azureidentitybindings"

	return i.filteredRawResourceGrp("apis/aadpodidentity.k8s.io/v1/azureidentitybindings", namespace, aibName)
}

func (i *IAM) filteredRawResourceGrp(apiGroup string, namespace string, resourceName string) (bool, error) {
	j, err := i.k.GetRawResourcesByGrp(apiGroup)

	if err != nil {
		return false, err
	}

	for _, r := range j.Items {
		if (r.Metadata["namespace"] == namespace) && (r.Metadata["name"] == resourceName) {
			return true, nil
		}
	}

	//false if none found in group g with key k and prefix f
	return false, nil
}

// CreateIAMProbePod creates a pod configured for IAM test cases.
func (i *IAM) CreateIAMProbePod(y []byte, useDefaultNS bool, aibName string, probe *summary.Probe) (*apiv1.Pod, error) {
	n := kubernetes.GenerateUniquePodName(i.probePodName)

	pod, err := i.k.CreatePodFromYaml(y, n, i.getNamespace(useDefaultNS),
		i.probeImage, *i.getAadPodIDBinding(useDefaultNS, aibName), true, probe)
	return pod, err
}

// DeleteIAMProbePod deletes the IAM test pod with the supplied name.
func (i *IAM) DeleteIAMProbePod(n string, useDefaultNS bool, e string) error {
	return i.k.DeletePod(n, i.getNamespace(useDefaultNS), e) //don't worry about waiting
}

// ExecuteVerificationCmd executes a verification command against the supplied pod name.
func (i *IAM) ExecuteVerificationCmd(pn string, cmd IAMProbeCommand, useDefaultNS bool) (*kubernetes.CmdExecutionResult, error) {
	c := cmd.String()
	ns := i.getNamespace(useDefaultNS)
	res := i.k.ExecCommand(c, ns, &pn)

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

func (i *IAM) getNamespace(useDefaultNS bool) string {
	if useDefaultNS {
		return "default"
	}
	return kubernetes.Namespace
}

func (i *IAM) getAadPodIDBinding(useDefaultNS bool, aibName string) *string {
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
		b = aibName
	}

	return &b
}

func beforeScenario(s *scenarioState, probeName string, gs *godog.Scenario) {
	s.name = gs.Name
	s.probe = summary.State.GetProbeLog(probeName)
	s.audit = summary.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	coreengine.LogScenarioStart(gs)
}
