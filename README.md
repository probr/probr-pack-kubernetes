<img src="assets/images/probr.png" width="200">

## Interactive Application Security Testing (IAST) for Cloud
Probr analyzes the complex behaviours and interactions in your cloud resources to enable engineers, developers and operations teams identify and fix security related flaws early and often, to assist in building secure platforms and reducing the number of defects discovered later in the development lifecycle.

Probr has been designed to test aspects of security and compliance that are otherwise challenging to assert using static code inspection or configuration inspection alone, providing a deeper level of confidence in the compliance of your cloud solutions.

### Control Specifications
Probr uses structured natural language to describe the behaviours of an adequately controlled set of cloud resources. These form the basis of control requirements without getting into the nitty gritty of how those controls should be implemented.  This leaves engineering teams the freedom to determine the best course of action to implement those behaviours. The implementation may change frequently, given the rapid feature velocity in the cloud and tooling ecosystem, without needing to update Probr. This differentiates Probr from policy-based tools, which are designed to look for specific features of resource implementation, so need to iterate in-line with changes to the underlying implementation approach.

### How it works
Probr deploys a series of probes to test the behaviours of the cloud resources in your code, returning a machine-readable set of structured results that can be integrated into the broader DevSecOps process for decision making.

## Quickstart Guide

### Requirements

The following elements are required to get started running the Probr Kubernetes service pack:

- A running Kubernetes cluster
- The kubeconfig file for the cluster you wish to probe
- Your cloud provider credentials (if probing the cloud provider)

### Get the executable

