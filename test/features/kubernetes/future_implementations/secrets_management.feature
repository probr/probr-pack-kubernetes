@tags
Feature: Ensure that secrets are stored as files instead of environment variables

As a Cloud Security Admin
I want to ensure that secret files are stored as seperate files instead of as environment variables
So that my organization can ensure that secrets are not easily discovered by reviewing the network logs

Rule: ...

@preventitive
Scenario Outline: Prevent secrets from being stored as environmental variables
Given an active Kubernetes cluster exists
And Some system exists to detect how and when secrets are being defined
When secrets are being defined 
And the user attempts to save the secrets as a <Secret Type>
Then communication will <Result> with error <Error Message>

Examples:
    | Secret Type            | Result  | Error Message                               |
    | File                   | Succeed | No error will show                          |
    | Environmental Variable | Fail    | Secrets should be stored as a separate file |
    
    