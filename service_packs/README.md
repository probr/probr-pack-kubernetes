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

1. Create a folder under the service_packs folder e.g. storage_packs/`storage`
   
1. Create a folder for each probe and insert code and feature files into the folder e.g. service_packs/`storage`/`encryption_in_flight` holds the following files for the encryption_in_flight probe:

      - encryption_in_flight.feature - the 'bdd' feature file storing the control requirements for this probe
      - encryption_in_flight.go - implementation code for the probe
      - encryption_in_flight_test.go - unit test code
  
1. Add code to the go file, in order to integrate the service pack probes with the GoDog handler
   - Define a ***ProbeStruct***:
      - type ProbeStruct struct{} 
      - var Probe ProbeStruct   // allows the probe to be added to the ProbeStore
      - func (p ProbeStruct) Name() string {return "probe name"}
   - Define ProbeInitialize(ctx *godog.TestSuiteContext) - code to initialize the Godog handler prior to the probe run:
      - e.g. initialize probe state
   - Define ScenarioInitialize(ctx *godog.ScenarioContext) - code to initialize the Godog handler prior to executing a scenario:
      - all steps defined in the scenario must be mapped to implementation code; the code will be called by the GoDog handler when executing the associated feature's scenario step
      - ctx.Step(***the scenario clause***, ***mapped function***) e.g. ctx.Step(`^the detective measure is enabled$`, `state.policyOrRuleAssigned`)
  
2. Add the service pack configuration variables to config/types.go
   - Define the service pack type e.g.
      - type ***Storage*** struct {
      - Excluded string `yaml:"Exclude"`
      - Probes   []Probe `yaml:"Probes"`
      - }
   - Add the type to the ServicePacks struct e.g. 
      - type ServicePacks struct {
	  - Kubernetes Kubernetes `yaml:"Kubernetes"`
	  - ***Storage*** Storage    `yaml:"Storage"`

1. Add the service pack and its probes to the init() function in service_packs/service_packs.go
	- packs["storage"] = []probe{
	-	encryption_in_flight.Probe,
	-	encryption_at_rest.Probe,
	-	access_whitelisting.Probe,
	- }

1. Add service pack exclusion logic to the handleProbeExclusions method in internal/config/config.go

1. Add any utilities under internal and import where required e.g. the storage probes use the azure connection utilities package, which is installed in the internal/azureutil folder