- **Option 1** - Download the latest Probr package by clicking the corresponding asset on our [release page](https://github.com/citihub/probr/releases).
- **Option 2** - You may build the edge version of Probr by using `go build -o probr.exe cmd/main.go` from the source code. This may also be necessary if an executable compatible with your system is not available in on the release page.
- **Option 3** - There is an example Dockerfile in [examples/docker](./examples/docker) which will build a Docker image with both Probr and [Cucumber HTML Reporter](https://www.npmjs.com/package/cucumber-html-reporter)

*Note: The usage docs refer to the executable as `probr` or `probr.exe` interchangeably. Use the former for unix/linux systems, and the latter package if you are working in Windows.*

### CLI Usage

1. If you will be using any custom files, move the downloaded executable to the associated working directory. Below are elements you may wish to add to your working directory:

      - **kubeconfig** - Required. Default location: `~/.kube/config`
      - **Probr config** - Not required, no default. Used to specify config options as code.
      - **output directory** - Not required *if* using output type of `INMEM`, which will simply print the scenario results to the terminal. Default directory still needs to be created, but path name can be modified via config. Default location: `./cucumber_output`

1. Set your configuration variables. For more on how to do this, see the config documentation further down on this page.

1. Run the probr executable via `./probr [OPTIONS]`.
    - Additional options can be seen via `./probr --help`
    - Review required variables by using `./probr show-requirements <SERVICE-PACK-NAME; optional>`

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

An example Vars file is available at [./examples/config.yml](./examples/config.yml).
You may have as many vars files as you wish in your codebase, which will enable you to maintain configurations for multiple environments in a single codebase.

The location of the vars file is passed as a CLI option e.g.

```
probr --varsFile=./config-dev.yml
```

### Probr Configuration Variables

These are general configuration variables.

| Variable | Description | CLI Option | Vars File | Env Var | Default |
|---|---|---|---|---|---|
|VarsFile|Config YAML File Path|yes|N/A|N/A|N/A|
|Silent|Disable visual runtime indicator|yes|no|N/A|false|
|NoSummary|Flag to switch off summary output|yes|no|N/A|false|
|WriteDirectory|Path to all output, including audit, cucumber results and other temp files|yes|yes|PROBR_WRITE_DIRECTORY|probr_output|
|Tags|Feature tag inclusions and exclusions|yes|yes|PROBR_TAGS| |
|LogLevel|Set log verbosity level|yes|yes|PROBR_LOG_LEVEL|ERROR|
|OutputType|"IO" will write to file, as is needed for CLI usage. "INMEM" should be used in non-CLI cases, where values should be returned in-memory instead|no|yes|PROBR_OUTPUT_TYPE|IO|
|AuditEnabled|Flag to switch on audit log|no|yes|PROBR_AUDIT_ENABLED|true|
|OverwriteHistoricalAudits|Flag to allow audit overwriting|no|yes|OVERWRITE_AUDITS|true|
|ContainerRegistry|Probe image container registry|no|yes|PROBR_CONTAINER_REGISTRY|docker.io|
|ProbeImage|Probe image name|no|probeImage|PROBR_PROBE_IMAGE|citihub/probr-probe|
|ContainerRequiredDropCapabilities|Container Required Drop Capabilities|no|ContainerRequiredDropCapabilities|PROBR_REQUIRED_DROP_CAPABILITIES|["NET_RAW"]|

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

## Tagging

A variety of tagging options are available to help you specify which probes should be included or excluded at runtime.

In the bullet-point lists below, each relationship level is represented by a `/` in the tag name. Examples are below the list of tags.

The available tags are as follows:

**Service Packs**

These tags describe `service_pack/probe/scenario` in progressive detail.

The first layer of the tag (`@probes`) is only an identifier, and serves no purpose by itself.

- probes
  - kubernetes
    - container_registry_access
      - `PROBR_VERSION`.`SCENARIO_ID`
    - iam
      - `PROBR_VERSION`.`SCENARIO_ID`
    - internet_access
      - `PROBR_VERSION`.`SCENARIO_ID`
    - general
      - `PROBR_VERSION`.`SCENARIO_ID`
    - pod_security_policy
      - `PROBR_VERSION`.`SCENARIO_ID`

_Examples:_

```
@probes/kubernetes  # all k8s probes and scenarios
@probes/kubernetes/iam  # scenarios for the k8s/iam prbe
@probes/kubernetes/pod_security_policy/1.0  # scenario 0 from the v1 Probr release
```

**Categories**

These "category" tags may target probes or scenarios with categorical similarities across multiple service packs. The first layer of the tag (`@category`) is only an identifier, and serves no purpose by itself.

- category
  - pod_security_policy
  - internet_access
  - iam

_Examples:_

```
@category/internet_access # targets internet access related probes from all service packs
```

**Standards**

These "standard" tags target probes and scenarios that validate a specific standard or control. The first layer of the tag (`@standard`) is only an identifier, and serves no purpose by itself.

Standards such as CIS can be drilled down by adding a `.` to drill down in accordance with the control identifiers.

- standard
  - cis
    - gke
      - `see description and examples`
  - citihub
      - `control ID`

_Examples:_

```
@standard/cis  # targets all CIS-compatible probes and scenarios
@standard/cis/gke  # as above, but only targets CIS GKE probes and scenarios
@standard/cis/gke/5  # as above, refined to a specific control
@standard/cis/gke/5.2  # More refined CIS control targeting
@standard/cis/gke/5.2.3  # More refined CIS control targeting
@standard/citihub/CHC2-IAM105  # targets probes related to this Citihub control
```

**Cloud Service Providers**
Some generalized probes target all providers, but others are only useful for specific providers. The first layer of the tag (`@csp`) is only an identifier, and serves no purpose by itself.

- csp
  - any
  - aws
  - azure
  - gke
  - openshift

```
@csp/any  # targets only probes and scenarios that are cloud agnostic
@csp/gke  # targets only GKE-compatible probes and scanarios
```

## Development & Contributing

Please see the [contributing docs](https://github.com/citihub/probr/blob/master/CONTRIBUTING.md) for information on how to develop and contribute to this repository as either a maintainer or open source contributor (the same rules apply for both).
