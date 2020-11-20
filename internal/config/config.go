package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/briandowns/spinner"
	"gopkg.in/yaml.v2"
)

// ConfigVars contains all possible config vars
type ConfigVars struct {
	// NOTE: Env and Defaults are ONLY available if corresponding logic is added to defaults.go and getters.go
	ServicePacks                  servicePacks   `yaml:"ServicePacks"`
	CloudProviders                cloudProviders `yaml:"CloudProviders"`
	OutputType                    string         `yaml:"OutputType"`
	CucumberDir                   string         `yaml:"CucumberDir"`
	AuditDir                      string         `yaml:"AuditDir"`
	AuditEnabled                  string         `yaml:"AuditEnabled"`
	LogLevel                      string         `yaml:"LogLevel"`
	OverwriteHistoricalAudits     string         `yaml:"OverwriteHistoricalAudits"`
	AuthorisedContainerRegistry   string         `yaml:"AuthorisedContainerRegistry"`
	UnauthorisedContainerRegistry string         `yaml:"UnauthorisedContainerRegistry"`
	ProbeImage                    string         `yaml:"ProbeImage"`
	Probes                        []Probe        `yaml:"Probes"`
	Tags                          string         `yaml:"Tags"`
	NoSummary                     bool           // set by flags only
	Silent                        bool           // set by flags only
	TagExclusions                 []string       // set programatically
}

type servicePacks struct {
	Kubernetes kubernetes `yaml:"Kubernetes"`
}

type cloudProviders struct {
	Azure azure `yaml:"Azure"`
}

type kubernetes struct {
	KubeConfigPath     string   `yaml:"KubeConfig"`
	KubeContext        string   `yaml:"KubeContext"`
	SystemClusterRoles []string `yaml:"SystemClusterRoles"`
}

type azure struct {
	SubscriptionID  string `yaml:"SubscriptionID"`
	ClientID        string `yaml:"ClientID"`
	ClientSecret    string `yaml:"ClientSecret"`
	TenantID        string `yaml:"TenantID"`
	LocationDefault string `yaml:"LocationDefault"`
	Identity        struct {
		DefaultNamespaceAI  string `yaml:"DefaultNamespaceAI"`
		DefaultNamespaceAIB string `yaml:"DefaultNamespaceAIB"`
	}
}

type Probe struct {
	Name          string `yaml:"Name"`
	Excluded      bool   `yaml:"Excluded"`
	Justification string `yaml:"Justification"`
}

// Vars is a singleton instance of ConfigVars
var Vars ConfigVars
var Spinner *spinner.Spinner

// GetTags parses Tags with TagExclusions
func (ctx *ConfigVars) GetTags() string {
	if ctx.Tags == "" {
		if len(ctx.Probes) > 0 {
			log.Printf("[WARN] Exclusions are being ignored due to Tags being set.")
		}
		for _, v := range ctx.Probes {
			if v.Excluded {
				ctx.HandleExclusion(v.Name, v.Justification)
			}
		}
	}
	return ctx.Tags
}

func (ctx *ConfigVars) HandleExclusion(name, justification string) {
	if name == "" {
		return
	}
	if justification == "" {
		log.Fatalf("[ERROR] A justification must be provided for the tag exclusion '%s'", name)
	}
	ctx.TagExclusions = append(ctx.TagExclusions, name) // Add exclusion to list
}

// Init will override config.Vars with the content retrieved from a filepath
func Init(configPath string) error {
	config, err := NewConfig(configPath)

	if err != nil {
		return err
	}
	Vars = config
	setFromEnvOrDefaults(&Vars) // Set any values not retrieved from file

	setLogFilter(Vars.LogLevel, os.Stderr) // Set the minimum log level obtained from Vars

	return nil
}

// NewConfig overrides the current config.Vars values
func NewConfig(c string) (ConfigVars, error) {
	// Create config structure
	config := ConfigVars{}
	if c == "" {
		return config, nil // No file path provided, return empty config
	}
	err := ValidateConfigPath(c)
	if err != nil {
		return config, err
	}
	// Open config file
	file, err := os.Open(c)
	if err != nil {
		return config, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return config, err
	}
	return config, nil
}

// ValidateConfigPath simply ensures the file exists
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func LogConfigState() {
	s, _ := json.MarshalIndent(Vars, "", "  ")
	log.Printf("[NOTICE] Config State: %s", s)
}
