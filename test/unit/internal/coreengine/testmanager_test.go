package coreengine_test

import (
	"fmt"
	"testing"
	"citihub.com/probr/internal/coreengine"
)

func TestGetAvailableTests(t *testing.T) {
	alltests := coreengine.GetAvailableTests()

	if alltests == nil {
		t.Errorf("No tests returned")
	}

	for _, ts := range alltests {
		fmt.Println(*ts.Category, *ts.Uuid)
	}

}