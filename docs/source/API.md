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

  Updates an existing Game or insert a new Game into database. `:id` must match `^[^-][a-zA-Z0-9-_]*$`.

  **Requires basic auth**.

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

  Lists all existing games.

  **Requires basic auth**.

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

## Offer Routes

  ### Create Offer
  `POST /offers`

  Inserts a new Offer into the database.

  **Requires basic auth**.

  * Payload
    ```
      {
        "name":      [string], // required, 255 characters max
        "productId": [string], // required, 255 characters max
        "gameId":    [string], // required, matches ^[^-][a-zA-Z0-9-_]*$
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
        "metadata":  [json],   // optional
        "filters":   [json]    // optional
      }
    ```

    * Field Descriptions
       - **name**:         Prettier game identifier to show on UI.  
       - **productId**:    Identifier of the item to be bought on PlayStore or AppStore.  
       - **gameId**:       ID of the game this template was made for (must exist on Games table on DB).  
       - **contents**:     What the offer provides (ex.: { "gem": 5, "gold": 100 }).  
       - **metadata**:     Any information the Front wants to access later.  
       - **period**:       Enable player to buy offer every x times, at most y times. <ul><li>every: decimal number with unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"</li><li>max: maximum number of times this offer can be bought by the player</li></ul>If "every" is an empty string, then the offer can be bought max times with no time restriction.  If "max" is 0, then the offer can be bought infinite times with time restriction.  They can't be "" and 0 at the same time.
       - **frequency**:    Enable player to see offer on UI x/unit of time, at most y times. <ul><li>every: decimal number with unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"</li><li>max: maximum number of times this offer can be seen by the player</li></ul>If "every" is an empty string, then the offer can be seen max times with no time restriction.  If "max" is 0, then the offer can be seen infinite times with time restriction.  They can't be "" and 0 at the same time.  
       - **trigger**:      Time when the offer is available.  
       - **filters**:      The filters for the offer, they can be of three different types for a given attribute: <ul><li>interval: the attribute must define the beginning and/or end of the interval with "geq" and "lt", the interval includes the beginning but not the end</li><li>equality: the attribute must define the "eq", the value that the filter expects the attribute to be equal to, it should be a string</li><li>difference: the attribute must define "neq", the value that the filter expects the attribute to be different from, it should be a string</ul>An example: "{ "intervalValue": { "geq": 0.0, "lt": 10.0 }, "equalValue": { "eq": "John" } }".  
       - **enabled**:      True if the offer is enabled.  
       - **placement**:    Where the offer is shown in the UI.  
       - **version**:      Offer current version.

  * Success Response
    * Code: `200`
    * Content:
      ```
        {
          "id":        [uuidv4],   // offer unique identifier
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
          "enabled":   true,       // created offer are enabled by default
          "version":   1,          // created offer version is 1 by default
          "filters":   [json]
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

  ### Update Offer
  `PUT /offers/:id`

  Updates the offer with given id in the database.

  **Requires basic auth**.

  * Payload
    ```
      {
        "name":      [string], // required, 255 characters max
        "productId": [string], // required, 255 characters max
        "gameId":    [string], // required, matches ^[^-][a-zA-Z0-9-_]*$
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
        "metadata":  [json],   // optional
        "filters":   [json]    // optional
      }
    ```

    * Field Descriptions
       - **name**:         Prettier game identifier to show on UI.  
       - **productId**:    Identifier of the item to be bought on PlayStore or AppStore.  
       - **gameId**:       ID of the game this template was made for (must exist on Games table on DB).  
       - **contents**:     What the offer provides (ex.: { "gem": 5, "gold": 100 }).  
       - **metadata**:     Any information the Front wants to access later.  
       - **period**:       Enable player to buy offer every x times, at most y times. <ul><li>every: decimal number with unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"</li><li>max: maximum number of times this offer can be bought by the player</li></ul>If "every" is an empty string, then the offer can be bought max times with no time restriction.  If "max" is 0, then the offer can be bought infinite times with time restriction.  They can't be "" and 0 at the same time.
       - **frequency**:    Enable player to see offer on UI x/unit of time, at most y times. <ul><li>every: decimal number with unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h"</li><li>max: maximum number of times this offer can be seen by the player</li></ul>If "every" is an empty string, then the offer can be seen max times with no time restriction.  If "max" is 0, then the offer can be seen infinite times with time restriction.  They can't be "" and 0 at the same time.  
       - **trigger**:      Time when the offer is available.  
       - **filters**:      The filters for the offer, they can be of three different types for a given attribute: <ul><li>interval: the attribute must define the beginning and/or end of the interval with "geq" and "lt", the interval includes the beginning but not the end</li><li>equality: the attribute must define the "eq", the value that the filter expects the attribute to be equal to, it should be a string</li><li>difference: the attribute must define "neq", the value that the filter expects the attribute to be different from, it should be a string</ul>An example: "{ "intervalValue": { "geq": 0.0, "lt": 10.0 }, "equalValue": { "eq": "John" } }".  
       - **enabled**:      True if the offer is enabled.  
       - **placement**:    Where the offer is shown in the UI.  
       - **version**:      Offer current version.

  * Success Response
    * Code: `200`
    * Content:
      ```
        {
          "id":      [uuidv4],   // offer unique identifier
          "version": [int]       // updated offer version
        }
      ```

  * Error response

    It will return an error if the offer with given id does not exist in the database

    * Code: `404`
    * Content:
      ```
        {
          "error": [string],       // error
          "code":  [string],       // error code
          "description": [string]  // error description
        }
      ```

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


  ### Enable offer
  `PUT /offers/:id/enable?game-id=<required-game-id>`

  Enables an offer. `:id` must be an `uuidv4`.

  **Requires basic auth**.

  * Success Response
    * Code: `200`
    * Content:
      ```
        {}
      ```

  * Error Response

    It will return status code 404 if the offer with given ID does not exist

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
  `PUT /offers/:id/disable?game-id=<required-game-id>`

  Disables an offer template. `:id` must be an `uuidv4`.

  **Requires basic auth**.

  * Success Response
    * Code: `200`
    * Content:
      ```
        {}
      ```

  * Error Response

    It will return status code 404 if the offer with given ID does not exist

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

  ### List Offers
  `GET /offers?game-id=<required-game-id>`

  Lists all game's offers.

  **Requires basic auth**.

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
        "enabled":   [bool],
        "version":   [int],
        "filters":   [json]
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

## Offer Request Routes

  There are the routes accessed by the offers lib.

  ### Get Available Offers
  `GET /available-offers?player-id=<required-player-id>&game-id=<required-game-id>&<attr1>=<val1>&...`

  Gets the available offers for a player of a game. An offer is available if it respects the frequency (last time player saw the offer), respects the period (last time player claimed the offer), is triggered (current time is between "from" and "to"), matches the filters of the offer for the parameters sent in the query string  and is enabled. The success response is a JSON where each key is a placement on the UI and the value is a list of available offers.  
  If an attribute sent in the query string doesn't exist in a filter it is ignored and the extra parameters for a filter are ignored if the request doesn't send a value for them. If the filter defines an interval the query string parameter value must be a number. There is no limit in the amount of attributes that can be sent to be used in the filters.

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
    * Header:
      A max-age header is sent to indicate how long the response returned by get available offers can be cached.
      ```
      Cache-Control: max-age=<seconds>
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
  `PUT /offers/claim`

  Claims a player's offer. Should only be called after payment confirmation. `:id` must be an `uuidv4`.

  * Payload
    ```
      {
        "gameId":   [string],      // required, matches ^[^-][a-zA-Z0-9-_]*$
        "playerId": [string],      // required, 255 characters max
        "productId": [string],     // required, 255 characters max
        "timestamp": [int64],      // required, unix timestamp of the purchase
        "transactionId": [string], // required, unique identifier of the purchase
        "id": [uuidv4]             // optional, the id of the offer being claimed
      }
    ```

    If the id of the offer being claimed is sent it will be used to find the offer, increment the claim counter and the timestamp of the last time this offer was claimed. If not, the other information present in the payload will be used to try to identify the offer that is being claimed.

  * Success Response
    * Code: `200`
    * Content:
      if the player can still see the offer:
      ```
        {
          "contents": [json]
          "nextAt": [int64]  // unix timestamp of the next time the offer can be shown
        }
      ```

      if the player can still see the offer
      ```
        {
          "contents": [json]
        }
      ```


  * If the player claimed an offer that he already claimed its contents are returned but with another status code, so the caller can decide whether to give the player the offer contents or not.
    * Code: `409`
    * Content:
      if the player can still see the offer:
      ```
        {
          "contents": [json]
          "nextAt": [int64]  // unix timestamp of the next time the offer can be shown
        }
      ```

      if the player can still see the offer
      ```
        {
          "contents": [json]
        }

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
  `PUT /offers/:id/impressions`

  Updates the time when the offer was last seen by the player and increments an impressions counter. `:id` must be an `uuidv4`.

  * Payload
    ```
      {
        "gameId":   [string],   // required, matches ^[^-][a-zA-Z0-9-_]*$
        "playerId": [string],   // required, 255 characters max
        "impressionId" [uuidv4] // required, unique identifier for this impression
      }
    ```

    The `impressionId` field is used so this request can be idempotent. If more than one request is sent with the same `impressionId` the counter and the last impression timestamp will not be updated.

  * Success Response
    * Code: `200`
    * Content:
      if the player can still see the offer:
      ```
        {
          "nextAt": [int64]  // unix timestamp of the next time the offer can be shown
        }
      ```

      if the player can still see the offer
      ```
        {}
      ```

  * Conflict Response (if the `impressionId` was already sent in a previous request):
    * Code: `200`
    * Content:
      if the player can still see the offer:
      ```
        {
          "nextAt": [int64]  // unix timestamp of the next time the offer can be shown
        }
      ```

      if the player can no longer see the offer:
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
