@tags
Feature: Use network policies to isolate traffic in your cluster network.

As a Cloud Security Admin
I want to ensure that resources in a kubernetes cluster are not in the default namespace
So that my organization can apply security contexts at that level

Rule: ...

@preventitive
Scenario Outline: Resources in the kubernetes cluster should not be uploaded into the default namespace
Given an active Kubernetes cluster exists which we can make changes to 
And some system exists which can detect if reasources are uploaded into the default namespace
When a resource is uploaded to the kubernetes cluster
And the system rules that the upload <NameSpace> uploaded to the default namespace
Then upload will <Result> with error <Error Message>

Examples:
    | NameSpace | Result  | Error Message                                           |
    | Was       | Fail    | A resource can't be uploaded with the default namespace |
    | Was Not   | Succeed | No error will show                                      |  