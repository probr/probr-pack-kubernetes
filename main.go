package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/probr/probr-pack-kubernetes/internal/config"
	"github.com/probr/probr-pack-kubernetes/internal/connection"
	"github.com/probr/probr-pack-kubernetes/internal/summary"
	"github.com/probr/probr-pack-kubernetes/pack"

	audit "github.com/probr/probr-sdk/audit"
	sdkConfig "github.com/probr/probr-sdk/config"
	"github.com/probr/probr-sdk/plugin"
	probeengine "github.com/probr/probr-sdk/probeengine"
	"github.com/probr/probr-sdk/utils"
)

var (
	// ServicePackName is the name for the service pack
	// TODO: Return binary name instead? Or add this to the interface to use instead of the binary name?
	ServicePackName = "Kubernetes"

	// Version is the main version number that is being run at the moment
	Version = "0.1.0"

	// Prerelease is a marker for the version. If this is "" (empty string)
	// then it means that it is a final release. Otherwise, this is a pre-release
	// such as "dev" (in development), "beta", "rc", etc.
	// This should only be modified thru ldflags in make file. See 'make release'
	Prerelease = "dev"

	// GitCommitHash references the commit id at build time
	// This should only be modified thru ldflags in make file. See 'make release'
	GitCommitHash = ""

	// BuiltAt is the build date
	// This should only be modified thru ldflags in make file. See 'make release'
	BuiltAt = ""
)

// ServicePack ...
type ServicePack struct {
}

// RunProbes is required to meet the Probr Service Pack interface
func (sp *ServicePack) RunProbes() error {
	return ProbrCoreLogic()
}

// main is executed when this file is called as a binary or `go run`
func main() {
	versionCmd, runCmd := setFlags()
	handleCommands(versionCmd, runCmd)
}

func setFlags() (versionCmd, runCmd *flag.FlagSet) {
	// > probr version [-v]
	versionCmd = flag.NewFlagSet("version", flag.ExitOnError)
	config.Vars.Verbose = *versionCmd.Bool("v", false, "Display extended version information") // TODO: Harness '-v' in the standard probr execution

	// > probr
	runCmd = flag.NewFlagSet("run", flag.ExitOnError)
	runCmd.StringVar(&config.Vars.VarsFile, "varsfile", "", "path to config file")
	runCmd.StringVar(&config.Vars.ServicePacks.Kubernetes.KubeConfigPath, "kubeconfig", "", "kube config file")
	return
}

func handleCommands(versionCmd, runCmd *flag.FlagSet) {
	subCommand := ""
	if len(os.Args) > 1 {
		subCommand = os.Args[1]
	}
	switch subCommand {
	case "version":
		versionCmd.Parse(os.Args[2:])
		printVersion(os.Stdout)

	case "debug": // Same cli args as run. Use this to bypass plugin and execute directly for debugging
		// Parse cli args
		runCmd.Parse(os.Args[2:]) // Skip first arg as it will be 'debug'
		ProbrCoreLogic()

	default:
		// Parse cli args
		runCmd.Parse(os.Args[1:])

		// Serve plugin
		spProbr := &ServicePack{}
		serveOpts := &plugin.ServeOpts{
			Pack: spProbr,
		}

		plugin.Serve(serveOpts)
	}
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
		defer sdkConfig.GlobalConfig.CleanupTmp()
		// TODO: Additional cleanup may be needed. For instance, any pods created during tests are not being dropped if aborted.
		os.Exit(0)
	}()
}

// ProbrCoreLogic ...
func ProbrCoreLogic() (err error) {
	defer sdkConfig.GlobalConfig.CleanupTmp()
	setupCloseHandler() // Sigterm protection

	config.Vars.Init()
	summary.State = audit.NewSummaryState(ServicePackName)

	connection.Connect()

	store := probeengine.NewProbeStore(ServicePackName, config.Vars.Tags(), &summary.State)
	s, err := store.RunAllProbes(pack.GetProbes())
	if err != nil {
		log.Printf("[ERROR] Error executing tests %v", err)
		return
	}
	log.Printf("[INFO] Overall test completion status: %v", s)
	summary.State.SetProbrStatus()

	summary.State.PrintSummary()
	summary.State.WriteSummary()

	if summary.State.ProbesPassed == 0 && summary.State.ProbesFailed == 0 {
		return utils.ReformatError("No probes ran, or all probes given statements were not met.")
	} else if summary.State.ProbesFailed > 0 {
		return utils.ReformatError("One or more probe scenarios were not successful. View the output logs for more details.")
	}
	return
}

func printVersion(w io.Writer) {

	if config.Vars.Verbose {
		fmt.Fprintf(w, "Service Pack : %s", ServicePackName)
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Version      : %s", getVersion())
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Commit       : %s", GitCommitHash)
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Built at     : %s", BuiltAt)
	} else {
		fmt.Fprintf(w, "Version      : %s", getVersion())
	}
}

func getVersion() string {
	if Prerelease != "" {
		return fmt.Sprintf("%s-%s", Version, Prerelease)
	}
	return Version
}
