package probr

import (
	"github.com/citihub/probr/internal/clouddriver/kubernetes"
	"github.com/citihub/probr/internal/coreengine"
	_ "github.com/citihub/probr/probes/clouddriver"
	k8s_probes "github.com/citihub/probr/probes/kubernetes"
)

//TODO: revise when interface this bit up ...
var kube = kubernetes.GetKubeInstance()

func RunAllProbes() (int, *coreengine.ProbeStore, error) {
	ts := coreengine.NewProbeStore() // get the test mgr

	for _, probe := range k8s_probes.Probes {
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
