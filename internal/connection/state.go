package connection

import (
	"log"

	"github.com/citihub/probr-pack-kubernetes/internal/config"
)

// TODO: Decide whether this 'connection.state' is the best naming convention

// State is a stateful Kubernetes API wrapper
var State *Conn

// Connect initializes connection.State using values from config.Vars
func Connect() {
	log.Printf("[DEBUG] Initializing connection with namespace '%s' and context '%s' using kubeconfig: %s",
		config.Vars.KubeConfigPath, config.Vars.KubeContext, config.Vars.ProbeNamespace)
	State = NewConnection(config.Vars.KubeConfigPath, config.Vars.KubeContext, config.Vars.ProbeNamespace)
}
