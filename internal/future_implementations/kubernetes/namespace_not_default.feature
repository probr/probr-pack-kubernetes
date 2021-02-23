@tags
Feature: Use network policies to isolate traffic in your cluster network.

  As a Cloud Security Admin
  I want to ensure that resources in a kubernetes cluster are not in the default namespace
  So that my organization can apply security contexts at that level

  Rule: ...

    #TODO PJITREVIEW pull into Service Pack and rename steps to match latest impls
    @preventative
    Scenario Outline: Resources in the kubernetes cluster should not be uploaded into the default namespace
      Given a Kubernetes cluster is deployed
      When a pod is deployed in the cluster
      And the system rules that the upload <SUCCEEDED> uploaded to the default namespace
      Then the operation will "<RESULT>" with an error "<ERRORMESSAGE>"

      Examples:
        | SUCCEEDED | RESULT  | ERRORMESSAGE                                            |
        | Was       | Fail    | A resource can't be uploaded with the default namespace |
        | Was Not   | Succeed | No error will show                                      |