package kubernetes_test

import (
	"fmt"
	"testing"

	"github.com/citihub/probr/internal/clouddriver/kubernetes"
)

func TestSetupContainerRegistyPod(t *testing.T) {

	pd, _, err := kubernetes.NewDefaultCRA().SetupContainerAccessProbePod("docker.io")

	fmt.Printf("pd: %v err: %v", pd, err)
}
