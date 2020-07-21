package kubernetes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"citihub.com/probr/internal/clouddriver/kubernetes"
)

func TestSetupNetworkAccessTestPod(t *testing.T) {

	p, err := kubernetes.SetupNetworkAccessTestPod()

	assert.Nil(t,err)
	assert.NotNil(t,p)
}

func TestIsURLAccessible(t *testing.T) {
	//need to dupe above .. fix?
	p, err := kubernetes.SetupNetworkAccessTestPod()

	assert.Nil(t,err)
	assert.NotNil(t,p)

	//now substance of this test
	url := "http://www.google.com"
	code, err := kubernetes.AccessURL(&url)

	assert.Nil(t, err)	

	//TODO: only true for now.  Ultimately network access should be locked down.
	assert.True(t, code==200)

}

func TestTeardownNetworkAccess(t *testing.T) {
	err := kubernetes.TeardownNetworkAccessTestPod()

	assert.Nil(t,err)
}