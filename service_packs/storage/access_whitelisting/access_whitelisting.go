package access_whitelisting

// EncryptionInFlight is an interface. For each CSP specific implementation
type accessWhitelisting interface {
	setup()
	cspSupportsWhitelisting() error
	examineStorageContainer(containerName string) error
	whitelistingIsConfigured() error
	checkPolicyAssigned() error
	provisionStorageContainer() error
	createWithWhitelist(ipPrefix string) error
	creationWill(result string) error
	teardown()
}
