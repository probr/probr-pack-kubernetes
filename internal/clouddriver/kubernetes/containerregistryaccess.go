package kubernetes

import (
	"strings"

	apiv1 "k8s.io/api/core/v1"
)

const (
	//TODO: default to these values for MVP - need to expose in future
	caNamespace   = "probr-container-access-test-ns"
	caTestImage   = "/busybox:latest"
	caContainer   = "container-access-test"
	caPodNameBase = "ca-test"
)

// ContainerRegistryAccess ...
type ContainerRegistryAccess interface {
	ClusterIsDeployed() *bool
	SetupContainerAccessTestPod(r *string) (*apiv1.Pod, error)
	TeardownContainerAccessTestPod(p *string) error
}

// CRA ...
type CRA struct {
	k Kubernetes
}

// NewCRA ...
func NewCRA(k Kubernetes) *CRA {
	c := &CRA{}
	c.k = k

	return c
}

// NewDefaultCRA ...
func NewDefaultCRA() *CRA {
	c := &CRA{}
	c.k = GetKubeInstance()

	return c
}

// ClusterIsDeployed ...
func (c *CRA) ClusterIsDeployed() *bool {
	return c.k.ClusterIsDeployed()
}

//SetupContainerAccessTestPod creates a pod with characteristics required for testing container access.
func (c *CRA) SetupContainerAccessTestPod(r *string) (*apiv1.Pod, error) {
	//full image is the repository + the caTestImage
	i := *r + caTestImage
	pname := GenerateUniquePodName(caPodNameBase + "-" + strings.ReplaceAll(*r, ".", "-"))
	ns, cname := caNamespace, caContainer
	// let caller handle result ...
	return c.k.CreatePod(&pname, &ns, &cname, &i, true, nil)
}

//TeardownContainerAccessTestPod ...
func (c *CRA) TeardownContainerAccessTestPod(p *string) error {
	ns := caNamespace
	err := c.k.DeletePod(p, &ns, false) //don't worry about waiting
	return err
}
