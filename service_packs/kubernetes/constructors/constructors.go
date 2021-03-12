// Package constructors provides functions to prepare new objects (as described by the name of the function)
// This implements factory pattern.
package constructors

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/citihub/probr/config"
	"github.com/citihub/probr/utils"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodSpec constructs a simple pod object
func PodSpec(baseName string, namespace string) *apiv1.Pod {
	name := strings.Replace(baseName, "_", "-", -1)
	podName := uniquePodName(name)
	containerName := fmt.Sprintf("%s-probe-pod", name)
	log.Printf(fmt.Sprintf("[DEBUG] Creating pod spec with podName=%s and containerName=%s", podName, containerName))

	annotations := make(map[string]string)
	annotations["seccomp.security.alpha.kubernetes.io/pod"] = "runtime/default"

	return &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "probr-probe",
			},
			Annotations: annotations,
		},
		Spec: apiv1.PodSpec{
			SecurityContext: DefaultPodSecurityContext(),
			Containers: []apiv1.Container{
				{
					Name:            containerName,
					Image:           DefaultProbrImageName(),
					ImagePullPolicy: apiv1.PullIfNotPresent,
					Command:         DefaultEntrypoint(),
					SecurityContext: DefaultContainerSecurityContext(),
				},
			},
			NodeSelector: map[string]string{
				"kubernetes.io/os": "linux",
			},
		},
	}
}

// DefaultContainerSecurityContext returns an SC with the drop capabilities specified in config vars
func DefaultContainerSecurityContext() *apiv1.SecurityContext {
	return &apiv1.SecurityContext{
		Privileged:               utils.BoolPtr(false),
		AllowPrivilegeEscalation: utils.BoolPtr(false),
		Capabilities: &apiv1.Capabilities{
			Drop: GetContainerDropCapabilitiesFromConfig(),
		},
	}
}

// DefaultPodSecurityContext returns a basic PSC
func DefaultPodSecurityContext() *apiv1.PodSecurityContext {
	return &apiv1.PodSecurityContext{
		RunAsUser:          utils.Int64Ptr(1000),
		FSGroup:            utils.Int64Ptr(2000),
		RunAsGroup:         utils.Int64Ptr(3000),
		SupplementalGroups: []int64{1},
	}
}

// DefaultProbrImageName joins the registry and image name specified in config vars
func DefaultProbrImageName() string {
	// Service pack will not start without these vars, so we can rely on them being present
	return fmt.Sprintf(
		"%s/%s",
		config.Vars.ServicePacks.Kubernetes.AuthorisedContainerRegistry,
		config.Vars.ServicePacks.Kubernetes.ProbeImage)
}

// DefaultEntrypoint is used by all default pods
func DefaultEntrypoint() []string {
	return []string{
		"sleep",
		"3600",
	}
}

// GetContainerDropCapabilitiesFromConfig returns a Capability object with the drop caps specified in config vars
func GetContainerDropCapabilitiesFromConfig() []apiv1.Capability {
	return CapabilityObjectList(config.Vars.ServicePacks.Kubernetes.ContainerRequiredDropCapabilities)
}

// CapabilityObjectList converts a list of strings into a list of capability objects
func CapabilityObjectList(capList []string) []apiv1.Capability {
	var capabilities []apiv1.Capability

	for _, cap := range capList {
		if cap != "" {
			capabilities = append(capabilities, apiv1.Capability(cap))
		}
	}

	return capabilities
}

func uniquePodName(baseName string) string {
	//take base and add some uniqueness
	t := time.Now()
	rand.Seed(t.UnixNano())
	uniq := fmt.Sprintf("%v-%v", t.Format("020106-150405"), rand.Intn(100))

	return fmt.Sprintf("%v-%v", baseName, uniq)
}
