package connection

import (
	"log"

	"github.com/probr/probr-pack-kubernetes/internal/config"
	"github.com/probr/probr-sdk/providers/kubernetes/connection"
)

// TODO: Decide whether this 'connection.state' is the best naming convention

// State is a stateful Kubernetes API wrapper
var State *connection.Conn

// Connect initializes connection.State using values from config.Vars.ServicePacks.Kubernetes
func Connect() {
	log.Printf("[DEBUG] Initializing connection with namespace '%s' and context '%s' using kubeconfig: %s",
		config.Vars.ServicePacks.Kubernetes.KubeConfigPath, config.Vars.ServicePacks.Kubernetes.KubeContext, config.Vars.ServicePacks.Kubernetes.ProbeNamespace)
	State = connection.NewConnection(config.Vars.ServicePacks.Kubernetes.KubeConfigPath, config.Vars.ServicePacks.Kubernetes.KubeContext, config.Vars.ServicePacks.Kubernetes.ProbeNamespace)
	log.Print("[DEBUG] Initialized Kubernetes API connection")
}
