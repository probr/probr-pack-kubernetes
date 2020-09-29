package kubernetes_test

import (
	"fmt"
	"testing"

	"github.com/citihub/probr/internal/clouddriver/kubernetes"
	"github.com/citihub/probr/internal/utils"
)

func TestSetupContainerRegistyPod(t *testing.T) {

	pd, err := kubernetes.NewDefaultCRA().SetupContainerAccessTestPod(utils.StringPtr("docker.io"))

	fmt.Printf("pd: %v err: %v", pd, err)
}
