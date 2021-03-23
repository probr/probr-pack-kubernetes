// Package utils provides general utility methods.  The '*Ptr' functions were borrowed/inspired by the kubernetes go-client.
package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/markbates/pkger"
)

func init() {
}

// BoolPtr returns a pointer to a bool
func BoolPtr(b bool) *bool {
	return &b
}

// StringPtr returns a pointer to the passed string.
func StringPtr(s string) *string {
	return &s
}

// Int64Ptr returns a pointer to an int64
func Int64Ptr(i int64) *int64 {
	return &i
}

// FindString searches a []string for a specific value.
// If found, returns the index of first occurrence, and True. If not found, returns -1 and False.
func FindString(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// CallerName retrieves the name of the function prior to the location it is called
// If using CallerName(0), the current function's name will be returned
// If using CallerName(1), the current function's parent name will be returned
// If using CallerName(2), the current function's parent's parent name will be returned
func CallerName(up int) string {
	s := strings.Split(CallerPath(up+1), ".") // split full caller path
	return s[len(s)-1]                        // select last element from caller path
}

// CallerPath checks the goroutine's stack of function invocation and returns the following:
// For up=0, return full caller path for caller function
// For up=1, returns full caller path for caller of caller
func CallerPath(up int) string {
	f := make([]uintptr, 1)
	runtime.Callers(up+2, f)                  // add full caller path to empty object
	return runtime.FuncForPC(f[0] - 1).Name() // get full caller path in string form
}

// CallerFileLine returns file name and line of invoker
// Similar to CallerName(1), but with file and line returned
func CallerFileLine() (string, int) {
	_, file, line, _ := runtime.Caller(2)
	return file, line
}

// ReformatError prefixes the error string ready for logging and/or output
func ReformatError(e string, v ...interface{}) error {
	var b strings.Builder
	b.WriteString("[ERROR] ")
	b.WriteString(e)

	s := fmt.Sprintf(b.String(), v...)

	return fmt.Errorf(s)
}

// ReadStaticFile returns the bytes for a given static file
// Path:
//  In most cases it will be ReadStaticFile(assetDir, fileName).
//  It could also be used as ReadStaticFile(assetDir, subfolder, filename)
func ReadStaticFile(path ...string) ([]byte, error) {

	// Validation for empty path
	if path == nil || len(path) == 0 {
		return nil, ReformatError("Path argument cannot be empty")
	}

	filename := path[len(path)-1]           // file name is the last string argument
	dirpathSlice := path[0:(len(path) - 1)] // folder path

	dirPath := ""
	for _, folder := range dirpathSlice {
		dirPath = filepath.Join(dirPath, folder)
	}
	if !filepath.IsAbs(dirPath) {
		dirPath = filepath.Join("/", dirPath) //Need the abs path (/) for pgker to work
	}

	filepath := filepath.Join(dirPath, filename)

	// If pkged.go file has been generated using pkger cli tool, this will open the file from bundled memory buffer.
	// Otherwise, this will read from file system
	// See: https://github.com/markbates/pkger

	f, err := pkger.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}

// ReplaceBytesValue replaces a substring with a new value for a given string in bytes
func ReplaceBytesValue(b []byte, old string, new string) []byte {
	newString := strings.Replace(string(b), old, new, -1)
	return []byte(newString)
}

// ReplaceBytesMultipleValues replaces multiple substring with a new value for a given string in bytes
func ReplaceBytesMultipleValues(b []byte, replacer *strings.Replacer) []byte {
	newString := replacer.Replace(string(b))
	return []byte(newString)
}

// AuditPlaceholders creates empty objects to reduce code repetition when auditing probe steps
func AuditPlaceholders() (strings.Builder, interface{}, error) {
	return *new(strings.Builder), *new(interface{}), *new(error)
}

// WriteAllowed determines whether a given filepath can be written, considering both permissions and overwrite flag
func WriteAllowed(path string, overwrite bool) bool {
	_, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil && overwrite == false {
		//log.Printf("[ERROR] OverwriteHistoricalAudits is set to false, preventing this from writing to file: %s", path)
		return false
	} else if os.IsPermission(err) {
		//log.Printf("[ERROR] Permissions prevent this from writing to file: %s", path)
		return false
	} else if err != nil {
		//log.Printf("[ERROR] Could not create or write to file: %s. Error: %s", path, err)
		return false
	}
	return true
}
