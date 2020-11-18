# Probes Engineering Notes

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
in the feature file as an english language statement, begining with a reserved
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
