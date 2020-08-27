@tags
Feature: Do not generally permit containers to be run with the hostIPC flag.

As a Cloud security admin..
I want to ensure that no containers run with the hostIPC flag.
So that my organization is not vulnerable to attacks on processes.

Rule: ...

	@preventative
	Scenario Outline: Prevent container running with the hostIPC flag.
		Given A kubernetes cluster exists which we can deploy into
		And some control exists to detect and prevent hostIPC access for kubernetes deployment
		When a Kubernetes deployment is applied to the active Kubernetes cluster 
		And hostIPC access request is marked <hostIPC access is requested> for the Kubernetes
		Then creation will <Result> with an error <Error Description>


	Example:
		| hostIPC access is requested | Result   | Error Description                        |
		| True                        | Fail     | Containers with hostIPC access can’t run |
		| False                       | Succeed  | No error would show                      |
		| Not defined                 | Succeed  | No error would show                      |
   

   
