    @k-iam-001
    Scenario Outline: Prevent cross namespace Azure Identities

        Security Standard References:
            - AZ-AAD-AI-1.0

        Given an "AzureIdentityBinding" called "probr-aib" exists in the namespace called "default"
        Then I succeed to create a simple pod in "<NAMESPACE>" namespace assigned with the "probr-aib" AzureIdentityBinding
        But an attempt to obtain an access token from that pod should "<RESULT>"

        Examples:
			| NAMESPACE     | RESULT  |
			| the probr     | Fail    |
			| the default   | Succeed |

    @k-iam-002
    Scenario: Prevent cross namespace Azure Identity Bindings

        Security Standard References:
            - AZ-AAD-AI-1.1

        Given an "AzureIdentity" called "probr-probe" exists in the namespace called "default"
        When I create an AzureIdentityBinding called "probr-aib" in the Probr namespace bound to the "probr-probe" AzureIdentity
        Then I succeed to create a simple pod in "the probr" namespace assigned with the "probr-aib" AzureIdentityBinding
        But an attempt to obtain an access token from that pod should "Fail"
