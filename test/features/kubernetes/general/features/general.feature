@kubernetes
@CIS-6.10
Feature: General Cluster Security Configurations
  As a Security Auditor
  I want to ensure that Kubernetes clusters have general security configurations in place
  So that no general cluster vulnerabilities can be exploited 

  Scenario: Ensure Kubernetes Web UI is disabled
    Given a Kubernetes cluster is deployed
    And the Kubernetes Web UI is disabled    
    Then I should not be able to access the Kubernetes Web UI    