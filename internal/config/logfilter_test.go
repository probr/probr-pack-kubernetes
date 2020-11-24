package config

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/citihub/probr/internal/utils"
)

func tmpLogger(test_string, level string) bytes.Buffer {
	defer func() {
		log.SetOutput(os.Stderr) // Return to normal Stderr handling after function
	}()
	var buf bytes.Buffer
	SetLogFilter(level, &buf) // Intercept expected Stderr output
	log.Printf(test_string)
	return buf
}

func bufferShouldLog(t *testing.T, test_string string, buf bytes.Buffer) {
	if len(buf.String()) < len(test_string) {
		file, line := utils.CallerFileLine()
		t.Logf("%v:%v:%s: Test string was not written to logs as expected: '%s'", file, line, utils.CallerName(0), buf.String())
		t.Fail()
	} else if len(buf.String()) == len(test_string) {
		file, line := utils.CallerFileLine()
		t.Logf("%v:%v:%s: Logger did not append timestamp to test string as expected: '%s'", file, line, utils.CallerName(0), buf.String())
		t.Fail()
	}
}

func bufferShouldNotLog(t *testing.T, test_string string, buf bytes.Buffer) {
	if len(buf.String()) > len(test_string) {
		file, line := utils.CallerFileLine()
		t.Logf("%v:%v:%s: Test string was written to logs, but not expected: '%s'", file, line, utils.CallerName(0), buf.String())
		t.Fail()
	}
}

func TestLog(t *testing.T) {
	test_string := "[ERROR] This should log an error"
	buf := tmpLogger(test_string, "ERROR")
	bufferShouldLog(t, test_string, buf)
}

func TestLogLevel(t *testing.T) {
	testString := "[NOTICE] this is a notice test string"

	// Validate lower than debug level does not print
	buf := tmpLogger(testString, "ERROR")
	bufferShouldNotLog(t, testString, buf)

	// Validate matching level prints
	buf = tmpLogger(testString, "NOTICE")
	bufferShouldLog(t, testString, buf)

	// Validate higer than debug level prints
	buf = tmpLogger(testString, "DEBUG")
	bufferShouldLog(t, testString, buf)

}
