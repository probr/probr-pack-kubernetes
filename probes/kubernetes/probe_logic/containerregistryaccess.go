package k8s_logic

import (
	"strings"

	"github.com/citihub/probr/internal/config"
	apiv1 "k8s.io/api/core/v1"
)

const (
	caNamespace   = "probr-container-access-test-ns"
	caContainer   = "container-access-test"
	caPodNameBase = "ca-test"
)

// ContainerRegistryAccess interface defines the methods to support container registry access tests.
type ContainerRegistryAccess interface {
	ClusterIsDeployed() *bool
	SetupContainerAccessProbePod(r string) (*apiv1.Pod, *PodAudit, error)
	TeardownContainerAccessProbePod(p *string, e string) error
}

// CRA implements the ContainerRegistryAccess interface.
type CRA struct {
	k Kubernetes
}

// NewCRA creates a new CRA with the supplied kubernetes instance.
func NewCRA(k Kubernetes) *CRA {
	c := &CRA{}
	c.k = k

	return c
}

// NewDefaultCRA creates a new CRA using the default kubernetes instance.
func NewDefaultCRA() *CRA {
	c := &CRA{}
	c.k = GetKubeInstance()

	return c
}

// ClusterIsDeployed verifies if a cluster is deployed.
func (c *CRA) ClusterIsDeployed() *bool {
	return c.k.ClusterIsDeployed()
}

//SetupContainerAccessProbePod creates a pod with characteristics required for testing container access.
func (c *CRA) SetupContainerAccessProbePod(r string) (*apiv1.Pod, *PodAudit, error) {
	//full image is the repository + the configured image
	i := r + "/" + config.Vars.ProbeImage
	pname := GenerateUniquePodName(caPodNameBase + "-" + strings.ReplaceAll(r, ".", "-"))
	ns, cname := caNamespace, caContainer
	// let caller handle result ...
	return c.k.CreatePod(pname, ns, cname, i, true, nil)
}

//TeardownContainerAccessProbePod deletes the supplied test pod in the container registry access namespace.
func (c *CRA) TeardownContainerAccessProbePod(p *string, e string) error {
	ns := caNamespace
	err := c.k.DeletePod(p, &ns, false, e) //don't worry about waiting
	return err
}
