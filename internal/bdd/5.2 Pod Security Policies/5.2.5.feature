@tags
Feature: Do not generally permit containers to be run with the allowPrivilegeEscalation flag.

As a a Cloud Security Admin ..
I want to ensure that no containers can run with the allowPrivilegeEscalation flag
So that my organization is not vulnerable to attacks on processes...

Rule: ...

	@preventative
	Scenario Outline: Prevent a deployment from running with the allowPrivilegeEscalation flag
		Given A kubernetes cluster exists which we can deploy into
		And some system exists to detect when a kubernetes deployment attempts to run on an existing kubernetes container using the allowPrivilegeEscalation flag
		When a kubernetes deployment is applied to an existing kubernetes cluster
		And the kubernetes deployment attempts to run with privileged escalation marked <allowPrivilegeEscalation requested>
		Then the operation will <RESULT> with an error <ERRORMESSAGE>

		Examples: 
			| AllowprivilegeEscalation requested | RESULT        | ERRORMESSAGE							                       |
			| True                               | Fail          | Containers cant run using the allowPrivilegeEscalation flag |
			| False                              | Succeed       | No error would show                                         |
			| Not Defined                        | Succeed       | No error would show                                         |
