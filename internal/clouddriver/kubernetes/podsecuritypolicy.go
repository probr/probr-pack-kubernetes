package kubernetes

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"citihub.com/probr/internal/clouddriver/azure"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
)

//PrivilegedAccess ...
type PrivilegedAccess int

//PrivilegedAccess enum
const (
	WithPrivilegedAccess PrivilegedAccess = iota
	WithoutPrivilegedAccess
)

const (
	//TODO: default to these values for MVP - need to expose in future
	//TODO: also, using network-access-test-ns here as there's an exclusion on the
	//container registry - needs to be tidied up ...
	pspTestNamespace = "probr-pod-security-test-ns"
	//NOTE: either the above namespace needs to be added to the exclusion list on the
	//container registry rule or busybox need to be available in the allowed (probably internal) registry
	pspTestImage     = "docker.io/busybox"
	pspTestContainer = "psp-test"
	pspTestPodName   = "psp-test-pod"
)

var securityPolicyProviders = []SecurityPolicyProvider{
	&KubeSecurityPolicyProvider{},
	azure.NewAzPolicyProvider()}

//ClusterHasPSP determines if the cluster has any Pod Security Policies set.
func ClusterHasPSP() (*bool, error) {
	// iterate over providers ...
	for _, p := range securityPolicyProviders {
		b, err := p.HasSecurityPolicies()
		if err != nil {
			return nil, err
		}
		if *b {
			//return on first 'true' - only trying to establish if we have any ...
			return b, nil
		}
	}

	//if we get to here, we haven't got any ...
	b := false
	return &b, nil
}

//PrivilegedAccessIsRestricted looks for a PodSecurityPolicy with 'Privileged' set to false (ie. NOT privileged).
func PrivilegedAccessIsRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range securityPolicyProviders {
		b, err := p.HasPrivilegedAccessRestriction()
		if err != nil {
			return nil, err
		}
		if b != nil && *b {
			//return on first 'true' - only trying to establish if we have any ...
			return b, nil
		}
	}

	//if we get to here, we haven't got any ...
	b := false
	return &b, nil
}

//HostPIDIsRestricted looks for a PodSecurityPolicy with 'HostPID' set to false (i.e. NO Access to HostPID ).
func HostPIDIsRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range securityPolicyProviders {
		b, err := p.HasHostPIDRestriction()
		if err != nil {
			return nil, err
		}
		if *b {
			//return on first 'true' - only trying to establish if we have any ...
			return b, nil
		}
	}

	//if we get to here, we haven't got any ...
	b := false
	return &b, nil

}

//CreatePODSettingPrivilegedAccess ...
func CreatePODSettingPrivilegedAccess(pa PrivilegedAccess) (*apiv1.Pod, error) {

	b := pa == WithPrivilegedAccess
	sc := apiv1.SecurityContext{
		Privileged:               &b,
		AllowPrivilegeEscalation: &b,
	}

	pname, ns, cname, image := pspTestPodName, pspTestNamespace, pspTestContainer, pspTestImage

	p, err := CreatePod(&pname, &ns, &cname, &image, true, &sc)

	if err != nil {
		return nil, err
	}

	return p, nil
}

//TeardownPodSecurityTestPod ...
func TeardownPodSecurityTestPod(p *string) error {
	ns := pspTestNamespace
	err := DeletePod(p, &ns, true)
	return err
}

func getPodSecurityPolicies() (*v1beta1.PodSecurityPolicyList, error) {
	c, err := GetClient()
	if err != nil {
		return nil, err
	}

	psp := c.PolicyV1beta1().PodSecurityPolicies()
	if psp == nil {
		return nil, fmt.Errorf("Pod Security Polices could not be obtained (nil returned)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pspList, err := psp.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	if pspList == nil {
		return nil, fmt.Errorf("Pod Security Polices list returned a nil list")
	}

	log.Printf("[NOTICE] There are %d psp policies in the cluster\n", len(pspList.Items))

	for _, e := range pspList.Items {
		log.Printf("[INFO] PSP: %v \n", e.GetName())
		log.Printf("[INFO] Spec: %+v \n", e.Spec)
	}

	return pspList, nil
}

// KubeSecurityPolicyProvider ...
type KubeSecurityPolicyProvider struct {
	psps     *v1beta1.PodSecurityPolicyList
	pspMutex sync.Mutex
}

func (p *KubeSecurityPolicyProvider) getPolicies() (*v1beta1.PodSecurityPolicyList, error) {
	p.pspMutex.Lock()
	defer p.pspMutex.Unlock()

	//already got them?
	if p.psps == nil {
		ps, err := getPodSecurityPolicies()
		if err != nil {
			return nil, err
		}
		p.psps = ps
	}
	return p.psps, nil
}

// HasSecurityPolicies ...
func (p *KubeSecurityPolicyProvider) HasSecurityPolicies() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	b := len(psps.Items) > 0
	return &b, nil
}

// HasPrivilegedAccessRestriction ...
func (p *KubeSecurityPolicyProvider) HasPrivilegedAccessRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least on of the PSPs should have Privileged set to false
	var res bool
	for _, e := range psps.Items {
		if !e.Spec.Privileged {
			log.Printf("[NOTICE] PASS: Privileged is set to %v on Policy: %v", e.Spec.Privileged, e.GetName())
			res = true
			break
		}
	}

	if !res {
		log.Printf("[NOTICE] FAIL: NO Policies found with Privileged set.\n")
	}

	return &res, nil

}

// HasHostPIDRestriction ...
func (p *KubeSecurityPolicyProvider) HasHostPIDRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least on of the PSPs should have HostPID set to false
	var res bool
	for _, e := range psps.Items {
		if !e.Spec.HostPID {
			log.Printf("[NOTICE] PASS: HostPID is set to %v on Policy: %v\n", e.Spec.HostPID, e.GetName())
			res = true
			break
		}
	}

	if !res {
		log.Printf("[NOTICE] FAIL: NO Policies found with HostPID set.\n")
	}

	return &res, nil

}
