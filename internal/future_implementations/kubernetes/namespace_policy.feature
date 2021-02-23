@tags
Feature: Ensure that Namespaces are given Network Policies

  As a Cloud Security Admin
  I want to ensure that all namespaces have an individual network policy
  So that my organization can ensure that my cluster runs smoothly 

  Rule: ...

    @preventative
    Scenario Outline: Prevent the creation of namespaces without an assigned network policy
      Given an active Kubernetes cluster exists which we can make changes to
      And Some system exists to detect whether created namespaces are given an individual network policy
      When a namespace is created
      And the system marks the network being assigned an individual namespace as <NetworkPolicy>
      Then communication will <Result> with error <Error Message>

      Examples:
        | NetworkPolicy | Result  | Error Message                                             |
        | False         | Fail    | A created namespace requires an individual network policy |
        | True          | Succeed | No error will show                                        |
    
    
    
    