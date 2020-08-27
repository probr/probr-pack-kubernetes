package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/citihub/probr/test/features"
)

func Test_LogAndReturnError(t *testing.T) {
	one := "one"
	two := "two"

	err := features.LogAndReturnError("error %v, %v", one, two)

	assert.NotNil(t, err)
	assert.Equal(t, "[ERROR] error one, two", err.Error())

}
