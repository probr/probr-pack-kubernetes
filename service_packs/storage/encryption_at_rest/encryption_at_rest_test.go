package encryption_at_rest

import (
	"flag"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

const csp = "CSP"

var opt = godog.Options{Output: colors.Colored(os.Stdout)}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("encryption_at_rest", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {
	//	logfilter.Setup()
	var state EncryptionAtRest

	cspEnv := strings.ToLower(os.Getenv(csp))
	switch cspEnv {
	case "azure":
		state = &EncryptionAtRestAzure{}
	case "aws":
		//		state = &EncryptionAtRestAWS{}
	default:
		log.Panicf("Environment variable CSP is defined as \"%s\"", cspEnv)
	}

	s.BeforeSuite(state.setup)

	s.Step(`^security controls that restrict data from being unencrypted at rest$`, state.securityControlsThatRestrictDataFromBeingUnencryptedAtRest)
	s.Step(`^we provision an Object Storage bucket$`, state.weProvisionAnObjectStorageBucket)
	s.Step(`^encryption at rest is "([^"]*)"$`, state.encryptionAtRestIs)
	s.Step(`^creation will "([^"]*)" with an error matching "([^"]*)"$`, state.creationWillWithAnErrorMatching)

	s.Step(`^there is a detective capability for creation of Object Storage without encryption at rest$`, state.policyOrRuleAvailable)
	s.Step(`^the capability for detecting the creation of Object Storage without encryption at rest is active$`, state.checkPolicyOrRuleAssignment)
	s.Step(`^the detective measure is enabled$`, state.policyOrRuleAssigned)
	s.Step(`^Object Storage is created with without encryption at rest$`, state.createContainerWithoutEncryption)
	s.Step(`^the detective capability detects the creation of Object Storage without encryption at rest$`, state.detectiveDetectsNonCompliant)
	s.Step(`^the detective capability enforces encryption at rest on the Object Storage Bucket$`, state.containerIsRemediated)
	s.AfterSuite(state.teardown)
}
