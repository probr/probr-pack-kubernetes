package kubernetes

// SecurityPolicyProvider ...
type SecurityPolicyProvider interface {
	HasSecurityPolicies() (*bool, error)
	HasPrivilegedAccessRestriction() (*bool, error)
	HasHostPIDRestriction() (*bool, error)
}
