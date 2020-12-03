@tags
Feature: Make sure that pods, containers and volumes have configured security contexts

As a Cloud Security Admin
I want to ensure that my pods, containers and volumes have configured security contexts
So that my organization kubernetes cluster has the right operating system on their containers

Rule: ...

@preventative
Scenario Outline: Prevent pods and containers from running if the security contexts hasn't been configured
Given an active kubernetes cluster exists which we can make changes to
And some system exists to detect when a pod attempts to run without configuring security contexts
When a pod or a container is deployed to the kubernetes cluster
And the security contexts have <Configuration> 
Then the upload will <Result> with error <Error Message>

Examples:
    | Configuration       | Result  | Error Message                                                           |
    | Been Configured     | Fail    | No security configuration prevent the pod/container from being uploaded |
    | Not been configured | Succeed | No error will show                                                      |

