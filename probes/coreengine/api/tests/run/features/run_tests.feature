Feature: Run tests

  As a developer
  I want to run compliance tests against my resources
  So that I can be confident those resources are compliant

  Scenario Outline: Run tests against cloud guardrails
    Given there are controls implemented around an account in the cloud
    And there are account credentials held in the system
    When I ask the system to test those controls for compliance
    And I provide the unique reference identifier for a specific set of account credentials
    And the account credential identifier is <credential_id>
    Then the system responds with <http_code>
    And the system responds with <response_body>
    And the tests are <tests_executed>
    ##TODO: we need to link the account to a test pack
    Examples:
      | credential_id | http_code | tests_executed | response_body                              |
      | found         | 200       | true           | a unique identifier for that specific test |
      | not found     | 404       | false          | error message "credentials not found"      |

  #MVP
  Scenario Outline:
    Given there are resources instantiated in the cloud of a given service type
    And there are account credentials held in the system
    When I ask the system to test instantiated resources for compliance
    And I provide a unique reference identifier for a set of account credentials
    And I provide cloud provider resource identifiers for the instantiated resources to be tested
    And I provide a unique identifier for the test pack to execute against the instantiated resources
    And the account credential identifier is <credential_id>
    And the test pack identifier is <testpack_id>
    And the resources based on the cloud provider resources identifiers can be <connected_to>
    Then the system responds with <http_code>
    And the tests are <tests_executed> against the resources based on the cloud provider resource identifiers using the account credentials behind the account credential identifier
    And the system responds with <response_body>

    Examples:
      | credential_id | testpack_id | connected_to     | http_code    | tests_executed | response_body                                             |
      | not found     | not found   | na               | 404          | not executed   | error message "credentials not found, testpack not found" |
      | not found     | found       | na               | 404          | not executed   | error message "credentials not found"                     |
      | found         | not found   | na               | 404          | not executed   | error message "test pack not found"                       |
      | found         | found       | connected to     | 200          | executed       | a unique identifier for that specific test                |
      | found         | found       | not connected to | 404          | not executed   | error message "resources cannot be connected to"          |

  #MVP
  Scenario Outline: Get the current status of a running test
    Given I have kicked off a test in the cloud
    When I ask the system for the status of a test
    And I provide the unique identifier for a specific test
    And the test unique identifier is <test_id>
    And the state of the test is <test_state>
    Then the system responds with <http_code>
    And the system responds with <response>

    Examples:
      | test_id   | test_state                            | response        | http_code |
      | not found | na                                    | test not found  | 404       |
      | found     | Waiting to execute                    | pending         | 200       |
      | found     | Test execution currently in progress  | running         | 200       |
      | found     | Error in test execution               | error           | 200       |
      | found     | All tests have completed              | completed       | 200       |

  Scenario Outline: Get result summary of completed test
    Given I have kicked off a test in the cloud
    When I ask the system for the summary results of a test
    And I provide the unique identifier for a specific test
    And the test unique identifier is <test_id>
    And the state of the test is <test_state>
    Then the system responds with <http_code>
    And the system responds with <results>

    Examples:
      | test_id   | test_state                            | results                              | http_code |
      | not found | na                                    | test not found                       | 404       |
      | found     | Waiting to execute                    | pending                              | 200       |
      | found     | Test execution currently in progress  | running                              | 200       |
      | found     | Error in test execution               | reason for error                     | 200       |
      | found     | All tests have completed              | array of [{"test_name": "pass\|fail"}] | 200       |


  #MVP
  Scenario Outline: Get detailed results of completed test
    Given I have kicked off a test in the cloud
    When I ask the system for the detailed results of a test
    And I provide the unique identifier for a specific test
    And the test unique identifier is <test_id>
    And the state of the test is <test_state>
    Then the system responds with <http_code>
    And the system responds with <results>

    Examples:
      | test_id   | test_state                            | results                               | http_code |
      | not found | na                                    | test not found                        | 404       |
      | found     | Waiting to execute                    | pending                               | 200       |
      | found     | Test execution currently in progress  | running                               | 200       |
      | found     | Error in test execution               | reason for error                      | 200       |
      | found     | All tests have completed              | array of [cucumber test json output]  | 200       |

