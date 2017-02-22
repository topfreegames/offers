Feature: Create an offer template

  Background:
    Given the server is up
    And a game with name "offer-template-game" exists

  Scenario Outline: offer template is created
    Given the server is up
    When an offer template is created in the "offer-template-game" game with:
      | name   | product_id | contents   | metadata   | period   | frequency   | trigger   | placement   |
      | <name> | <pid>      | <contents> | <metadata> | <period> | <frequency> | <trigger> | <placement> |
    Then an offer template with name "<name>" exists in game "offer-template-game"

    Examples:
      | name | pid            | contents   | metadata   | period       | frequency          | trigger                                 | placement |
      | oc1  | com.tfg.sample | { 'x': 1 } | { 'y': 2 } | { 'max': 1 } | { 'every': '50s' } | { 'from': 1486678078, 'to': 1486678079} | 'popup'   |

  Scenario: can't create offer template with same name
    Given an offer template exists with name "oc2" in game "offer-template-game"
    When an offer template is created in the "offer-template-game" game with:
      | name | product_id     | contents   | metadata   | period       | frequency          | trigger                                 | placement |
      | oc2  | com.tfg.sample | { 'x': 1 } | { 'y': 2 } | { 'max': 1 } | { 'every': '50s' } | { 'from': 1486678078, 'to': 1486678079} | 'popup'   |
    Then the last request returned status code 409
    And the last error is "OFF-03" with message "There's already an offer template with the same name"
    And an offer template with name "oc2" does not exist in game "offer-template-game"

  Scenario Outline: can't create offer template with invalid payload
    Given the server is up
    When an offer template is created in the "offer-template-game" game with:
      | name   | product_id | contents   | metadata   | period   | frequency   | trigger   | placement   |
      | <name> | <pid>      | <contents> | <metadata> | <period> | <frequency> | <trigger> | <placement> |
    Then the last request returned status code 400
    And the last error is "OFF-04" with message "<error>"
    And an offer template with name "oc3" does not exist in game "offer-template-game"

    Examples:
      | name | pid            | contents   | metadata   | period       | frequency          | trigger                                 | placement | error                                                      |
      |      | com.tfg.sample | { 'x': 1 } | { 'y': 2 } | { 'max': 1 } | { 'every': '50s' } | { 'from': 1486678078, 'to': 1486678079} | 'popup'   | The name is required to create a new offer template.       |
      | oc3  |                | { 'x': 1 } | { 'y': 2 } | { 'max': 1 } | { 'every': '50s' } | { 'from': 1486678078, 'to': 1486678079} | 'popup'   | The product id is required to create a new offer template. |
      | oc3  | com.tfg.sample |            | { 'y': 2 } | { 'max': 1 } | { 'every': '50s' } | { 'from': 1486678078, 'to': 1486678079} | 'popup'   | The contents are required to create a new offer template.  |
      | oc3  | com.tfg.sample | { 'x': 1 } | { 'y': 2 } |              | { 'every': '50s' } | { 'from': 1486678078, 'to': 1486678079} | 'popup'   | The period is required to create a new offer template.     |
      | oc3  | com.tfg.sample | { 'x': 1 } | { 'y': 2 } | { 'max': 1 } |                    |                                         | 'popup'   | The trigger is required to create a new offer template.    |
