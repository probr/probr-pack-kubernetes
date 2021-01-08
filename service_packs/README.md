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
A _scenario_ is executed as a sequence of _steps_- each of which is described
in the feature file as an english language statement, beginning with a reserved
word (Given, And, When, Then). The functions that define each step can be found
in the `.go` file that has the name of the associated probe. (For example,
"container registry access" steps are defined in `container_registry_access.go`)

Within the probe's go file, functions must be defined to execute each of the
steps specified in the corresponding feature file.

For example, the `container_registry_access.feature` file specifies a step
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

1. Create a subdirectory that will import all of your newly created probes. Directory must be named "pack" (e.g. `storage_packs/storage/pack`) and must contain a function named `GetProbes`.

   ```go
      // storage_packs/storage/pack/pack.go
      func GetProbes() []coreengine.Probe {
         ...
         return []coreengine.Probe{
            access_whitelisting.Probe,
            encryption_at_rest.Probe,
            encryption_in_flight.Probe,
         }
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

1. Optional: Add an `IsExcluded` function to `internal/config` with at least one required variable. Without this, you will need some other logic in `GetProbes()` to avoid service pack being run by default (which is undesired behavior). See the storage service pack `pack.go` for an example of how to conditionally include probes without using `IsExcluded`.

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
         storage_pack "github.com/citihub/probr/service_packs/storage/pack"
      )
      ...
      func packs() (packs map[string][]coreengine.Probe) {
         ...
         packs["storage"] = storage_pack.GetProbes()
         return
      }
      ```
