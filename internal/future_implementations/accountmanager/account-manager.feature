Feature: Cloud Driver Account Manager

  As an administrator
  I want configure target entities (AWS Account, Azure Subscription, Google Cloud Project etc)
  So that I can run tests against them
  
  Scenario Outline: Set up of credentials to connect to cloud account
    Given I am configuring a "<Cloud Account>" Account
    And "<Cloud Credential>" Credential with access to the "<Cloud Account>" Account is already configured in the system
    When I add the "<Cloud Account>" Account details to the system
    And I link the "<Cloud Credential>" Credential to the "<Cloud Account>" Account
    Then a resource deployment will "<Result>" with the message "<Error Description>"

    Examples:
      | Cloud Account | Cloud Credential | Result  | Error Description                                   |
      | AWS           | AWS Credential   | Succeed | Success                                                    |
      | AWS           | Azure Credential | Fail    | Cannot connect to account using supplied credential |