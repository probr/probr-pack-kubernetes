package pod_security_policy

// SecurityPolicyProvider defines a set of methods for interrogating the security policies set on the kubernetes cluster.
type SecurityPolicyProvider interface {
	HasSecurityPolicies() (*bool, error)
	HasPrivilegedAccessRestriction() (*bool, error)
	HasHostPIDRestriction() (*bool, error)
	HasHostIPCRestriction() (*bool, error)
	HasHostNetworkRestriction() (*bool, error)
	HasAllowPrivilegeEscalationRestriction() (*bool, error)
	HasRootUserRestriction() (*bool, error)
	HasNETRAWRestriction() (*bool, error)
	HasAllowedCapabilitiesRestriction() (*bool, error)
	HasAssignedCapabilitiesRestriction() (*bool, error)
	HasHostPortRestriction() (*bool, error)
	HasVolumeTypeRestriction() (*bool, error)
	HasSeccompProfileRestriction() (*bool, error)
}
