@probes/kubernetes
@probes/kubernetes/iam
@category/iam
@standard/citihub
@standard/citihub/CHC2-IAM105
Feature: Ensure stringent authentication and authorisation
    As a Security Auditor
    I want to ensure that stringent authentication and authorisation policies are applied to my organisation's Kubernetes clusters
    So that only approve actors have ability to perform sensitive operations in order to prevent malicious attacks on my organization

    @probes/kubernetes/iam/AZ-AAD-AI-1.0 @control_type/preventative @csp/azure
    Scenario Outline: Prevent cross namespace Azure Identities
        Given a Kubernetes cluster exists which we can deploy into
        And an AzureIdentityBinding called "probr-aib" exists in the namespace called "default"
        When I create a simple pod in "<NAMESPACE>" namespace assigned with the "probr-aib" AzureIdentityBinding
        Then the pod is deployed successfully
        But an attempt to obtain an access token from that pod should "<RESULT>"

        Examples:
			| NAMESPACE     | RESULT  |
			| the probr     | Fail    |
			| the default   | Succeed |

    @probes/kubernetes/iam/AZ-AAD-AI-1.1 @control_type/preventative @csp/azure
    Scenario: Prevent cross namespace Azure Identity Bindings
        Given a Kubernetes cluster exists which we can deploy into
        And the namespace called "default" has an AzureIdentity called "probr-probe"
        When I create an AzureIdentityBinding called "probr-aib" in the Probr namespace bound to the "probr-probe" AzureIdentity
        And I deploy a pod assigned with the "probr-aib" AzureIdentityBinding into the Probr namespace
        Then the pod is deployed successfully
        But an attempt to obtain an access token from that pod should fail

    @probes/kubernetes/iam/AZ-AAD-AI-1.2 @control_type/preventative @csp/azure
    Scenario: Prevent access to AKS credentials via Azure Identity Components
        Given a Kubernetes cluster exists which we can deploy into
        And the cluster has managed identity components deployed
        When I execute the command "cat /etc/kubernetes/azure.json" against the MIC pod
        Then Kubernetes should prevent me from running the command
