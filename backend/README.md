# Bermudia Backend

The Bermudia backend is a robust Go application that powers the gamified learning platform. It provides RESTful APIs, WebSocket connections for real-time features, and integrates with external services for enhanced functionality.

## 🏗️ Architecture

The backend follows clean architecture principles with clear separation of concerns:

```text
backend/
├── main.go                    # Application entry point
├── adminbot/                  # Telegram/Bale bot integration
│   └── bot.go
├── api/                       # API layer
│   ├── handler/              # HTTP request handlers
│   │   ├── handler.go        # Router setup, middleware, server lifecycle
│   │   ├── auth.go           # Authentication endpoints/middleware
│   │   ├── admin.go          # Admin API endpoints
│   │   ├── events.go         # WebSocket event streaming
│   │   ├── island.go         # Island-related endpoints
│   │   ├── player.go         # Player management
│   │   ├── territory.go      # Territory endpoints
│   │   └── response.go       # Response utilities
│   └── hub/                  # WebSocket hub
│       └── hub.go
└── internal/                 # Internal packages
    ├── config/               # Configuration management (koanf-based)
    ├── domain/               # Domain models and business rules
    ├── mock/                 # Mock game content seeded in dev mode
    ├── repository/           # Data access layer
    └── service/              # Business logic layer
```

## 🚀 Getting Started

### Prerequisites

