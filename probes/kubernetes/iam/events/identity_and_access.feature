@kubernetes
@iam_control
@CIS-6.8
@AZ-AAD-AI
Feature: Ensure stringent authentication and authorisation
  As a Security Auditor
  I want to ensure that stringent authentication and authorisation policies are applied to my organisation's Kubernetes clusters
  So that only approve actors have ability to perform sensitive operations in order to prevent malicious attacks on my organization

  #There will be CIS control here, for now, straight into Azure AAD Managed Identity ...

  @preventative @AZ-AAD-AI-1.0
  Scenario Outline: Prevent cross namespace Azure Identities
    Given a Kubernetes cluster exists which we can deploy into
    And the default namespace has an AzureIdentityBinding
    When I create a simple pod in "<namespace>" namespace assigned with that AzureIdentityBinding
    Then the pod is deployed successfully
    But an attempt to obtain an access token from that pod should "<RESULT>"

    Examples: 
			| namespace     | RESULT  |
			| a non-default | Fail    |
			| the default   | Succeed |

  @preventative @AZ-AAD-AI-1.1
  Scenario: Prevent cross namespace Azure Identity Bindings
    Given a Kubernetes cluster exists which we can deploy into
    And the default namespace has an AzureIdentity
    When I create an AzureIdentityBinding called "probr-aib" in a non-default namespace
    And I deploy a pod assigned with the "probr-aib" AzureIdentityBinding into the same namespace as the "probr-aib" AzureIdentityBinding 
    Then the pod is deployed successfully
    But an attempt to obtain an access token from that pod should fail

  @preventative @AZ-AAD-AI-1.2
  Scenario: Prevent access to AKS credentials via Azure Identity Components
    Given a Kubernetes cluster exists which we can deploy into
    And the cluster has managed identity components deployed
    When I execute the command "cat /etc/kubernetes/azure.json" against the MIC pod
    Then Kubernetes should prevent me from running the command
