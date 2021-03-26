package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	pack "github.com/citihub/probr-pack-kubernetes"
	cliflags "github.com/citihub/probr-pack-kubernetes/cmd/cli_flags"
	"github.com/citihub/probr-sdk/audit"
	"github.com/citihub/probr-sdk/config"
	"github.com/citihub/probr-sdk/logging"
	"github.com/citihub/probr-sdk/plugin"
	"github.com/citihub/probr-sdk/probeengine"
)

// ServicePack ...
type ServicePack struct {
}

// Greet ...
func (sp *ServicePack) Greet() string {
	log.Printf("[DEBUG] message from ServicePack_Probr.Greet")
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
func ProbrCoreLogic() string {
	log.Printf("[INFO] message from ProbCoreLogic: %s", "Start")

	// Setup for handling SIGTERM (Ctrl+C)
	//setupCloseHandler()

	err := config.Init("") // Create default config
	if err != nil {
		log.Printf("[ERROR] error returned from config.Init: %v", err)
		os.Exit(2)
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
		os.Exit(2) // Exit 2+ is for logic/functional errors
	}
	log.Printf("[INFO] Overall test completion status: %v", s)
	audit.State.SetProbrStatus()

	out := probeengine.GetAllProbeResults(ts)
	if out == nil || len(out) == 0 {
		audit.State.Meta["no probes completed"] = fmt.Sprintf(
			"Probe results not written to file, possibly due to all being excluded or permissions on the specified output directory: %s",
			config.Vars.CucumberDir(),
		)
	}
	audit.State.PrintSummary()
	audit.State.WriteSummary()

	log.Printf("[INFO] message from ProbCoreLogic: %s", "End")

	return "Hello Probr!"
}
