Feature: Configure access to a cloud account

  As a developer
  I want to set up credentials for my cloud account
  So that the system can execute tests against my cloud account

  Scenario Outline: Add a cloud account
    Given there is a <account_nomenclature> set up in <cloud_provider>
    When I provide <account_nomenclature> information and access credentials for a specific <account_nomenclature>
    Then the system can successfully connect to the <account_nomenclature>
    And the <account_nomenclature> information and credentials are stored in the system for future use
    And the system returns a unique reference identifier for the <account_nomenclature> information and credentials

    Examples:
      | account_nomenclature | cloud_provider |
      | account              | AWS            |
      | subscription         | Azure          |
      | project              | GCP            |

  # an attempt to express that passwords need to be obtained from a secrets vault every time
  Scenario: Handling of sensitive credential information
    Given there is a need to communicate with a cloud account to complete a particular activity
    When the system attempts to connect to the cloud account
    Then the sensitive credential information (e.g. passwords, certificates, access token) must always be pulled from the secrets vault
    And the sensitive credential information is never held in memory beyond the scope of the activity
    And the sensitive credential information is never persisted to disk