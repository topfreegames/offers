Feature: Claiming an offer for a given player

  Background:
    Given the server is up
    And a game with name "offer-claim-game" exists

  Scenario Outline: A player requests offers that can be bought only once
    Given the following offer templates exist in the "offer-claim-game" game:
      | game | name | product_id     | contents   | metadata   | period             | trigger                  |
      | org  | oc1  | com.tfg.sample | { "x": 1 } | { "y": 4 } | { "type": "once" } | { "from": 0, "to": 5 }   |
      | org  | oc2  | com.tfg.sample | { "x": 2 } | { "y": 5 } | { "type": "once" } | { "from": 6, "to": 10 }  |
      | org  | oc3  | com.tfg.sample | { "x": 3 } | { "y": 6 } | { "type": "once" } | { "from": 20, "to": 25 } |
    And the following players exist in the "offer-claim-game" game:
      | id        | claimed-offers |
      | joseph    |                |
      | john      |            oc1 |
      | michael   |       oc1, oc2 |
      | jane      |            oc3 |
      | mary      |       oc2, oc3 |
      | christine |  oc1, oc2, oc3 |
    When the current time is <time>
    And player "<player>" claims offer "<offer>" in game "offer-claim-game"
    Then the last request returned status code "<code>" and body "<body>"

    Examples:
      | time | player | offer | code | body                                                                            |
      |    3 | joseph | oc1   |  200 | { "claimedAt": 3, "contents": { "x": 1 } }                                      |
      |    6 | joseph | oc1   |  200 | { "claimedAt": 6, "contents": { "x": 1 } }                                      |
      | 9999 | joseph | oc1   |  200 | { "claimedAt": 6, "contents": { "x": 1 } }                                      |
      |    0 | joseph | oc3   |    ? | ?                                                                               |
      |    3 | john   | oc1   |  409 | { "code": "OFF-05", "reason": "Offer oc1 has already been claimed by player." } |
      |    3 | joseph | oc1   |  200 | { "claimedAt": 3, "contents": { "x": 1 } }                                      |
