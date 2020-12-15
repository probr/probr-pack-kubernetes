@probes/storage
@probes/storage/encryption_at_rest
@standard/citihub/CHC2-SVD001
@standard/citihub/CHC2-AGP140
@standard/citihub/CHC2-EUC001
Feature: Object Storage Encryption at Rest

  As a Cloud Security Architect
  I want to ensure that suitable security controls are applied to Object Storage
  So that my organisation is protected against data leakage due to misconfiguration

    #Rule: CHC2-AGP140 - Ensure cryptographic controls are in place to protect the confidentiality and integrity of data in-transit, stored, generated and processed in the cloud

    @probes/storage/encryption_at_rest/1.0 @control_type/preventative @csp/azure
    Scenario Outline: Prevent Creation of Object Storage Without Encryption at Rest
      Given security controls that restrict data from being unencrypted at rest
      When we provision an Object Storage bucket
      And encryption at rest is "<Encryption Option>"
      Then creation will "<Result>" with an error matching "<Error Description>"

      Examples:
        | Encryption Option | Result  | Error Description                                                      |
        | enabled           | Fail    | Storage Buckets must not be created without encryption as rest enabled |
        | disabled          | Succeed |                                                                        |

    @probes/storage/encryption_at_rest/1.1 @control_type/detective @csp/azure
    Scenario: Detect creation of Object Storage Without Encryption at Rest
      Given there is a detective capability for creation of Object Storage without encryption at rest
      And the capability for detecting the creation of Object Storage without encryption at rest is active
      When Object Storage is created with without encryption at rest
      Then the detective capability detects the creation of Object Storage without encryption at rest
      And the detective capability enforces encryption at rest on the Object Storage Bucket