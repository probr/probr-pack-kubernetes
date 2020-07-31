@tags
Feature: Do not generally permit containers to be run with the host Network flag.

As a a Cloud Security Admin ..
I want to ensure that no containers can run using the host Network flag
So that my organization is not vulnerable to attacks on processes...

Rule: ...

	@detective
	Scenario Outline: Detect when a deployment attempts to run with the host Network flag
		Given some system exists to detect when a kubernetes deployment attempts to run on an existing kubernetes container using the host Network
		When a kubernetes deployment is created 
		And the kubernetes deployment attempts to run using the host Network
		Then the creation will <FAIL> with an error <ERRORMESSAGE>
		
		Examples:
			| Container runs with host Network | Result   | Error Description                             |
			| True                             | Fail     | You cannnot run Containers using host Network |
			| False                            | Succeed  |                                               |
			| Not Defined                      | Succeed  | No error would show                           |

			