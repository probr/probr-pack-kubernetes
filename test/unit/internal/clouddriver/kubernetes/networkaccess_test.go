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
	n := p.GetObjectMeta().GetName()
	code, err := kubernetes.AccessURL(&n, &url)


	log.Printf("[NOTICE] URL: %v Result: %v", url, code)
}