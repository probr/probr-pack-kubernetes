@kubernetes
@CIS-6.1
Feature: Protect image container registries
  As a Security Auditor
  I want to ensure that containers image registries are secured in my organisation's Kubernetes clusters
  So that only approved software can be run in our cluster in order to prevent malicious attacks on my organization

  #Rule: CHC2-APPDEV135 - Ensure software release and deployment is managed through a formal, controlled process

  @preventative @CIS-6.1.4
  Scenario Outline: Ensure only authorised container registries are allowed
    Given a Kubernetes cluster is deployed
    When a user attempts to deploy a container from "<auth>" registry "<registry>"
    Then the deployment attempt is "<result>"

    Examples:
      | auth          | registry          | result  |
      | unauthorised  | docker.io         | denied  |
      | unauthorised  | gcr.io            | denied  |
      | authorised    | mcr.microsoft.com | allowed |
      | authorised    | allowed-registry  | allowed |