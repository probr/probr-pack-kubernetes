package coreengine_test

import (
	"testing"

	"github.com/citihub/probr/internal/coreengine"
	"github.com/stretchr/testify/assert"
)

func TestTestStatus(t *testing.T) {
	ts := coreengine.Pending
	assert.Equal(t, ts.String(), "Pending")

	ts = coreengine.Running
	assert.Equal(t, ts.String(), "Running")

	ts = coreengine.CompleteSuccess
	assert.Equal(t, ts.String(), "CompleteSuccess")

	ts = coreengine.CompleteFail
	assert.Equal(t, ts.String(), "CompleteFail")

	ts = coreengine.Error
	assert.Equal(t, ts.String(), "Error")
}