- Go 1.25 or higher
- PostgreSQL 15+ or SQLite 3+
- Git

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/Rastaiha/bermudia.git
   cd bermudia/backend
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Set up environment variables**

   Configuration is loaded via [koanf](https://github.com/knadh/koanf) from environment variables, not from a `.env` file. Every variable is prefixed with `BERMUDIA__`, and `__` is used as the nesting separator for nested structs (e.g. `Config.Postgres.Enable` becomes `BERMUDIA__POSTGRES__ENABLE`). See the [Environment Variables](#environment-variables) table below for the full list.

   Example:

   ```bash
   export BERMUDIA__DEV_MODE=true
   export BERMUDIA__POSTGRES__ENABLE=true
   export BERMUDIA__POSTGRES__HOST=localhost
   export BERMUDIA__POSTGRES__PORT=5432
   export BERMUDIA__POSTGRES__USER=postgres
   export BERMUDIA__POSTGRES__PASS=password
   export BERMUDIA__POSTGRES__DB=bermudia
   export BERMUDIA__TOKEN_SIGNING_KEY=$(echo -n "your-secret-key" | base64)
   export BERMUDIA__BOT_TOKEN=your-bale-bot-token
   ```

   If `BERMUDIA__POSTGRES__ENABLE` is not `true`, the app falls back to SQLite (see `internal/repository`).

4. **Run the application**

   ```bash
   go run main.go
   ```

   The server listens on the hardcoded port `8080` (`http://localhost:8080`); there is currently no environment variable to change it.

## 🐳 Docker Deployment

### Build Docker Image

```bash
docker build -t bermudia-backend .
```

### Run with Docker

The `Dockerfile` builds the binary in a `golang:1.25` stage (with `CGO_ENABLED=1`, required by the SQLite driver) and runs it from a `debian:12.11` image, exposing port `8080`.

```bash
docker run -p 8080:8080 \
  -e BERMUDIA__POSTGRES__ENABLE=true \
  -e BERMUDIA__POSTGRES__HOST=host \
  -e BERMUDIA__POSTGRES__PORT=5432 \
  -e BERMUDIA__POSTGRES__USER=user \
  -e BERMUDIA__POSTGRES__PASS=password \
  -e BERMUDIA__POSTGRES__DB=bermudia \
  -e BERMUDIA__TOKEN_SIGNING_KEY=your-base64-secret-key \
  -e BERMUDIA__BOT_TOKEN=your-bale-bot-token \
  bermudia-backend
```

### Docker Compose

```yaml
version: '3.8'

services:
  backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      - BERMUDIA__POSTGRES__ENABLE=true
      - BERMUDIA__POSTGRES__HOST=db
      - BERMUDIA__POSTGRES__PORT=5432
      - BERMUDIA__POSTGRES__USER=postgres
      - BERMUDIA__POSTGRES__PASS=password
      - BERMUDIA__POSTGRES__DB=bermudia
      - BERMUDIA__TOKEN_SIGNING_KEY=your-base64-secret-key
      - BERMUDIA__BOT_TOKEN=your-bale-bot-token
    depends_on:
      - db
  
  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=bermudia
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## 📡 API Documentation

All player-facing routes are mounted under `/api/v1`. There is a separate, independently authenticated `/admin` tree (see below), and a plain `GET /health` check.

### Authentication

#### Login

```http
POST /api/v1/login
Content-Type: application/json

{
  "username": "player1",
  "password": "password123"
}
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

There is no self-service registration endpoint; users are created via the admin API (`POST /admin/users`).

Login checks the password with bcrypt and issues a JWT (`HS512`) containing `user_id` and `iat`, signed with `Config.TokenSigningKey` (base64-decoded). Tokens are accepted for 16 hours after issuance, checked against `iat` in `ValidateToken` (see `internal/service/auth.go`).

### WebSocket Endpoints

These are unauthenticated at the route level (auth, if any, happens in the handler itself) and live directly under `/api/v1`:

```text
GET /api/v1/events         # player update events
GET /api/v1/trade/events    # market/trade events
GET /api/v1/inbox/events    # inbox events
```

### Authenticated Endpoints (not paused)

All routes below require `Authorization: Bearer <token>` and are unaffected by the game-pause state.

```text
GET  /api/v1/me
GET  /api/v1/territories/{territoryID}
GET  /api/v1/territories/{territoryID}/players
GET  /api/v1/islands/{islandID}
GET  /api/v1/player
POST /api/v1/travel_check
POST /api/v1/refuel_check
POST /api/v1/anchor_check
POST /api/v1/migrate_check
POST /api/v1/unlock_treasure_check
POST /api/v1/trade/make_offer_check
GET  /api/v1/trade/offers
POST /api/v1/invest_check
GET  /api/v1/inbox/messages
```

`GET /api/v1/player` returns the player's `Player`/`FullPlayer` domain object, e.g.:

```json
{
  "atTerritory": "territory1",
  "atIsland": "island_start",
  "anchored": true,
  "fuel": 15,
  "fuelCap": 15,
  "coin": 100,
  "blueKey": 0,
  "redKey": 0,
  "goldenKey": 0,
  "masterKey": 0,
  "knowledgeBars": [],
  "books": []
}
```

### Authenticated + Paused-Gated Endpoints

These require the same auth as above, and additionally return `423 Locked` if the game is currently paused (`h.pauseCheckMiddleware`):

```text
POST /api/v1/answer/{inputID}
GET  /api/v1/answer/{inputID}/help
POST /api/v1/travel
POST /api/v1/refuel
POST /api/v1/anchor
POST /api/v1/migrate
POST /api/v1/unlock_treasure
POST /api/v1/trade/make_offer
POST /api/v1/trade/accept_offer
POST /api/v1/trade/delete_offer
POST /api/v1/invest
```

Example bodies (see `api/handler/player.go` for the full set):

```http
POST /api/v1/travel
Content-Type: application/json

{
  "fromIsland": "island_a",
  "toIsland": "island_b"
}
```

```http
POST /api/v1/refuel
Content-Type: application/json

{
  "amount": 5
}
```

```http
POST /api/v1/trade/make_offer
Content-Type: application/json

{
  "offered": { "items": [{ "type": "coin", "amount": 100 }] },
  "requested": { "items": [{ "type": "blueKey", "amount": 2 }] }
}
```

`answer/{inputID}` is submitted as `multipart/form-data` (a text answer or a file), not JSON — see `api/handler/island.go`.

### Admin API

A completely separate route tree under `/admin`, authenticated independently of the player JWT flow (backed by `Config.AdminUsername` / `Config.AdminPassword`, not the player user store):

```text
POST /admin/login
GET  /admin/territories
POST /admin/territories
GET  /admin/territories/{territoryID}/island_bindings
POST /admin/territories/{territoryID}/island_bindings
GET  /admin/books/{bookID}
GET  /admin/islands/{islandID}
POST /admin/islands/{islandID}/book
GET  /admin/pools
POST /admin/pools/{poolID}/books
GET  /admin/users
POST /admin/users
GET  /admin/game_state
POST /admin/game_state
POST /admin/broadcast
GET  /admin/connections
```

`POST /admin/territories` upserts a territory by `id` (`INSERT ... ON CONFLICT
(id) DO UPDATE`), so the same endpoint both creates and updates. `SetTerritory`
validates that the territory has a `startIsland` present in its island list,
that every island has an id and name, and that all edges / refuel / terminal /
prerequisite references point to islands in the list.

#### Game controls (`game_state`, `broadcast`, `connections`)

These expose over HTTP what previously lived only in the Telegram/Bale admin
bot, for the web admin panel:

```http
GET  /admin/game_state          → { "isPaused": bool }
POST /admin/game_state          body { "isPaused": bool } → new state
POST /admin/broadcast           body { "message": string } → { "sentTo": n }
GET  /admin/connections         → { "players": n, "market": n, "inbox": n }
```

`game_state` reads/writes the shared `GameStateStore` that the
`pauseCheckMiddleware` consults — pausing makes the paused-gated player
endpoints return `423 Locked`. `broadcast` sends an announcement to every
player's inbox. `connections` returns the live WebSocket hub counts.

#### Removal / deletion semantics

The admin API has **no general delete endpoints**, and the two removals it does
perform behave differently — worth knowing before relying on the admin panel to
"clean up" content:

- **Question / treasure inside a book** — truly deleted. `POST /admin/islands/{id}/book`
  (and the pool equivalent) diffs the submitted book against the stored one and
  runs `DELETE FROM questions` / `DELETE FROM treasures` for anything dropped.
  A book with zero questions is valid; the book row itself is never deleted.
  Caveat: the delete does **not** cascade to players' existing `answers` rows,
  so answers/corrections tied to a removed question are orphaned.
- **Island removed from a territory** — partial. `POST /admin/territories`
  overwrites the territory's island list, so the island disappears from the
  map, but a row remains in the `islands` table (islands are also tracked there
  via `ReserveIDForTerritory`, which has no delete counterpart). The id stays
  reserved to that territory: re-adding the same id to the **same** territory
  reclaims it, but reusing it in a **different** territory is rejected.
- **Book from a pool / a whole book / a whole territory** — not possible. There
  is no `RemoveBookFromPool`, no book delete, and no territory delete anywhere
  in the backend. A book can be moved between pools (re-`POST` to a different
  pool), but not removed from all pools.

### Health Check

```http
GET /health
```

Returns `200 OK` with body `OK`.

## 🔧 Configuration

Configuration is defined in `internal/config/config.go` and loaded in `internal/config/load.go` using [koanf](https://github.com/knadh/koanf)'s env provider. Every variable must be prefixed with `BERMUDIA__`, and nested fields (currently only `Postgres`) use `__` as the separator, e.g. `Config.Postgres.Host` -> `BERMUDIA__POSTGRES__HOST`.

### Environment Variables

| Variable | Config field | Description | Default |
|----------|-------------|--------------|---------|
| `BERMUDIA__DEV_MODE` | `DevMode` | Enables request logging middleware and allows mock data creation | `false` |
| `BERMUDIA__POSTGRES__ENABLE` | `Postgres.Enable` | Use PostgreSQL instead of SQLite | `false` |
| `BERMUDIA__POSTGRES__HOST` | `Postgres.Host` | PostgreSQL host | - |
| `BERMUDIA__POSTGRES__PORT` | `Postgres.Port` | PostgreSQL port | - |
| `BERMUDIA__POSTGRES__USER` | `Postgres.User` | PostgreSQL user | - |
| `BERMUDIA__POSTGRES__PASS` | `Postgres.Pass` | PostgreSQL password | - |
| `BERMUDIA__POSTGRES__DB` | `Postgres.DB` | PostgreSQL database name | - |
| `BERMUDIA__POSTGRES__SSL_MODE` | `Postgres.SSLMode` | PostgreSQL SSL mode | `disable` |
| `BERMUDIA__TOKEN_SIGNING_KEY` | `TokenSigningKey` | Base64-encoded key used to sign/verify JWTs | - |
| `BERMUDIA__MOCK_USERS_PASSWORD` | `MockUsersPassword` | Password assigned to generated mock users (dev mode only) | - |
| `BERMUDIA__BOT_TOKEN` | `BotToken` | Telegram/Bale bot API token (currently pointed at the Bale server, but Bale's bot API is Telegram-Bot-API-compatible) | - |
| `BERMUDIA__MIN_CORRECTION_DELAY` | `MinCorrectionDelay` | Minimum delay before a submitted answer is corrected | `10s` |
| `BERMUDIA__CORRECTION_JOB_INTERVAL` | `CorrectionJobInterval` | Interval of the background correction job | `10s` |
| `BERMUDIA__DEFAULT_CORRECTION_GROUP` | `DefaultCorrectionGroup` | Fallback chat ID for correction notifications | - |
| `BERMUDIA__CORRECTION_GROUPS` | `CorrectionGroupsStr` | Comma-separated `chatId:territory` pairs mapping territories to correction chat groups | - |
| `BERMUDIA__CONTENT_FILE_ID` | `ContentFileID` | File ID used by the bot for content delivery | - |
| `BERMUDIA__CORRECTION_REVERT_WINDOW` | `CorrectionRevertWindow` | Time window during which a correction can be reverted | - |
| `BERMUDIA__CREATE_MOCK` | `CreateMock` | If true (and `DevMode` is true), seeds mock game content on startup | `false` |
| `BERMUDIA__ADMINS_GROUP` | `AdminsGroup` | Chat ID for admin notifications | - |
| `BERMUDIA__ADMIN_USERNAME` | `AdminUsername` | Username for the `/admin` login | - |
| `BERMUDIA__ADMIN_PASSWORD` | `AdminPassword` | Password for the `/admin` login | - |

There is no variable to change the listening port; it is hardcoded to `:8080` in `api/handler/handler.go`.

### Database Configuration

The app connects to PostgreSQL when `BERMUDIA__POSTGRES__ENABLE=true` (see `repository.ConnectToPostgres`), otherwise it falls back to a local SQLite database (`repository.ConnectToSqlite`). There is no `scripts/init.sql`; schema setup is handled by the repository layer itself at startup.

## 🧪 Testing

There are currently no automated tests in the backend (no `*_test.go` files exist in this module). `go test ./...` will run successfully but exercises nothing. Contributions adding tests are welcome.

## 📦 Domain Models

### User

Represents a user account in the system.

### Player

Game state for a user, including resources and progress.

### Territory

A game zone containing multiple islands.

### Island

Individual locations with challenges and treasures.

### Treasure

Collectible items that require keys to unlock.

### Question

Educational challenges presented to players.

### Market Offer

Player-created trade proposals.

## 🛠️ Development

### Project Structure Explained

- **api/handler**: HTTP request handlers for each domain
- **api/hub**: WebSocket hub for managing real-time connections
- **internal/config**: Configuration loading and management
- **internal/domain**: Core business entities and interfaces
- **internal/repository**: Database access layer (DAL)
- **internal/service**: Business logic layer
- **internal/mock**: Mock game content seeded on startup when `BERMUDIA__DEV_MODE` and `BERMUDIA__CREATE_MOCK` are both true
- **adminbot**: Integration with messaging platforms

### Adding New Features

1. Define domain models in `internal/domain/`
2. Create repository interface and implementation in `internal/repository/`
3. Implement business logic in `internal/service/`
4. Add HTTP handlers in `api/handler/`
5. Update router in `main.go`

### Code Style

Follow Go conventions:

- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Write tests for new features
- Document exported functions

## 🔒 Security

- Player passwords are hashed using bcrypt (`internal/domain/user.go`)
- JWT (HS512) tokens are used for player authentication, valid for 16 hours after issuance (`internal/service/auth.go`)
- The `/admin` tree uses a separate, simpler auth scheme backed by a single configured username/password (`Config.AdminUsername` / `Config.AdminPassword`), not JWTs
- CORS is currently wide open: `corsMiddleware` in `api/handler/handler.go` sets `Access-Control-Allow-Origin: *` unconditionally — there is no `CORS_ORIGIN` (or similar) environment variable to restrict it
- Input validation on endpoints is minimal and handled ad hoc in each handler
- SQL injection prevention through parameterized queries in the repository layer

## 🐛 Troubleshooting

### Database Connection Issues

```bash
# Check PostgreSQL is running
systemctl status postgresql

# Test connection
psql -h localhost -U user -d bermudia -c "SELECT 1"
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

The port is hardcoded to `8080` and cannot currently be changed via configuration.

### JWT Token Issues

Ensure `BERMUDIA__TOKEN_SIGNING_KEY` is properly set (it must be valid base64) and stays the same across restarts/deployments, since it is used both to sign and verify tokens.

## 📚 Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Chi Router](https://github.com/go-chi/chi)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [JWT Introduction](https://jwt.io/introduction)

## 🤝 Contributing

Contributions are welcome! Please ensure:

- Code builds (`go build ./...`)
- New features include tests where practical (note: the module currently has no test suite)
- Code is properly formatted (`go fmt`)
- Commit messages are descriptive

## 📄 License

No LICENSE file is currently present in this repository.

---

For questions or support, contact the development team at Rasta.
