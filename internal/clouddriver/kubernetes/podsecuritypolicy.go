package kubernetes

import (
	"context"
	"fmt"
	"log"
	"time"

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
	pspTestNamespace = "network-access-test-ns"
	pspTestImage     = "curlimages/curl"
	pspTestContainer = "curlimages"
	pspTestPodName   = "psp-test-pod"
)

//ClusterHasPSP determines if the cluster has any Pod Security Policies set.
func ClusterHasPSP() (*bool, error) {
	psps, err := getPodSecurityPolicies()
	if err != nil {
		return nil, err
	}

	b := len(psps.Items) > 0
	return &b, nil
}

//PrivilegedAccessIsRestricted looks for a PodSecurityPolicy with 'Privileged' set to false (ie. NOT privileged).
func PrivilegedAccessIsRestricted() (*bool, error) {
	psps, err := getPodSecurityPolicies()
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

//HostPIDIsRestricted looks for a PodSecurityPolicy with 'HostPID' set to false (i.e. NO Access to HostPID ).
func HostPIDIsRestricted() (*bool, error) {
	psps, err := getPodSecurityPolicies()
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

//CreatePODSettingPrivilegedAccess ...
func CreatePODSettingPrivilegedAccess(pa PrivilegedAccess) (*apiv1.Pod, error) {
	
	
	b := pa == WithPrivilegedAccess
	sc := apiv1.SecurityContext{
		Privileged: &b,		
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
