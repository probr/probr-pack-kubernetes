Feature: Minimize the admission of privileged containers (5.2.1)

    As a Cloud Security Admin ..
    I want to ensure that no containers run with privileged access
    So that my organisation is not venerable to attacks for processes...

    Scenario Outline: Prevent container running with privileged access
	    Given control exists to prevent privileged access 
	    When a deployment is created
	    And "<Privileged>" access is requested
	    Then creation will "<Result>" with a message "<Error Description>"

        Examples:
      | Privileged |  Result | Error Description        |
      | No         | Succeed |                          |
      | Yes        | Fail    | Cannot create deployment |