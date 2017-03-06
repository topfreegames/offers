Feature: Claiming an offer for a given player

  Background:
    Given the server is up
    And a game with name "offer-claim-game" exists

  Scenario Outline: A player claims an offer
    Given the following offer templates exist in the "offer-claim-game" game:
      | game | name    | product_id     | contents   | placement | period       | frequency          | trigger                  | key                                  |
      | org  | oclaim1 | com.tfg.sample | { 'x': 1 } | popup     | { 'max': 1 } | { 'every': '50s' } | { 'from': 0, 'to': 5 }   | 6597e909-ee8e-4a0c-82dd-533568f86aa6 |
      | org  | oclaim2 | com.tfg.sample | { 'x': 2 } | popup     | { 'max': 1 } | { 'every': '50s' } | { 'from': 20, 'to': 25 } | dbe912cd-d248-44ae-af4b-8c6e4ef9ad8f |
    And the following players claimed in the "offer-claim-game" game:
      | id     | claimed-offers | last-seen-offer-at |
      | joseph | -              | -                  |
      | jane   | -              | -                  |
      | mary   | -              | -                  |
      | jack   | -              | -                  |
      | john   | oclaim1        | 1                  |
    When the current time is "<request_time>"
    And the game "offer-claim-game" requests offers for player "<player>" in "<placement>"
    And the current time is "<claim_time>"
    And player "<player>" claims offer "<offer>" in game "offer-claim-game"
    Then the last request returned status code "<code>" and body "<body>"

    Examples:
      | request_time | claim_time | player | offer    | placement | code | body       |
      | 0            | 3          | joseph | oclaim1  | popup     | 200  | { 'x': 1 } |
      | 1            | 6          | jane   | oclaim1  | popup     | 200  | { 'x': 1 } |
      | 2            | 7d         | mary   | oclaim1  | popup     | 200  | { 'x': 1 } |
      | 21           | 0          | joseph | oclaim2  | popup     | 200  | { 'x': 2 } |
      | 4            | 3          | john   | oclaim1  | popup     | 409  | { 'x': 1 } |
      | 5            | 3          | jack   | oclaim1  | popup     | 200  | { 'x': 1 } |
      | 0            | 3          | jack   | oclaim15 | popup     | 301  |            |
