package utils

import (
	"testing"

	"github.com/citihub/probr/probes"
	"github.com/stretchr/testify/assert"
)

func Test_LogAndReturnError(t *testing.T) {
	one := "one"
	two := "two"

	err := probes.LogAndReturnError("error %v, %v", one, two)

	assert.NotNil(t, err)
	assert.Equal(t, "[ERROR] error one, two", err.Error())

}
