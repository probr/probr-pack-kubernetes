# Probr
### Your Zero Trust Compliance Toolbox

Probr is intended to provide flexible "probing" of your cloud provider and Kubernetes cluster to ensure that the _result_ of your security controls have properly taken effect.

Instead of reading configurations or scanning to validate that specific policies are in place, Probr attempts to perform specific tasks that should or shouldn't be able to occur from specific roles.

Probr may be used by **security professionals** to audit or demonstrate the need for specific policies and remidiation, or Probr may be used by **engineering teams** to display that necessary regulations are being met.

## Quickstart Guide

### Requirements

The following elements are required to get started with Probr:

- A running Kubernetes cluster
- The kubeconfig file for the cluster you wish to probe
- Your cloud provider credentials (if probing the cloud provider)

### Get the executable

- **Option 1** - Download the latest Probr package by clicking the corresponding asset on our [release page](https://github.com/citihub/probr/releases).
- **Option 2** - You may build the edge version of Probr by using `go build cmd/probr-cli/*.go` from the source code. This may also be necessary if an executable compatible with your system is not available in on the release page.

*Note: The usage docs refer to the executable as `probr` but on the release page it will have the version number in its name. You can use that name for execution, or simply change the package's name after you download it.*

### CLI Usage

1. If you will be using any custom files, move the downloaded executable to the associated working directory. Below are elements you may wish to add to your working directory:

      - **kubeconfig** - Required. Default location: `~/. kube/config`
      - **Probr config** - Not required, no default. Used to specify config options as code.
      - **output directory** - Not required *if* using output type of `INMEM`, which will simply print the probe results to the terminal. Default directory still needs to be created, but path name can be modified via config. Default location: `./testoutput`

1. Set your configuration variables. For more on how to do this, see the [config documentation](https://github.com/citihub/probr/blob/master/internal/config/README.md)

1. Run the probr executable via `./probr [OPTIONS]`. Additional options can be seen via `./probr --help`
