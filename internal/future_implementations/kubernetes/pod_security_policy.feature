Feature: Pod Security Policy additions

	@probes/kubernetes/pod_security_policy/1.9 @control_type/preventative @standard/none/PSP-0.1
	Scenario Outline: Prevent deployments from accessing unapproved port range
		Given a Kubernetes cluster exists which we can deploy into
		And some system exists to prevent Kubernetes deployments with unapproved port range from being deployed to an existing Kubernetes cluster
        When a Kubernetes deployment is applied to an existing Kubernetes cluster
		And an "<requested>" port range is requested for the Kubernetes deployment
		Then the operation will "<RESULT>" with an error "<ERRORMESSAGE>"
		But I should not be able to perform a command that access an unapproved port range
		And I should be able to perform an allowed command

		Examples:
			| requested 	| RESULT 	| ERRORMESSAGE							|
			| unapproved  	| Fail  	| Cannot access unapproved port range	|
			| approved		| Succeed	|									  	|
			| not defined	| Succeed	|										|
