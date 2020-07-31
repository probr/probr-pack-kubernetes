@tags
Feature: Do not generally permit containers to be run with the NET_RAW capability.

As a Cloud Security Admin
I want to ensure that no containers run with the NET_RAW capability.
So that my organization is not vulnerable to attacks by malicious containers.

Rule: ...

	@preventitive
	Scenario Outline: Prevent deployments from running with the NET_RAW capability.
		Given A kubernetes cluster exists which we can deploy into
		And some control exists to detect and prevent containers with NET_RAW capability for being deployed to an active kubernetes cluster
		When a deployment is applied to an active kubernetes cluster
		And  the deployment has the net_raw request marked <Container has NET_RAW capabilities> for the Kubernetes deployment
		Then the operation will <RESULT> with an error <ERRORMESSAGE>

		Examples: 
			| HostPID requested           | RESULT        | ERRORMESSAGE							      |
			| True                        | Fail          | Containers cant run with NET_RAW capabilities |           
			| False                       | Succeed       |                                               |
			| Not Defined                 | Succeed       |                                               |