@tags
Feature: Create boundaries between resources using namespaces

As a Cloud Security Admin
I want to ensure that that users have limited permissions 
To limit the scope of malicious activities/ mistakes

Rule: ...

@detective
Scenario Outline: detect when users attempt to access permissions outside those which they are allowed
Given an active kubernetes cluster exists which we can make changes to
And some system exists which can detect when a user attempts to access denied permissions
When a user attempts to access <PermissionsGranted> permissions
Then the system will <Result> with an error <ErrorMessage>

Examples:
    | PermissionsGranted | Result  | ErrorMessage                                        |
    | Denied             | Fail    | Your role does not have access to these permissions |
    | Allowed            | Succeed | No error would show                                 |