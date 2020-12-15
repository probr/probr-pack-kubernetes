package access_whitelisting

import (
	"flag"
	"log"
	"os"
	"strings"
	"testing"

	//	"citihub.com/compliance-as-code/internal/logfilter"
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

	status := godog.RunWithOptions("access_whitelisting_test", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {
	//	logfilter.Setup()
	var state accessWhitelisting

	csp := strings.ToLower(os.Getenv("CSP"))
	switch csp {
	case "azure":
		state = &accessWhitelistingAzure{}
	case "aws":
		//		state = &accessWhitelistingAWS{}
	default:
		log.Panicf("Cloud Provider '%s' not supported - set environment variable 'CSP'", csp)
	}

	s.BeforeSuite(state.setup)

	s.Step(`^the CSP provides a whitelisting capability for Object Storage containers$`, state.cspSupportsWhitelisting)
	s.Step(`^we examine the Object Storage container in environment variable "([^"]*)"$`, state.examineStorageContainer)
	s.Step(`^whitelisting is configured with the given IP address range or an endpoint$`, state.whitelistingIsConfigured)
	s.Step(`^security controls that Prevent Object Storage from being created without network source address whitelisting are applied$`, state.checkPolicyAssigned)
	s.Step(`^we provision an Object Storage container$`, state.provisionStorageContainer)
	s.Step(`^it is created with whitelisting entry "([^"]*)"$`, state.createWithWhitelist)
	s.Step(`^creation will "([^"]*)"$`, state.creationWill)

	s.AfterSuite(state.teardown)
}
