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
          "error": "DatabaseError",
          "code": "OFF-000",
          "description": [string]    // error description
        }
      ```

## Game Routes
  ### Upsert Game
  `PUT /games/:id`

  Update an existing Game or insert a new Game into database. `:id` must match `^[^-][a-z0-9-]*$`.

  * Payload
    ```
    {
      "name":     [string], // required, 255 characters max
      "metadata": [json]    // optional
    }
    ```
    * Field Descriptions
      - **id**:       Unique ID that identifies the game              
      - **name**:     Prettier game identifier to show on UI          
      - **metadata**: Any additional information one would like to access later

  * Success Response
    * Code: `200`
    * Content:
      ```
        {
          "gameId": [string]
        }
      ```

  * Error response
    * If missing or invalid arguments
      * Code: `422`
      * Content:
          ```
            {
              "error": [string],       // error
              "code":  [string],       // error code
              "description": [string]  // error description
            }
        ```
    * It will return an error if the query on db (upsert) failed
      * Code: `500`
      * Content:
        ```
          {
            "error": [string],       // error
            "code":  [string],       // error code
            "description": [string]  // error description
          }
        ```

  ### List Games
  `GET /games`

  List all existing games.

  * Success Response
    * Code: `200`
    * Content:

    ```
    [    
      {
        "id":       [string],
        "name":     [string],
        "metadata": [json]
      },
      ...
    ]
    ```

  * Error response

    It will return an error if the query on db failed

    * Code: `500`
    * Content:
      ```
        {
          "error": [string],       // error
          "code":  [string],       // error code
          "description": [string]  // error description
        }
      ```

## Offer Template Routes

  ### Insert Offer Template
  `POST /templates`

  Insert a new Offer Template into database.

  * Payload
    ```
      {
        "name":      [string], // required, 255 characters max
        "key":       [uuidv4], // required
        "productId": [string], // required, 255 characters max
        "gameId":    [string], // required, matches ^[^-][a-z0-9-]*$
        "contents":  [json],   // required
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
        },
        "metadata":  [json]    // optional
      }
    ```

    * Field Descriptions
       - **name**:         Prettier game identifier to show on UI.  
       - **key**:          Identifies an offer template. It is common between the templates versions, meaning it keeps the same when an offer is updated.
       - **productId**:    Identifier of the item to be bought on PlayStore or AppStore.  
       - **gameId**:       ID of the game this template was made for (must exist on Games table on DB).  
       - **contents**:     What the offer provides (ex.: { "gem": 5, "gold": 100 }).  
       - **metadata**:     Any information the Front wants to access later.  
       - **period**:       Enable player to buy offer every x times, at most y times. <ul><li>every: decimal number with unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"</li><li>max: maximum number of times this offer can be bought by the player</li></ul>If "every" is an empty string, then the offer can be bought max times with no time restriction.  If "max" is 0, then the offer can be bought infinite times with time restriction.  They can't be "" and 0 at the same time.
       - **frequency**:    Enable player to see offer on UI x/unit of time, at most y times. <ul><li>every: decimal number with unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"</li><li>max: maximum number of times this offer can be seen by the player</li></ul>If "every" is an empty string, then the offer can be seen max times with no time restriction.  If "max" is 0, then the offer can be seen infinite times with time restriction.  They can't be "" and 0 at the same time.  
       - **trigger**:      Time when the offer is available.  
       - **enabled**:      True if the offer is enabled.  
       - **placement**:    Where the offer is shown in the UI.  

  * Success Response
    * Code: `200`
    * Content:
      ```
        {
          "id":        [uuidv4],   // offer template unique identifier
          "key":       [uuidv4],
          "name":      [string],
          "productId": [string],
          "gameId":    [string],
          "contents":  [json],  
          "metadata":  [json],  
          "placement": [string],
          "period":    {        
            "every": [string],  
            "max":   [int]      
          },   
          "frequency": {        
            "every": [string],  
            "max":   [int]      
          },   
          "trigger":   {        
            "from":  [int],     
            "to":    [int]      
          },
          "enabled": true     // created templates are enabled by default
        }
      ```

  * Error response

    It will return an error if the request has missing or invalid arguments

    * Code: `422`
    * Content:
      ```
        {
          "error": [string],       // error
          "code":  [string],       // error code
          "description": [string]  // error description
        }
      ```

    It will return an error if the query on db (insert) failed

    * Code: `500`
    * Content:
      ```
        {
          "error": [string],       // error
          "code":  [string],       // error code
          "description": [string]  // error description
        }
      ```

  ### Enable offer template
  `PUT /templates/:id/enable`

  Enable an offer template. `:id` must be an `uuidv4`.

  * Success Response
    * Code: `200`
    * Content:
      ```
        {
          "id": [uuidv4]
        }
      ```

  * Error Response

    It will return status code 404 not found if the ID doesn't exist

    * Code: `404`
    * Content:
      ```
        {
          "error": [string],       // error
          "code":  [string],       // error code
          "description": [string]  // error description
        }
      ```

    It will return status code 500 internal error occurred

    * Code: `500`
    * Content:
      ```
        {
          "error": [string],       // error
          "code":  [string],       // error code
          "description": [string]  // error description
        }
      ```

  ### Disable offer template
  `PUT /templates/:id/disable`

  Disable an offer template. `:id` must be an `uuidv4`.

  * Success Response
    * Code: `200`
    * Content:
      ```
        {
          "id": [uuidv4]
        }
      ```

  * Error Response

    It will return status code 404 not found if the ID doesn't exist

    * Code: `404`
    * Content:
      ```
        {
          "error": [string],       // error
          "code":  [string],       // error code
          "description": [string]  // error description
        }
      ```

    It will return status code 500 internal error occurred

    * Code: `500`
    * Content:
      ```
        {
          "error": [string],       // error
          "code":  [string],       // error code
          "description": [string]  // error description
        }
      ```

  ### List Offer Templates
  `GET /templates?game-id=<required-game-id>`

  List all game's offer templates.

  * Success Response
    * Code: `200`
    * Content:

    ```
    [    
      {
        "id":        [uuidv4],   // offer template unique identifier
        "key":       [uuidv4],
        "name":      [string],
        "productId": [string],
        "gameId":    [string],
        "contents":  [json],  
        "metadata":  [json],  
        "placement": [string],
        "period":    {        
          "every": [string],  
          "max":   [int]      
        },   
        "frequency": {        
          "every": [string],  
          "max":   [int]      
        },   
        "trigger":   {        
          "from":  [int],     
          "to":    [int]      
        },
        "enabled": [bool]
      },
      ...
    ]
    ```

  * Error response

    It will return an error if the query on db failed

    * Code: `500`
    * Content:
      ```
        {
          "error": [string],       // error
          "code":  [string],       // error code
          "description": [string]  // error description
        }
      ```

## Offer Routes

  ### Get Available Offers
  `GET /offers?player-id=<required-player-id>&game-id=<required-game-id>`

  Get the available offers for a player of a game. An offer is available if it respects the frequency (last time player saw the offer), respects the period (last time player claimed the offer), is triggered (current time is between "from" and "to") and is enabled. The success response is a JSON where each key is a placement on the UI and the value is a list of available offers.

  * Success Response
    * Code: `200`
    * Content:
      ```
        {
          "placement-1": [
            {
                "id":                   [uuidv4], // offer id
                "productId":            [string], // required, 255 characters max
                "contents":             [json],   // offer contents as registered in the offer template
                "metadata":             [json],   // offer metadata as registered in the offer template
                "remainingPurchases":   [int],    // if the template has a max period, how many purchases are still available for this offer
                "remainingImpressions": [int],    // if the template has a max frequency, how many purchases are still available for this offer
                "expireAt":             [int64]   // timestamp (seconds since epoch) until when the offer is valid
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
        "error": [string],       // error
        "code":  [string],       // error code
        "description": [string]  // error description
      }
      ```

  ### Claim Offer
  `PUT /offers/:id/claim`

  Claim a player's offer. Should only be called after payment confirmation. `:id` must be an `uuidv4`.

  * Payload
    ```
      {
        "gameId":   [string], // required, matches ^[^-][a-z0-9-]*$
        "playerId": [string]  // required, 255 characters max
      }
    ```

  * Success Response (if the player can still see the offer)
    * Code: `200`
    * Content:
      ```
        {
          "contents": [json],
          "nextAt": [int64]  // unix timestamp of the next time the offer can be shown
        }
      ```

  * Success Response (if the player can no longer see the offer)
    * Code: `200`
    * Content:
      ```
        {
          "contents": [json]
        }
      ```

  * If the player claimed an offer that he already claimed its contents are returned but with another status code, so the caller can decide whether to give the player the offer contents or not.
    * Code: `409`
    * Content:
      ```
        {
          "contents": [json]
          "nextAt": [int64]  // unix timestamp of the next time the offer can be shown
        }
      ```

  * Error Response
    * If a offer with id, gameId and playerId was not found in database.
      * Code: `404`
      * Content:
        ```
        {
          "error": [string],       // error
          "code":  [string],       // error code
          "description": [string]  // error description
        }
        ```

    * If any internal error occurred.
      * Code: `500`
      * Content:
        ```
          {
            "error": [string],       // error
            "code":  [string],       // error code
            "description": [string]  // error description
          }
        ```

  ### Offer Impressions
  `POST /offers/:id/impressions`

  Updates the time when the offer was last seen by the player. `:id` must be an `uuidv4`.

  * Payload
    ```
      {
        "gameId":   [string], // required, matches ^[^-][a-z0-9-]*$
        "playerId": [string]  // required, 255 characters max
      }
    ```

  * Success Response (if the player can still see the offer)
    * Code: `200`
    * Content:
      ```
        {
          "nextAt": [int64]  // unix timestamp of the next time the offer can be shown
        }
      ```

  * Success Response (if the player can no longer see the offer)
    * Code: `200`
    * Content:
      ```
        {}
      ```

  * Error Response
    * If missing or invalid arguments.
      * Code: `422`
      * Content:
        ```
          {
            "error": [string],       // error
            "code":  [string],       // error code
            "description": [string]  // error description
          }
        ```

    * If a offer with id, gameId and playerId was not found in database.
      * Code: `404`
      * Content:
        ```
          {
            "error": [string],       // error
            "code":  [string],       // error code
            "description": [string]  // error description
          }
        ```

    * If any internal error occurred.
      * Code: `500`
      * Content:
        ```
        {
          "error": [string],       // error
          "code":  [string],       // error code
          "description": [string]  // error description
        }
        ```
