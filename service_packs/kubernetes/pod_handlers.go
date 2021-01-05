// Package kubernetes provides functions for interacting with Kubernetes and
// is built using the kubernetes client-go (https://github.com/kubernetes/client-go).
package kubernetes

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//needed for authentication against the various GCPs
	k8s "k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// PodCreationErrorReason provides an CSP agnostic reason for errors encountered when creating pods.
type PodCreationErrorReason int

// enum values for PodCreationErrorReason
const (
	UndefinedPodCreationErrorReason PodCreationErrorReason = iota
	PSPNoPrivilege
	PSPNoPrivilegeEscalation
	PSPAllowedUsersGroups
	PSPContainerAllowedImages
	PSPHostNamespace
	PSPHostNetwork
	PSPAllowedCapabilities
	PSPAllowedPortRange
	PSPAllowedVolumeTypes
	PSPSeccompProfile
	ImagePullError
	Blocked
	Unauthorized
)

func (r PodCreationErrorReason) String() string {
	return [...]string{"podcreation-error: undefined",
		"podcreation-error: psp-container-no-privilege",
		"podcreation-error: psp-container-no-privilege-escalation",
		"podcreation-error: psp-allowed-users-groups",
		"podcreation-error: psp-container-allowed-images",
		"podcreation-error: psp-host-namespace",
		"podcreation-error: psp-host-network",
		"podcreation-error: psp-allowed-capabilities",
		"podcreation-error: psp-allowed-portrange",
		"podcreation-error: psp-allowed-volume-types-profile",
		"podcreation-error: psp-allowed-seccomp-profile",
		"podcreation-error: image-pull-error",
		"podcreation-error: blocked"}[r]
}

// PodCreationError encapsulates the underlying pod creation error along with a map of platform agnostic
// PodCreationErrorReason codes.  Note that there could be more that one PodCreationErrorReason.  For
// example a pod may fail due to a 'psp-container-no-privilege' error and 'psp-host-network', in which
// case there would be two entries in the ReasonCodes map.
type PodCreationError struct {
	err         error
	ReasonCodes map[PodCreationErrorReason]*PodCreationErrorReason
}

type PodAudit struct {
	PodName         string
	Namespace       string
	ContainerName   string
	Image           string
	SecurityContext *apiv1.SecurityContext
}

func (p *PodCreationError) Error() string {
	return fmt.Sprintf("pod creation error: %v %v", p.ReasonCodes, p.err)
}

func getPods(c *k8s.Clientset, ns string) (*apiv1.PodList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pods, err := c.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}
	if pods == nil {
		return nil, fmt.Errorf("pod list returned nil")
	}

	log.Printf("[DEBUG] There are %d pods in the cluster\n", len(pods.Items))

	for i := 0; i < len(pods.Items); i++ {
		log.Printf("[DEBUG] Pod: %v %v\n", pods.Items[i].GetNamespace(), pods.Items[i].GetName())
	}

	return pods, nil
}

// GenerateUniquePodName creates a unique pod name based on the format: 'baseName'-'nanosecond time'-'random int'.
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

func waitForDelete(c *k8s.Clientset, ns string, n string) error {
	// Currently unused, not deleting yet in case it is useful elsewhere
	ps := c.CoreV1().Pods(ns)

	w, err := ps.Watch(context.Background(), metav1.ListOptions{})

	if err != nil {
		return err
	}

	log.Printf("[INFO] *** Waiting for DELETE on pod %v ...", n)

	for e := range w.ResultChan() {
		log.Printf("[DEBUG] Watch Probe Type: %v", e.Type)
		p, ok := e.Object.(*apiv1.Pod)
		if !ok {
			log.Printf("[WARN] Unexpected Watch Probe Type received for pod %v - skipping", p.GetObjectMeta().GetName())
			break
		}
		log.Printf("[INFO] Watch Container phase: %v", p.Status.Phase)
		log.Printf("[DEBUG] Watch Container status: %+v", p.Status.ContainerStatuses)

		if e.Type == "DELETED" {
			log.Printf("[DEBUG] DELETED probe received for pod %v", p.GetObjectMeta().GetName())
			break
		}

	}

	log.Printf("[INFO] *** Completed waiting for DELETE on pod %v", n)

	return nil
}
