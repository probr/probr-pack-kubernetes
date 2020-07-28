@tags
Feature: Do not generally permit kubernetes deployment to be run with capabilities beyond the default set.

As a (....)
I want to ensure that no kubernetes deployment can run with capabilities beyond the default set.
So that my organization is not vulnerable to attacks by malicious kubernetes deployments.

Rule: Insert tags here.

@preventitive
Scenario Outline: Prevent container running with capabilities beyond the default set.
		Given A kubernetes cluster exists which we can deploy into
		And some control exists to prevent kubernetes deployments with capabilities beyond the default set for being deployed to an active kubernetes cluster
		When a container is deployed to an active kubernetes deployment
		And the <additional capabilities requested> by the kubernetes deployment.
		Then the operation will <RESULT> with an error <ERRORMESSAGE>

		Examples: 
			| additional capabilities requested | RESULT        | ERRORMESSAGE							           |
			| True                              | Fail          | Containers cant run with privileged capabilities |
			| False                             | Succeed       |                                                  |
			| Not Defined                       | Succeed       |                                                  |