package probr

import (
	"github.com/citihub/probr/internal/clouddriver/kubernetes"
	"github.com/citihub/probr/internal/coreengine"
	_ "github.com/citihub/probr/probes/clouddriver"
	_ "github.com/citihub/probr/probes/kubernetes"
)

//TODO: revise when interface this bit up ...
var kube = kubernetes.GetKubeInstance()

func addTest(tm *coreengine.TestStore, n string, g coreengine.Group, c coreengine.Category) {
	td := coreengine.TestDescriptor{
		Group:    g,
		Category: c,
		Name:     n,
	}
	tm.AddTest(td)
}

func RunAllTests() (int, *coreengine.TestStore, error) {
	tm := coreengine.NewTestManager() // get the test mgr

	//add some tests and add them to the TM - we need to tidy this up!
	addTest(tm, "container_registry_access", coreengine.Kubernetes, coreengine.ContainerRegistryAccess)
	addTest(tm, "internet_access", coreengine.Kubernetes, coreengine.InternetAccess)
	addTest(tm, "pod_security_policy", coreengine.Kubernetes, coreengine.PodSecurityPolicies)
	addTest(tm, "account_manager", coreengine.CloudDriver, coreengine.General)
	addTest(tm, "general", coreengine.Kubernetes, coreengine.General)
	addTest(tm, "iam_control", coreengine.Kubernetes, coreengine.IAM)

	s, err := tm.ExecAllTests() // Executes all added (queued) tests
	return s, tm, err
}

//GetAllTestResults ...
func GetAllTestResults(ts *coreengine.TestStore) (map[string]string, error) {
	out := make(map[string]string)
	for name := range ts.Tests {
		r, n, err := ReadTestResults(ts, name)
		if err != nil {
			return nil, err
		}
		if r != "" {
			out[n] = r
		}
	}
	return out, nil
}

//ReadTestResults ...
func ReadTestResults(ts *coreengine.TestStore, name string) (string, string, error) {
	t, err := ts.GetTest(name)
	test := t
	if err != nil {
		return "", "", err
	}
	r := test.Results
	n := test.TestDescriptor.Name
	if r != nil {
		b := r.Bytes()
		return string(b), n, nil
	}
	return "", "", nil
}
