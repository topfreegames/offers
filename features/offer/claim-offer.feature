Feature: Claiming an offer for a given player

  Background:
    Given the server is up
    And a game with id "offer-claim-game" exists

  Scenario Outline: A player claims an offer
    Given the following offers exist in the "offer-claim-game" game:
      | game | name    | product_id     | contents   | placement | period       | frequency         | trigger                  |
      | org  | oclaim1 | com.tfg.sample | { 'x': 1 } | popup     | { 'max': 1 } | { 'every': '1s' } | { 'from': 0, 'to': 5 }   |
      | org  | oclaim2 | com.tfg.sample | { 'x': 2 } | popup     | { 'max': 1 } | { 'every': '1s' } | { 'from': 20, 'to': 25 } |
    And the following players claimed in the "offer-claim-game" game:
      | id     | claimed-offers | last-seen-offer-at |
      | joseph | -              | -                  |
      | jane   | -              | -                  |
      | mary   | -              | -                  |
      | jack   | -              | -                  |
      | john   | oclaim1        | 1                  |
    When the current time is "<request_time>"
    And the game "offer-claim-game" requests offer instances for player "<player>" in "<placement>"
    And the current time is "<claim_time>"
    And player "<player>" claims offer instance "<offer>" in game "offer-claim-game"
    Then the last request returned status code "<code>" and body "<body>"

    Examples:
      | request_time | claim_time | player | offer    | placement | code | body                                                                                                                                                                                                               |
      | 1            | 3          | joseph | oclaim1  | popup     | 200  | {'contents':{'x':1}}                                                                                                                                                                                               |
      | 2            | 6          | jane   | oclaim1  | popup     | 200  | {'contents':{'x':1}}                                                                                                                                                                                               |
      | 3            | 7d         | mary   | oclaim1  | popup     | 200  | {'contents':{'x':1}}                                                                                                                                                                                               |
      | 21           | 1          | joseph | oclaim2  | popup     | 200  | {'contents':{'x':2}}                                                                                                                                                                                               |
      | 4            | 3          | john   | oclaim1  | popup     | 200  | {'contents':{'x':1}}                                                                                                                                                                                               |
      | 5            | 3          | jack   | oclaim1  | popup     | 200  | {'contents':{'x':1}}                                                                                                                                                                                               |
      | 1            | 3          | jack   | oclaim15 | popup     | 404  | {'code':'OFF-001','description':'offerInstance was not found with specified filters.','error':'offerInstanceNotFoundError','filters':{'GameID':'offer-claim-game','PlayerID':'jack','ProductID':'com.tfg.sample'}} |
