package encryption_in_flight

// EncryptionInFlight is an interface. For each CSP specific implementation
type EncryptionInFlight interface {
	setup()
	anAzureResourceGroupExists() error
	weProvisionAnObjectStorageBucket() error
	httpAccessIs(arg1 string) error
	httpsAccessIs(arg1 string) error
	creationWillWithAnErrorMatching(result, errDescription string) error

	detectObjectStorageUnencryptedTransferAvailable() error
	detectObjectStorageUnencryptedTransferEnabled() error
	createUnencryptedTransferObjectStorage() error
	detectsTheObjectStorage() error
	encryptedDataTrafficIsEnforced() error
	teardown()
}
