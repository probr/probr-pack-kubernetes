// Package kubernetes provides functions for interacting with Kubernetes and
// is built using the kubernetes client-go (https://github.com/kubernetes/client-go).
package kubernetes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/citihub/probr/audit"
	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/utils"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"

	executil "k8s.io/client-go/util/exec"

	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	//needed for authentication against the various GCPs
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// K8SJSON encapsulates the response from a raw/rest call to the Kubernetes API
type K8SJSON struct {
	APIVersion string
	Items      []struct {
		Kind     string
		Metadata map[string]string
	}
}

// CmdExecutionResult encapsulates the result from an exec call to the kubernetes cluster.  This includes 'stdout',
// 'stderr', 'exit code' and any error details in the case of a non-zero exit code.
type CmdExecutionResult struct {
	Stdout string
	Stderr string

	Err      error
	Code     int
	Internal bool
}

func (e *CmdExecutionResult) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("cmd execution result: code=%v stdout=%v stderr=%v", e.Code, e.Stdout, e.Stderr))

	if e.Err != nil {
		b.WriteString(fmt.Sprintf(" || error: internal=%t msg=%v", e.Internal, e.Err))
	}

	return b.String()
}

// Kubernetes interface defines the methods available to interact with the kubernetes cluster.
type Kubernetes interface {
	ClusterIsDeployed() *bool
	GetClient() (*kubernetes.Clientset, error)
	GetPods(ns string) (*apiv1.PodList, error)
	CreatePod(pname string, ns string, cname string, image string, w bool, sc *apiv1.SecurityContext, probe *audit.Probe) (*apiv1.Pod, *PodAudit, error)
	CreatePodFromObject(pod *apiv1.Pod, podName string, ns string, wait bool, probe *audit.Probe) (*apiv1.Pod, error)
	CreatePodFromYaml(y []byte, pname string, ns string, image string, aadpodidbinding string, w bool, probe *audit.Probe) (*apiv1.Pod, error)
	GetPodObject(pname string, ns string, cname string, image string, sc *apiv1.SecurityContext) *apiv1.Pod
	ExecCommand(cmd string, ns string, pn *string) *CmdExecutionResult
	DeletePod(pname string, ns string, e string) error
	DeleteNamespace(ns *string) error
	CreateConfigMap(n *string, ns string) (*apiv1.ConfigMap, error)
	DeleteConfigMap(name string) error
	GetConstraintTemplates(prefix string) (*map[string]interface{}, error)
	GetRawResourcesByGrp(g string) (*K8SJSON, error)
	GetClusterRolesByResource(r string) (*[]rbacv1.ClusterRole, error)
	GetClusterRoles() (*rbacv1.ClusterRoleList, error)
}

// Kube provides an implementation of Kubernetes.
type Kube struct {
	kubeClient                   *kubernetes.Clientset
	clientMutex                  sync.Mutex
	cspErrorToProbrCreationError map[string]PodCreationErrorReason

	k8statusToPodCreationError map[string]PodCreationErrorReason
}

var instance *Kube
var once sync.Once

