package kubernetes

import (
	"os"
	"flag"
	"log"
	"context"
	"sync"
	"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/api/core/v1"
	
)

var (
	kubeConfigFile *string
	kubeClient	*kubernetes.Clientset
	clientMutex		sync.Mutex
)


//SetKubeConfigFile - explict/full path to kube config .. TODO: hmm not sure I like this
func SetKubeConfigFile(fullyQualifiedKubeConfig *string) {
	kubeConfigFile = fullyQualifiedKubeConfig
}


//GetClient ... just for initial dev - will change TODO ... will probably remove this from export
func GetClient() (*kubernetes.Clientset, error) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	if kubeClient != nil {
		return kubeClient, nil		
	}

	var kubeconfig *string	

	//prefer kube config path if it's been supplied
	if kubeConfigFile != nil && *kubeConfigFile != "" {
		kubeconfig = flag.String("kubeconfig", *kubeConfigFile, "fully qualified and supplied absolute path to the kubeconfig file")		
	} else if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		//TODO: DA, not sure we really need to panic here :-)
		panic(err.Error())
	}

	// create the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	kubeClient = clientSet

	return kubeClient, nil		
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
		log.Printf("P: %v %v\n", pods.Items[i].GetNamespace() , pods.Items[i].GetName() )			
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