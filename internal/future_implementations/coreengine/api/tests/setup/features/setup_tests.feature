Feature: setup test
  As a developer
  I want to configure a test pack
  So that I can customise the tests that are executed

  #MVP
  Scenario: Create a test pack
    Given The system is running
    When I ask the system to create a new test pack
    Then the system creates an empty test pack
    And the system responds with a unique identifier for the test pack

  #MVP - we'll only have one service pack = Kubernetesl; or possibly two (AKS + EKS)
  Scenario: Add a service pack to a test pack
    Given There are service packs in the system
    And There are test packs in the system
    When I request the system adds a service pack to a test pack
    And I provide the unique identifier for the service pack
    And I provide the unique identifier for the test pack
    And #validation stuff
    Then the system responds with <response>

  Scenario: Configure a service pack in a test pack
    Given there is a test pack with one or more service packs added to it
    When I ask the system to configure the service pack in the test pack
    And I provide the unique identifier for the test pack
    And I provide the unique identifier for the service pack
    And I provide a configuration object which is specific to the service pack
    And the configuration object is <a valid schema>
    Then the system configures the service pack behind the service pack unique identifier only within the test pack behind the test pack unique identifier
    And the system responds with <response>

  Scenario: Remove a test from a test pack
    Given x
    When y
    Then z

  Scenario: Show full configuration of a service pack


  Scenario: List all test packs