// GetKubeInstance returns a singleton instance of Kube.
func GetKubeInstance() *Kube {
	//TODO: revise use of singleton here ...
	once.Do(func() {
		instance = &Kube{}

		//This is brittle!:
		//Map error message strings to common creation error types
		//unfortunately there is no alternative mechanism to interpret the reason for
		//pod creation failure.
		//('azurepolicy' messages are from AKS via Azure Policy constraints; 'securityContext' are from
		//EKS via underlying PSP)
		instance.cspErrorToProbrCreationError = make(map[string]PodCreationErrorReason, 7)
		instance.cspErrorToProbrCreationError["azurepolicy-container-no-privilege"] = PSPNoPrivilege
		instance.cspErrorToProbrCreationError["securityContext.privileged: Invalid value: true"] = PSPNoPrivilege

		instance.cspErrorToProbrCreationError["azurepolicy-psp-container-no-privilege-esc"] = PSPNoPrivilegeEscalation
		instance.cspErrorToProbrCreationError["securityContext.allowPrivilegeEscalation: Invalid value: true"] = PSPNoPrivilegeEscalation

		instance.cspErrorToProbrCreationError["azurepolicy-psp-allowed-users-groups"] = PSPAllowedUsersGroups
		instance.cspErrorToProbrCreationError["securityContext.runAsUser: Invalid value: 0"] = PSPAllowedUsersGroups

		instance.cspErrorToProbrCreationError["azurepolicy-container-allowed-images"] = PSPContainerAllowedImages

		instance.cspErrorToProbrCreationError["azurepolicy-psp-host-namespace"] = PSPHostNamespace
		instance.cspErrorToProbrCreationError["securityContext.hostPID: Invalid value: true"] = PSPHostNamespace
		instance.cspErrorToProbrCreationError["securityContext.hostIPC: Invalid value: true"] = PSPHostNamespace

		instance.cspErrorToProbrCreationError["azurepolicy-psp-host-network"] = PSPHostNetwork
		instance.cspErrorToProbrCreationError["securityContext.hostNetwork: Invalid value: true"] = PSPHostNetwork

		instance.cspErrorToProbrCreationError["azurepolicy-container-allowed-capabilities"] = PSPAllowedCapabilities
		instance.cspErrorToProbrCreationError["securityContext.capabilities.add: Invalid value: \"NET_RAW\""] = PSPAllowedCapabilities
		instance.cspErrorToProbrCreationError["securityContext.capabilities.add: Invalid value: \"NET_ADMIN\""] = PSPAllowedCapabilities

		instance.cspErrorToProbrCreationError["azurepolicy-psp-host-network-ports"] = PSPAllowedPortRange
		instance.cspErrorToProbrCreationError["hostPort: Invalid value"] = PSPAllowedPortRange

		instance.cspErrorToProbrCreationError["azurepolicy-psp-volume-types"] = PSPAllowedVolumeTypes

		instance.cspErrorToProbrCreationError["azurepolicy-psp-seccomp"] = PSPSeccompProfile
		instance.cspErrorToProbrCreationError["not an allowed seccomp profile"] = PSPSeccompProfile

		instance.k8statusToPodCreationError = make(map[string]PodCreationErrorReason, 2)
		instance.k8statusToPodCreationError["ErrImagePull"] = ImagePullError
		instance.k8statusToPodCreationError["Blocked"] = Blocked
	})

	return instance
}

// ClusterIsDeployed verifies if a cluster is deployed that can be contacted based on the current
// kubernetes config and context.
func (k *Kube) ClusterIsDeployed() *bool {
	kc, err := k.GetClient()
	if err != nil {
		utils.ReformatError("Error raised when getting Kubernetes client: %v", err)
		return nil
	}

	t, f := true, false
	if kc == nil {
		return &f
	}

	return &t
}

//GetPods returns a collection of pods on the target kubernetes cluster.
func (k *Kube) GetPods(ns string) (*apiv1.PodList, error) {
	c, err := k.GetClient()
	if err != nil {
		return nil, err
	}

	pl, err := getPods(c, ns)
	if err != nil {
		return nil, err
	}

	return pl, nil
}

// CreatePod creates a pod with the supplied parameters.  A true value for 'wait' indicates that
// the function should wait (block) until the pod is in a running state.
func (k *Kube) CreatePod(podName string, ns string, containerName string, image string, wait bool, sc *apiv1.SecurityContext, probe *audit.Probe) (*apiv1.Pod, *PodAudit, error) {
	//create Pod Object ...
	p := k.GetPodObject(podName, ns, containerName, image, sc)
	audit := &PodAudit{podName, "probr-general-test-ns", containerName, image, sc}

	pod, err := k.CreatePodFromObject(p, podName, ns, wait, probe)
	return pod, audit, err
}

// CreatePodFromYaml creates a pod for the supplied yaml.  A true value for 'w' indicates that the function
// should wait (block) until the pod is in a running state.
func (k *Kube) CreatePodFromYaml(y []byte, pname string, ns string, image string, aadpodidbinding string, w bool, probe *audit.Probe) (*apiv1.Pod, error) {
	vars := config.Vars.ServicePacks.Kubernetes
	approvedImage := vars.AuthorisedContainerRegistry + "/" + vars.ProbeImage
	podSpec := utils.ReplaceBytesValue(y, "{{ probr-compatible-image }}", approvedImage)
	o, _, err := scheme.Codecs.UniversalDeserializer().Decode(podSpec, nil, nil)
	if err != nil {
		log.Printf("[ERROR] %s: could not create pod from yaml asset, %v", utils.CallerName(2), err)
	}

	p := o.(*apiv1.Pod)
	//update the name to the one that's supplied
	p.SetName(pname)
	//also update the image (which could have been supplied via the env)
	//(only expecting one container, but loop in case of many)
	if image != "" {
		for _, c := range p.Spec.Containers {
			c.Image = image
		}
	}

	if aadpodidbinding != "" {
		if p.Labels == nil {
			p.Labels = make(map[string]string)
		}
		p.Labels["aadpodidbinding"] = aadpodidbinding
	}
	return k.CreatePodFromObject(p, pname, ns, w, probe)
}

