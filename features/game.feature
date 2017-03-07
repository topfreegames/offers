Feature: Game Management

  Scenario Outline: Game is created that did not exist before
    Given the server is up
    When a game with id "<id>" is created with a name "<name>"
    Then the game "<id>" exists
    And the game "<id>" has name of "<name>"

    Examples:
      | id    | name           |
      | game1 | game 1 example |
      | game2 | game 2 example |

  Scenario Outline: Game is updated, created if not exists
    Given the server is up
    When a game with id "<id>" is updated with a name "<name>"
    Then the game "<id>" exists
    And the game "<id>" has name of "<name>"

    Examples:
      | id    | name           |
      | game3 | game 3 example |
      | game4 | game 4 example |
      | game4 | game 5 example |

  Scenario Outline: The game is not created with invalid information
    Given the server is up
    When a game with id "<id>" is updated with a name "<name>"
    Then the last request returned status code <status>
    And the last error is "<error_code>" with message "<error_message>"
    And the game "<id>" does not exist

    Examples:
      | id            | name        | status | error_code | error_message                               |
      |               | game name x | 404    | OFF-002    | 404 page not found                          |
      | game5         |             | 422    | OFF-002    | Name: non zero value required;              |
      | asd*@!3[1249  | game name y | 422    | OFF-002    | ID: non zero value required;                |
      | @VeryBigText@ | game name z | 422    | OFF-002    | *does not validate as stringlength(1\|255); |
