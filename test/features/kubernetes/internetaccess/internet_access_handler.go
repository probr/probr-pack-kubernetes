package internetaccess

import (	
	"os"
	"fmt"
	"path/filepath"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"

	"citihub.com/probr/internal/coreengine"
	"citihub.com/probr/test/features"
)

//this is the "TEST HANDLER" impl  and will get called when probr is invoked from the CLI or API
//all we do here is set the godog args based on what has been supplied (e.g. output path) 
//and call to the "feature" implementation (i.e the same impl when godog / go test is invoked)

//Init ...
func init() {
	n, c := "internet_access", coreengine.InternetAccess
	td := coreengine.TestDescriptor{Category: c, Name: n}

	coreengine.TestHandleFunc(td,TH)
}

//TH ...
func TH() (int, error) {
	probrRoot, err := features.GetProbrRoot()
	
	if err != nil {
		return -1, fmt.Errorf("unable to determine probr root - not able to perform tests")
	}
	
	featPath := filepath.Join(probrRoot,"test","features","kubernetes","internetaccess","features")
	
	//TODO: FIX - get this from env/arg
	outPath := filepath.Join(probrRoot,"testoutput")
	os.Mkdir(outPath,os.ModeDir)
		
	f, err := os.Create(filepath.Join(outPath,"iatestout.json"))
	if err != nil {
		return -2, err
	}

	opts := godog.Options{
		Format:    "cucumber",		
		Output: 	colors.Colored(f), 	
		Paths:     []string{featPath},	
	}

	status := godog.TestSuite{
		Name: "internetaccess",
		TestSuiteInitializer: TestSuiteInitialize,
		ScenarioInitializer:  ScenarioInitialize,
		Options: &opts,
	}.Run()
		
	return status, nil
}