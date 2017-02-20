Feature: Game Management

  Scenario Outline: Game is created that did not exist before
    Given the server is up
    When a game named "<name>" is created with bundle id of "<bundle_id>"
    Then the game "<name>" exists
    And the game "<name>" has bundle id of "<bundle_id>"

    Examples:
      | name  | bundle_id                |
      | game1 | com.topfreegames.example |
      | game2 | random-bundle-id         |

  Scenario Outline: Game is updated, created if not exists
    Given the server is up
    When a game named "<name>" is updated with bundle id of "<bundle_id>"
    Then the game "<name>" exists
    And the game "<name>" has bundle id of "<bundle_id>"

    Examples:
      | name  | bundle_id                |
      | game3 | com.topfreegames.example |
      | game4 | random-bundle-id         |
      | game4 | com.topfreegames.other   |

  Scenario Outline: The game is not created with invalid information
    Given the server is up
    When a game named "<name>" is updated with bundle id of "<bundle_id>"
    Then the last request returned status code <status>
    And the last error is "<error_code>" with message "<error_message>"
    And the game "<name>" does not exist

    Examples:
      | name          | bundle_id                | status | error_code | error_message                               |
      |               | com.topfreegames.example | 422    | OFF-002    | ID: non zero value required;                |
      | game5         |                          | 422    | OFF-002    | BundleID: non zero value required;          |
      | asd*@!3[1249  | com.topfreegames.example | 422    | OFF-002    | ID: non zero value required;                |
      | @VeryBigText@ | com.topfreegames.example | 422    | OFF-002    | *does not validate as stringlength(1\|255); |
