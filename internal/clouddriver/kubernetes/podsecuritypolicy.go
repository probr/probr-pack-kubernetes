package kubernetes

import (	
	"log"
	"context"
	"time"
	"fmt"
		
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/api/policy/v1beta1"
)

//ClusterHasPSP - determines if the cluster has any Pod Security Policies set
func ClusterHasPSP() (bool, error) {
	psps, err := getPodSecurityPolicies()
	if err != nil {
		return false, err
	}

	return len(psps.Items)>0, nil
}

//PrivilegedAccessIsRestricted - look for a PodSecurityPolicy with 'Privileged' set to false (ie. NOT privileged)
func PrivilegedAccessIsRestricted() (bool, error) {
	psps, err := getPodSecurityPolicies()
	if err != nil {
		return false, err
	}

	//at least on of the PSPs should have Privileged set to false
	for i := 0; i < len(psps.Items); i++ {
		if !psps.Items[i].Spec.Privileged {
			log.Printf("PASS: Privileged is set to %v on Policy: %v", psps.Items[i].Spec.Privileged, psps.Items[i].GetName())
			return true, nil
		}						
	}

	log.Printf("FAIL: NO Policy's found with Privileged set.\n")

	return false, nil
}

//HostPIDIsRestricted - look for a PodSecurityPolicy with 'HostPID' set to false (i.e. NO Access to HostPID )
func HostPIDIsRestricted() (bool, error) {
	psps, err := getPodSecurityPolicies()
	if err != nil {
		return false, err
	}

	//at least on of the PSPs should have HostPID set to false
	for i := 0; i < len(psps.Items); i++ {
		if !psps.Items[i].Spec.HostPID {
			log.Printf("PASS: HostPID is set to %v on Policy: %v\n", psps.Items[i].Spec.HostPID, psps.Items[i].GetName())
			return true, nil
		}						
	}

	log.Printf("FAIL: NO Policy's found with HostPID set.\n")

	return false, nil
}

//getPodSecurityPolicies ...
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
		
	log.Printf("There are %d psp policies in the cluster\n", len(pspList.Items))

	for i := 0; i < len(pspList.Items); i++ {		
		log.Printf("PSP: %v \n", pspList.Items[i].GetName())
		log.Printf("Spec: %+v \n", pspList.Items[i].Spec)		
	}
	 	
	return pspList, nil
}
