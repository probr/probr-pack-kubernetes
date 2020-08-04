package kubernetes

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
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

	executil "k8s.io/client-go/util/exec"

	apiv1 "k8s.io/api/core/v1"
)

//PodCreationErrorReason ... TODO: not sure if this is the correct name for this?
type PodCreationErrorReason int

//PodCreationErrorReason enum
const (
	UndefinedPodCreationErrorReason PodCreationErrorReason = iota
	PSPNoPrivilege
	PSPNoPrivilegeEscalation
	PSPAllowedUsersGroups
	PSPContainerAllowedImages
)

func (r PodCreationErrorReason) String() string {
	return [...]string{"podcreation-error: undefined",
		"podcreation-error: psp-container-no-privilege",
		"podcreation-error: psp-container-no-privilege-escalation",
		"podcreation-error: psp-allowed-users-groups",
		"podcreation-error: psp-container-allowed-images"}[r]
}

//PodCreationError ...
type PodCreationError struct {
	err        error
	ReasonCode PodCreationErrorReason
}

func (p *PodCreationError) Error() string {
	return fmt.Sprintf("pod creation error: %v %v", p.ReasonCode, p.err)
}

var (
	kubeConfigFile *string
	kubeClient     *kubernetes.Clientset
	clientMutex    sync.Mutex

	azErrorToPodCreationError = make(map[string]PodCreationErrorReason)
)

func init() {
	azErrorToPodCreationError["azurepolicy-container-no-privilege"] = PSPNoPrivilege
	azErrorToPodCreationError["azurepolicy-psp-container-no-privilege-escalation"] = PSPNoPrivilegeEscalation
	azErrorToPodCreationError["azurepolicy-psp-allowed-users-groups"] = PSPAllowedUsersGroups
	azErrorToPodCreationError["azurepolicy-container-allowed-images"] = PSPContainerAllowedImages
}

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

// CreatePod creates a pod with the following parameters:
// pname - pod name
// ns - namespace
// cname - container name
// image - image
// w - indicates whether or not to wait for the pod to be running
// sc - security context
func CreatePod(pname *string, ns *string, cname *string, image *string, w bool, sc *apiv1.SecurityContext) (*apiv1.Pod, error) {
	//create Pod Objet ...
	p := GetPodObject(*pname, *ns, *cname, *image, sc)

	return CreatePodFromObject(p, pname, ns, w)
}

// CreatePodFromObject creates a pod from the supplied pod object in the gievn namespace
func CreatePodFromObject(p *apiv1.Pod, pname *string, ns *string, w bool) (*apiv1.Pod, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := pc.Create(ctx, p, metav1.CreateOptions{})
	if err != nil {
		if isAlreadyExists(err) {
			log.Printf("[NOTICE] POD %v already exists. Returning existing.", *pname)
			res, _ := pc.Get(ctx, *pname, metav1.GetOptions{})

			//return it and nil out err
			return res, nil
		} else if isForbidden(err) {
			log.Printf("[NOTICE] Creation of POD %v is forbidden: %v", *pname, err)
			//return a specific error:
			return nil, &PodCreationError{err, toPodCreationErrorCode(err)}
		}
		return nil, err
	}

	log.Printf("[NOTICE] POD %q created.\n", res.GetObjectMeta().GetName())

	if w {
		//wait:
		waitForPhase(apiv1.PodRunning, c, ns, pname)
	}

	return res, nil
}

// GetPodObject ...
func GetPodObject(pname string, ns string, cname string, image string, sc *apiv1.SecurityContext) *apiv1.Pod {

	a := make(map[string]string)
	a["seccomp.security.alpha.kubernetes.io/pod"] = "runtime/default"

	if sc == nil {
		sc = defaultContainerSecurityContext()
	}

	return &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pname,
			Namespace: ns,
			Labels: map[string]string{
				"app": "demo",
			},
			Annotations: a,
		},
		Spec: apiv1.PodSpec{
			SecurityContext: defaultPodSecurityContext(),
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
					SecurityContext: sc,
				},
			},
		},
	}
}

//GenerateUniquePodName ...
func GenerateUniquePodName(baseName string) string {
	//take base and add some uniqueness
	t := time.Now()
	rand.Seed(t.UnixNano())
	uniq := fmt.Sprintf("%v-%v", t.Format("020106-150405"), rand.Intn(100))

	return fmt.Sprintf("%v-%v", baseName, uniq)
}

func defaultPodSecurityContext() *apiv1.PodSecurityContext {
	var user, grp, fsgrp int64
	user, grp, fsgrp = 1000, 3000, 2000

	return &apiv1.PodSecurityContext{
		RunAsUser:          &user,
		RunAsGroup:         &grp,
		FSGroup:            &fsgrp,
		SupplementalGroups: []int64{1},
	}
}

func defaultContainerSecurityContext() *apiv1.SecurityContext {
	b := false

	return &apiv1.SecurityContext{
		Privileged:               &b,
		AllowPrivilegeEscalation: &b,
	}
}

