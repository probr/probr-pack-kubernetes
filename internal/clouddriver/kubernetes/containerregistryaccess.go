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

//SetupContainerAccessTestPod creates a pod with characteristics required for testing container access.
func SetupContainerAccessTestPod(r *string) (*apiv1.Pod, error) {
	//full image is the repository + the caTestImage
	i := *r + caTestImage
	pname := caPodNameBase + "-" + strings.ReplaceAll(*r, ".", "-")
	ns, cname := caNamespace, caContainer
	p, err := CreatePod(&pname, &ns, &cname, &i, true, nil)

	if err != nil {
		return nil, err
	}

	return p, nil
}

//TeardownContainerAccessTestPod ...
func TeardownContainerAccessTestPod(p *string) error {
	ns := caNamespace
	err := DeletePod(p, &ns, true)
	return err
}
