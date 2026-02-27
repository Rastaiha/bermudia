# Admin API Documentation

## Base URL

```
{{protocol}}://bermudia-api-internal.darkube.app/admin
```

## Response Format

All endpoints return responses in the same format as the main API:

```json
{
  "ok": true,
  "error": "string",
  "result": {}
}
```

- `ok`: Boolean indicating success/failure
- `error`: Error message (only present when `ok=false`)
- `result`: Response data (only present when `ok=true`)

## Authentication

All admin endpoints (except [Admin Login](#admin-login)) require a valid admin JWT token.

The token must be sent in the `Authorization` request header:

```
Authorization: Bearer <token>
```

Admin tokens are valid for **6 hours** from the time of issuance.

---

## Endpoints

### Admin Login

Authenticates an admin user and returns a JWT token.

**Endpoint:** `POST /admin/login`

**Request Body:**

| Field    | Type   | Required | Description        |
|----------|--------|----------|--------------------|
| username | string | yes      | Admin username     |
| password | string | yes      | Admin password     |

**Response:** [`AdminLoginResult`](#adminloginresult)

---

### Get Territories

_**Requires authentication.**_

Returns a list of all territories in the game.

**Endpoint:** `GET /admin/territories`

**Response:** Array of [`Territory`](#territory)

---

### Set Territory

_**Requires authentication.**_

Creates or updates a territory. The territory is identified by its `id` field. If a territory with the given `id` already exists, it will be fully replaced. Islands listed in the territory that do not yet exist in the database will be created (their IDs reserved).

**Endpoint:** `POST /admin/territories`

**Request Body:** [`Territory`](#territory)

**Response:** [`Territory`](#territory) — the territory as submitted

**Validation rules:**
- `id` is required
- `startIsland` is required and must refer to one of the islands in the `islands` array
- All islands in `islands` must have a non-empty `id` and `name`
- All edges must refer to islands in the `islands` array
- All `refuelIslands` and `terminalIslands` must refer to islands in the `islands` array
- All keys in `islandPrerequisites` and their prerequisite values must refer to islands in the `islands` array

---

### Get Territory Island Bindings

_**Requires authentication.**_

Returns the current island bindings and pool settings for a territory. Shows which islands have a directly-assigned book, which are assigned from a pool, and which are empty.

**Endpoint:** `GET /admin/territories/{territoryID}/island_bindings`

**Path Parameters:**

| Parameter   | Description                         |
|-------------|-------------------------------------|
| territoryID | The unique identifier of the territory |

**Response:** [`TerritoryIslandBindings`](#territoryislandbindings)

---

### Set Territory Island Bindings

_**Requires authentication.**_

Updates the island bindings and pool settings for a territory. This determines which islands draw their content from a pool of books (assigned randomly per player) and which are empty.

**Endpoint:** `POST /admin/territories/{territoryID}/island_bindings`

**Request Body:** [`TerritoryIslandBindings`](#territoryislandbindings)

**Response:** [`TerritoryIslandBindings`](#territoryislandbindings)

**Note:** The total count of islands in `pooledIslands` must equal `poolSettings.easy + poolSettings.medium + poolSettings.hard`, otherwise the request is rejected.

---

### Get Island Header

_**Requires authentication.**_

Returns the header metadata for a single island, including whether it has a book assigned and whether that book is from a pool.

**Endpoint:** `GET /admin/islands/{islandID}`

**Path Parameters:**

| Parameter | Description                       |
|-----------|-----------------------------------|
| islandID  | The unique identifier of the island |

**Response:** [`IslandHeader`](#islandheader)

---

### Set Book and Bind to Island

_**Requires authentication.**_

Creates or updates a book and directly binds it to the given island. The operation is upsert-style: if `bookId` is omitted from the request body, a new book is created and a new ID is generated and returned. If a `bookId` from a previous response is provided, the existing book is updated in place.

The same logic applies to individual components within the book:
- If a question's `id` field is omitted, a new question is created with a generated ID.
- If a question's `id` (e.g., `qst_...`) from a previous response is included, that question is updated.
- If a treasure's `id` field is omitted, a new treasure is created.
- If a treasure's `id` from a previous response is included, that treasure entry is updated.

**Endpoint:** `POST /admin/islands/{islandID}/book`

**Path Parameters:**

| Parameter | Description                       |
|-----------|-----------------------------------|
| islandID  | The unique identifier of the island |

**Request Body:** [`BookInput`](#bookinput)

**Response:** [`BookInput`](#bookinput) — the full book as saved, with all IDs populated

---

### Get Book

_**Requires authentication.**_

Returns the full contents of a book, including all its components (iframes and questions) and treasures. The returned object can be modified and resubmitted to update the book.

**Endpoint:** `GET /admin/books/{bookID}`

**Path Parameters:**

| Parameter | Description                    |
|-----------|--------------------------------|
| bookID    | The unique identifier of the book |

**Response:** [`BookInput`](#bookinput)

---

### Get Pools

_**Requires authentication.**_

Returns all book pools (`easy`, `medium`, `hard`) and the list of book IDs assigned to each pool.

**Endpoint:** `GET /admin/pools`

**Response:** Array of [`PoolOutput`](#pooloutput)

---

### Set Book and Bind to Pool

_**Requires authentication.**_

Creates or updates a book and adds it to the given pool. The same upsert behavior as [Set Book and Bind to Island](#set-book-and-bind-to-island) applies: omitting `bookId` creates a new book, while including an existing `bookId` updates that book.

**Endpoint:** `POST /admin/pools/{poolID}/books`

**Path Parameters:**

| Parameter | Description                                              |
|-----------|----------------------------------------------------------|
| poolID    | The pool to bind the book to: `easy`, `medium`, or `hard` |

**Request Body:** [`BookInput`](#bookinput)

**Response:** [`BookInput`](#bookinput) — the full book as saved, with all IDs populated

---

### Get Users

_**Requires authentication.**_

Returns a list of all registered users. Passwords are not included in the response.

**Endpoint:** `GET /admin/users`

**Response:** Array of [`UserOutput`](#useroutput)

---

### Create or Update User

_**Requires authentication.**_

Creates a new user, or updates an existing user if the `username` already exists. When updating, all provided fields are overwritten. If `password` is omitted, a random password is generated and returned in the response — **this is the only time the password is visible**, so store it immediately.

A player record is initialized for the user in their `startingTerritory` upon creation. If the user already exists, the player record is not re-initialized.

**Endpoint:** `POST /admin/users`

**Request Body:** [`UserInput`](#userinput)

**Response:** [`UserInput`](#userinput) — includes the plaintext `password` (generated or as provided)

---

## Data Types

### AdminLoginResult

| Field | Type   | Description                              |
|-------|--------|------------------------------------------|
| token | string | JWT token to use for authenticated requests |

---

### Territory

| Field               | Type                                            | Description                                                                   |
|---------------------|-------------------------------------------------|-------------------------------------------------------------------------------|
| id                  | string                                          | Unique identifier of the territory                                            |
| name                | string                                          | Display name                                                                  |
| backgroundAsset     | string                                          | Asset key for the background image                                            |
| startIsland         | string                                          | ID of the island where players start when entering this territory             |
| islands             | [`Island`](#island)[]                           | All islands in this territory                                                 |
| edges               | [`Edge`](#edge)[]                               | Connections between islands that players can travel along                     |
| refuelIslands       | [`RefuelIsland`](#refuelisland)[]               | Islands where players can purchase fuel                                       |
| terminalIslands     | [`TerminalIsland`](#terminalisland)[]           | Islands from which players can migrate to another territory                   |
| islandPrerequisites | [`IslandPrerequisites`](#islandprerequisites)   | Map of island ID → list of island IDs that must be completed before access    |

---

### Island

| Field      | Type   | Description                         |
|------------|--------|-------------------------------------|
| id         | string | Unique identifier                   |
| name       | string | Display name                        |
| x          | float  | X position on the territory map     |
| y          | float  | Y position on the territory map     |
| width      | float  | Width on the territory map          |
| height     | float  | Height on the territory map         |
| iconAsset  | string | Asset key for the island icon       |

---

### Edge

| Field | Type   | Description                            |
|-------|--------|----------------------------------------|
| from  | string | ID of the island at one end of the edge |
| to    | string | ID of the island at the other end      |

Edges are **bidirectional** — a player can travel in either direction along an edge.

---

### RefuelIsland

| Field | Type   | Description               |
|-------|--------|---------------------------|
| id    | string | ID of the refuel island   |

---

### TerminalIsland

| Field | Type   | Description                |
|-------|--------|----------------------------|
| id    | string | ID of the terminal island  |

---

### IslandPrerequisites

A JSON object mapping an island ID (string) to an array of island IDs (strings) that must be answered/completed before the player can access it.

```json
{
  "island_math2": ["island_math1"],
  "island_math3": ["island_math1", "island_math2"]
}
```

---

### IslandHeader

| Field       | Type    | Description                                                         |
|-------------|---------|---------------------------------------------------------------------|
| id          | string  | Unique identifier of the island                                     |
| name        | string  | Display name                                                        |
| territory_id | string | ID of the territory this island belongs to                          |
| bookId      | string  | ID of the directly-assigned book, if any (empty string if none)     |
| fromPool    | boolean | True if this island's content is assigned from a pool               |

---

### BookInput

The primary structure for creating or updating a book. This same structure is used for both requests and responses.

| Field      | Type                                              | Required (on write) | Description                                                                                               |
|------------|---------------------------------------------------|---------------------|-----------------------------------------------------------------------------------------------------------|
| bookId     | string                                            | no                  | ID of the book (`bok_...` prefix). Omit to create a new book; include to update an existing book.         |
| components | [`BookInputComponent`](#bookinputcomponent)[]     | yes                 | Ordered list of the book's content components                                                             |
| treasures  | [`BookTreasureComponent`](#booktreasurecomponent)[] | no                | Treasure chests embedded in this book                                                                     |

---

### BookInputComponent

Exactly one of `iframe` or `question` must be present.

| Field    | Type                                              | Description                       |
|----------|---------------------------------------------------|-----------------------------------|
| iframe   | [`IslandIFrame`](#islandiframe)?                  | An embedded iframe component      |
| question | [`IslandInputQuestion`](#islandinputquestion)?    | A question/answer input component |

---

### IslandIFrame

| Field | Type   | Description            |
|-------|--------|------------------------|
| url   | string | URL to embed in the iframe |

---

### IslandInputQuestion

| Field               | Type     | Required | Description                                                                                                       |
|---------------------|----------|----------|-------------------------------------------------------------------------------------------------------------------|
| id                  | string   | no       | ID of the question (`qst_...` prefix). Omit to create; include a previously returned ID to update.               |
| text                | string   | yes      | The question prompt shown to the player                                                                           |
| inputType           | string   | yes      | Type of answer input: `text` or `file`                                                                            |
| inputAccept         | string[] | conditional | Required when `inputType` is `file`. Accepted MIME types or file extensions.                                 |
| knowledgeAmount     | int      | yes      | Amount of knowledge points awarded upon correct answer (must be ≥ 0)                                              |
| rewardSource        | string   | no       | Reward tier for correct answers. Valid values: `edu1`, `edu2`, `edu3`, `edu4`, `edu5`, `edu6`, `final`. Omit for no reward. |
| correctionHintMessage | string | no      | Optional hint/context message shown to correctors                                                                 |

---

### BookTreasureComponent

| Field | Type   | Description                                                                                      |
|-------|--------|--------------------------------------------------------------------------------------------------|
| id    | string | ID of the treasure (`trs_...` prefix). Omit to create; include a previously returned ID to update. |

---

### PoolOutput

| Field | Type     | Description                                       |
|-------|----------|---------------------------------------------------|
| id    | string   | Pool identifier: `easy`, `medium`, or `hard`      |
| books | string[] | List of book IDs currently in this pool           |

---

### TerritoryIslandBindings

| Field         | Type                                                    | Description                                                                         |
|---------------|---------------------------------------------------------|-------------------------------------------------------------------------------------|
| territoryId   | string                                                  | ID of the territory                                                                 |
| emptyIslands  | string[]                                                | IDs of islands with no book assigned (not from pool, no direct book)                |
| pooledIslands | string[]                                                | IDs of islands that draw their content from a pool (assigned per-player at runtime) |
| poolSettings  | [`TerritoryPoolSettings`](#territorypoolsettings)       | How many pooled books belong to each difficulty pool                                |

---

### TerritoryPoolSettings

| Field  | Type | Description                                       |
|--------|------|---------------------------------------------------|
| easy   | int  | Number of pooled islands drawing from the `easy` pool   |
| medium | int  | Number of pooled islands drawing from the `medium` pool |
| hard   | int  | Number of pooled islands drawing from the `hard` pool   |

The sum `easy + medium + hard` must equal the number of IDs in `pooledIslands`.

---

### UserInput

| Field             | Type   | Required | Description                                                                             |
|-------------------|--------|----------|-----------------------------------------------------------------------------------------|
| username          | string | yes      | Login username (case-insensitive). If a user with this username already exists, that user is updated. |
| name              | string | no       | Display name of the user                                                                |
| password          | string | no       | Plaintext password. If omitted, a random password is generated and returned.            |
| startingTerritory | string | yes      | ID of the territory where the player will start. Must be an existing territory.         |
| meetLink          | string | no       | URL for the user's meet/video call link                                                 |

---

### UserOutput

Same fields as [`UserInput`](#userinput), but **`password` is always omitted**.

| Field     | Type   | Description                        |
|-----------|--------|------------------------------------|
| username  | string | Login username                     |
| name      | string | Display name                       |
| meetLink  | string | URL for the user's meet link       |

---

## Error Responses

| HTTP Status | When                                                                 |
|-------------|----------------------------------------------------------------------|
| 400         | Invalid request body or failed validation (e.g., missing required fields, mismatched counts) |
| 401         | Missing or invalid/expired auth token                                |
| 404         | Referenced resource not found (territory, island, book, etc.)        |
| 409         | Rule violation (e.g., conflicting state)                             |
| 500         | Internal server error                                                |