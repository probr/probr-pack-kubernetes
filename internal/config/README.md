# Config

Internal code to manage config, including Cloud Driver parameters and Test Packs

## Log Filter

Probr Log Levels:
- **ERROR** - Behavior that is a result of a definite misconfiguration or code failure
- **WARN** - Behavior that is likely due to a misconfiguration, but is not fatal
- **NOTICE** - (1) User config information to prevent confusion, or (2) behavior that could result from a misconfiguration but also may be intentional
- **INFO** - Non-verbose information that doesn't fit the above criteria
- **DEBUG** - Any potentially helpful information that doesn't fit the above criteria

Multi-line logs should be formatted prior to `log.Printf(...)`. By using this command multiple times, each line will get a separate timestamp and will appear to be separate entries.

For example, `Results: ` could be read as if an empty string was being output.

However, by misusing `log.Printf` we may cause a similar appearance:

```
log.Printf("[NOTICE] Results:")
log.Printf("[NOTICE] %s", myVar)
// Prints:
// 2020/09/28 11:18:01 [NOTICE] Results:
// 2020/09/28 11:18:01 [NOTICE] {"some": "information"}
```


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
