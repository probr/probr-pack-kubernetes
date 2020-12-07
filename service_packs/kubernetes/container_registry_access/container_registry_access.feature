@probes/kubernetes
@probes/kubernetes/container_registry_access
@standard/cis/gke/6
@standard/cis/gke/6.1
@standard/citihub/CHC2-APPDEV135
@standard/citihub/CHC2-ITS120
@csp/any
Feature: Protect image container registries
    As a Security Auditor
    I want to ensure that containers image registries are secured in my organisation's Kubernetes clusters
    So that only approved software can be run in our cluster in order to prevent malicious attacks on my organization

    #Rule: CHC2-APPDEV135 - Ensure software release and deployment is managed through a formal, controlled process

    @probes/kubernetes/container_registry_access/1.0 @control_type/preventative @standard/cis/gke/6.1.3
    Scenario: Ensure container image registries are read-only
        Given a Kubernetes cluster is deployed
        And I am authorised to pull from a container registry
        When I attempt to push to the container registry using the cluster identity
        Then the push request is rejected due to authorization

    @probes/kubernetes/container_registry_access/1.1 @control_type/preventative @standard/cis/gke/6.1.4
    Scenario: Ensure deployment from an authorised container registry is allowed
        Given a Kubernetes cluster is deployed
        When a user attempts to deploy a container from an authorised registry
        Then the deployment attempt is allowed

    @probes/kubernetes/container_registry_access/1.2 @control_type/preventative @standard/cis/gke/6.1.5
    Scenario: Ensure deployment from an unauthorised container registry is denied
        Given a Kubernetes cluster is deployed
        When a user attempts to deploy a container from an unauthorised registry
        Then the deployment attempt is denied
