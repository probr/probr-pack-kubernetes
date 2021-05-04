// Package constructors provides functions to prepare new objects (as described by the name of the function)
// This implements factory pattern.
package constructors

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/citihub/probr-sdk/utils"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PersistentVolumeClaimConfig holds the state of the PVC
type PersistentVolumeClaimConfig struct {
	Name       string // Name of the PVC. If set, overrides NamePrefix
	NamePrefix string // NamePrefix defaults to "pvc-" if unspecified
	ClaimSize  string // ClaimSize must be specified in the Quantity format. Defaults to 2Gi if unspecified

	AccessModes      []apiv1.PersistentVolumeAccessMode // AccessModes defaults to RWO if unspecified
	Annotations      map[string]string
	Selector         *metav1.LabelSelector
	StorageClassName *string

	VolumeMode *apiv1.PersistentVolumeMode // VolumeMode defaults to nil if unspecified or specified as the empty string
}

// PodSpec constructs a simple pod object
func PodSpec(baseName, namespace, image string) *apiv1.Pod {
	name := strings.Replace(baseName, "_", "-", -1)
	podName := uniquePodName(name)
	containerName := fmt.Sprintf("%s-probe-pod", name)
	log.Printf("[DEBUG] Creating pod spec with podName=%s and containerName=%s", podName, containerName)

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
					Image:           image,
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
			Drop: CapabilityObjectList([]string{"NET_RAW"}),
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

// DefaultEntrypoint is used by all default pods
func DefaultEntrypoint() []string {
	return []string{
		"sleep",
		"3600",
	}
}

// DynamicPersistentVolumeClaim constructs a simple Dynamic PersistentVolumeClaim
func DynamicPersistentVolumeClaim(baseName, namespace, storageClass string) *apiv1.PersistentVolumeClaim {

	name := strings.Replace(baseName, "_", "-", -1)
	pvcName := uniquePodName(name)

	config := PersistentVolumeClaimConfig{
		Name:             pvcName,
		ClaimSize:        "1Gi",
		AccessModes:      []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
		StorageClassName: &storageClass,
	}

	return &apiv1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.Name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": "probr-probe",
			},
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: config.AccessModes,
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					apiv1.ResourceStorage: resource.MustParse(config.ClaimSize),
				},
			},
			StorageClassName: config.StorageClassName,
		},
	}
}

// AddPVCToPod adds a PersistentVolumeClaim to a Pod
func AddPVCToPod(pod *apiv1.Pod, pvc *apiv1.PersistentVolumeClaim) {
	pvcSource := apiv1.PersistentVolumeClaimVolumeSource{
		ClaimName: pvc.ObjectMeta.Name,
	}

	volume := apiv1.Volume{
		Name: "probr",
		VolumeSource: apiv1.VolumeSource{
			PersistentVolumeClaim: &pvcSource,
		},
	}

	volumeMount := apiv1.VolumeMount{
		Name:      "probr",
		MountPath: "/probr",
	}

	pod.Spec.Volumes = append(pod.Spec.Volumes, volume)
	pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, volumeMount)
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
	uniq := fmt.Sprintf("%v-%v%v", t.Format("020106-150405"), rand.Intn(100), rand.Intn(100))

	return fmt.Sprintf("%v-%v", baseName, uniq)
}
