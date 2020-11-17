package probr

import (
	"github.com/citihub/probr/internal/coreengine"
	_ "github.com/citihub/probr/probes/clouddriver"
	"github.com/citihub/probr/probes/kubernetes"
	k8s_logic "github.com/citihub/probr/probes/kubernetes/probe_logic"
)

//TODO: revise when interface this bit up ...
var kube = k8s_logic.GetKubeInstance()

func RunAllProbes() (int, *coreengine.ProbeStore, error) {
	ts := coreengine.NewProbeStore() // get the test mgr

	for _, probe := range kubernetes.Probes {
		ts.AddProbe(probe.GetGodogProbe())
	}

	s, err := ts.ExecAllProbes() // Executes all added (queued) tests
	return s, ts, err
}

//GetAllProbeResults ...
func GetAllProbeResults(ps *coreengine.ProbeStore) (map[string]string, error) {
	out := make(map[string]string)
	for name := range ps.Probes {
		r, n, err := ReadProbeResults(ps, name)
		if err != nil {
			return nil, err
		}
		if r != "" {
			out[n] = r
		}
	}
	return out, nil
}

//ReadProbeResults ...
func ReadProbeResults(ps *coreengine.ProbeStore, name string) (string, string, error) {
	p, err := ps.GetProbe(name)
	probe := p
	if err != nil {
		return "", "", err
	}
	r := probe.Results
	n := probe.ProbeDescriptor.Name
	if r != nil {
		b := r.Bytes()
		return string(b), n, nil
	}
	return "", "", nil
}
