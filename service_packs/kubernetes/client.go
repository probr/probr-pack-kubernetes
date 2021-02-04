package kubernetes

import (
	"log"

	"github.com/citihub/probr/config"

	k8s "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	//needed for authentication against the various GCPs
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

//GetClient gets a client connection to the Kubernetes cluster specified via config.Vars.ServicePacks.Kubernetes.KubeConfigPath
func (k *Kube) GetClient() (*k8s.Clientset, error) {
	k.clientMutex.Lock()
	defer k.clientMutex.Unlock()

	if k.kubeClient != nil {
		// Singleton
		return k.kubeClient, nil
	}

	clientConfig, err := getClientConfig()
	if err != nil {
		return nil, err
	}

	// create the clientset (note: assigned to global "kubeClient")
	k.kubeClient, err = k8s.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}
	return k.kubeClient, nil
}

// getClientConfig retrieves config and changes context if needed prior to returning
func getClientConfig() (*restclient.Config, error) {
	// Adapted from clientcmd.BuildConfigFromFlags:
	// https://github.com/kubernetes/client-go/blob/5ab99756f65dbf324e5adf9bd020a20a024bad85/tools/clientcmd/client_config.go#L606

	k := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: config.Vars.ServicePacks.Kubernetes.KubeConfigPath},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: ""}})
	rawConfig, _ := k.RawConfig()
	if config.Vars.ServicePacks.Kubernetes.KubeContext == "" {
		log.Printf("[NOTICE] Initializing client with default context from provided config")
	} else {
		log.Printf("[NOTICE] Initializing kube config with non-default context: %v", config.Vars.ServicePacks.Kubernetes.KubeContext)
		modifyContext(rawConfig)
	}
	return k.ClientConfig()
}

func modifyContext(rawConfig clientcmdapi.Config) {
	ctx := config.Vars.ServicePacks.Kubernetes.KubeContext
	if rawConfig.Contexts[ctx] == nil {
		log.Fatalf("Required context does not exist in provided kubeconfig: %v", ctx)
	}
	rawConfig.CurrentContext = ctx
	err := clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(), rawConfig, true)
	if err != nil {
		log.Fatalf("[ERROR] Failed to modify context in kubeconfig: %v", ctx)
	}
}
