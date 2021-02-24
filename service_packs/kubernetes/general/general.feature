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

    @k-gen-004
    Scenario Outline: Test outgoing connectivity of a deployed pod
    Ensure that containers running inside Kubernetes clusters cannot directly access the Internet
    So that Internet traffic can be inspected and controlled

        Given a Kubernetes cluster is deployed
        When a pod is deployed in the cluster
        And a process inside the pod establishes a direct http(s) connection to "<url>"
        Then access is "<result>"

        Examples:
            | url               | result  |
            | www.google.com    | blocked |
            | www.microsoft.com | blocked |
            | www.ubuntu.com    | blocked |