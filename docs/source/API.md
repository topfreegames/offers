Offers API
==========

## Healthcheck Routes

  ### Healthcheck

  `GET /healthcheck`

  Validates if the app is still up, including its database connection.

  * Success Response
    * Code: `200`
    * Content:
      ```
        {
          "healthy": true
        }
      ```
  * Error Response

    It will return an internal error if it failed to connect to the database.

    * Code: `500`
    * Content:
      ```
        {
          "healthy": false
        }
      ```

## Game Routes

  ### Upsert Game
  `PUT /games`

  Update an existing Game or insert a new Game into database.

  * Payload
    ```
    {
      "id":       [string], // required, matches ^[^-][a-z0-9-]*$
      "name":     [string], // required, 255 characters max
      "bundleId": [string], // required, 255 characters max
      "metadata": [json],   // optional
    }
    ```

  | Parameter   | Description                                     |
  | ----------- | -------------                                   |
  | id          | Unique ID that identifies the game              |
  | name        | Prettier game identifier to show on UI          |
  | bundleId    | App identifier on PlayStore or AppStore         |
  | metadata    | Any information the Front wants to access later |

  * Success Response
    * Code: `200`
    * Content:
      ```
        {
          "gameId": [string]
        }
      ```

  * Error response

    It will return an error if the query on db (upsert) failed

    * Code: `500`
    * Content:
      ```
        {
          "reason": [string]
        }
      ```

## Offer Template Routes

  ### Insert Offer Template
  `PUT /offer-template`

  Insert a new Offer Template into database.

  * Payload
    ```
      {
        "id":        [uuidv4], // required
        "name":      [string], // required, 255 characters max
        "productId": [string], // required, 255 characters max
        "gameId":    [string], // required, matches ^[^-][a-z0-9-]*$
        "contents":  [json],   // required
        "metadata":  [json],   // optional
        "enabled":   [bool],   // optional
        "placement": [string], // required, 255 characters max
        "period":    {         // required
          "every": [string],   // required
          "max":   [int]       // required
        },   
        "frequency": {         // required
          "every": [string],   // required
          "max":   [int]       // required
        },   
        "trigger":   {         // required
          "from":  [int],      // required
          "to":    [int]       // required
        }
      }
    ```

  | Parameter   | Description                                                                                                                                                                                                                                                                                               |
  | ----------- | -------------                                                                                                                                                                                                                                                                                             |
  | id          | Unique ID that identifies the offer template                                                                                                                                                                                                                                                              |
  | name        | Prettier game identifier to show on UI                                                                                                                                                                                                                                                                    |
  | productId   | Identifier of the item to be bought on PlayStore or AppStore                                                                                                                                                                                                                                              |
  | gameId      | ID of the game this template was made for (must exist on Games table on DB)                                                                                                                                                                                                                               |
  | contents    | What the offer provides (ex.: { "gem": 5, "gold": 100 })                                                                                                                                                                                                                                                  |
  | metadata    | Any information the Front wants to access later                                                                                                                                                                                                                                                           |
  | period      | Enable player to buy offer every x times, at most y times. <ul><li>every: decimal number with unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"</li><li>max: maximum number of times this offer can be bought by the player</li></ul>If "every" is an empty string, then the offer can be bought max times with no time restriction.  If "max" is 0, then the offer can be bought infinite times with time restriction.  They can`t be "" and 0 at the same time.      |
  | frequency   | Enable player to see offer on UI x/unit of time, at most y times. <ul><li>every: decimal number with unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"</li><li>max: maximum number of times this offer can be seen by the player</li></ul>If "every" is an empty string, then the offer can be seen max times with no time restriction.  If "max" is 0, then the offer can be seen infinite times with time restriction.  They can`t be "" and 0 at the same time. |
  | trigger     | Time when the offer is available                                                                                                                                                                                                                                                                          |
  | enabled     | True if the offer is enabled                                                                                                                                                                                                                                                                              |
  | placement   | Where the offer is shown in the UI                                                                                                                                                                                                                                                                        |

  * Success Response
    * Code: `200`
    * Content:
      ```
        {
          "OfferTemplateId": [string]
        }
      ```

  * Error response

    It will return an error if the query on db (insert) failed

    * Code: `500`
    * Content:
      ```
        {
          "reason": [string]
        }
      ```

## Offer Routes

  ### Get Available Offers
  `GET /offers?player-id=<required-player-id>&game-id=<required-game-id>`

  Get the available offers for a player of a game. An offer is available if it respects the frequency (last time player saw the offer), respects the period (last time player claimed the offer), is triggered (current time is between "from" and "to") and is enabled. The success response is a JSON where each key is a placement on the UI and the value is a list of OfferTemplates.

  * Success Response
    * Code: `200`
    * Content:
      ```
        {
          "placement-1": [
            {
                "id":        [uuidv4], // required
                "name":      [string], // required, 255 characters max
                "productId": [string], // required, 255 characters max
                "gameId":    [string], // required, matches ^[^-][a-z0-9-]*$
                "contents":  [json],   // required
                "metadata":  [json],   // optional
                "enabled":   [bool],   // optional
                "placement": [string], // required, 255 characters max
                "period":    {         // required
                  "every": [string],   // required
                  "max":   [int]       // required
                },   
                "frequency": {         // required
                  "every": [string],   // required
                  "max":   [int]       // required
                },   
                "trigger":   {         // required
                  "from":  [int],      // required
                  "to":    [int]       // required
                }
            },
            ...
          ]
          "placement-2": [
            ...
          ],
          ...
        }
      ```

  * Error Response
    * Code: `400`, if player-id is not informed
    * Code: `400`, if game-id is not informed
    * Code: `500`, if server failed in any other way
    * Content:
      ```
        {
          "reason": [string]
        }
      ```

  ### Claim Offer
  `PUT /offer/claim`

  Claim a player's offer. Should only be called after payment confirmation.

  * Payload
    ```
      {
        "id":       [uuidv4], // required
        "gameId":   [string], // required, matches ^[^-][a-z0-9-]*$
        "playerId": [string]  // required, 255 characters max
      }
    ```

  * Success Response
    * Code: `200`
    * Content: JSON with the offer contents, got from the key contents in OfferTemplate

  * Error Response

    * If a offer with id, gameId and playerId was not found in database.
      * Code: `404`
      * Content:
        ```
          {
            "reason": [string]
          }
        ```

    * If any internal error occurred.
      * Code: `500`
      * Content:
        ```
          {
            "reason": [string]
          }
        ```

    * If the player claimed an offer before the time defined in OfferTemplate["period"]["every"] or more than what is defined in OfferTemplate["period"]["max"], the offer is not claimed anymore but its contents are returned.
      * Code: `409`
      * Content:
        ```
          {
            "reason":   [string],
            "contents": [json]
          }
        ```

  ### Update Offer Last Seen At
  `PUT /offer/last-seen-at`

  Updates the time when the offer was last seen in the UI by the player

  * Payload
    ```
      {
        "id":       [uuidv4], // required
        "gameId":   [string], // required, matches ^[^-][a-z0-9-]*$
        "playerId": [string]  // required, 255 characters max
      }
    ```

  * Success Response
    * Code: `200`

  * Error Response
    * If a offer with id, gameId and playerId was not found in database.
      * Code: `404`
      * Content:
        ```
          {
            "reason": [string]
          }
        ```

    * If any internal error occurred.
      * Code: `500`
      * Content:
        ```
          {
            "reason": [string]
          }
        ```