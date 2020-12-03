package utils

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestReformatError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf) // Intercept expected Stderr output
	defer func() {
		log.SetOutput(os.Stderr) // Return to normal Stderr handling after function
	}()

	longString := "Verify that this somewhat long string remains unchanged in the output after being handled"
	err := ReformatError(longString)
	errContainsString := strings.Contains(err.Error(), longString)
	if !errContainsString {
		t.Logf("Test string was not properly included in retured error")
		t.Fail()
	}
}
