# Service Pack Engineering Notes

This directory contains subdirectories for each _service pack_ within Probr.
A service pack contains multiple probes that are designed to verify that the
service is in compliance with various controls.

## Feature Files

Each probe describes a set of related controls, which are together defined in
a single Cucumber BDD `.feature` file.

For more information on the _gherkin_ language and how a BDD "feature file" is
written, please review the
[Cucumber documentation](https://cucumber.io/docs/gherkin/reference/)

### Feature Scenarios and Steps

Each control that is to be validated takes the form of a Cucumber _scenario_.

A _scenario_ is executed as a sequence of _steps_- each of which is present
in the feature file as an english language statement, beginning with a reserved
word (Given, And, When, Then).

The functions that define each step can be found in the `.go` file
in the same directory as the `.feature` file. (For example,
"container registry access" steps are defined in `container_registry_access.go`)

Within the probe's go file, functions must be defined to execute each of the
steps specified in the corresponding feature file.

For example, the `container_registry_access.feature` file may specify a step
`When I attempt to push to the container registry using the cluster identity`.
Within the `container_registry_access.go` file, a ScenarioState method
`iAttemptToPushToTheContainerRegistryUsingTheClusterIdentity` is defined and
mapped to the string 'I attempt to push to the container registry using the
cluster identity'.

Mapping a step to a golang function is registered in the ScenarioContext of
the 'godog test handler', so that when the test handler runs a probe for a
feature file, it executes the appropriate go code for each defined step.
See any `ScenarioInitialize` function for an example of the registration
of step functions.

## Adding a new service pack

1. Create a directory for your new pack under the service_packs folder (e.g. `storage_packs/storage`)

1. Create a folder for each probe and insert code and feature files into the folder e.g. `service_packs/storage/encryption_in_flight` holds the following files for the encryption_in_flight probe:

      - encryption_in_flight.feature - the 'bdd' feature file storing the control requirements for this probe
      - encryption_in_flight.go - implementation code for the probe
      - encryption_in_flight_test.go - unit test code

1. Create a file at the top level of the pack that will import all of your newly created probes. The package should be named after the service pack, and must contain a function named `GetProbes`.

   ```go
      // storage_packs/storage/storage.go
      package storage
      import (
         "storage_packs/storage/access_whitelisting"
         "storage_packs/storage/encryption_at_rest"
         "storage_packs/storage/encryption_in_flight"
      )
      func GetProbes() []coreengine.Probe {
         ...
         return []coreengine.Probe{
            access_whitelisting.Probe,
            encryption_at_rest.Probe,
            encryption_in_flight.Probe,
         }
      }
   ```

   - A function named `init` must also be added to the `packname.go` file, containing logic to include all feature files in pkger bundle.
   Logic has been added to the Makefile to aid in packaging any files that are listed in this way (`make binary`).
   For more information about pkger, please review the [official pkger docs](https://github.com/markbates/pkger)
   ```go
      func init() {
         // This line will ensure that all static files are bundled into pkged.go file when using pkger cli tool
         // See: https://github.com/markbates/pkger
         pkger.Include("/service_packs/kubernetes/general/general.feature")
         pkger.Include("/service_packs/kubernetes/podsecurity/podsecurity.feature")
         pkger.Include("/service_packs/kubernetes/iam/iam.feature")
      }
   ```

1. Add code to the go file, in order to integrate the service pack probes with the GoDog handler
   - Define a ***ProbeStruct*** - follow an existing example to ensure proper implementation

      ```go
         // storage_packs/my_pack/my_probe/my_probe.go
         type ProbeStruct struct{} 
         var Probe ProbeStruct   // allows the probe to be added to the ProbeStore
         func (p ProbeStruct) Name() string {return "my-probe-name"} // Used in storage_packs/storage_packs.go
         func (p ProbeStruct) Path() string { return coreengine.GetFeaturePath("service_packs", "kubernetes", p.Name()) } // Allows for custom pack file structure
         func ProbeInitialize(ctx *godog.TestSuiteContext) {} // required by the Godog handler
         func ScenarioInitialize(ctx *godog.ScenarioContext) {} // defines each step, required by the Godog handler
      ```

1. Add the service pack configuration variables to `internal/config/types.go`, allowing users to specify the inclusion of your service pack.
   - Define the service pack type. Example:

      ```go
        // internal/config/types.go
        type Storage struct {
          Excluded string `yaml:"Exclude"`
          Probes   []Probe `yaml:"Probes"`
        }
      ```

   - Add the type to the ServicePacks struct. Example: 

      ```go
         // internal/config/types.go
         type ServicePacks struct {
           Kubernetes Kubernetes `yaml:"Kubernetes"`
           Storage    `yaml:"Storage"`
         }
      ```

1. Add an `IsExcluded` function to `internal/config` with at least one required variable. This must be used in `GetProbes()` or (1) `probr run <SERVICE-PACK>` will not function properly, and (2) the service pack will be run every time by default (which is undesired behavior). 

   ```go
      // internal/config/config.go
      // Log and return exclusion configuration
      func (k Kubernetes) IsExcluded() bool {
         return validatePackRequirements("Kubernetes", k)
      }
   ```

   ```go
      // internal/config/requirements.go
      var Requirements = map[string][]string{
         "Storage":    []string{"Provider"},
         "Kubernetes": []string{"AuthorisedContainerRegistry", "UnauthorisedContainerRegistry"},
      }

   ```

1. Add the service pack and its probes to `service_packs/service_packs.go`

   ```go
      // service_packs/service_packs.go
      import (
         ...   
         "github.com/citihub/probr/service_packs/storage"
      )
      ...
      func packs() (packs map[string][]coreengine.Probe) {
         ...
         packs["storage"] = storage.GetProbes()
         return
      }
      ```