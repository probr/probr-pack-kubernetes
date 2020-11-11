package kubernetes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/citihub/probr/internal/config"
	"k8s.io/client-go/kubernetes/scheme"

	apiv1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

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
	CreateAIB(y []byte, ai string, n string, ns string) (bool, error)
	CreateIAMProbePod(y []byte, useDefaultNS bool) (*apiv1.Pod, error)
	DeleteIAMProbePod(n string, useDefaultNS bool, e string) error
	ExecuteVerificationCmd(pn string, cmd IAMProbeCommand, useDefaultNS bool) (*CmdExecutionResult, error)
	GetAccessToken(pn string, useDefaultNS bool) (*string, error)
}

// IAM implements the IdentityAccessManagement interface.
type IAM struct {
	k Kubernetes

	probeNamespace string
	probeImage     string
	probeContainer string
	probePodName   string

	testAzureIdentityBinding string
}

// IAMVerification provides an IAM specific type wrapper extending PSPVerificationProbe.
type IAMVerification struct {
	PSPVerificationProbe
}

// NewDefaultIAM creates a new IAM instance using the default kubernetes provider.
func NewDefaultIAM() *IAM {
	i := &IAM{}
	i.k = GetKubeInstance()

	i.setenv()
	return i
}

func (i *IAM) setenv() {
	//just default these for now (not sure we'll ever want to supply):
	i.probeNamespace = defaultIAMProbeNamespace
	i.probeContainer = defaultIAMProbeContainer
	i.probePodName = defaultIAMProbePodName

	// Extract registry and image info from config
	i.probeImage = config.Vars.ContainerRegistry + "/" + config.Vars.ProbeImage

	i.testAzureIdentityBinding = "probr-specificns-aib"
}

//CreateAIB creates an AzureIdentityBinding to the supplied AzureIdentity
//ai - name of the AzureIdentity
//n - name of AzureIdentityBinding
//ns - namespace in which to create the AIB
func (i *IAM) CreateAIB(y []byte, ai string, n string, ns string) (bool, error) {

	i.createFromYaml(y, nil, &ns, nil, false)
	return false, nil
}

func (i *IAM) createFromYaml(y []byte, pname *string, ns *string, image *string, w bool) (*apiv1.Pod, error) {

	// g := schema.GroupVersionKind{
	// 	Group:   "aadpodidentity.k8s.io",
	// 	Kind:    "AzureIdentityBinding",
	// 	Version: "v1",
	// }

	sch := runtime.NewScheme()
	// sch.Recognizes(g)
	_ = scheme.AddToScheme(sch)

	// decode := serializer.NewCodecFactory(sch).UniversalDeserializer().Decode
	// decode := scheme.Codecs.UniversalDeserializer().Decode

	codec := scheme.Codecs.LegacyCodec(scheme.Scheme.PrioritizedVersionsAllGroups()...)

	unst := unstructured.Unstructured{}
	err := runtime.DecodeInto(codec, y, &unst)

	// o, k, err := decode(y, &g, nil)

	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	log.Printf("unst is %v", unst)
	// p := o.(*apiv1.Pod)
	// //update the name to the one that's supplied
	// p.SetName(*pname)
	// //also update the image (which could have been supplied via the env)
	// //(only expecting one container, but loop in case of many)
	// for _, c := range p.Spec.Containers {
	// 	c.Image = *image
	// }

	// return i.k.CreatePodFromObject(p, pname, ns, w)

	c, _ := i.k.GetClient()

	apiNS := apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}
	c.CoreV1().Namespaces().Create(context.TODO(), &apiNS, metav1.CreateOptions{})

	// r := c.CoreV1().RESTClient().Post().NamespaceIfScoped("test", true).Body(unst)
	r := c.CoreV1().RESTClient().Post().Body(unst)

	res := r.Do(context.TODO())
	log.Printf("[DEBUG] RAW Result: %+v", res)

	if res.Error() != nil {
		return nil, res.Error()
	}

	b, _ := res.Raw()
	bs := string(b)
	log.Printf("[DEBUG] STRING result: %v", bs)

	j := K8SJSON{}
	json.Unmarshal(b, &j)

	log.Printf("[DEBUG] JSON result: %+v", j)

	return nil, nil
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
func (i *IAM) CreateIAMProbePod(y []byte, useDefaultNS bool) (*apiv1.Pod, error) {
	n := GenerateUniquePodName(i.probePodName)

	return i.k.CreatePodFromYaml(y, n, *i.getNamespace(useDefaultNS), i.probeImage,
		*i.getAadPodIDBinding(useDefaultNS), true)
}

// DeleteIAMProbePod deletes the IAM test pod with the supplied name.
func (i *IAM) DeleteIAMProbePod(n string, useDefaultNS bool, e string) error {
	return i.k.DeletePod(&n, i.getNamespace(useDefaultNS), false, e) //don't worry about waiting
}

// ExecuteVerificationCmd executes a verification command against the supplied pod name.
func (i *IAM) ExecuteVerificationCmd(pn string, cmd IAMProbeCommand, useDefaultNS bool) (*CmdExecutionResult, error) {
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
		b = config.Vars.Azure.Identity.DefaultNamespaceAIB
	} else {
		//if not the default namespace, then we are testing a specific
		//identity binding set up as part of the probr run.
		b = i.testAzureIdentityBinding
	}

	return &b
}