// CreatePodFromObject creates a pod from the supplied pod object with the given pod name and namespace.  A true value for 'w' indicates that the function
// should wait (block) until the pod is in a running state.
func (k *Kube) CreatePodFromObject(pod *apiv1.Pod, podName string, ns string, wait bool, probe *audit.Probe) (*apiv1.Pod, error) {
	if pod == nil || podName == "" || ns == "" {
		return nil, fmt.Errorf("one or more of pod (%v), podName (%v) or namespace (%v) is nil - cannot create POD", pod, podName, ns)
	}

	log.Printf("[INFO] Creating pod %v in namespace %v", podName, ns)
	log.Printf("[DEBUG] Pod details: %+v", *pod)

	c, err := k.GetClient()
	if err != nil {
		return nil, err
	}

	//create the namespace for the POD (noOp if already present)
	_, err = k.createNamespace(&ns)
	if err != nil {
		return nil, err
	}

	//now do pod ...
	pc := c.CoreV1().Pods(ns)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := pc.Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		if isAlreadyExists(err) {
			log.Printf("[INFO] POD %v already exists. Returning existing.", podName)
			res, _ := pc.Get(ctx, podName, metav1.GetOptions{})

			//return it and nil out err
			return res, nil
		} else if isForbidden(err) {
			log.Printf("[NOTICE] Creation of POD %v is forbidden: %v", podName, err)
			//return a specific error:
			return nil, &PodCreationError{err, *k.toPodCreationErrorCode(err)}
		}
		return nil, err
	}

	log.Printf("[INFO] POD %q creation started.", res.GetObjectMeta().GetName())

	if wait {
		err = k.waitForPhase(apiv1.PodRunning, c, &ns, &podName)
		if err != nil {
			return res, err
		}
	}
	probe.CountPodCreated(podName)
	log.Printf("[INFO] POD %q creation completed. Pod is up and running.", res.GetObjectMeta().GetName())
	return res, nil
}

// GetPodObject constructs a simple pod object using kubernetes API types.
func (k *Kube) GetPodObject(pname string, ns string, cname string, image string, sc *apiv1.SecurityContext) *apiv1.Pod {

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
					},
					SecurityContext: sc,
				},
			},
		},
	}
}

// CreateConfigMap creates a config map with the supplied name in the given namespace.
func (k *Kube) CreateConfigMap(n *string, ns string) (*apiv1.ConfigMap, error) {
	c, err := k.GetClient()
	if err != nil {
		return nil, err
	}

	//create the namespace for the POD (noOp if already present)
	_, err = k.createNamespace(&ns)
	if err != nil {
		return nil, err
	}

	//now do config map ...
	cms := c.CoreV1().ConfigMaps(ns)

	cm := apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      *n,
			Namespace: ns,
			Labels: map[string]string{
				"app": "demo",
			},
		},
		Data: map[string]string{
			"key": "value",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := cms.Create(ctx, &cm, metav1.CreateOptions{})

	if err != nil {
		log.Printf("[ERROR] Error creating ConfigMap %q: %v", res.GetObjectMeta().GetName(), err)
		return nil, err
	}

	log.Printf("[INFO] ConfigMap %q created.", res.GetObjectMeta().GetName())

	return res, nil
}

// DeleteConfigMap deletes the named config map in the given namespace.
func (k *Kube) DeleteConfigMap(name string) error {
	c, err := k.GetClient()
	if err != nil {
		return err
	}

	cms := c.CoreV1().ConfigMaps(Namespace)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = cms.Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	log.Printf("[INFO] ConfigMap %v deleted.", name)

	return nil
}

