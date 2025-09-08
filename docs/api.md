# Educational Game API Documentation

## Base URL

### Development

```
{{protocol}}://bermudia-api-internal.darkube.app/api/v1
```

protocols: `https`, `wss`

## Response Format

All API endpoints return responses in a consistent format:

```json
{
  "ok": true,
  "error": "string",
  "result": {}
}
```

- `ok`: Boolean indicating success/failure
- `error`: Error message (only present when ok=false)
- `result`: Response data (only present when ok=true)

## Authentication

Some endpoints are authenticated and their request is authenticated done via JWT.

User must [log in](#login) to receive a JWT.

The generated JWT must be but in the `Authorization` request header for authenticated HTTP endpoints.

For websocket endpoints, the JWT is received in `token` query parameter.

## Endpoints

### Login

Checks user credentials and generates a jwt token.

Receives a [LoginRequest](#loginrequest) in body.

Returns a [LoginResult](#loginresult) in response.

**Endpoint:** `POST /login`

---

### Get Me

_This endpoint **is authenticated** and needs an auth token for access._

Returns [Me](#me) in response.

This endpoint is just a way to test whether the JWT is correctly set and still valid.

**Endpoint:** `GET /me`

---

### Get Territory

Retrieves the complete data for a specific territory, including all islands and their connections.

Returns a [Territory](#territory) in response.

**Endpoint:** `GET /territories/{territoryID}`

**Parameters:**

- `territoryID` (path parameter, required): The unique identifier of the territory

---

### Get Island

_This endpoint **is authenticated** and needs an auth token for access._

Retrieves the content of the island.

Returns a [IslandContent](#islandcontent) in response.

**Preconditions:** The player must be at the island and also must have anchored.

**Endpoint:** `GET /islands/{islandID}`

**Parameters:**

- `islandID` (path parameter, required): The unique identifier of the island

---

### Submit Answer

_This endpoint **is authenticated** and needs an auth token for access._

Receives the input that the user enters into a [IslandInput](#islandinput) component.

Returns the updated [SubmissionState](#submissionstate) field of the [IslandInput](#islandinput) component in response.

**Preconditions:** The player must be at the corresponding island and also must have anchored.

**Endpoint:** `POST /answer/{inputID}`

**Parameters:**

- `inputID` (path parameter, required): The _id_ of the [IslandInput](#islandinput) component.
- `data` (body parameter, required): The user data. Its type depends on the _type_ field in [IslandInput](#islandinput); If _type_ is `file` , pass the file, otherwise pass the plain text in this field.

**Note:** Request's `Content-Type` must be `multipart/form-data`

```shell
curl --request POST \
  --url https://bermudia-api-internal.darkube.app/api/v1/answer/ans_29C12F3C7D089666 \
  --header 'Authorization: TOKEN' \
  --header 'Content-Type: multipart/form-data' \
  --form data=@/path/to/file.txt
```

---

### Stream Events

_This endpoint **is authenticated** and needs an auth token for access._

A **websocket** endpoint for receiving realtime events.

Type of messages is text; JSON encoding of [Event](#event).

**Endpoint:** `/events?token=TOKEN`

---

### Get Player

_This endpoint **is authenticated** and needs an auth token for access._

Returns the [Player](#player) object.

**Endpoint:** `POST /player`

---

### Travel Check

_This endpoint **is authenticated** and needs an auth token for access._

Checks whether the specified travel is possible.

Receives a [TravelCheckRequest](#travelcheckrequest) in body.

Returns [TravelCheckResult](#travelcheckresult) in response.

**Endpoint:** `POST /travel_check`

```shell
curl --request POST \
  --url https://bermudia-api-internal.darkube.app/api/v1/travel_check \
  --header 'Authorization: TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{"fromIsland": "island_final","toIsland": "island_math2"}'
```

---

### Travel

_This endpoint **is authenticated** and needs an auth token for access._

Makes the player travel from the current island to another island.

Receives a [TravelRequest](#travelrequest) in body.

Returns an empty object in response.

**Endpoint:** `POST /travel`

```shell
curl --request POST \
  --url https://bermudia-api-internal.darkube.app/api/v1/travel \
  --header 'Authorization: TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{"fromIsland": "island_final","toIsland": "island_math2"}'
```

---

### Refuel Check

_This endpoint **is authenticated** and needs an auth token for access._

Used to check whether refuel is possible in the current state.

Does not receive anything.

Returns a [RefuelCheckResult](#refuelcheckresult) in response.

**Endpoint:** `POST /refuel_check`

```shell
curl --request POST \
  --url https://bermudia-api-internal.darkube.app/api/v1/refuel_check \
  --header 'Authorization: TOKEN'
```

---

### Refuel

_This endpoint **is authenticated** and needs an auth token for access._

Used to refuel player's vehicle.

Receives [RefuelRequest](#refuelrequest) in body.

Returns an empty object in response.

**Endpoint:** `POST /refuel`

```shell
curl --request POST \
  --url https://bermudia-api-internal.darkube.app/api/v1/refuel \
  --header 'Authorization: TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{"amount": 5}'
```

---

### Anchor Check

_This endpoint **is authenticated** and needs an auth token for access._

Checks whether anchoring on the specified island is possible.

Receives a [AnchorCheckRequest](#anchorcheckrequest) in body.

Returns [AnchorCheckResult](#anchorcheckresult) in response.

**Endpoint:** `POST /anchor_check`

```shell
curl --request POST \
  --url https://bermudia-api-internal.darkube.app/api/v1/anchor_check \
  --header 'Authorization: TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{"island": "island_final"}'
```

---

### Anchor

_This endpoint **is authenticated** and needs an auth token for access._

Anchors the player on the specified island.

Receives a [AnchorRequest](#anchorrequest) in body.

Returns an empty object in response.

**Endpoint:** `POST /anchor`

```shell
curl --request POST \
  --url https://bermudia-api-internal.darkube.app/api/v1/anchor \
  --header 'Authorization: TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{"island": "island_final"}'
```

---

### Migrate Check

_This endpoint **is authenticated** and needs an auth token for access._

Checks migration requirements for all available territories. Shows which territories the player can migrate to and what conditions must be met.

Does not receive anything.

Returns [MigrateCheckResult](#migratecheckresult) in response.

**Endpoint:** `POST /migrate_check`

```shell
curl --request POST \
  --url https://bermudia-api-internal.darkube.app/api/v1/migrate_check \
  --header 'Authorization: TOKEN'
```

---

### Migrate

_This endpoint **is authenticated** and needs an auth token for access._

Migrates the player to a different territory.

Receives a [MigrateRequest](#migraterequest) in body.

Returns an empty object in response.

**Endpoint:** `POST /migrate`

```shell
curl --request POST \
  --url https://bermudia-api-internal.darkube.app/api/v1/migrate \
  --header 'Authorization: TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{"toTerritory": "territory_math"}'
```

---

### Unlock Treasure Check

_This endpoint **is authenticated** and needs an auth token for access._

Checks whether unlocking a specific treasure is possible.

Receives a [UnlockTreasureCheckRequest](#unlocktreasurecheckrequest) in body.

Returns [UnlockTreasureCheckResult](#unlocktreasurecheckresult) in response.

**Endpoint:** `POST /unlock_treasure_check`

```shell
curl --request POST \
  --url https://bermudia-api-internal.darkube.app/api/v1/unlock_treasure_check \
  --header 'Authorization: TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{"treasureID": "trs_C0B869257687459"}'
```

---

### Unlock Treasure

_This endpoint **is authenticated** and needs an auth token for access._

Unlocks a treasure for the player.

Receives a [UnlockTreasureRequest](#unlocktreasurerequest) in body.

Returns the updated [IslandTreasure](#islandtreasure) in response.

**Preconditions:** The player must be at the island where the treasure is located and also must have anchored.

**Endpoint:** `POST /unlock_treasure`

```shell
curl --request POST \
  --url https://bermudia-api-internal.darkube.app/api/v1/unlock_treasure \
  --header 'Authorization: TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{"treasureID": "trs_C0B869257687459"}'
```
---

### Make Trade Offer Check

_This endpoint **is authenticated** and needs an auth token for access._

Checks the limitation of making a new trade offer.

Does not receive anything.

Returns [MakeOfferCheckResult](#makeoffercheckresult) in response.

**Endpoint:** `POST /trade/make_offer_check`

```shell
curl --request POST \
  --url https://bermudia-api-internal.darkube.app/api/v1/trade/make_offer_check \
  --header 'Authorization: TOKEN'
```

---

### Make Trade Offer

_This endpoint **is authenticated** and needs an auth token for access._

Creates a new trade offer in the marketplace. The offered items are immediately deducted from the player's inventory and will be returned if the offer is deleted or accepted.

Receives a [MakeOfferRequest](#makeofferrequest) in body.

Returns an empty object in response.

**Endpoint:** `POST /trade/make_offer`

```shell
curl --request POST \
  --url https://bermudia-api-internal.darkube.app/api/v1/trade/make_offer \
  --header 'Authorization: TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{"offered": {"items": [{"type": "coin", "amount": 100}]}, "requested": {"items": [{"type": "blueKey", "amount": 1}]}}'
```

---

### Accept Trade Offer

_This endpoint **is authenticated** and needs an auth token for access._

Accepts an existing trade offer from another player. The trade is completed immediately if the accepting player has the required items.

Receives an [AcceptOfferRequest](#acceptofferrequest) in body.

Returns an empty object in response.

**Endpoint:** `POST /trade/accept_offer`

```shell
curl --request POST \
  --url https://bermudia-api-internal.darkube.app/api/v1/trade/accept_offer \
  --header 'Authorization: TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{"offerID": "tof_C0B869257687459"}'
```

---

### Delete Trade Offer

_This endpoint **is authenticated** and needs an auth token for access._

Deletes the player's own trade offer from the marketplace. The offered items are returned to the player's inventory.

Receives a [DeleteOfferRequest](#deleteofferrequest) in body.

Returns an empty object in response.

**Endpoint:** `POST /trade/delete_offer`

```shell
curl --request POST \
  --url https://bermudia-api-internal.darkube.app/api/v1/trade/delete_offer \
  --header 'Authorization: TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{"offerID": "tof_C0B869257687459"}'
```

---

### Get Trade Offers

_This endpoint **is authenticated** and needs an auth token for access._

Retrieves a paginated list of active trade offers from the marketplace, showing which offers the current player can accept.

**Parameters**:

- `page` (int, query param): Page number for pagination (0-based, default: 0)
- `limit` (int, query param): Number of offers per page (default: 5, max: 100)
- `by` (string, query param): Filters list based on the offerer. One of `me`, `others`. If empty, returns all offers.

Returns an array of [TradeOfferView](#tradeofferview) in response.

**Endpoint:** `GET /trade/offers`

```shell
curl --request POST \
  --url https://bermudia-api-internal.darkube.app/api/v1/trade/offers?page=1&limit=20 \
  --header 'Authorization: TOKEN'
```

---

## Data Models

### LoginRequest

| Field    | Type   | Description                 |
|----------|--------|-----------------------------|
| username | string | Unique username of the user |
| password | string | Password of user            |

### LoginResult

| Field | Type   | Description                                                |
|-------|--------|------------------------------------------------------------|
| token | string | JWT to be put in future requests. It is valid for 16 hours |


### TravelCheckRequest

| Field      | Type   | Description                                                                                         |
|------------|--------|-----------------------------------------------------------------------------------------------------|
| fromIsland | string | The current island of player (it is received by server to prevent travel in case of state mismatch) |
| toIsland   | string | The destination island                                                                              |


### TravelRequest

| Field      | Type   | Description                                                                                         |
|------------|--------|-----------------------------------------------------------------------------------------------------|
| fromIsland | string | The current island of player (it is received by server to prevent travel in case of state mismatch) |
| toIsland   | string | The destination island                                                                              |


### RefuelRequest

| Field  | Type   | Description                                                                                                                  |
|--------|--------|------------------------------------------------------------------------------------------------------------------------------|
| amount | int    | Amount of fuel to buy. Must be positive and not bigger that _maxAvailableAmount_ in [RefuelCheckResult](#refuelcheckresult). |


### AnchorCheckRequest

| Field  | Type   | Description                            |
|--------|--------|----------------------------------------|
| island | string | The id of island you want to anchor on |


### AnchorRequest

| Field  | Type   | Description                            |
|--------|--------|----------------------------------------|
| island | string | The id of island you want to anchor on |


### MigrateRequest

| Field       | Type   | Description                           |
|-------------|--------|---------------------------------------|
| toTerritory | string | The ID of the territory to migrate to |


### UnlockTreasureCheckRequest

| Field      | Type   | Description                 |
|------------|--------|-----------------------------|
| treasureID | string | The ID of the treasure      |


### UnlockTreasureRequest

| Field      | Type   | Description                 |
|------------|--------|-----------------------------|
| treasureID | string | The ID of the treasure      |


### MakeOfferRequest

| Field     | Type          | Description                                                     |
|-----------|---------------|-----------------------------------------------------------------|
| offered   | [Cost](#cost) | The items the player is offering in the trade                   |
| requested | [Cost](#cost) | The items the player is requesting in exchange for their offer  |


### AcceptOfferRequest

| Field   | Type   | Description                        |
|---------|--------|------------------------------------|
| offerID | string | The unique identifier of the offer |


### DeleteOfferRequest

| Field   | Type   | Description                        |
|---------|--------|------------------------------------|
| offerID | string | The unique identifier of the offer |


### Me

| Field    | Type   | Description                 |
|----------|--------|-----------------------------|
| id       | int    | Unique numeric id of user   |
| username | string | Unique username of the user |


### Player

| Field         | Type                            | Description                                                    |
|---------------|---------------------------------|----------------------------------------------------------------|
| atTerritory   | string                          | Current territory of player                                    |
| atIsland      | string                          | Current island of player                                       |
| anchored      | boolean                         | `true` if player has anchored at _atIsland_, `false` otherwise |
| fuel          | string                          | Current fuel level of player's vehicle                         |
| fuelCap       | int                             | Current fuel capacity of player's vehicle                      |
| coins         | int                             | Current number of coins of player                              |
| blueKeys      | int                             | Current number of blue keys of player                          |
| redKeys       | int                             | Current number of red keys of player                           |
| goldenKeys    | int                             | Current number of golden keys of player                        |
| knowledgeBars | [KnowledgeBar](#knowledgebar)[] | Current state of player's knowledge in each territory          |
| books         | [Book](#book)[]                 | Player's books, sorted in the order they were achieved         |

### Territory

| Field               | Type                                | Description                                                                                                                                                       |
|---------------------|-------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id                  | string                              | Unique identifier for the territory                                                                                                                               |
| name                | string                              | Display name of the territory                                                                                                                                     |
| backgroundAsset     | string                              | Asset file name for background                                                                                                                                    |
| islands             | [Island](#island)[]                 | Array of islands in this territory                                                                                                                                |
| edges               | [Edge](#edge)[]                     | Array of connections between islands                                                                                                                              |
| refuelIslands       | [RefuelIsland](#refuelisland)[]     | Array of refuel islands                                                                                                                                           |
| terminalIslands     | [TerminalIsland](#terminalisland)[] | Array of terminal islands. To call [Migrate](#migrate) you must be in one of these islands                                                                        |
| islandPrerequisites | Map<string, string[]>               | Map of island id to a list of prerequisite island ids.<br />Player can't travel to an island if they have not answered all questions in the prerequisite islands. |


### Island

| Field     | Type   | Description                                                                                     |
|-----------|--------|-------------------------------------------------------------------------------------------------|
| id        | string | Unique identifier for the island                                                                |
| name      | string | Display name of the island                                                                      |
| x         | number | X coordinate on the map (between 0 and 1, a coefficient of territory's background asset width)  |
| y         | number | Y coordinate on the map (between 0 and 1, a coefficient of territory's background asset height) |
| width     | number | Width of icon (between 0 and 1, a coefficient of territory's background asset width)            |
| height    | number | Height of icon (between 0 and 1, a coefficient of territory's background asset height)          |
| iconAsset | string | Asset file name for background.                                                                 |

### Edge

| Field | Type   | Description                  |
|-------|--------|------------------------------|
| from  | string | ID of the source island      |
| to    | string | ID of the destination island |

**Note:** Edges represent bidirectional connections. If there's an edge from A to B, players can travel both ways.

### RefuelIsland

| Field | Type   | Description  |
|-------|--------|--------------|
| id    | string | ID of island |

### TerminalIsland

| Field | Type   | Description  |
|-------|--------|--------------|
| id    | string | ID of island |


### IslandContent

| Field      | Type                                  | Description                                   |
|------------|---------------------------------------|-----------------------------------------------|
| components | [IslandComponent](#islandcomponent)[] | Array of components to be displayed in the UI |
| treasures  | [IslandTreasure](#islandtreasure)[]   | Array of treasures available on the island    |


### IslandComponent

| Field  | Type                            | Description                  |
|--------|---------------------------------|------------------------------|
| iframe | [IslandIFrame](#islandiframe)[] | The data of iframe component |
| input  | [IslandInput](#islandinput)[]   | The data of input component  |

**Note:** Exactly one of these fields will be present in a component object

### IslandIFrame

| Field | Type   | Description                  |
|-------|--------|------------------------------|
| url   | string | The source url of the iframe |


### IslandTreasure

| Field    | Type           | Description                                                                             |
|----------|----------------|-----------------------------------------------------------------------------------------|
| id       | string         | Unique identifier for the treasure                                                      |
| unlocked | bool           | Whether the treasure has been unlocked by the player                                    |
| reward   | [Cost](#cost)? | If _unlocked_ is true, shows the achieved reward. Length of _reward.items_ can be zero. |


### IslandInput

| Field           | Type                                | Description                                                                                                                                                                                                  |
|-----------------|-------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id              | string                              | The unique id of this input, to be used in [Submit Answer](#submit-answer)                                                                                                                                   |
| type            | string                              | Type of the data this input receives. One of [HTML Input Element Types](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/input#input_types) (usually one of `text`, `number` and `file`) |
| accept          | []string?                           | If type is `file`, this field is present and contains the accepted MIME types.                                                                                                                               |
| description     | string                              | Description of the input to be shown to user                                                                                                                                                                 |
| submissionState | [SubmissionState](#submissionstate) | The current submission state of this input.                                                                                                                                                                  |


### SubmissionState

| Field       | Type    | Description                                                                                                          |
|-------------|---------|----------------------------------------------------------------------------------------------------------------------|
| submittable | boolean | True if the a new answer can be submitted, false otherwise.                                                          |
| status      | string  | The status of answer; one of `empty`, `pending` (in process of correction) , `correct`, `wrong`                      |
| filename    | string? | If _status_ is not `empty` and [IslandInput](#islandinput) _type_ is `file`, the name of the last submitted file.    |
| value       | string? | If _status_ is not `empty` and [IslandInput](#islandinput) _type_ is not `file`, the last submitted plain text value |
| submittedAt | string? | If _status_ is not `empty`, the time of last submission in Unix milliseconds.                                        |


### Book

| Field       | Type   | Description                          |
|-------------|--------|--------------------------------------|
| territoryId | string | ID of territory this book belongs to |
| islandId    | string | ID of island this book belongs to    |
| islandName  | string | Name of island this book belongs to  |


### KnowledgeBar

| Field       | Type   | Description                                            |
|-------------|--------|--------------------------------------------------------|
| territoryId | string | ID of territory this knowledge bar belongs to          |
| value       | int    | Player's knowledge in the territory                    |
| total       | int    | Total amount of knowledge that exists in the territory |

### Event

| Field        | Type                                    | Description                                                                          |
|--------------|-----------------------------------------|--------------------------------------------------------------------------------------|
| playerUpdate | [PlayerUpdateEvent](#playerupdateevent) | If event is a player update event, this field is present.                            |
| timestamp    | string                                  | Time of event emission in Unix milliseconds. Can be used to discard very old events. |


### Cost

| Field | Type                    | Description                                                    |
|-------|-------------------------|----------------------------------------------------------------|
| items | [CostItem](#costitem)[] | An array of [CostItem](#costitem) determining the needed items |


### CostItem

| Field  | Type   | Description                                                                      |
|--------|--------|----------------------------------------------------------------------------------|
| type   | string | Type of the needed item. One of `fuel`, `coin`, `blueKey`, `redKey`, `goldenKey` |
| amount | int    | The number of items needed of this type                                          |


### TravelCheckResult

| Field      | Type          | Description                                                   |
|------------|---------------|---------------------------------------------------------------|
| feasible   | boolean       | True if the travel can be done, false otherwise.              |
| travelCost | [Cost](#cost) | An object representing the needed items for the travel        |
| reason     | string?       | If _feasible_ is false, this field is present and reports why |


### RefuelCheckResult

| Field              | Type   | Description                                                                     |
|--------------------|--------|---------------------------------------------------------------------------------|
| maxAvailableAmount | int    | Maximum amount of fuel that can be bought depending on conditions. Can be zero. |
| coinCostPerUnit    | int    | The number of coins needed to buy a unit of fuel.                               |
| maxReason          | string | A description for the cause of _maxAvailableAmount_                             |


### AnchorCheckResult

| Field         | Type          | Description                                                   |
|---------------|---------------|---------------------------------------------------------------|
| feasible      | boolean       | True if the travel can be done, false otherwise               |
| anchoringCost | [Cost](#cost) | An object representing the needed items for anchoring         |
| reason        | string?       | If _feasible_ is false, this field is present and reports why |


### UnlockTreasureCheckResult

| Field    | Type          | Description                                                   |
|----------|---------------|---------------------------------------------------------------|
| feasible | boolean       | True if the treasure can be unlocked, false otherwise         |
| cost     | [Cost](#cost) | An object representing the needed items for unlocking         |
| reason   | string?       | If _feasible_ is false, this field is present and reports why |


### MigrateCheckResult

| Field                      | Type                                                    | Description                                                               |
|----------------------------|---------------------------------------------------------|---------------------------------------------------------------------------|
| knowledgeCriteriaTerritory | string                                                  | The territory used as criteria for knowledge-based migration requirements |
| knowledgeValue             | int                                                     | Player's current knowledge value in the criteria territory                |
| minAcceptableKnowledge     | int                                                     | Minimum knowledge required in the criteria territory for free migration   |
| territoryMigrationOptions  | [TerritoryMigrationOption](#territorymigrationoption)[] | Array of migration options for all territories                            |


### TerritoryMigrationOption

| Field                  | Type          | Description                                                                         |
|------------------------|---------------|-------------------------------------------------------------------------------------|
| territoryId            | string        | ID of the territory                                                                 |
| territoryName          | string        | Name of the territory                                                               |
| status                 | string        | Migration status: `resident` (current), `visited` (previously visited), `untouched` |
| migrationCost          | [Cost](#cost) | Cost required for migration (when mustPayCost is true)                              |
| mustPayCost            | boolean       | Whether the player must pay the migration cost                                      |
| feasible               | boolean       | Whether migration to this territory is possible                                     |
| reason                 | string?       | If _feasible_ is false, explanation of why migration is not possible                |


### MakeOfferCheckResult

| Field         | Type          | Description                                                                      |
|---------------|---------------|----------------------------------------------------------------------------------|
| feasible      | boolean       | True if player can make a new offer, false otherwise                             |
| tradableItems | [Cost](#cost) | Shows a list of tradable items and how many of each can the user put in a offer. |
| reason        | string?       | If _feasible_ is false, this field is present and reports why                    |


### TradeOfferView

| Field      | Type          | Description                                         |
|------------|---------------|-----------------------------------------------------|
| id         | string        | Unique identifier of the trade offer                |
| by         | string        | Name of the player who created the offer            |
| offered    | [Cost](#cost) | The items being offered by the creator              |
| requested  | [Cost](#cost) | The items being requested in exchange               |
| createdAt  | string        | Time when the offer was created (Unix milliseconds) |
| acceptable | boolean       | Whether the current player can accept this offer    |


### PlayerUpdateEvent

| Field  | Type              | Description                                                                                                                                                                                                      |
|--------|-------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| reason | string            | The reason for change in player state. One of `initial`, `travel`, `refuel`, `correction`, `anchor`, `migration`, `unlockTreasure`, `newBook`, `makeOffer`, `acceptOffer`, `ownOfferAccepted`, `ownOfferDeleted` |
| player | [Player](#player) | The new value of player object.                                                                                                                                                                                  |