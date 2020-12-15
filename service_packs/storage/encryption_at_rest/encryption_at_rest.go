package encryption_at_rest

// EncryptionAtRest is an interface. For each CSP specific implementation
type EncryptionAtRest interface {
	setup()
	securityControlsThatRestrictDataFromBeingUnencryptedAtRest() error
	weProvisionAnObjectStorageBucket() error
	encryptionAtRestIs(encryptionOption string) error
	creationWillWithAnErrorMatching(result string) error
	policyOrRuleAvailable() error
	checkPolicyOrRuleAssignment() error
	policyOrRuleAssigned() error
	prepareToCreateContainer() error
	createContainerWithoutEncryption() error
	detectiveDetectsNonCompliant() error
	containerIsRemediated() error
	teardown()
}