// ExecCommand executes the supplied command on the given pod name in the specified namespace.
func (k *Kube) ExecCommand(cmd string, ns string, pn *string) (s *CmdExecutionResult) {
	if cmd == "" {
		return &CmdExecutionResult{Err: fmt.Errorf("command string is nil - nothing to execute"), Internal: true}
	}
	log.Printf("[DEBUG] Executing command: \"%s\" on POD '%s' in namespace '%s'", cmd, *pn, ns)
	c, err := k.GetClient()
	if err != nil {
		return &CmdExecutionResult{Err: err, Internal: true}
	}

	req := c.CoreV1().RESTClient().Post().Resource("pods").
		Name(*pn).Namespace(ns).SubResource("exec")

	scheme := runtime.NewScheme()
	if err := apiv1.AddToScheme(scheme); err != nil {
		return &CmdExecutionResult{Err: fmt.Errorf("error adding to scheme: %v", err), Internal: true}
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	options := apiv1.PodExecOptions{
		Command: strings.Fields(cmd),
		// Container: containerName, //specify if more than one container
		Stdout: true,
		Stderr: true,
		TTY:    false,
	}

	req.VersionedParams(&options, parameterCodec)

	log.Printf("[DEBUG] %s.%s: ExecCommand Request URL: %v", utils.CallerName(2), utils.CallerName(1), req.URL().String())

	config, err := clientcmd.BuildConfigFromFlags("", config.Vars.ServicePacks.Kubernetes.KubeConfigPath)
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return &CmdExecutionResult{Err: fmt.Errorf("error while creating Executor: %v", err), Internal: true}
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		if ce, ok := err.(executil.CodeExitError); ok {
			//the command has been executed on the container, but the underlying command raised an error
			//this is an 'external' error and represents a successful communication with the cluster
			return &CmdExecutionResult{Stdout: stdout.String(), Stderr: stderr.String(), Code: ce.Code, Err: fmt.Errorf("error raised on cmd execution: %v", err)}
		}
		return &CmdExecutionResult{Stdout: stdout.String(), Stderr: stderr.String(), Err: fmt.Errorf("error in Stream: %v", err), Internal: true}
	}

	//all good:
	return &CmdExecutionResult{Stdout: stdout.String(), Stderr: stderr.String()}
}

