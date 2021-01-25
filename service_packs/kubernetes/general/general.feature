@k-gen
@probes/kubernetes/general
Feature: General Cluster Security Configurations
    As a Security Auditor
    I want to ensure that Kubernetes clusters have general security configurations in place
    So that no general cluster vulnerabilities can be exploited

    @k-gen-001
    Scenario Outline: Minimise wildcards in Roles and Cluster Roles
        Given a Kubernetes cluster is deployed
        When I inspect the "<rolelevel>" that are configured
        Then I should only find wildcards in known and authorised configurations

        Examples:
            | rolelevel     |
            | Roles         |
            | Cluster Roles |

    @k-gen-002
    Scenario: Ensure Security Contexts are enforced
        Given a Kubernetes cluster is deployed
        When I attempt to create a deployment which does not have a Security Context
        Then the deployment is rejected

    @k-gen-003
    Scenario: Ensure Kubernetes Web UI is disabled
        Given a Kubernetes cluster is deployed
        And the Kubernetes Web UI is disabled
        Then I should not be able to access the Kubernetes Web UI
