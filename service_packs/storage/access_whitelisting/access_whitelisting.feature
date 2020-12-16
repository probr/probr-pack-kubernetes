@probes/storage
@probes/storage/access_whitelisting
@standard/citihub/CHC2-CHC2-SVD030
@csp.aws
@csp.azure
Feature: Object Storage Has Network Whitelisting Measures Enforced

  As a Cloud Security Architect
  I want to ensure that suitable security controls are applied to Object Storage
  So that my organisation's data can only be accessed from whitelisted IP addresses

  #Rule: CHC2-SVD030 - protect cloud service network access by limiting access from the appropriate source network only

    @probes/storage/access_whitelisting/1.0 @control_type/detective
    Scenario: Check Object Storage is Configured With Network Source Address Whitelisting
      Given the CSP provides a whitelisting capability for Object Storage containers
      And a specified azure resource group exists
      When we examine the Object Storage container in environment variable "TARGET_STORAGE_CONTAINER"
      Then whitelisting is configured with the given IP address range or an endpoint

    @probes/storage/access_whitelisting/1.1 @control_type/preventative
    Scenario Outline: Prevent Object Storage from Being Created Without Network Source Address Whitelisting
      Given security controls that Prevent Object Storage from being created without network source address whitelisting are applied
      And a specified azure resource group exists
      When we provision an Object Storage container
      And it is created with whitelisting entry "<Whitelist Entry>"
      Then creation will "<Result>"

      Examples:
        | Whitelist Entry | Result  |
        | 219.79.19.0/24  | Success |
        | 219.79.19.1     | Fail    |
        | 219.108.32.1    | Fail    |
        | 170.74.231.168  | Success |
        | nil             | Fail    |
