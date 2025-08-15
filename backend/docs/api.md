# Educational Game API Documentation

## Base URL

For development purposes:

```
http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1
```

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

The generated JWT must be but in the `Authorization` request header for authenticated endpoints.

## Endpoints

### Login

Checks user credentials and generates a jwt token.

Receives a [LoginRequest](#loginrequest) in body.

Returns a [LoginResult](#loginresult) in response.

**Endpoint:** `POST /login`

### Get Me (authenticated)

Returns [Me](#me) in response.

This endpoint is just a way to test whether the JWT is correctly set and still valid.

**Endpoint:** `GET /me`

### Get Territory

Retrieves the complete data for a specific territory, including all islands and their connections.

Returns a [Territory](#territory) in response.

**Endpoint:** `GET /territories/{territoryID}`

**Parameters:**

- `territoryID` (path parameter, required): The unique identifier of the territory

### Get Island

Retrieves the content of the island.

Returns a [IslandContent](#islandcontent) in response.

**Endpoint:** `GET /islands/{islandID}`

**Parameters:**

- `islandID` (path parameter, required): The unique identifier of the island

### Submit Answer

Receives the input of a [IslandInput](#islandinput) components.

Returns an empty object on response.

**Endpoint:** `POST /answer/{inputID}`

**Parameters:**

- `data` (body parameter, required): The user data. Its type depends on the [IslandInput](#islandinput) `type`.

**Note:** Request's `Content-Type` should be `multipart/form-data`

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

### Me

| Field    | Type   | Description                 |
|----------|--------|-----------------------------|
| id       | int    | Unique numeric id of user   |
| username | string | Unique username of the user |

### Territory

| Field           | Type                | Description                          |
|-----------------|---------------------|--------------------------------------|
| id              | string              | Unique identifier for the territory  |
| name            | string              | Display name of the territory        |
| backgroundAsset | string              | Asset file name for background       |
| islands         | [Island](#island)[] | Array of islands in this territory   |
| edges           | [Edge](#edge)[]     | Array of connections between islands |

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

| Field       | Type      | Description                                                                                                                                                                                                  |
|-------------|-----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id          | string    | The unique id of this input, to be used in [Submit Answer](#submit-answer)                                                                                                                                   |
| type        | string    | Type of the data this input receives. One of [HTML Input Element Types](https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/input#input_types) (usually one of `text`, `number` and `file`) |
| accept      | []string? | If type is `file`, this field is present and contains the accepted MIME types.                                                                                                                               |
| description | string    | Description of the input to be shown to user                                                                                                                                                                 |
