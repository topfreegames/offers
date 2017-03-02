Feature: Request offers for a player

  Background:
    Given the server is up
    And a game with name "offer-request-game" exists

  Scenario Outline: A player requests offers that can be bought only once to be shown every minute
    Given the following offer templates exist in the "offer-request-game" game:
      | game | name    | product_id     | contents   | placement | period       | frequency          | trigger                   |
      | org  | oronce1 | com.tfg.sample | { "x": 1 } | popup     | { "max": 1 } | { "every": "60s" } | { "from": 0, "to": 5 }    |
      | org  | oronce2 | com.tfg.sample | { "x": 2 } | popup     | { "max": 1 } | { "every": "60s" } | { "from": 6, "to": 10 }   |
      | org  | oronce3 | com.tfg.sample | { "x": 3 } | popup     | { "max": 1 } | { "every": "60s" } | { "from": 20, "to": 800 } |
      | org  | oronce4 | com.tfg.sample | { "x": 4 } | storoncee | { "max": 1 } | { "every": "60s" } | { "from": 20, "to": 800 } |
    And the following players exist in the "offer-request-game" game:
      | id        | claimed-offers            | last-seen-offer-at |
      | joseph    | -                         | -, -, -, -         |
      | john      | oronce1                   | 3, -, -, -         |
      | michael   | oronce1, oronce2          | 3, 8, -, -         |
      | jane      | oronce3                   | -, -, 23, -        |
      | mary      | oronce2, oronce3          | -, 8, 23, -        |
      | christine | oronce1, oronce2, oronce3 | 3, 8, 23, 23       |
    When the current time is "<current_time>"
    And the game "offer-request-game" requests offers for player "<player_id>" in "<placement>"
    Then an offer with name "<offer>" is returned

    Examples:
      | current_time | player_id | offer   | placement |
      | 3            | joseph    | oronce1 | popup     |
      | 8            | joseph    | oronce2 | popup     |
      | 18           | joseph    | -       | popup     |
      | 23           | joseph    | oronce3 | popup     |
      | 23           | joseph    | oronce4 | store     |
      | 3            | john      | -       | popup     |
      | 8            | john      | oronce2 | popup     |
      | 18           | john      | -       | popup     |
      | 23           | john      | oronce3 | popup     |
      | 3            | michael   | -       | popup     |
      | 8            | michael   | -       | popup     |
      | 18           | michael   | -       | popup     |
      | 23           | michael   | oronce3 | popup     |
      | 3            | jane      | oronce1 | popup     |
      | 8            | jane      | oronce2 | popup     |
      | 18           | jane      | -       | popup     |
      | 23           | jane      | -       | popup     |
      | 23           | jane      | oronce4 | store     |
      | 3            | mary      | oronce1 | popup     |
      | 8            | mary      | -       | popup     |
      | 18           | mary      | -       | popup     |
      | 23           | mary      | -       | popup     |
      | 23           | mary      | oronce4 | store     |
      | 3            | christine | -       | popup     |
      | 8            | christine | -       | popup     |
      | 18           | christine | -       | popup     |
      | 23           | christine | -       | popup     |
      | 23           | christine | -       | store     |
      | 83           | christine | oronce4 | store     |


  Scenario Outline: After a player has seen an offer there is a track of it
    Given the following offer templates exist in the "offer-request-game" game:
      | game | name    | product_id     | contents   | placement | period       | frequency          | trigger                 |
      | org  | otosee1 | com.tfg.sample | { "x": 1 } | popup     | { "max": 1 } | { "every": "60s" } | { "from": 0, "to": 5 }  |
      | org  | otosee2 | com.tfg.sample | { "x": 2 } | popup     | { "max": 1 } | { "every": "60s" } | { "from": 6, "to": 10 } |
    When the current time is "<current_time>"
    And the game "offer-request-game" requests offers for player "<player_id>" in "popup"
    And the player "<player_id>" of game "offer-request-game" sees offer with name "<seen_offer>"
    Then player "<player_id>" of game "offer-request-game" has seen offer "<seen_offer>"
    And player "<player_id>" of game "offer-request-game" has not seen offer "<unseen_offer>"

    Examples:
      | current_time | player_id | seen_offer | unseen_offer |
      | 3            | jack      | otosee1    | otosee2      |
      | 8            | jenniffer | otosee2    | otosee1      |

  Scenario: When a player sees offers subsequently they are both tracked
    Given the following offer templates exist in the "offer-request-game" game:
      | game | name   | product_id     | contents   | placement | period       | frequency          | trigger                 |
      | org  | oseen1 | com.tfg.sample | { "x": 1 } | popup     | { "max": 1 } | { "every": "60s" } | { "from": 0, "to": 5 }  |
      | org  | oseen2 | com.tfg.sample | { "x": 2 } | popup     | { "max": 1 } | { "every": "60s" } | { "from": 6, "to": 10 } |
    When the current time is 3
    And the game "offer-request-game" requests offers for player "Henry" in "popup"
    And the player "Henry" of game "offer-request-game" sees offer in "popup"
    And the current time is 8
    And the game "offer-request-game" requests offers for player "Henry" in "popup"
    And the player "Henry" of game "offer-request-game" sees offer in "popup"
    Then player "Henry" of game "offer-request-game" has seen offer "oseen1"
    And player "Henry" of game "offer-request-game" has seen offer "oseen2"
