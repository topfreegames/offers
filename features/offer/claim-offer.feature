Feature: Claiming an offer for a given player

  Background:
    Given the server is up
    And a game with name "offer-claim-game" exists

  Scenario Outline: A player requests offers that can be bought only once
    Given the following offer templates exist in the "offer-claim-game" game:
      | game | name | product_id     | contents   | period             | frequency       | trigger                  |
      | org  | oc1  | com.tfg.sample | { "x": 1 } | { "type": "once" } | { "every": 50 } | { "from": 0, "to": 5 }   |
      | org  | oc2  | com.tfg.sample | { "x": 2 } | { "type": "once" } | { "every": 50 } | { "from": 20, "to": 25 } |
    And the following players exist in the "offer-claim-game" game:
      | id        | claimed-offers |
      | joseph    |                |
      | john      |            oc1 |
    When the current time is <time>
    And player "<player>" claims offer "<offer>" in game "offer-claim-game"
    Then the last request returned status code <status> and body "<body>"

    Examples:
      | time | player | offer | code | body                                                                            |
      | 3    | joseph | oc1   | 200  | { 'claimedAt': 3, 'contents': { 'x': 1 } }                                      |
      | 6    | joseph | oc1   | 200  | { 'claimedAt': 6, 'contents': { 'x': 1 } }                                      |
      | 7d   | joseph | oc1   | 200  | { 'claimedAt': 6, 'contents': { 'x': 1 } }                                      |
      | 0    | joseph | oc2   | 200  | { 'claimedAt': 0, 'contents': { 'x': 2 } }                                      |
      | 3    | john   | oc1   | 409  | { 'code': 'OFF-05', 'reason': 'Offer oc1 has already been claimed by player.' } |
      | 3    | joseph | oc1   | 200  | { 'claimedAt': 3, 'contents': { 'x': 1 } }                                      |
