Feature: Claiming an offer for a given player

  Background:
    Given the server is up
    And a game with name "offer-claim-game" exists

  Scenario Outline: A player claims an offer
    Given the following offer templates exist in the "offer-claim-game" game:
      | game | name    | product_id     | contents  | placement | period       | frequency          | trigger                  |
      | org  | oclaim1 | com.tfg.sample | { 'x': 1 }| popup     | { 'max': 1 } | { 'every': '50s' } | { 'from': 0, 'to': 5 }   |
      | org  | oclaim2 | com.tfg.sample | { 'x': 2 }| popup     | { 'max': 1 } | { 'every': '50s' } | { 'from': 20, 'to': 25 } |
    And the following players exist in the "offer-claim-game" game:
      | id     | claimed-offers | last-seen-offer-at |
      | joseph | -              | -                  |
      | jane   | -              | -                  |
      | mary   | -              | -                  |
      | jack   | -              | -                  |
      | john   | oclaim1        | 1                  |
    When the current time is "<time>"
    And the game "offer-claim-game" requests offers for player "<player>" in "<placement>"
    And player "<player>" claims offer "<offer>" in game "offer-claim-game"
    Then the last request returned status code "<code>" and body "<body>"

    Examples:
      | time | player | offer    | placement | code | body                                                                                     |
      | 3    | joseph | oclaim1  | popup     | 200  | { 'x': 1 }                                                                               |
      | 6    | jane   | oclaim1  | popup     | 200  | { 'x': 1 }                                                                               |
      | 7d   | mary   | oclaim1  | popup     | 200  | { 'x': 1 }                                                                               |
      | 0    | joseph | oclaim2  | popup     | 200  | { 'x': 2 }                                                                               |
      | 3    | john   | oclaim1  | popup     | ?    | { 'code': 'OFF-05', 'description': 'Offer oclaim1 has already been claimed by player.' } |
      | 3    | jack   | oclaim1  | popup     | 200  | { 'x': 1 }                                                                               |
      | 3    | jack   | oclaim15 | popup     | 404  | { 'code': 'OFF-06', 'description': 'Offer oclaim15 was not found.' }                     |
