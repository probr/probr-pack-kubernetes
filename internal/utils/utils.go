// Package utils provides general utility methods.  The '*Ptr' functions were borrowed/inspired by the kubernetes go-client.
package utils

import (
	"runtime"
	"strings"
)

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

// GetCallerName retrieves the name of the function prior to the location it is called
func GetCallerName(up int) string {
	f := make([]uintptr, 1)
	runtime.Callers(up+2, f)                   // add full caller path to empty object
	step := runtime.FuncForPC(f[0] - 1).Name() // get full caller path in string form
	s := strings.Split(step, ".")              // split full caller path
	return s[len(s)-1]                         // select last element from caller path
}

func GetCallerFileLine() (string, int) {
	_, file, line, _ := runtime.Caller(2)
	return file, line
}
