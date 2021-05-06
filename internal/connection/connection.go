package connection

import (
	"log"

	"github.com/citihub/probr-pack-kubernetes/internal/config"
	"github.com/citihub/probr-sdk/providers/kubernetes/connection"
)

// TODO: Decide whether this 'connection.state' is the best naming convention

// State is a stateful Kubernetes API wrapper
var State *connection.Conn

// Connect initializes connection.State using values from config.Vars.Kube
func Connect() {
	log.Printf("[DEBUG] Initializing connection with namespace '%s' and context '%s' using kubeconfig: %s",
		config.Vars.Kube.KubeConfigPath, config.Vars.Kube.KubeContext, config.Vars.Kube.ProbeNamespace)
	State = connection.NewConnection(config.Vars.Kube.KubeConfigPath, config.Vars.Kube.KubeContext, config.Vars.Kube.ProbeNamespace)
	log.Print("[DEBUG] Initialized Kubernetes API connection")
}
