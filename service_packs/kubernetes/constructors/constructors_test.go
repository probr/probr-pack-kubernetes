// Package constructors provides functions to prepare new objects (as described by the name of the function)

// This implements factory pattern.

package constructors

import (
	"reflect"
	"strings"
	"testing"

	"github.com/citihub/probr/config"
	"github.com/citihub/probr/utils"
	apiv1 "k8s.io/api/core/v1"
)

func Test_uniquePodName(t *testing.T) {
	tests := []struct {
		testName string
		arg      string
	}{
		{
			testName: "Unique Pod Name Contains Base Name",
			arg:      "basename1",
		},
		{
			testName: "Unique Pod Name Contains Base Name",
			arg:      "base-name-2",
		},
		{
			testName: "Unique Pod Name Contains Base Name",
			arg:      "base_name_3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got := uniquePodName(tt.arg)
			if !strings.Contains(got, tt.arg) || len(got) <= len(tt.arg) {
				t.Errorf("uniquePodName() = %v, want %v", got, tt.arg)
			}
		})
	}
}

func TestPodSpec(t *testing.T) {
	type args struct {
		baseName                 string
		namespace                string
		containerSecurityContext *apiv1.SecurityContext
	}
	tests := []struct {
		name string
		args args
		want func(gotPod *apiv1.Pod, args args, t *testing.T)
	}{
		{
			name: "Pod's security context is always the default security context",
			args: args{
				baseName:  "pod",
				namespace: "pod",
			},
			want: func(gotPod *apiv1.Pod, args args, t *testing.T) {
				if !reflect.DeepEqual(gotPod.Spec.SecurityContext, DefaultPodSecurityContext()) {
					t.Errorf("PodSpec() should set the pod's security context using DefaultPodSecurityContext()")
				}
			},
		},
		{
			name: "Pod has at least one container",
			args: args{
				baseName:  "pod1",
				namespace: "pod1",
			},
			want: func(gotPod *apiv1.Pod, want args, t *testing.T) {
				gotContainers := gotPod.Spec.Containers
				if !(len(gotContainers) > 0) {
					t.Error("PodSpec() did not create a container object, but wanted at least one")
				}
			},
		},
		{
			name: "Container uses a unique pod name",
			args: args{
				baseName:                 "pod3",
				namespace:                "pod3",
				containerSecurityContext: nil,
			},
			want: func(gotPod *apiv1.Pod, args args, t *testing.T) {
				gotContainerName := gotPod.Spec.Containers[0].Name
				if len(gotContainerName) <= len(args.baseName) || !strings.Contains(gotContainerName, args.baseName) {
					t.Errorf("PodSpec() got container name '%s', but wanted: '%s'", gotContainerName, args.baseName)
				}
			},
		},
		{
			name: "Container uses the default probr image name",
			args: args{
				baseName:  "pod4",
				namespace: "pod4",
			},
			want: func(gotPod *apiv1.Pod, args args, t *testing.T) {
				gotImageName := gotPod.Spec.Containers[0].Image
				wantImageName := DefaultProbrImageName()
				if strings.Compare(gotImageName, wantImageName) != 0 {
					t.Errorf("PodSpec() got image name '%s', but wanted: '%s'", gotImageName, wantImageName)
				}
			},
		},
		{
			name: "Pod uses the provided name",
			args: args{
				baseName:  "pod5",
				namespace: "pod5",
			},
			want: func(gotPod *apiv1.Pod, want args, t *testing.T) {
				gotName := gotPod.ObjectMeta.Name
				if !strings.Contains(gotName, want.baseName) {
					t.Errorf("PodSpec() got name '%s', but wanted it to include: '%s'", gotName, want.baseName)
				}
			},
		},
		{
			name: "Pod uses the provided namespace",
			args: args{
				baseName:  "pod6",
				namespace: "pod6",
			},
			want: func(gotPod *apiv1.Pod, want args, t *testing.T) {
				gotNamespace := gotPod.ObjectMeta.Namespace
				if strings.Compare(gotNamespace, want.namespace) != 0 {
					t.Errorf("PodSpec() got namespace '%s', but wanted: '%s'", gotNamespace, want.namespace)
				}
			},
		},
		{
			name: "Pod Labels is not nil",
			args: args{
				baseName:                 "pod7",
				namespace:                "pod7",
				containerSecurityContext: nil,
			},
			want: func(gotPod *apiv1.Pod, want args, t *testing.T) {
				if gotPod.Labels == nil {
					t.Error("PodSpec() did not create a Labels object, but wanted at least an empty map")
				}
			},
		},
		{
			name: "Node Selector contains 'kubernetes.io/os':'linux'",
			args: args{
				baseName:                 "pod8",
				namespace:                "pod8",
				containerSecurityContext: nil,
			},
			want: func(gotPod *apiv1.Pod, want args, t *testing.T) {
				if gotPod.Spec.NodeSelector["kubernetes.io/os"] != "linux" {
					t.Error("PodSpec() did not include a node selector 'kubernetes.io/os':'linux'")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PodSpec(tt.args.baseName, tt.args.namespace)
			tt.want(got, tt.args, t)
		})
	}
}

func TestDefaultContainerSecurityContext(t *testing.T) {
	tests := []struct {
		name string
		want *apiv1.SecurityContext
	}{
		{
			name: "Very strict test to enforce expectations of DefaultContainerSecurityContext",
			want: &apiv1.SecurityContext{
				Privileged:               utils.BoolPtr(false),
				AllowPrivilegeEscalation: utils.BoolPtr(false),
				Capabilities: &apiv1.Capabilities{
					Drop: GetContainerDropCapabilitiesFromConfig(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultContainerSecurityContext(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultContainerSecurityContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultPodSecurityContext(t *testing.T) {
	tests := []struct {
		name string
		want *apiv1.PodSecurityContext
	}{
		{
			name: "Very strict test to enforce expectations of DefaultContainerSecurityContext",
			want: &apiv1.PodSecurityContext{
				RunAsUser:          utils.Int64Ptr(1000),
				FSGroup:            utils.Int64Ptr(2000),
				RunAsGroup:         utils.Int64Ptr(3000),
				SupplementalGroups: []int64{1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultPodSecurityContext(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultPodSecurityContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultProbrImageName(t *testing.T) {
	tests := []struct {
		name     string
		registry string
		image    string
		want     string
	}{
		{
			name:     "Ensure probr image name is built by joining the registry and image specified in config vars",
			registry: "string1",
			image:    "string2",
			want:     "string1/string2",
		},
		{
			name:     "Ensure probr image name is built by joining the registry and image specified in config vars",
			registry: "registryName",
			image:    "imageName",
			want:     "registryName/imageName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Vars.ServicePacks.Kubernetes.AuthorisedContainerRegistry = tt.registry
			config.Vars.ServicePacks.Kubernetes.ProbeImage = tt.image
			if got := DefaultProbrImageName(); got != tt.want {
				t.Errorf("DefaultProbrImageName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetContainerDropCapabilitiesFromConfig(t *testing.T) {
	tests := []struct {
		name         string
		capabilities []string
		want         []apiv1.Capability
	}{
		{
			name:         "Ensure list of capabilities is populated from convig vars",
			capabilities: []string{"value1", "value2"},
			want: []apiv1.Capability{
				apiv1.Capability("value1"),
				apiv1.Capability("value2"),
			},
		},
		{
			name:         "Ensure list of capabilities is populated from convig vars",
			capabilities: []string{"cap1", "cap2"},
			want: []apiv1.Capability{
				apiv1.Capability("cap1"),
				apiv1.Capability("cap2"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Vars.ServicePacks.Kubernetes.ContainerRequiredDropCapabilities = tt.capabilities
			if got := GetContainerDropCapabilitiesFromConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetContainerDropCapabilitiesFromConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCapabilityObjectList(t *testing.T) {
	tests := []struct {
		name         string
		capabilities []string
		want         []apiv1.Capability
	}{
		{
			name:         "Ensure list strings is converted to list of capability objects",
			capabilities: []string{"value1", "value2"},
			want: []apiv1.Capability{
				apiv1.Capability("value1"),
				apiv1.Capability("value2"),
			},
		},
		{
			name:         "Ensure list strings is converted to list of capability objects",
			capabilities: []string{"cap1", "cap2"},
			want: []apiv1.Capability{
				apiv1.Capability("cap1"),
				apiv1.Capability("cap2"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CapabilityObjectList(tt.capabilities); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CapabilityObjectList() = %v, want %v", got, tt.want)
			}
		})
	}
}
