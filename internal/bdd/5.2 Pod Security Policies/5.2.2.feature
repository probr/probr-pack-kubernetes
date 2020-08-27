@tags
Feature: Do not generally permit containers to be run with the hostPID .

As a a Cloud Security Admin ..
I want to ensure that no containers can run using the host PID 
So that my organization is not vulnerable to attacks on processes...

Rule: ...

	@preventative
	Scenario Outline: Prvent a deployment from running with the hostPID 
		Given A kubernetes cluster exists which we can deploy into
		And some system exists to prevent a Kubernetes container from running using the hostPID on the active kubernetes cluster
		When a Kubernetes deployment is applied to an existing Kubernetes cluster 
		And hostPID request has been marked <HostPID requested> for the kubernetes deployment.
		Then the operation will <RESULT> with an error <ERRORMESSAGE>

		Examples: 
			| HostPID requested           | RESULT        | ERRORMESSAGE							     |
			| True                        | Fail          | Containers cant run using hostPID            |
			| False                       | Succeed       |                                              |
			| Not Defined                 | Succeed       |                                              |
		