package probr

import (
	"github.com/citihub/probr/internal/coreengine"
	"github.com/citihub/probr/service_packs"
	"github.com/citihub/probr/service_packs/kubernetes"
)

//TODO: revise when interface this bit up ...
var kube = kubernetes.GetKubeInstance()

func RunAllProbes() (int, *coreengine.ProbeStore, error) {
	ts := coreengine.NewProbeStore() // get the test mgr

	for _, probe := range service_packs.GetAllProbes() {
		ts.AddProbe(probe)
	}

	s, err := ts.ExecAllProbes() // Executes all added (queued) tests
	return s, ts, err
}

//GetAllProbeResults maps ProbeStore results to strings
func GetAllProbeResults(ps *coreengine.ProbeStore) map[string]string {
	out := make(map[string]string)
	for name := range ps.Probes {
		results, name, err := readProbeResults(ps, name)
		if err != nil {
			out[name] = err.Error()
		} else {
			out[name] = results
		}
	}
	return out
}

func readProbeResults(ps *coreengine.ProbeStore, name string) (probeResults, probeName string, err error) {
	p, err := ps.GetProbe(name)
	if err != nil {
		return
	}
	probeResults = p.Results.String()
	probeName = p.ProbeDescriptor.Name
	return
}
