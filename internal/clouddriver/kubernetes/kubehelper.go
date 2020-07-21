package kubernetes

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"

	apiv1 "k8s.io/api/core/v1"
	// appsv1 "k8s.io/api/apps/v1"
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
	n := "kubeconfig"
	f := flag.Lookup(n)
	if f != nil {
		return &f.DefValue
	}

	var c *string
	//prefer kube config path if it's been supplied
	if kubeConfigFile != nil && *kubeConfigFile != "" {
		log.Printf("[NOTICE] Setting Kube Config to: %v", *kubeConfigFile)
		c = flag.String("kubeconfig", *kubeConfigFile, "fully qualified and supplied absolute path to the kubeconfig file")
	} else if e := getConfigPathFromEnv(); e != "" {
		log.Printf("[NOTICE] Setting Kube Config to: %v", e)
		c = flag.String("kubeconfig", e, "(optional) absolute path to the kubeconfig file")
	} else if home := homeDir(); home != "" {
		p := filepath.Join(home, ".kube", "config")
		log.Printf("[NOTICE] Setting Kube Config to: %v", p)
		c = flag.String("kubeconfig", p, "(optional) absolute path to the kubeconfig file")
	} else {
		c = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	return c
}

//GetPods ...
func GetPods() (*apiv1.PodList, error) {
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

func getPods(c *kubernetes.Clientset) (*apiv1.PodList, error) {
	pods, err := c.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	if pods == nil {
		return nil, fmt.Errorf("pod list returned nil")
	}

	log.Printf("[NOTICE] There are %d pods in the cluster\n", len(pods.Items))

	for i := 0; i < len(pods.Items); i++ {
		log.Printf("[INFO] Pod: %v %v\n", pods.Items[i].GetNamespace(), pods.Items[i].GetName())
	}

	return pods, nil
}

//CreatePod ...
func CreatePod(pname *string, ns *string, cname *string, image *string) (*apiv1.Pod, error) {
	c, err := GetClient()
	if err != nil {
		return nil, err
	}

	//create the namespace for the POD (noOp if already present)
	_, err = CreateNamespace(ns)
	if err != nil {
		return nil, err
	}

	//now do pod ...
	pc := c.CoreV1().Pods(*ns)
	p := getPodObject(*pname, *ns, *cname, *image)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := pc.Create(ctx, p, metav1.CreateOptions{})
	if err != nil {
		if isAlreadyExists(err) {
			log.Printf("[INFO] POD %v already exists. Returning existing.\n", *pname)
			//return it and nil out err
			return res, nil
		}
		return nil, err
	}

	log.Printf("[NOTICE] POD %q created.\n", res.GetObjectMeta().GetName())

	//wait:
	waitForRunning(c, ns, pname)

	return res, nil
}

func getPodObject(pname string, ns string, cname string, image string) *apiv1.Pod {
	return &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pname,
			Namespace: ns,
			Labels: map[string]string{
				"app": "demo",
			},
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:            cname,
					Image:           image,
					ImagePullPolicy: apiv1.PullIfNotPresent,
					Command: []string{
						"sleep",
						"3600",
						// "/bin/sh", "-c",
					},
				},
			},
		},
	}
}

//ExecCommand ...
func ExecCommand(cmd, ns, pn *string) (string, string, error) {
	if cmd == nil {
		return "", "", fmt.Errorf("command string is nil - nothing to execute")
	}
	logCmd(cmd, pn, ns)

	c, err := GetClient()
	if err != nil {
		return "", "", err
	}

	req := c.CoreV1().RESTClient().Post().Resource("pods").
		Name(*pn).Namespace(*ns).SubResource("exec")

	scheme := runtime.NewScheme()
	if err := apiv1.AddToScheme(scheme); err != nil {
		return "", "", fmt.Errorf("error adding to scheme: %v", err)
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	options := apiv1.PodExecOptions{
		Command: strings.Fields(*cmd),
		// Container: containerName, //specify if more than one container
		Stdout: true,
		Stderr: true,
		TTY:    false,
	}

	req.VersionedParams(&options, parameterCodec)

	log.Printf("[INFO] ExecCommand Request URL: %v", req.URL().String())

	config, err := clientcmd.BuildConfigFromFlags("", *setConfigPath())
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("error while creating Executor: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return stdout.String(), stderr.String(), fmt.Errorf("error in Stream: %v", err)
	}

	return stdout.String(), stderr.String(), nil
}

//DeletePod ...
func DeletePod(pname *string, ns *string) error {
	c, err := GetClient()
	if err != nil {
		return err
	}

	pc := c.CoreV1().Pods(*ns)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = pc.Delete(ctx, *pname, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	log.Printf("[NOTICE] POD %v deleted.", *pname)

	return nil
}

//CreateNamespace ...
func CreateNamespace(ns *string) (*apiv1.Namespace, error) {
	c, err := GetClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//try and create ...
	apiNS := apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *ns,		
		},
	}
	n, err := c.CoreV1().Namespaces().Create(ctx, &apiNS, metav1.CreateOptions{})
	if err != nil {
		if isAlreadyExists(err) {
			log.Printf("[INFO] Namespace %v already exists. Returning existing.", *ns)
			//return it and nil out the err
			return n, nil
		}
		return nil, err
	}

	log.Printf("[NOTICE] Namespace %q created.", n.GetObjectMeta().GetName())

	return n, nil
}

//DeleteNamespace ...
func DeleteNamespace(ns *string) error {
	c, err := GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = c.CoreV1().Namespaces().Delete(ctx, *ns, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	log.Printf("[NOTICE] Namespace %v deleted.", *ns)

	return nil
}

func isAlreadyExists(err error) bool {
	if se, ok := err.(*errors.StatusError); ok {
		//409 is "already exists"
		return se.ErrStatus.Code == 409
	}
	return false
}

func waitForRunning(c *kubernetes.Clientset, ns *string, n *string) (bool, error) {

	ps := c.CoreV1().Pods(*ns)

	w, err := ps.Watch(context.Background(), metav1.ListOptions{})

	if err != nil {
		return false, nil
	}

	go func() {
		for e := range w.ResultChan() {
			log.Printf("[INFO] Watch Event Type: %v", e.Type)
			p, ok := e.Object.(*apiv1.Pod)
			if !ok {
				log.Printf("[WARNING] Unexpected Watch Event Type - skipping")
				break
			}			
			log.Printf("[INFO] Watch Container phase: %v", p.Status.Phase)
			log.Printf("[DEBUG] Watch Container status: %+v", p.Status.ContainerStatuses)

		}
	}()
	time.Sleep(5 * time.Second)

	return false, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func getConfigPathFromEnv() string {	
	return os.Getenv("KUBE_CONFIG")	
}

func logCmd(c *string, p *string, n *string) {
	log.Printf("[NOTICE] Executing command: \"%v\" on POD '%v' in namespace '%v'", *c, *p, *n)	
}
