# probr

Probr is intended to provide flexible "probing" of your cloud provider and Kubernetes cluster to ensure that the _result_ of your security policies is in place.

Instead of scanning to validate that specific measures have been taken, Probr attempts to perform specific tasks that it should or _shouldn't_ be able to perform.

Probr may be used by **security professionals** to demonstrate the need for specific policies and remidiation, or it may be used by **engineering teams** to display that necessary regulations are being met.

## Getting Started

### Requirements

The following elements are required to get started with Probr:

- A running Kubernetes cluster
- The kubeconfig file for the cluster you wish to probe
- Your cloud provider credentials (if probing the cloud provider)

### CLI Usage

1. Pull the version of the code you wish to run from `github.com/citihub/probr` and package the code using the following command:

`go build -o probr probr/cmd/probr-cli/main.go`

This will make a `probr` or `probr.exe` binary in your current working directory.

2. Set your configuration variables. For more on this, see the [config documentation](https://github.com/citihub/probr/blob/master/internal/config/README.md)

3. Run the probr executable via `./probr [OPTIONS]` on *nix or `probr.exe [OPTIONS]` on Windows.

Additional options, as seen using `./probr --help`:
```
$ ./probr --help
Usage of ./probr:
  -integrationTest
        run integration tests
  -kubeConfig string
        kube config file
  -outputDir string
        output directory
  -outputType string
        output defaults to write in memory, if 'IO' will write to specified output directory (default "INMEM")
  -varsFile string
        path to config file
```

### Dockerized Usage

If you would like to compile the latest code, but don't have the requirements to do so, you may run the provided Dockerfile to build and run the executable.

** Build **

To build and run probr using Docker:
```
cd probr
docker build -t probr .
docker run probr
```

The build step does not need to be repeated each time you run the container.

** Run **

To run with variables:

```
docker run -e <VAR_NAME>=<VALUE> -e <VAR2_NAME>=<VALUE2> probr
```

For example, you may pass in your Azure credentials at this point:

```
docker run -e AZURE_CLIENT_ID=<YOUR_VALUE> -e AZURE_CLIENT_SECRET=<YOUR_VALUE>
```

For additional environment variable options, including env files, view details on the (Official Docker Docs)[https://docs.docker.com/engine/reference/commandline/run/#set-environment-variables--e---env---env-file]

** Output **

The docker container will not be able to modify your filesystem by default. If you would like to get the output from the tests, you may mount a volume.

(To be continued...)