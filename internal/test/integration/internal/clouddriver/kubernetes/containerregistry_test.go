package kubernetes_test

import (
	"fmt"
	"testing"

	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"
	"gitlab.com/citihub/probr/internal/utils"
)

func TestSetupContainerRegistyPod(t *testing.T) {

	pd, err := kubernetes.NewDefaultCRA().SetupContainerAccessTestPod(utils.StringPtr("docker.io"))

	fmt.Printf("pd: %v err: %v", pd, err)
}
