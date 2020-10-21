package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// ConfigVars contains all possible config vars
type ConfigVars struct {
	// NOTE: Env and Defaults are ONLY available if corresponding logic is added to defaults.go and getters.go
	KubeConfigPath            string `yaml:"kubeConfig"`
	KubeContext               string `yaml:"kubeContext"`
	OutputType                string `yaml:"outputType"`
	CucumberDir               string `yaml:"outputDir"`
	AuditDir                  string `yaml:"auditDir"`
	SummaryEnabled            string `yaml:"summaryEnabled"`
	AuditEnabled              string `yaml:"auditEnabled"`
	OverwriteHistoricalAudits string `yaml:"overwriteHistoricalAudits"`
	Images                    struct {
		Repository string `yaml:"repository"`
		Curl       string `yaml:"curl"`
		BusyBox    string `yaml:"busyBox"`
	} `yaml:"images"`
	Azure struct {
		SubscriptionID  string `yaml:"subscriptionID"`
		ClientID        string `yaml:"clientID"`
		ClientSecret    string `yaml:"clientSecret"`
		TenantID        string `yaml:"tenantID"`
		LocationDefault string `yaml:"locationDefault"`
		Identity        struct {
			DefaultNamespaceAI  string `yaml:"defaultNamespaceAI"`
			DefaultNamespaceAIB string `yaml:"defaultNamespaceAIB"`
		} `yaml:"azureIdentity"`
	} `yaml:"azure"`
	Events             []Event  `yaml:"events"`
	SystemClusterRoles []string `yaml:"systemClusterRoles"`
	Tags               string   `yaml:"tags"`
	TagExclusions      []string // not from yaml
}

type Event struct {
	Name          string  `yaml:"name"`
	Excluded      bool    `yaml:"excluded"`
	Justification string  `yaml:"justification"`
	Probes        []Probe `yaml:"probes"`
}

type Probe struct {
	Name          string `yaml:"name"`
	Excluded      bool   `yaml:"excluded"`
	Justification string `yaml:"justification"`
}

// Vars is a singleton instance of ConfigVars
var Vars ConfigVars

// GetTags parses Tags with TagExclusions
func (ctx *ConfigVars) GetTags() string {
	for _, v := range ctx.Events {
		if v.Excluded {
			ctx.HandleExclusion(v.Name, v.Justification)
		} else {
			ctx.HandleProbeExclusions(&v)
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
	r := "@" + name + ","
	ctx.Tags = strings.Replace(ctx.Tags, r, "", -1)     // Remove exclusion from tags
	ctx.TagExclusions = append(ctx.TagExclusions, name) // Add exclusion to list
}

func (ctx *ConfigVars) HandleProbeExclusions(e *Event) {
	for _, v := range e.Probes {
		if v.Excluded {
			ctx.HandleExclusion(v.Name, v.Justification)
		}
	}
}

func init() {
	//create a defaulted config
	Init("")
}

// Init will override config.Vars with the content retrieved from a filepath
func Init(configPath string) error {
	config, err := NewConfig(configPath)
	if err != nil {
		return err
	}
	Vars = config
	setFromEnvOrDefaults(&Vars) // Set any values not retrieved from file
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
