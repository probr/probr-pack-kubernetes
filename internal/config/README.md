# Config

Internal code to manage config, including Cloud Driver parameters and Test Packs

## Available Config Vars

The struct `config.Vars` uses the `Config` struct presented below.

```
type Config struct {
	KubeConfigPath string `yaml:"kubeConfig"`
	OutputType     string `yaml:"outputType"`
	Images         struct {
		Repository string `yaml:"repository"`
		Curl       string `yaml:"curl"`
		BusyBox    string `yaml:"busyBox"`
	} `yaml:"Images"`
	Azure struct {
		SubscriptionID  string `yaml:"subscriptionID"`
		ClientID        string `yaml:"clientID"`
		ClientSecret    string `yaml:"clientSecret"`
		TenantID        string `yaml:"tenantID"`
		LocationDefault string `yaml:"locationDefault"`
	} `yaml:"azure"`
}
```

## How the Config Works

These values are populated in one of three ways, with the value being taken from the highest priority entry.

1. Default values; found in `config/defaults.go` (lowest priority)
2. OS environment variables; set locally prior to probr execution (mid priority)
3. Vars file; yaml (highest priority)

### Default Values

Most variables pertain to environment specific elements such as image repository information or cloud provider credentials. As such, most defaults are empty values.

### Environment Variables

If you would like to handle logic differently per environment, env vars may be useful. An example of how to set an env var is as follows:

`export KUBE_CONFIG=./path/to/config`

### Vars File

You may have as many vars files as you wish in your codebase, which will enable you to maintain configurations for multiple environments in a single codebase.

An example of how to a vars file is as follows:

```
probr --varsFile=./config-dev.yml
```

**IMPORTANT:** Remember to encrypt your config file if it contains secrets.
