# Educational Game API Documentation

## Base URL

### Development

```
{{protocol}}://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1
```

protocols: `http`, `ws`

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
  --url http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1/answer/ans_29C12F3C7D089666 \
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
  --url http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1/travel_check \
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
  --url http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1/travel \
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
  --url http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1/refuel_check \
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
  --url http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1/refuel \
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
  --url http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1/anchor_check \
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
  --url http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1/anchor \
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
  --url http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1/migrate_check \
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
  --url http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1/migrate \
  --header 'Authorization: TOKEN' \
  --header 'Content-Type: application/json' \
  --data '{"toTerritory": "territory_math"}'
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
| knowledgeBars | [KnowledgeBar](#knowledgebar)[] | Current state of player's knowledge in each territory          |

### Territory

| Field           | Type                                | Description                                                                                |
|-----------------|-------------------------------------|--------------------------------------------------------------------------------------------|
| id              | string                              | Unique identifier for the territory                                                        |
| name            | string                              | Display name of the territory                                                              |
| backgroundAsset | string                              | Asset file name for background                                                             |
| islands         | [Island](#island)[]                 | Array of islands in this territory                                                         |
| edges           | [Edge](#edge)[]                     | Array of connections between islands                                                       |
| refuelIslands   | [RefuelIsland](#refuelisland)[]     | Array of refuel islands                                                                    |
| terminalIslands | [TerminalIsland](#terminalisland)[] | Array of terminal islands. To call [Migrate](#migrate) you must be in one of these islands |


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

| Field  | Type   | Description                                    |
|--------|--------|------------------------------------------------|
| type   | string | Type of the needed item. One of `fuel`, `coin` |
| amount | int    | The number of items needed of this type        |


### TravelCheckResult

| Field      | Type          | Description                                                    |
|------------|---------------|----------------------------------------------------------------|
| feasible   | boolean       | True if the travel can be done, false otherwise.               |
| travelCost | [Cost](#cost) | An object representing the needed items for the travel         |
| reason     | string?       | If _feasible_ is false, this field is presents and reports why |


### RefuelCheckResult

| Field              | Type   | Description                                                                     |
|--------------------|--------|---------------------------------------------------------------------------------|
| maxAvailableAmount | int    | Maximum amount of fuel that can be bought depending on conditions. Can be zero. |
| coinCostPerUnit    | int    | The number of coins needed to buy a unit of fuel.                               |
| maxReason          | string | A description for the cause of _maxAvailableAmount_                             |


### AnchorCheckResult

| Field         | Type          | Description                                                    |
|---------------|---------------|----------------------------------------------------------------|
| feasible      | boolean       | True if the travel can be done, false otherwise                |
| anchoringCost | [Cost](#cost) | An object representing the needed items for anchoring          |
| reason        | string?       | If _feasible_ is false, this field is presents and reports why |


### MigrateCheckResult

| Field                      | Type                                                    | Description                                                               |
|----------------------------|---------------------------------------------------------|---------------------------------------------------------------------------|
| knowledgeCriteriaTerritory | string                                                  | The territory used as criteria for knowledge-based migration requirements |
| knowledgeValue             | int                                                     | Player's current knowledge value in the criteria territory                |
| territoryMigrationOptions  | [TerritoryMigrationOption](#territorymigrationoption)[] | Array of migration options for all territories                            |


### TerritoryMigrationOption

| Field                  | Type          | Description                                                                         |
|------------------------|---------------|-------------------------------------------------------------------------------------|
| territoryId            | string        | ID of the territory                                                                 |
| territoryName          | string        | Name of the territory                                                               |
| status                 | string        | Migration status: `resident` (current), `visited` (previously visited), `untouched` |
| minAcceptableKnowledge | int           | Minimum knowledge required for free migration to this territory                     |
| migrationCost          | [Cost](#cost) | Cost required for migration (when mustPayCost is true)                              |
| mustPayCost            | boolean       | Whether the player must pay the migration cost                                      |
| feasible               | boolean       | Whether migration to this territory is possible                                     |
| reason                 | string?       | If _feasible_ is false, explanation of why migration is not possible                |


### PlayerUpdateEvent

| Field  | Type              | Description                                                                                                      |
|--------|-------------------|------------------------------------------------------------------------------------------------------------------|
| reason | string            | The reason for change in player state. One of `initial`, `travel`, `refuel`, `correction`, `anchor`, `migration` |
| player | [Player](#player) | The new value of player object.                                                                                  |
