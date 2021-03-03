package errors

import "k8s.io/apimachinery/pkg/api/errors"

// IsStatusCode validates whether an error is a StatusError with a specific status code
func IsStatusCode(expected int32, err error) bool {
	if se, ok := err.(*errors.StatusError); ok {
		return se.ErrStatus.Code == expected
	}
	return false
}
