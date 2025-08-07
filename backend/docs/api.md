# Educational Game API Documentation

## Base URL

```
http://localhost:8080/api/v1
```

## Endpoints

### Get Territory

Retrieves the complete data for a specific territory, including all islands and their connections.

**Endpoint:** `GET /territories/{territoryID}`

**Parameters:**

- `territoryID` (path parameter, required): The unique identifier of the territory

**Response Format:**

# Educational Game API Documentation

## Base URL

```
http://localhost:8080/api/v1
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

## Endpoints

### Get Territory

Retrieves the complete data for a specific territory, including all islands and their connections.

**Endpoint:** `GET /territories/{territoryID}`

**Parameters:**

- `territoryID` (path parameter, required): The unique identifier of the territory

**Success Response (200 OK):**

```json
{
  "ok": true,
  "result": {
    "id": "territory1",
    "name": "Mystic Archipelago",
    "description": "A magical collection of islands where mathematics and science come alive",
    "islands": [
      {
        "id": "island_start",
        "name": "Welcome Harbor",
        "x": 100,
        "y": 150,
        "assetFile": "harbor_icon.png",
        "description": "The starting point for all new adventurers"
      }
    ],
    "edges": [
      {
        "from": "island_start",
        "to": "island_math1"
      }
    ]
  }
}
```

**Error Responses:**

**404 Not Found:**

```json
{
  "ok": false,
  "error": "Territory not found"
}
```

**400 Bad Request:**

```json
{
  "ok": false,
  "error": "Territory ID is required"
}
```

**408 Request Timeout:**

```json
{
  "ok": false,
  "error": "Request timeout"
}
```

**500 Internal Server Error:**

```json
{
  "ok": false,
  "error": "Failed to retrieve territory"
}
```

## Data Models

### Territory

| Field       | Type     | Description                           |
|-------------|----------|---------------------------------------|
| id          | string   | Unique identifier for the territory   |
| name        | string   | Display name of the territory         |
| description | string   | Optional description of the territory |
| islands     | Island[] | Array of islands in this territory    |
| edges       | Edge[]   | Array of connections between islands  |

### Island

| Field       | Type   | Description                        |
|-------------|--------|------------------------------------|
| id          | string | Unique identifier for the island   |
| name        | string | Display name of the island         |
| x           | number | X coordinate on the map            |
| y           | number | Y coordinate on the map            |
| assetFile   | string | Filename of the icon asset         |
| description | string | Optional description of the island |

### Edge

| Field | Type   | Description                  |
|-------|--------|------------------------------|
| from  | string | ID of the source island      |
| to    | string | ID of the destination island |

**Note:** Edges represent bidirectional connections. If there's an edge from A to B, players can travel both ways.

## Example Usage

### Get Territory Data

```bash
curl -X GET "http://localhost:8080/api/v1/territories/territory1"
```

### Health Check

```bash
curl -X GET "http://localhost:8080/health"
```

## Error Handling

The API uses proper HTTP status codes and consistent error response format. All errors include:

- `ok`: Always `false` for errors
- `error`: Human-readable error message

The API also handles context timeouts and cancellations properly, returning appropriate timeout errors when requests
take too long.

## CORS

The API includes CORS headers to allow cross-origin requests from frontend applications.

## Data Models

### Territory

| Field       | Type     | Description                           |
|-------------|----------|---------------------------------------|
| id          | string   | Unique identifier for the territory   |
| name        | string   | Display name of the territory         |
| islands     | Island[] | Array of islands in this territory    |
| edges       | Edge[]   | Array of connections between islands  |

### Island

| Field       | Type   | Description                        |
|-------------|--------|------------------------------------|
| id          | string | Unique identifier for the island   |
| name        | string | Display name of the island         |
| x           | number | X coordinate on the map            |
| y           | number | Y coordinate on the map            |
| assetFile   | string | Filename of the icon asset         |
| description | string | Optional description of the island |

### Edge

| Field | Type   | Description                  |
|-------|--------|------------------------------|
| from  | string | ID of the source island      |
| to    | string | ID of the destination island |

**Note:** Edges represent bidirectional connections. If there's an edge from A to B, players can travel both ways.

## Example Usage

### Get Territory Data

```bash
curl -X GET "http://localhost:8080/api/v1/territories/territory1"
```

### Health Check

```bash
curl -X GET "http://localhost:8080/health"
```

## CORS

The API includes CORS headers to allow cross-origin requests from frontend applications.
