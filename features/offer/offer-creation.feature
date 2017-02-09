Feature: Create an offer

  Background:
    Given the server is up
    And a game with name "offer-game" exists

  Scenario Outline: offer is created
    Given the server is up
    When an offer is created with:
      | name   | product_id | contents   | metadata   | trigger   | period   |
      | <name> | <pid>      | <contents> | <metadata> | <trigger> | <period> |
    Then an offer with name "<name>" exists in game "offer-game"

    Examples:
      | name | pid            | contents   | metadata   | trigger           | period                                  |
      | oc1  | com.tfg.sample | { "x": 1 } | { "y": 2 } | { "daily": true } | { "from": 1486678078, "to": 1486678079} |

  Scenario: can't create offer with same name
    Given an offer exists with name "oc2" in game "offer-game"
    When an offer is created with:
      | name | product_id     | contents   | metadata   | trigger           | period                                  |
      | oc2  | com.tfg.sample | { "x": 1 } | { "y": 2 } | { "daily": true } | { "from": 1486678078, "to": 1486678079} |
    Then the last request returned status code 409
    And the last error is "OFF-03" with message "There's already an offer with the same name"
    And an offer with name "oc2" does not exist in game "offer-game"

  Scenario Outline: can't create offer with invalid payload
    Given the server is up
    When an offer is created with:
      | name   | product_id | contents   | metadata   | trigger   | period   |
      | <name> | <pid>      | <contents> | <metadata> | <trigger> | <period> |
    Then the last request returned status code 400
    And the last error is "OFF-04" with message "<error>"
    And an offer with name "oc3" does not exist in game "offer-game"

    Examples:
      | name | pid            | contents   | metadata   | trigger           | period                                  | error                                                  |
      |      | com.tfg.sample | { "x": 1 } | { "y": 2 } | { "daily": true } | { "from": 1486678078, "to": 1486678079} | The name is required to create a new offer.            |
      | oc3  |                | { "x": 1 } | { "y": 2 } | { "daily": true } | { "from": 1486678078, "to": 1486678079} | The product id is required to create a new offer.      |
      | oc3  | com.tfg.sample |            | { "y": 2 } | { "daily": true } | { "from": 1486678078, "to": 1486678079} | The offer contents are required to create a new offer. |
      | oc3  | com.tfg.sample | { "x": 1 } | { "y": 2 } |                   | { "from": 1486678078, "to": 1486678079} | The trigger is required to create a new offer.         |
      | oc3  | com.tfg.sample | { "x": 1 } | { "y": 2 } | { "daily": true } |                                         | The period is required to create a new offer.          |
