@kubernetes
@csp.azure
@CCO:CHC2-APPDEV135
Feature: Protect software deployment using authorised container registry
  As a Security Auditor
  I want to ensure that only containers from approved container registries can be run in my organisation's Kubernetes clusters
  So that only approved software can be run in our cluster

  #Rule: CHC2-APPDEV135 - Ensure software release and deployment is managed through a formal, controlled process

  Scenario Outline: Test only authorised container registry is allowed
    Given a Kubernetes cluster is deployed
    When a user attempts to deploy a container from "<registry>"
    Then the deployment attempt is "<result>"

    Examples:
      | registry          | result  |
      | docker.io         | denied  |
      | gcr.io            | denied  |
      | mcr.microsoft.com | allowed |
      | allowed-registry  | allowed |