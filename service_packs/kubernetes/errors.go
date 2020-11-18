package kubernetes

import "k8s.io/apimachinery/pkg/api/errors"

func isAlreadyExists(err error) bool {
	if se, ok := err.(*errors.StatusError); ok {
		//409 is "already exists"
		return se.ErrStatus.Code == 409
	}
	return false
}

func isForbidden(err error) bool {
	if se, ok := err.(*errors.StatusError); ok {
		//403 is "forbidden"
		return se.ErrStatus.Code == 403
	}
	return false
}
