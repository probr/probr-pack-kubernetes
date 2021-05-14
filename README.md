# Probr Kubernetes Service Pack

The [Probr](https://github.com/probr/probr) Kubernetes Service pack provides a variety of provider-agnostic compliance checks.

Get the latest stable version [here](https://github.com/probr/probr-pack-kubernetes/releases/latest).

## To Build from Source

The following will build a binary named "kubernetes":

```sh
git clone https://github.com/probr/probr-pack-kubernetes.git
cd probr-pack-kubernetes
make binary
```

Move the `kubernetes` binary into your probr service pack location (default is `${HOME}/probr/binaries`)

## Pre-Requisites

You will need:

1. [Probr Core](https://github.com/probr/probr) to execute this service pack.
1. A Kubernetes Cluster
1. An active kubeconfig against the cluster, that can deploy into the probe namespace (see config below. Default is probr-general-test-ns)

## Configuration

### Minimum configuration

The minimum required additions to your Probr runtime configuration is as follows:

```yaml
Run:
  - "kubernetes"
ServicePacks:
  Kubernetes:
    AuthorisedContainerImage: "yourprivateregistry.io/citihub/probr-probe"
```

### Full configuration

If you don't want to use the defaults you can add the following to your Probr config.yml:

```yaml
Run:
  - "aks"
ServicePacks:
  Kubernetes:
    KubeConfig: "location of your kubeconfig if not the default"
    KubeContext: "specific kubecontext if not the current context"
    AuthorisedContainerImage: "yourprivateregistry.io/citihub/probr-probe"
    ProbeNamespace: "namespace Probr deploys into. Defaults to 'probr-general-test-ns'"
CloudProviders:
  Azure:
    TenantID: "UUID of your tenant"
    SubscriptionID: "UUID of your subscription"
    ClientID: "Client ID UUID of your service principle"
    ClientSecret: "Recommend leaving this blank and using envvar"
```

## Running the Service Pack

If all of the instructions above have been followed, then you should be able to run `./probr` and the service pack will run.
