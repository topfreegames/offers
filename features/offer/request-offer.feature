Feature: Request offers for a player

  Background:
    Given the server is up
    And a game with name "offer-request-game" exists

  Scenario Outline: A player requests offers that can be bought only once to be shown every minute
    Given the following offer templates exist in the "offer-request-game" game:
      | game | name | product_id     | contents   | placement | period          | frequency                          | trigger                   |
      | org  | or1  | com.tfg.sample | { 'x': 1 } | popup     | { 'amount': 1 } | { 'every': 60, 'unit': 'seconds' } | { 'from': 0, 'to': 5 }    |
      | org  | or2  | com.tfg.sample | { 'x': 2 } | popup     | { 'amount': 1 } | { 'every': 60, 'unit': 'seconds' } | { 'from': 6, 'to': 10 }   |
      | org  | or3  | com.tfg.sample | { 'x': 3 } | popup     | { 'amount': 1 } | { 'every': 60, 'unit': 'seconds' } | { 'from': 20, 'to': 800 } |
      | org  | or4  | com.tfg.sample | { 'x': 4 } | store     | { 'amount': 1 } | { 'every': 60, 'unit': 'seconds' } | { 'from': 20, 'to': 800 } |
    And the following players exist in the "offer-request-game" game:
      | id        | claimed-offers | last-seen-offer-at |
      | joseph    | -              | -, -, -, -         |
      | john      | or1            | 3, -, -, -         |
      | michael   | or1, or2       | 3, 8, -, -         |
      | jane      | or3            | -, -, 23, -        |
      | mary      | or2, or3       | -, 8, 23, -        |
      | christine | or1, or2, or3  | 3, 8, 23, 23       |
    When the current time is <current_time>
    And the game "offer-request-game" requests offers for player "<player_id>" in "<placement>"
    Then an offer with name "<offer>" is returned

    Examples:
      | current_time | player_id | offer | placement |
      | 3            | joseph    | or1   | popup     |
      | 8            | joseph    | or2   | popup     |
      | 18           | joseph    | -     | popup     |
      | 23           | joseph    | or3   | popup     |
      | 23           | joseph    | or4   | store     |
      | 3            | john      | -     | popup     |
      | 8            | john      | or2   | popup     |
      | 18           | john      | -     | popup     |
      | 23           | john      | or3   | popup     |
      | 3            | michael   | -     | popup     |
      | 8            | michael   | -     | popup     |
      | 18           | michael   | -     | popup     |
      | 23           | michael   | or3   | popup     |
      | 3            | jane      | or1   | popup     |
      | 8            | jane      | or2   | popup     |
      | 18           | jane      | -     | popup     |
      | 23           | jane      | -     | popup     |
      | 23           | jane      | or4   | store     |
      | 3            | mary      | or1   | popup     |
      | 8            | mary      | -     | popup     |
      | 18           | mary      | -     | popup     |
      | 23           | mary      | -     | popup     |
      | 23           | mary      | or4   | store     |
      | 3            | christine | -     | popup     |
      | 8            | christine | -     | popup     |
      | 18           | christine | -     | popup     |
      | 23           | christine | -     | popup     |
      | 23           | christine | -     | store     |
      | 83           | christine | or4   | store     |


  Scenario Outline: After a player has seen an offer there is a track of it
    Given the following offer templates exist in the "offer-request-game" game:
      | game | name | product_id     | contents   | placement | period             | frequency       | trigger                   |
      | org  | ot1  | com.tfg.sample | { 'x': 1 } | popup     | { 'type': 'once' } | { 'every': 60 } | { 'from': 0, 'to': 5 }    |
      | org  | ot2  | com.tfg.sample | { 'x': 2 } | popup     | { 'type': 'once' } | { 'every': 60 } | { 'from': 6, 'to': 10 }   |
    When the current time is <current_time>
    And the game "offer-request-game" requests offers for player "<player_id>" in "popup"
    Then player "<player_id>" has seen offer "<seen_offer>"
    And player "<player_id"> has not seen offer "<unseen_offer>"

    Examples:
      | current_time | player_id | seen_offer | unseen_offer |
      | 3            | jack      | ot1        | ot2          |
      | 8            | jenniffer | ot2        | ot1          |

  Scenario: When a player sees offers subsequently they are both tracked
    Given the following offer templates exist in the "offer-request-game" game:
      | game | name   | product_id     | contents   | placement | period             | frequency       | trigger                 |
      | org  | oseen1 | com.tfg.sample | { 'x': 1 } | popup     | { 'type': 'once' } | { 'every': 60 } | { 'from': 0, 'to': 5 }  |
      | org  | oseen2 | com.tfg.sample | { 'x': 2 } | popup     | { 'type': 'once' } | { 'every': 60 } | { 'from': 6, 'to': 10 } |
    When the current time is 3
    And the game "offer-request-game" requests offers for player "Henry" in "popup"
    And the current time is 8
    And the game "offer-request-game" requests offers for player "Henry" in "popup"
    Then player "Henry" has seen offer "oseen1"
    And player "Henry" has seen offer "oseen2"
