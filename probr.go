package probr

import (
	"github.com/citihub/probr/internal/clouddriver/kubernetes"
	"github.com/citihub/probr/internal/coreengine"
	_ "github.com/citihub/probr/probes/clouddriver"
	k8s_probes "github.com/citihub/probr/probes/kubernetes"
)

//TODO: revise when interface this bit up ...
var kube = kubernetes.GetKubeInstance()

func RunAllTests() (int, *coreengine.TestStore, error) {
	ts := coreengine.NewTestManager() // get the test mgr

	for _, probe := range k8s_probes.Probes {
		ts.AddTest(probe.GetGodogTest())
	}

	s, err := ts.ExecAllTests() // Executes all added (queued) tests
	return s, ts, err
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
