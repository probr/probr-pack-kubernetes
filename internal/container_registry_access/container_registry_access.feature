@k-cra
@probes/kubernetes/container_registry_access
Feature: Protect image container registries
    As a Security Auditor
    I want to ensure that containers image registries are secured in my organisation's Kubernetes clusters
    So that only approved software can be run in our cluster in order to prevent malicious attacks on my organization

    Security Standard References:
        CHC2-APPDEV135 - Ensure software release and deployment is managed through a formal, controlled process

    Background:
        Given a Kubernetes cluster exists which we can deploy into

    @k-cra-003
    Scenario: Ensure deployment from an unauthorised container registry is denied
        Then pod creation "succeeds" with container image from "authorized" registry
        And pod creation "is denied" with container image from "unauthorized" registry
