@tags

Feature: Ensure that the CNI supports Network Policies

As a Cloud Security Admin
I want to ensure that kubernetes deployments do not run without Network policies enabled
So that my organization can ensure that the CNI is secure

Rule: ...

    @preventative
    Scenario Outline: Ensure that the CNI supports Network Policies
        Given a created kubernetes cluster has a working CNI which can be intereacted with
        When a kubernetes cluster is created 
        And and the cluster's CNI has Network Policies <Network Policy Used>
        Then the cluster's creation will <Result> with an error <Error Message>
    
    Examples:
        | Network Policy Used | Result   | Error Message                                |
        | enabled             | Succeed  |                                              |
        | disabled            | Fail     | Network Policies must be enabled on this CNI |