//ExecCommand TODO: fix error codes
func ExecCommand(cmd, ns, pn *string) (string, string, int, error) {
	if cmd == nil {
		return "", "", -1, fmt.Errorf("command string is nil - nothing to execute")
	}
	logCmd(cmd, pn, ns)

	c, err := GetClient()
	if err != nil {
		return "", "", -2, err
	}

	req := c.CoreV1().RESTClient().Post().Resource("pods").
		Name(*pn).Namespace(*ns).SubResource("exec")

	scheme := runtime.NewScheme()
	if err := apiv1.AddToScheme(scheme); err != nil {
		return "", "", -3, fmt.Errorf("error adding to scheme: %v", err)
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
		return "", "", -4, fmt.Errorf("error while creating Executor: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		if ce, ok := err.(executil.CodeExitError); ok {
			return stdout.String(), stderr.String(), ce.Code, fmt.Errorf("error in Stream: %v", err)
		}
		return stdout.String(), stderr.String(), -5, fmt.Errorf("error in Stream: %v", err)
	}

	return stdout.String(), stderr.String(), 0, nil
}

// DeletePod deletes the pod with the following parameters:
// pname - pod name
// ns - namespace
// w - indicates whether or not to wait on the deletion
func DeletePod(pname *string, ns *string, w bool) error {
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

	if w {
		//wait:
		waitForDelete(c, ns, pname)
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
			Name: *ns,
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

func isForbidden(err error) bool {
	if se, ok := err.(*errors.StatusError); ok {
		//40 is "already exists"
		return se.ErrStatus.Code == 403
	}
	return false
}

func toPodCreationErrorCode(err error) PodCreationErrorReason {
	if se, ok := err.(*errors.StatusError); ok {
		//get the reason
		r := se.ErrStatus.Reason
		m := se.ErrStatus.Message

		log.Printf("[INFO] *** reason: %v", r)
		log.Printf("[INFO] *** message: %v", m)
		//map this to the pod creation code
		for k, e := range azErrorToPodCreationError {
			if strings.Contains(m, k) {
				//take the element
				return e
			}
		}
	}

	return UndefinedPodCreationErrorReason
}

func waitForPhase(ph apiv1.PodPhase, c *kubernetes.Clientset, ns *string, n *string) error {

	ps := c.CoreV1().Pods(*ns)

	w, err := ps.Watch(context.Background(), metav1.ListOptions{})

	if err != nil {
		return err
	}

	log.Printf("[NOTICE] *** Waiting for phase %v on pod %v ...", ph, *n)

	for e := range w.ResultChan() {
		log.Printf("[INFO] Watch Event Type: %v", e.Type)
		p, ok := e.Object.(*apiv1.Pod)
		if !ok {
			log.Printf("[WARNING] Unexpected Watch Event Type - skipping")
			break
		}
		log.Printf("[INFO] Watch Container phase: %v", p.Status.Phase)
		log.Printf("[DEBUG] Watch Container status: %+v", p.Status.ContainerStatuses)

		// don't wait if we're getting errors:
		if podInErrorState(p) {
			break
		}

		if p.Status.Phase == ph {
			break
		}

	}

	log.Printf("[NOTICE] *** Completed waiting for phase %v on pod %v", ph, *n)

	return nil
}

func podInErrorState(p *apiv1.Pod) bool {

	// check the container statuses for error conditions:
	if len(p.Status.ContainerStatuses) > 0 {
		if p.Status.ContainerStatuses[0].State.Waiting != nil {
			n := p.GetObjectMeta().GetName()
			r := p.Status.ContainerStatuses[0].State.Waiting.Reason
			log.Printf("[INFO] Pod: %v Waiting reason: %v", n, r)

			//TODO: other error states?
			if r == "ErrImagePull" {
				log.Printf("[WARN] Giving up waiting on pod %v . Error reason: %v", n, r)
				return true
			}
		}
	}

	return false
}

func waitForDelete(c *kubernetes.Clientset, ns *string, n *string) error {

	ps := c.CoreV1().Pods(*ns)

	w, err := ps.Watch(context.Background(), metav1.ListOptions{})

	if err != nil {
		return err
	}

	log.Printf("[NOTICE] *** Waiting for DELETE on pod %v ...", *n)

	for e := range w.ResultChan() {
		log.Printf("[INFO] Watch Event Type: %v", e.Type)
		p, ok := e.Object.(*apiv1.Pod)
		if !ok {
			log.Printf("[WARNING] Unexpected Watch Event Type received for pod %v - skipping", p.GetObjectMeta().GetName())
			break
		}
		log.Printf("[INFO] Watch Container phase: %v", p.Status.Phase)
		log.Printf("[DEBUG] Watch Container status: %+v", p.Status.ContainerStatuses)

		if e.Type == "DELETED" {
			log.Printf("[NOTICE] DELETED event received for pod %v", p.GetObjectMeta().GetName())
			break
		}

	}

	log.Printf("[NOTICE] *** Completeed waiting for DELETE on pod %v", *n)

	return nil
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
