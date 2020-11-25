# Probr

## Your Zero Trust Compliance Toolbox

Probr is intended to provide flexible "probing" of your cloud provider and Kubernetes cluster to ensure that the result of your security controls have properly taken effect.

Instead of reading configurations or scanning to validate that specific policies are in place, Probr attempts to perform specific tasks that should or shouldn't be able to occur from specific roles.

Probr may be used by **security professionals** to audit or demonstrate the need for specific policies and remediation, or Probr may be used by **engineering teams** to display that necessary regulations are being met.

## Quickstart Guide

### Requirements

The following elements are required to get started with Probr:

- A running Kubernetes cluster
- The kubeconfig file for the cluster you wish to probe
- Your cloud provider credentials (if probing the cloud provider)

### Get the executable

- **Option 1** - Download the latest Probr package by clicking the corresponding asset on our [release page](https://github.com/citihub/probr/releases).
- **Option 2** - You may build the edge version of Probr by using `go build -o probr.exe cmd/main.go` from the source code. This may also be necessary if an executable compatible with your system is not available in on the release page.

*Note: The usage docs refer to the executable as `probr` or `probr.exe` interchangably. Use the former for unix/linux systems, and the latter package if you are working in Windows.*

### CLI Usage

1. If you will be using any custom files, move the downloaded executable to the associated working directory. Below are elements you may wish to add to your working directory:

      - **kubeconfig** - Required. Default location: `~/.kube/config`
      - **Probr config** - Not required, no default. Used to specify config options as code.
      - **output directory** - Not required *if* using output type of `INMEM`, which will simply print the scenario results to the terminal. Default directory still needs to be created, but path name can be modified via config. Default location: `./cucumber_output`

1. Set your configuration variables. For more on how to do this, see the config documentation further down on this page.

1. Run the probr executable via `./probr [OPTIONS]`. Additional options can be seen via `./probr --help`

*Note: Feature files are not included in the binary. In this present state, Probr must be executed from the top level directory of the source code.*

## Configuration

### How the Config Works

Configuration variables can be populated in one of four ways, with the value being taken from the highest priority entry.

1. Default values; found in `internal/config/defaults.go` (lowest priority)
1. OS environment variables; set locally prior to probr execution (mid priority)
1. Vars file; yaml (highest non-CLI priority)
1. CLI flags; see `./probr --help` for available flags (highest priority)

_Note: See `internal/config/README.md` for engineering notes regarding configuration._

### Environment Variables

If you would like to handle logic differently per environment, env vars may be useful. An example of how to set an env var is as follows:

`export KUBE_CONFIG=./path/to/config`

### Vars File

An example Vars file is available at `probr/examples/config.yml`
You may have as many vars files as you wish in your codebase, which will enable you to maintain configurations for multiple environments in a single codebase.

The location of the vars file is passed as a CLI option e.g.

```
probr --varsFile=./config-dev.yml
```

**IMPORTANT:** Remember to encrypt your config file if it contains secrets.

### Probr Configuration Variables

These are general configuration variables.

| Variable | Description | CLI Option | Vars File | Env Var | Default |
|---|---|---|---|---|---|
|VarsFile|Config YAML File Path|yes|N/A|N/A|N/A|
|Silent|Disable visual runtime indicator|yes|no|N/A|false|
|NoSummary|Flag to switch off summary output|yes|no|N/A|false|
|CucumberDir|Path to cucumber output dir if applicable|yes|yes|PROBR_CUCUMBER_DIR|cucumber_output|
|Tags|Feature tag inclusions and exclusions|yes|yes|PROBR_TAGS| |
|LogLevel|Set log verbosity level|yes|yes|PROBR_LOG_LEVEL|ERROR|
|OutputType|"IO" will write to file, as is needed for CLI usage. "INMEM" should be used in non-CLI cases, where values should be returned in-memory instead|no|yes|PROBR_OUTPUT_TYPE|IO|
|AuditEnabled|Flag to switch on audit log|no|yes|PROBR_AUDIT_ENABLED|true|
|AuditDir|Path to audit dir|no|yes|PROBR_AUDIT_DIR|audit_output|
|OverwriteHistoricalAudits|Flag to allow audit overwriting|no|yes|OVERWRITE_AUDITS|true|
|ContainerRegistry|Probe image container regsitry|no|yes|PROBR_CONTAINER_REGISTRY|docker.io|
|ProbeImage|Probe image name|no|probeImage|PROBR_PROBE_IMAGE|citihub/probr-probe|

### Service Pack Configuration Variables

Variables that are specific to a service pack. May be configured in the Vars file via embedded tags under ServicePacks.

| Variable | Description | CLI Flag | VarsFile | Env Var | Default |
|---|---|---|---|---|---|
|Kubernetes.KubeConfig|Path to kubernetes config|yes|yes|KUBE_CONFIG|~/.kube/config|
|Kubernetes.KubeContext|Kubernetes context|no|yes|KUBE_CONTEXT| |
|Kubernetes.SystemClusterRoles|Cluster names|no|yes|N/A|{"system:", "aks", "cluster-admin", "policy-agent"}|

### Cloud Provider Configuration Variables

Variables that are specific to a cloud service provider and can be configured in the Vars file via embedded tags under CloudProviders.

| Variable | Description | CLI Flag | VarsFile | Env Var | Default |
|---|---|---|---|---|---|
|Azure.SubscriptionID|Azure subscription|no|yes|AZURE_SUBSCRIPTION_ID| |
|Azure.ClientId|Azure client id|no|yes|AZURE_CLIENT_ID| |
|Azure.ClientSecret|Azure client secret|no|yes|AZURE_CLIENT_SECRET| |
|Azure.TenantID|Azure tenant id|no|yes|AZURE_TENANT_ID| |
|Azure.LocationDefault|Azure location default|no|yes|AZURE_LOCATION_DEFAULT| |
|Azure.AzureIdentity.DefaultNamespaceAI|Azure namespace|no|yes|DEFAULT_NS_AZURE_IDENTITY|probr-defaultns-ai|
|Azure.AzureIdentity.DefaultNamespaceAIB|Azure namespace|no|yes|DEFAULT_NS_AZURE_IDENTITY_BINDING|probr-defaultns-aib|

### Probes Configuration Variables

Variables used to configure which probes and their associated controls & scenarios are to be run. Probes can be excluded via the ProbeExclusions config file variable. Specific controls/scenarios can be excluded via the TagExclusions config file variable.
Tags may also be specified via a command line parameter to control which tagged probes/controls/scenarios are to be run. This takes precedence over the TagExclusions variable

| Variable | Description | CLI Flag | VarsFile | Env Var | Default |
|---|---|---|---|---|---|
|Tags|Specify tags for probes/controls/scenarios to be included or excluded|yes|no| | |
|ProbeExclusions|Specify names of probes to be excluded and provide justification|no|yes| | |
|TagExclusions|Specify the tags for controls/scenarios to be excluded|no|yes| | |

## Development & Contributing

Please see the [contributing docs](https://github.com/citihub/probr/blob/master/CONTRIBUTING.md) for information on how to develop and contribute to this repository as either a maintainer or open source contributor (the same rules apply for both).
