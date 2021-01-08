package utils

import (
	"bytes"
	"fmt"
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

func TestFindString(t *testing.T) {

	var tests = []struct {
		slice         []string
		val           string
		expectedIndex int
		expectedFound bool
	}{
		{[]string{"the", "and", "for", "so", "go"}, "and", 1, true},
		{[]string{"the", "and", "for", "so", "go"}, "for", 2, true},
		{[]string{"the", "and", "for", "so", "go"}, "in", -1, false},
	}

	for _, c := range tests {

		testName := fmt.Sprintf("FindString(%q,%q) - Expected:%d,%v", c.slice, c.val, c.expectedIndex, c.expectedFound)

		t.Run(testName, func(t *testing.T) {
			actualPosition, actualFound := FindString(c.slice, c.val)

			if actualPosition != c.expectedIndex || actualFound != c.expectedFound {
				t.Errorf("\nCall: FindString(%q,%q)\nResult: %d,%v\nExpected: %d,%v", c.slice, c.val, actualPosition, actualFound, c.expectedIndex, c.expectedFound)
			}
		})
	}
}

func TestReplaceBytesValue(t *testing.T) {

	var tests = []struct {
		bytes          []byte
		oldValue       string
		newValue       string
		expectedResult []byte
	}{
		{[]byte("oldstringhere"), "old", "new", []byte("newstringhere")},                       //Replace a word with no spaces
		{[]byte("oink oink oink"), "k", "ky", []byte("oinky oinky oinky")},                     //Replace a character
		{[]byte("oink oink oink"), "oink", "moo", []byte("moo moo moo")},                       //Replace a word with spaces
		{[]byte("nothing to replace"), "www", "something", []byte("nothing to replace")},       //Nothing to replace due to no match
		{[]byte(""), "a", "b", []byte("")},                                                     //Empty string
		{[]byte("Unicode character: ㄾ"), "Unicode", "Unknown", []byte("Unknown character: ㄾ")}, //String with unicode character
		{[]byte("Unicode character: ㄾ"), "ㄾ", "none", []byte("Unicode character: none")},       //Replace unicode character
	}

	for _, c := range tests {

		testName := fmt.Sprintf("ReplaceBytesValue(%q,%q,%q) - Expected:%q", string(c.bytes), c.oldValue, c.newValue, string(c.expectedResult))

		t.Run(testName, func(t *testing.T) {
			actualResult := ReplaceBytesValue(c.bytes, c.oldValue, c.newValue)

			if string(actualResult) != string(c.expectedResult) {
				t.Errorf("\nCall: ReplaceBytesValue(%q,%q,%q)\nResult: %q\nExpected: %q", string(c.bytes), c.oldValue, c.newValue, string(actualResult), string(c.expectedResult))
			}
		})
	}
}
