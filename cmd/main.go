package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/citihub/probr"
	"github.com/citihub/probr/audit"
	cliflags "github.com/citihub/probr/cmd/cli_flags"
	"github.com/citihub/probr/config"

	"github.com/citihub/probr-sdk/plugin"
	"github.com/hashicorp/go-hclog"
	hcplugin "github.com/hashicorp/go-plugin"
)

// ServicePack ...
type ServicePack struct {
	logger hclog.Logger
}

// Greet ...
func (g *ServicePack) Greet() string {
	g.logger.Debug("message from ServicePack_Probr.Greet")
	//g.logger.Debug("args...", os.Args)

	//return "Hello Probr!"
	return ProbrCoreLogic(g.logger)
}

// handshakeConfigs are used to just do a basic handshake between
// a hcplugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad hcplugins or executing a hcplugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = hcplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "probr.servicepack.kubernetes",
}

func main() {

	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	spProbr := &ServicePack{
		logger: logger,
	}
	// hcpluginMap is the map of hcplugins we can dispense.
	var hcpluginMap = map[string]hcplugin.Plugin{
		"kubernetes": &plugin.ServicePackPlugin{Impl: spProbr},
	}

	logger.Debug("message from Probr hcplugin", "foo", "bar")

	hcplugin.Serve(&hcplugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         hcpluginMap,
	})

	probr.Logger = &logger
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
		////log.Printf("Execution aborted - %v", "SIGTERM")
		probr.CleanupTmp()
		// TODO: Additional cleanup may be needed. For instance, any pods created during tests are not being dropped if aborted.
		os.Exit(0)
	}()
}

// ProbrCoreLogic ...
func ProbrCoreLogic(logger hclog.Logger) string {
	logger.Debug("message from ProbCoreLogic", "Start")

	// Setup for handling SIGTERM (Ctrl+C)
	//setupCloseHandler()

	err := config.Init("") // Create default config
	if err != nil {
		////log.Printf("[ERROR] error returned from config.Init: %v", err)
		os.Exit(2)
	}
	if len(os.Args[1:]) > 0 {
		////log.Printf("[DEBUG] Checking for CLI options or flags")
		cliflags.HandleRequestForRequiredVars()
		////log.Printf("[DEBUG] Handle pack option")
		cliflags.HandlePackOption()
		cliflags.HandleFlags()
	}

	config.Vars.LogConfigState()

	_, ts, err := probr.RunAllProbes()
	if err != nil {
		//log.Printf("[ERROR] Error executing tests %v", err)
		os.Exit(2) // Exit 2+ is for logic/functional errors
	}
	//log.Printf("[INFO] Overall test completion status: %v", s)
	audit.State.SetProbrStatus()

	out := probr.GetAllProbeResults(ts)
	if out == nil || len(out) == 0 {
		audit.State.Meta["no probes completed"] = fmt.Sprintf(
			"Probe results not written to file, possibly due to all being excluded or permissions on the specified output directory: %s",
			config.Vars.CucumberDir(),
		)
	}
	audit.State.PrintSummary()
	audit.State.WriteSummary()

	logger.Debug("message from ProbCoreLogic", "Complete")

	return "Hello Probr!"
}
