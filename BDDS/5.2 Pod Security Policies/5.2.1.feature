@tags
Feature: Minimize the admission of privileged containers (5.2.1)

As a Cloud Security Admin ..
I want to ensure that no containers run with privileged access
So that my organization is not vulnerable to attacks on processes...

Rule: Insert tags here.

	@preventative
	Scenario Outline: prevent deployments running with privileged access
		Given A kubernetes cluster exists which we can deploy into
		And some control exists to prevent privileged access for kubernetes deployments to an active kubernetes cluster
		When a Kubernetes deployment is applied to the active Kubernetes cluster
		And privileged access request is marked <privileged access requested> for the Kubernetes deployment
		Then the operation will <RESULT> with an error <ERRORMESSAGE>

		Examples: 
			| privileged access requested | RESULT        | ERRORMESSAGE							     |
			| True                        | Fail          | Containers with privileged access can’t run  |
			| False                       | Succeed       | No error would show                          |
			| Not Defined                 | Succeed       | No error would show                          |
			