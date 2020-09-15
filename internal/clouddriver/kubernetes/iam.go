package kubernetes

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"gitlab.com/citihub/probr/internal/config"
	"gitlab.com/citihub/probr/internal/utils"
	"k8s.io/client-go/kubernetes/scheme"

	apiv1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	//default values.  Overrides can be supplied via the environment.
	defaultIAMTestNamespace = "probr-rbac-test-ns"
	//NOTE: either the above namespace needs to be added to the exclusion list on the
	//container registry rule or busybox need to be available in the allowed (probably internal) registry
	defaultIAMImageRepository = "docker.io"
	defaultIAMTestImage       = "busybox"
	defaultIAMTestContainer   = "iam-test"
	defaultIAMTestPodName     = "iam-test-pod"
)

//IdentityAccessManagement encapsulates functionality for querying and probing Identity and Access Management setup
type IdentityAccessManagement interface {
	AzureIdentityExists(ns string) (bool, error)
	AzureIdentityBindingExists(ns string) (bool, error)
	CreateAIB(y []byte, ai string, n string, ns string) (bool, error)
	CreateIAMTestPod(y []byte, identityBinding string) (*apiv1.Pod, error)
	DeleteIAMTestPod(n string) error
	ExecuteVerificationCmd(pn string, cmd PSPTestCommand) (*CmdExecutionResult, error)
}

//IAM implements the IdentityAccessManagement interface
type IAM struct {
	k Kubernetes

	testNamespace string
	testImage     string
	testContainer string
	testPodName   string

	testAzureIdentityBinding string
}

// IAMVerification ...
type IAMVerification struct {
	PSPVerificationProbe
}

// NewDefaultIAM ...
func NewDefaultIAM() *IAM {
	i := &IAM{}
	i.k = GetKubeInstance()

	i.setenv()
	return i
}

func (i *IAM) setenv() {
	//just default these for now (not sure we'll ever want to supply):
	i.testNamespace = defaultIAMTestNamespace
	i.testContainer = defaultIAMTestContainer
	i.testPodName = defaultIAMTestPodName

	// image repository + busy box from config
	// but default if not supplied
	ig := config.Vars.Images.Repository
	if len(ig) < 1 {
		ig = defaultIAMImageRepository
	}
	b := config.Vars.Images.BusyBox
	if len(b) < 1 {
		b = defaultIAMTestImage
	}

	i.testImage = ig + "/" + b

	//TODO: possibly externalise this
	i.testAzureIdentityBinding = "probr-aib"
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

//AzureIdentityExists gets the AzureIdentityBindings and filter for namespace (if supplied)
func (i *IAM) AzureIdentityExists(ns string) (bool, error) {
	//need to make a 'raw' call to get the AIBs
	//the AIB's are in the API group: "apis/aadpodidentity.k8s.io/v1/azureidentity"

	return i.filteredRawResourceGrp("apis/aadpodidentity.k8s.io/v1/azureidentities", "namespace", ns)
}

//AzureIdentityBindingExists gets the AzureIdentityBindings and filter for namespace (if supplied)
func (i *IAM) AzureIdentityBindingExists(ns string) (bool, error) {
	//need to make a 'raw' call to get the AIBs
	//the AIB's are in the API group: "apis/aadpodidentity.k8s.io/v1/azureidentitybindings"

	return i.filteredRawResourceGrp("apis/aadpodidentity.k8s.io/v1/azureidentitybindings", "namespace", ns)	
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

//CreateIAMTestPod ...
func (i *IAM) CreateIAMTestPod(y []byte, identityBinding string) (*apiv1.Pod, error) {
	n := GenerateUniquePodName(i.testPodName)

	//TODO: pass a nil image in for now and take it from the yaml
	return i.k.CreatePodFromYaml(y, &n, utils.StringPtr(i.testNamespace),
		nil, &identityBinding, true)	
}

//DeleteIAMTestPod ...
func (i *IAM) DeleteIAMTestPod(n string) error {
	return i.k.DeletePod(&n, &i.testNamespace, false) //don't worry about waiting
}

// ExecuteVerificationCmd ...
func (i *IAM) ExecuteVerificationCmd(pn string, cmd PSPTestCommand) (*CmdExecutionResult, error) {
	c := cmd.String()
	// ns := i.testNamespace
	ns := "default"
	res := i.k.ExecCommand(&c, &ns, &pn)

	log.Printf("[NOTICE] ExecPSPTestCmd: %v stdout: %v exit code: %v (error: %v)", cmd, res.Stdout, res.Code, res.Err)
	
	return res, nil

}