// DeletePod deletes the given pod in the specified namespace.
func (k *Kube) DeletePod(podName string, ns string, probeName string) error {
	_, err := k.PodStatus(podName, ns)
	if err != nil {
		return err // If pod does not exist, it cannot be deleted
	}

	c, err := k.GetClient()
	if err != nil {
		return err
	}

	pc := c.CoreV1().Pods(ns)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = pc.Delete(ctx, podName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	audit.State.GetProbeLog(probeName).CountPodDestroyed()
	log.Printf("[INFO] POD %v deleted.", podName)

	return nil
}

func (k *Kube) PodStatus(name, ns string) (apiv1.PodStatus, error) {
	client, err := k.GetClient()
	if err != nil {
		return apiv1.PodStatus{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pod, err := client.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
	return pod.Status, err
}

func (k *Kube) createNamespace(ns *string) (*apiv1.Namespace, error) {
	c, err := k.GetClient()
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

	log.Printf("[INFO] Namespace %q created.", n.GetObjectMeta().GetName())

	return n, nil
}

// DeleteNamespace deletes the supplied namespace.
func (k *Kube) DeleteNamespace(ns *string) error {
	c, err := k.GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = c.CoreV1().Namespaces().Delete(ctx, *ns, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	log.Printf("[INFO] Namespace %v deleted.", *ns)

	return nil
}

// GetConstraintTemplates returns the constraint templates associated with the active cluster.
func (k *Kube) GetConstraintTemplates(prefix string) (*map[string]interface{}, error) {
	return k.getAPIResourcesByGrp("constraints", prefix)
}

// GetIdentityBindings returns the identity bindings associated with the active cluster.
func (k *Kube) GetIdentityBindings(prefix string) (*map[string]interface{}, error) {
	return k.getAPIResourcesByGrp("aadpodidentity", prefix)
}

func (k *Kube) getAPIResourcesByGrp(grp string, nPrefix string) (*map[string]interface{}, error) {
	c, err := k.GetClient()
	if err != nil {
		return nil, err
	}

	_, arl, err := c.ServerGroupsAndResources()

	if err != nil {
		return nil, err
	}

	var con = make(map[string]interface{})
	for _, ar := range arl {
		if ar == nil {
			continue
		}
		g := ar.GroupVersion
		log.Printf("[DEBUG] API Resource Group %v", g)
		if len(grp) > 0 && !strings.HasPrefix(g, grp) {
			continue
		}

		for _, a := range ar.APIResources {
			log.Printf("[DEBUG] API Resource %+v", a)
			log.Printf("[DEBUG] API Resource - Group: %v Name: %v Kind: %v", g, a.Name, a.Kind)

			//skip if it doesn't pass the prefix filter (if one has been supplied):
			if len(nPrefix) > 0 && !strings.HasPrefix(a.Name, nPrefix) {
				continue
			}
			//treat it like a set ...
			_, exists := con[a.Name]
			if !exists {
				con[a.Name] = a
			}
		}
	}

	return &con, nil
}

// GetRawResourcesByGrp makes a 'raw' REST call to k8s to get the resources specified by the
// supplied group string, e.g. "apis/aadpodidentity.k8s.io/v1/azureidentitybindings".  This
// is required to support resources that are not supported by typed API calls (e.g. "pods").
func (k *Kube) GetRawResourcesByGrp(g string) (*K8SJSON, error) {
	c, err := k.GetClient()
	if err != nil {
		return nil, err
	}

	r := c.RESTClient().Get().AbsPath(g)
	log.Printf("[DEBUG] REST request: %+v", r)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res := r.Do(ctx)
	if res.Error() != nil {
		return nil, res.Error()
	}

	b, _ := res.Raw()
	bs := string(b)
	log.Printf("[DEBUG] STRING result: %v", bs)

	j := K8SJSON{}
	json.Unmarshal(b, &j)

	log.Printf("[DEBUG] JSON result: %+v", j)

	return &j, nil
}

// GetClusterRolesByResource returns a collection of cluster roles filtered by
// the supplied resource type.
func (k *Kube) GetClusterRolesByResource(r string) (*[]rbacv1.ClusterRole, error) {
	var crs []rbacv1.ClusterRole

	crl, err := k.GetClusterRoles()
	if err != nil {
		return &crs, err
	}

	for _, cr := range crl.Items {
		log.Printf("[DEBUG] ClusterRole: %+v", cr)
		if k.meetsResourceFilter(r, &cr.ObjectMeta, &cr.Rules) {
			//add to results
			log.Printf("[DEBUG] ClusterRole meets resource filter (%v): %+v", r, cr)
			crs = append(crs, cr)
		}
	}

	return &crs, nil
}

// GetRolesByResource returns a collection of roles filtered by
// the supplied resource type.
func (k *Kube) GetRolesByResource(r string) (*[]rbacv1.Role, error) {
	var ros []rbacv1.Role

	rl, err := k.GetRoles()
	if err != nil {
		return &ros, err
	}

	for _, ro := range rl.Items {
		log.Printf("[DEBUG] Role: %+v", ro)
		if k.meetsResourceFilter(r, &ro.ObjectMeta, &ro.Rules) {
			//add to results
			log.Printf("[DEBUG] Role meets resource filter (%v): %+v", r, ro)
			ros = append(ros, ro)
		}
	}

	return &ros, nil
}

func (k *Kube) meetsResourceFilter(f string, m *v1.ObjectMeta, p *[]rbacv1.PolicyRule) bool {

	//skip system/known roles
	if k.skipSystemRole(m) {
		return false
	}

	for _, ru := range *p {
		log.Printf("[DEBUG] PolicyRule: %+v", ru)
		var b bool

		for _, res := range ru.Resources {
			if strings.HasPrefix(res, f) {
				log.Printf("[DEBUG] PolicyRule meets filter %v", f)
				//meets filter
				//can also break out of the rules loop as
				//we want to add full role to results if one rule
				//passes filter
				b = true
				break
			}
		}
		if b {
			return true
		}
	}
	return false
}

func (k *Kube) skipSystemRole(m *v1.ObjectMeta) bool {
	//first check for known system namespaces:
	if strings.HasPrefix(m.Namespace, "kube") || strings.HasPrefix(m.Namespace, "gatekeeper") {
		return true
	}

	//next, check to see if the role name is on the list of system roles
	for _, r := range config.Vars.ServicePacks.Kubernetes.SystemClusterRoles {
		//use a prefix check:
		if strings.HasPrefix(m.Name, r) {
			return true
		}
	}

	return false
}

// GetClusterRoles retrieves all cluster roles associated with the active cluster.
func (k *Kube) GetClusterRoles() (*rbacv1.ClusterRoleList, error) {
	c, err := k.GetClient()
	if err != nil {
		// return nil, err
	}

	cr := c.RbacV1().ClusterRoles()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return cr.List(ctx, metav1.ListOptions{LabelSelector: "gatekeeper.sh/system!=yes"})
}

//GetRoles retrieves all roles associated with the active cluster.
func (k *Kube) GetRoles() (*rbacv1.RoleList, error) {
	c, err := k.GetClient()
	if err != nil {
		// return nil, err
	}

	r := c.RbacV1().Roles("")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return r.List(ctx, metav1.ListOptions{LabelSelector: "gatekeeper.sh/system!=yes"})
}

func (k *Kube) toPodCreationErrorCode(err error) *map[PodCreationErrorReason]*PodCreationErrorReason {
	//try and map the error codes within the error message issued by the service provider
	//to known error codes (return a map so they can be easily accessed)

	var pcErr = make(map[PodCreationErrorReason]*PodCreationErrorReason)
	if se, ok := err.(*errors.StatusError); ok {
		//get the reason
		r := se.ErrStatus.Reason
		m := se.ErrStatus.Message

		log.Printf("[INFO] *** reason: %v", r)
		log.Printf("[INFO] *** message: %v", m)
		//map this to the pod creation code

		for key, e := range k.cspErrorToProbrCreationError {
			if strings.Contains(m, key) {
				//take the element
				pcErr[e] = &e
			}
		}
	}

	return &pcErr
}

func (k *Kube) waitForPhase(ph apiv1.PodPhase, c *kubernetes.Clientset, ns *string, n *string) error {

	ps := c.CoreV1().Pods(*ns)

	//don't wait for more than 1 min ...
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	w, err := ps.Watch(ctx, metav1.ListOptions{})

	if err != nil {
		return err
	}

	log.Printf("[INFO] *** Waiting for phase %v on pod %v ...", ph, *n)

	for e := range w.ResultChan() {
		log.Printf("[DEBUG] Watch Probe Type: %v", e.Type)
		p, ok := e.Object.(*apiv1.Pod)
		if !ok {
			log.Printf("[WARN] Unexpected Watch Probe Type - skipping")
			//check for timeout
			if ctx.Err() != nil {
				log.Printf("[WARN] Context error received while waiting on pod %v. Error: %v", *n, ctx.Err())
				return ctx.Err()
			}
			// break
			continue
		}
		if p.GetName() != *n {
			log.Printf("[NOTICE] Probe received for pod %v which we're not waiting on. Skipping.", p.GetName())
			continue
		}

		log.Printf("[INFO] Pod %v Phase: %v", p.GetName(), p.Status.Phase)
		for _, con := range p.Status.ContainerStatuses {
			log.Printf("[DEBUG] Container Status: %+v", con)
		}

		// don't wait if we're getting errors:
		b, err := k.podInErrorState(p)
		if b {
			log.Printf("[ERROR] Giving up waiting on pod creation. Error: %v", err)
			return err
		}

		if p.Status.Phase == ph {
			break
		}

	}

	log.Printf("[INFO] *** Completed waiting for phase %v on pod %v", ph, *n)

	return nil
}

func (k *Kube) podInErrorState(p *apiv1.Pod) (bool, *PodCreationError) {

	// check the container statuses for error conditions:
	if len(p.Status.ContainerStatuses) > 0 {
		if p.Status.ContainerStatuses[0].State.Waiting != nil {
			n := p.GetObjectMeta().GetName()
			r := p.Status.ContainerStatuses[0].State.Waiting.Reason
			log.Printf("[DEBUG] Pod: %v Waiting reason: %v", n, r)

			//TODO: other error states? Also need to tidy up the error creation
			pe, exists := k.k8statusToPodCreationError[r]

			if exists {
				log.Printf("[ERROR] Giving up waiting on pod %v . Error reason: %v", n, r)
				pcErr := make(map[PodCreationErrorReason]*PodCreationErrorReason, 1)

				pcErr[pe] = &pe
				return true, &PodCreationError{fmt.Errorf("giving up waiting on pod %v . Error reason: %v", n, r), pcErr}
			}
		}
	}

	return false, nil
}
