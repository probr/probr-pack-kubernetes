package encryption_in_flight

import (
	"flag"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var opt = godog.Options{Output: colors.Colored(os.Stdout)}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("encryption_in_flight", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {
	// logfilter.Setup()
	var state EncryptionInFlight
	csp := os.Getenv("CSP")

	switch strings.ToLower(csp) {
	case "azure":
		state = &EncryptionInFlightAzure{}
	case "aws":
	//	state = &EncryptionInFlightAWS{}
	default:
		log.Panicf("Environment variable CSP is defined as \"%s\"", csp)
	}

	s.BeforeSuite(state.setup)

	s.Step(`^security controls that restrict data from being unencrypted in flight$`, state.securityControlsThatRestrictDataFromBeingUnencryptedInFlight)
	s.Step(`^we provision an Object Storage bucket$`, state.weProvisionAnObjectStorageBucket)
	s.Step(`^http access is "([^"]*)"$`, state.httpAccessIs)
	s.Step(`^https access is "([^"]*)"$`, state.httpsAccessIs)
	s.Step(`^creation will "([^"]*)" with an error matching "([^"]*)"$`, state.creationWillWithAnErrorMatching)

	s.Step(`^there is a detective capability for creation of Object Storage with unencrypted data transfer enabled$`, state.detectObjectStorageUnencryptedTransferAvailable)
	s.Step(`^the capability for detecting the creation of Object Storage with unencrypted data transfer enabled is active$`, state.detectObjectStorageUnencryptedTransferEnabled)
	s.Step(`^Object Storage is created with unencrypted data transfer enabled$`, state.createUnencryptedTransferObjectStorage)
	s.Step(`^the detective capability detects the creation of Object Storage with unencrypted data transfer enabled$`, state.detectsTheObjectStorage)
	s.Step(`^the detective capability enforces encrypted data transfer on the Object Storage Bucket$`, state.encryptedDataTrafficIsEnforced)

	s.AfterSuite(state.teardown)
}
