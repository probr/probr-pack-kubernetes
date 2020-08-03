@tags
Feature: Do not generally permit containers to be run as the root user

As a a Cloud Security Admin ..
I want to ensure that no containers can run as the root user
So that my organization is not vulnerable to attacks on processes...

Rule: ...

	@preventative 
	Scenario Outline: Prevent a deployment from running as the root user
		Given A kubernetes cluster exists which we can deploy into
		And some system exists to detect when a kubernetes deployment attempts to run on an existing kubernetes container as the root user
		When a kubernetes deployment is applied to the existing kubernetes cluster
		And the deployment has marked the additional capabilities request as <additional capabilities requested>
		Then the operation will <RESULT> with an error <ERRORMESSAGE>

		Examples: 
			| additional capabilities requested | RESULT        | ERRORMESSAGE							           |
			| True                              | Fail          | Containers cant run with privileged capabilities |
			| False                             | Succeed       |                                                  |
			| Not Defined                       | Succeed       |                                                  |