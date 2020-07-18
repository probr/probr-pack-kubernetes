@tags
Feature: Do not generally permit containers to be run with assigned capabilities.

As a Cloud Security Admin
I want to ensure that no kubernetes deployments run with assigned capabilities
So that my organization is not vulnerable to attacks by malicious deployments.

Rule: ...

    @preventitive
    Scenario Outline: Prevent deployments from running with assigned capabilities.
        Given a kubernetes cluster exists which we can deploy into.
		And some system exists to prevent a Kubernetes container from running using assigned capabilities on the active kubernetes cluster
        When a deployment is created to an active kubernetes cluster
		And container has assigned capabilities is marked true <Container has assigned capabilities> for the Kubernetes
        Then creation will <Result> with an error <ERRORMESSAGE>

    Examples: 
        | Container has assigned capabilities | Result   | ERRORMESSAGE                                            |
        | True                                | Fail     | Containers with assigned capabilities can’t be deployed |
        | False                               | Succeed  |                                                         |
		| Not defined                         | Succeed  |                                                         |