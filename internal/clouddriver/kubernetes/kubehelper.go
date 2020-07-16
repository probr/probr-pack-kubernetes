package kubernetes

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	v1 "k8s.io/api/core/v1"
)

var (
	kubeConfigFile *string
	kubeClient     *kubernetes.Clientset
	clientMutex    sync.Mutex
)

//SetKubeConfigFile sets the fully qualified path to the Kubernetes config file.
func SetKubeConfigFile(f *string) {
	kubeConfigFile = f
}

//GetClient gets a client connection to the Kubernetes cluster specifed via @SetKubeConfigFile or from home directory.
func GetClient() (*kubernetes.Clientset, error) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	if kubeClient != nil {
		return kubeClient, nil
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *setConfigPath())
	if err != nil {
		return nil, err
	}

	// create the clientset (note: assigned to global "kubeClient")
	kubeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return kubeClient, nil
}

func setConfigPath() *string {
	var c *string

	//prefer kube config path if it's been supplied
	if kubeConfigFile != nil && *kubeConfigFile != "" {
		c = flag.String("kubeconfig", *kubeConfigFile, "fully qualified and supplied absolute path to the kubeconfig file")
	} else if home := homeDir(); home != "" {
		c = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		c = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	return c
}

//GetPods ...
func GetPods() (*v1.PodList, error) {
	c, err := GetClient()
	if err != nil {
		return nil, err
	}

	pl, err := getPods(c)
	if err != nil {
		return nil, err
	}

	return pl, nil
}

func getPods(c *kubernetes.Clientset) (*v1.PodList, error) {
	pods, err := c.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	if pods == nil {
		return nil, fmt.Errorf("pod list returned nil")
	}

	log.Printf("There are %d pods in the cluster\n", len(pods.Items))

	for i := 0; i < len(pods.Items); i++ {
		log.Printf("P: %v %v\n", pods.Items[i].GetNamespace(), pods.Items[i].GetName())
	}

	return pods, nil
}

//CreatePod ...
func CreatePod() {

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
