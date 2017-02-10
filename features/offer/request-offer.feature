Feature: Request offers for a player

  Background:
    Given the server is up
    And a game with name "offer-request-game" exists

  Scenario Outline: A player requests offers that can be bought only once
    Given the following offer templates exist in the "offer-request-game" game:
      | game | name | product_id     | contents   | metadata   | period             | trigger                  |
      | org  | or1  | com.tfg.sample | { "x": 1 } | { "y": 4 } | { "type": "once" } | { "from": 0, "to": 5 }   |
      | org  | or2  | com.tfg.sample | { "x": 2 } | { "y": 5 } | { "type": "once" } | { "from": 6, "to": 10 }  |
      | org  | or3  | com.tfg.sample | { "x": 3 } | { "y": 6 } | { "type": "once" } | { "from": 20, "to": 25 } |
    And the following players exist in the "offer-request-game" game:
      | id        | claimed-offers |
      | joseph    |                |
      | john      |            or1 |
      | michael   |       or1, or2 |
      | jane      |            or3 |
      | mary      |       or2, or3 |
      | christine |  or1, or2, or3 |
    When the current time is <current_time>
    And the game "offer-request-game" requests offers for player "<player_id>"
    Then an offer with name "<offer_name>" is returned

    Examples:
      | current_time | player_id | offer_name |
      |            3 | joseph    | or1        |
      |            8 | joseph    | or2        |
      |           18 | joseph    |            |
      |           23 | joseph    | or3        |
      |            3 | john      |            |
      |            8 | john      | or2        |
      |           18 | john      |            |
      |           23 | john      | or3        |
      |            3 | michael   |            |
      |            8 | michael   |            |
      |           18 | michael   |            |
      |           23 | michael   | or3        |
      |            3 | jane      | or1        |
      |            8 | jane      | or2        |
      |           18 | jane      |            |
      |           23 | jane      |            |
      |            3 | mary      | or1        |
      |            8 | mary      |            |
      |           18 | mary      |            |
      |           23 | mary      |            |
      |            3 | christine |            |
      |            8 | christine |            |
      |           18 | christine |            |
      |           23 | christine |            |
