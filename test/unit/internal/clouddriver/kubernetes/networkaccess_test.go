package kubernetes_test

import (
	"log"
	"testing"

	"citihub.com/probr/internal/clouddriver/kubernetes"
	"github.com/stretchr/testify/assert"
)

func TestSetupNetworkAccessTestPod(t *testing.T) {

	p, err := kubernetes.SetupNetworkAccessTestPod()

	assert.Nil(t, err)
	assert.NotNil(t, p)
}

func TestIsURLAccessible(t *testing.T) {
	//need to dupe above .. fix?
	p, err := kubernetes.SetupNetworkAccessTestPod()

	assert.Nil(t, err)
	assert.NotNil(t, p)

	//now substance of this test
	url := "http://www.google.com"
	code, err := kubernetes.AccessURL(&url)

	assert.Nil(t, err)

	log.Printf("[NOTICE] URL: %v Result: %v", url, code)
}

func TestTeardownNetworkAccess(t *testing.T) {
	err := kubernetes.TeardownNetworkAccessTestPod()

	assert.Nil(t, err)
}
