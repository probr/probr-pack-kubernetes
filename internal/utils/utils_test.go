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

	long_string := "Verify that this somewhat long string remains unchanged in the output after being handled"
	err := ReformatError(long_string)
	err_contains_string := strings.Contains(err.Error(), long_string)
	if !err_contains_string {
		t.Logf("Test string was not properly included in retured error")
		t.Fail()
	}
}
