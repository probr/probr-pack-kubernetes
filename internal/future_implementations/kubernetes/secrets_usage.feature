@tags
Feature: Secrets are stored through an external secret provider

As a Cloud Security Admin
I want to ensure that secret objects are stored in an external secrets providers
So that access to secrets for my organization is limited

Rule: ...

@preventative
Scenario Outline: Prevent secrets from being accessed inside the kubernetes cluster
Given an active kubernetes cluster exists
And some system exists to detect when someone attempts to access secrets
When someone tries to access secrets inside the kubernetes cluster
And the access limited is <Access> 
Then the attempt will <Result> with error <Error Message>

Examples:
    | Access    | Result      | Error Message                                            |
    | True      | Fail        | Secrets can't be accessed inside the kubernetes cluster |
    | False     | Succeed     | No error will show                                       |