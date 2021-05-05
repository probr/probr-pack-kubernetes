package connection

import (
	"log"

	"github.com/citihub/probr-pack-kubernetes/internal/config"
	"github.com/citihub/probr-sdk/providers/kubernetes/connection"
)

// TODO: Decide whether this 'connection.state' is the best naming convention

// State is a stateful Kubernetes API wrapper
var State *connection.Conn

// Connect initializes connection.State using values from config.Vars
func Connect() {
	log.Printf("[DEBUG] Initializing connection with namespace '%s' and context '%s' using kubeconfig: %s",
		config.Vars.KubeConfigPath, config.Vars.KubeContext, config.Vars.ProbeNamespace)
	State = connection.NewConnection(config.Vars.KubeConfigPath, config.Vars.KubeContext, config.Vars.ProbeNamespace)
	log.Print("[DEBUG] Initialized Kubernetes API connection")
}
