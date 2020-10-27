package kubernetes_test

import (
	"log"
	"testing"

	"github.com/citihub/probr/internal/clouddriver/kubernetes"
	"github.com/stretchr/testify/assert"
)

//TODO: this will be revised when the unit/integration tests are refactored to be properly mocked
var na = kubernetes.NewDefaultNA()

func TestSetupNetworkAccessProbePod(t *testing.T) {

	p, _, err := na.SetupNetworkAccessProbePod()

	assert.Nil(t, err)
	assert.NotNil(t, p)
}

func TestIsURLAccessible(t *testing.T) {
	//need to dupe above .. fix?
	p, _, err := na.SetupNetworkAccessProbePod()

	assert.Nil(t, err)
	assert.NotNil(t, p)

	//now substance of this test
	url := "http://www.google.com"
	n := p.GetObjectMeta().GetName()
	code, err := na.AccessURL(&n, &url)

	log.Printf("[INFO] URL: %v Result: %v Error: %v", url, code, err)
}
