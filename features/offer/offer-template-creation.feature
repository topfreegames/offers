Feature: Create an offer

  Background:
    Given the server is up
    And a game with id "offer-game" exists

  Scenario Outline: offer is created
    Given the server is up
    When an offer is created in the "offer-game" game with name "<name>" pid "<pid>" contents "<contents>" metadata "<metadata>" period "<period>" freq "<frequency>" trigger "<trigger>" place "<placement>"
    Then an offer with name "<name>" exists in game "offer-game"

    Examples:
      | name | pid            | contents   | metadata   | period       | frequency          | trigger                                 | placement |
      | oc1  | com.tfg.sample | { 'x': 1 } | { 'y': 2 } | { 'max': 1 } | { 'every': '50s' } | { 'from': 1486678078, 'to': 1486678079} | popup     |

  Scenario Outline: can't create offer with invalid payload
    Given the server is up
    When an offer is created in the "offer-game" game with name "<name>" pid "<pid>" contents "<contents>" metadata "<metadata>" period "<period>" freq "<frequency>" trigger "<trigger>" place "<placement>"
    Then the last request returned status code 422
    And the last error is "OFF-002" with message "<error>"
    And an offer with name "oc3" does not exist in game "offer-game"

    Examples:
      | name | pid            | contents   | metadata   | period       | frequency          | trigger                                 | placement | error                                             |
      |      | com.tfg.sample | { 'x': 1 } | { 'y': 2 } | { 'max': 1 } | { 'every': '50s' } | { 'from': 1486678078, 'to': 1486678079} | popup     | The name is required to create a new offer.       |
      | oc3  |                | { 'x': 1 } | { 'y': 2 } | { 'max': 1 } | { 'every': '50s' } | { 'from': 1486678078, 'to': 1486678079} | popup     | The product id is required to create a new offer. |
      | oc3  | com.tfg.sample |            | { 'y': 2 } | { 'max': 1 } | { 'every': '50s' } | { 'from': 1486678078, 'to': 1486678079} | popup     | The contents are required to create a new offer.  |
      | oc3  | com.tfg.sample | { 'x': 1 } | { 'y': 2 } |              | { 'every': '50s' } | { 'from': 1486678078, 'to': 1486678079} | popup     | The period is required to create a new offer.     |
      | oc3  | com.tfg.sample | { 'x': 1 } | { 'y': 2 } | { 'max': 1 } |                    | { 'from': 1486678078, 'to': 1486678079} | popup     | The frequency is required to create a new offer.  |
      | oc3  | com.tfg.sample | { 'x': 1 } | { 'y': 2 } | { 'max': 1 } | { 'every': '50s' } |                                         | popup     | The trigger is required to create a new offer.    |
      | oc3  | com.tfg.sample | { 'x': 1 } | { 'y': 2 } | { 'max': 1 } | { 'every': '50s' } | { 'from': 1486678078, 'to': 1486678079} |           | The placement is required to create a new offer.  |
