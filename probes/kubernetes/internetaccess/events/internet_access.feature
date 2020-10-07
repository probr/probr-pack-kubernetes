@kubernetes
@internet_access
@csp.azure
@CCO:CHC2-SVD030
Feature: Egress control of a kubernetes cluster
  As a Security Auditor
  I want to ensure that containers running inside Kubernetes clusters cannot directly access the Internet
  So that Internet traffic can be inspected and controlled

#  Rule: CHC2-SVD030 - protect cloud service network access by limiting access from the appropriate source network only

  Scenario Outline: Test outgoing connectivity of a deployed pod
    Given a Kubernetes cluster is deployed
    And a pod is deployed in the cluster
    When a process inside the pod establishes a direct http(s) connection to "<url>"
    Then access is "<result>"

    Examples:
      | url               | result  |
      | www.google.com    | blocked |
      | www.microsoft.com | blocked |
      | www.ubuntu.com    | blocked |