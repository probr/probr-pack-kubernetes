package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	pack "github.com/citihub/probr-pack-kubernetes"
	"github.com/citihub/probr-sdk/audit"
	cliflags "github.com/citihub/probr-sdk/cli_flags"
	"github.com/citihub/probr-sdk/config"
	"github.com/citihub/probr-sdk/logging"
	"github.com/citihub/probr-sdk/plugin"
	"github.com/citihub/probr-sdk/probeengine"
	"github.com/citihub/probr-sdk/utils"
)

// ServicePack ...
type ServicePack struct {
}

// RunProbes ...
func (sp *ServicePack) RunProbes() error {
	log.Printf("[DEBUG] message from ServicePack_Probr.RunProbes")
	log.Printf("[DEBUG] args... %v", os.Args)

	return ProbrCoreLogic()
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "debug" {
		ProbrCoreLogic()
		return
	}
	spProbr := &ServicePack{}
	serveOpts := &plugin.ServeOpts{
		Pack: spProbr,
	}

	plugin.Serve(serveOpts)
}

// setupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
// Ref: https://golangcode.com/handle-ctrl-c-exit-in-terminal/
func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Printf("Execution aborted - %v", "SIGTERM")
		probeengine.CleanupTmp()
		// TODO: Additional cleanup may be needed. For instance, any pods created during tests are not being dropped if aborted.
		os.Exit(0)
	}()
}

// ProbrCoreLogic ...
func ProbrCoreLogic() (err error) {
	log.Printf("[INFO] message from ProbCoreLogic: %s", "Start")
	defer probeengine.CleanupTmp()
	setupCloseHandler()

	// Setup for handling SIGTERM (Ctrl+C)
	//setupCloseHandler()

	err = config.Init("") // Create default config
	if err != nil {
		log.Printf("[ERROR] error returned from config.Init: %v", err)
		return
	}
	if len(os.Args[1:]) > 0 {
		log.Printf("[DEBUG] Checking for CLI options or flags")
		cliflags.HandleRequestForRequiredVars()
		log.Printf("[DEBUG] Handle pack option")
		cliflags.HandlePackOption()
		cliflags.HandleFlags()
	}

	config.Vars.LogConfigState()

	logWriter := logging.ProbrLoggerOutput()
	log.SetOutput(logWriter) // TODO: This is a temporary patch, since logger output is being overritten while loading config vars

	s, ts, err := probeengine.RunAllProbes("kubernetes", pack.GetProbes())
	if err != nil {
		log.Printf("[ERROR] Error executing tests %v", err)
		return
	}
	log.Printf("[INFO] Overall test completion status: %v", s)
	audit.State.SetProbrStatus()

	_, success := probeengine.GetAllProbeResults(ts) // TODO: Use the results provided here
	audit.State.PrintSummary()
	audit.State.WriteSummary()

	log.Printf("[INFO] message from ProbCoreLogic: %s", "End")

	if !success {
		return utils.ReformatError("One or more probe scenarios were not successful. View the output logs for more details.")
	}
	return
}
