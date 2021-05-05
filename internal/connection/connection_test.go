package connection

import (
	"testing"

	"github.com/citihub/probr-pack-kubernetes/internal/config"
)

func TestConnect(t *testing.T) {
	tests := []struct {
		name           string
		kubeConfigPath string
		kubeContext    string
		probeNamespace string
	}{
		{
			name: "Validate Connect() updates connection.State",
		},
		{
			name:           "Validate Connect() updates connection.State",
			kubeConfigPath: "fakeConfigPath",
		},
		{
			name:           "Validate Connect() updates connection.State",
			kubeConfigPath: "fakeConfigPath2",
			kubeContext:    "fakeContext",
		},
		{
			name:           "Validate Connect() updates connection.State",
			kubeConfigPath: "fakeConfigPath3",
			kubeContext:    "fakeContext2",
			probeNamespace: "fakeNamespace",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			State = nil
			config.Vars.KubeConfigPath = tt.kubeConfigPath
			config.Vars.KubeContext = tt.kubeContext
			config.Vars.ProbeNamespace = tt.probeNamespace
			Connect()
			if State == nil {
				t.Error("Connect failed to set connection.State")
			}
			if State.KubeConfigPath != tt.kubeConfigPath ||
				State.KubeContext != tt.kubeContext ||
				State.Namespace != tt.probeNamespace {
				t.Error("Connect failed to use config.Vars")
			}
		})
	}
}
