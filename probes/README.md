# Probes Engineering Notes

This folder includes feature files and functionality associated with probes. A probe has functions defined to allow it to test the required control behaviour of a 'system feature', which is defined in a single BDD feature file.
Probes are grouped into 'categories' and placed within a category folder. For example, the testing functionality for the 'container_registry_probe' is defined in the container_registry_probe.go file within the kubernetes folder, because it is part of the kubernetes category. The associated container_registry_access.feature file is located within the kubernetes/probe_specifications folder.

## Feature Files

A feature file defines behaviour using the gherkin language
For more information on how a _feature file_ is written, please review the [Cucumber documentation](https://cucumber.io/docs/gherkin/reference/)

### Feature Scenarios and Steps

Feature behaviour is specified in the form of a set of _scenarios_, each of which are tested by the associated probe.
A scenario is executed as a sequence of _steps_, each of which is described in the feature file as an english language statement, begining with a reserved word (Given, And, When, Then).
The status of a scenario execution is managed by a ScenarioState struct, which is defined within the k8s_probes.go file

## Probe functionality definition

Within the probe's go file, functions must be defined to execute each of the steps specified in the corresponding feature file.
For example, the container_registry_access.feature file specifies a step 'When I attempt to push to the container registry using the cluster identity'. Within the container_registry_access.go file, a ScenarioState method _iAttemptToPushToTheContainerRegistryUsingTheClusterIdentity_ is defined and mapped to the string 'I attempt to push to the container registry using the cluster identity'.
The mappings from step specification to go function must be registered in the ScenarioContext of the 'godog test handler', so that when the test handler runs a probe for a feature file, it executes the appropriate go code for each defined step. See the craScenarioInitialize function for an example of the registration of step functions.

For information regarding the godog test handler, refer to the README.MD in the internal/coreengine folder.

Note that some steps may be used across multiple probes within the kubernetes category, and are found in `probes/kubernetes/k8s_probes.go.
