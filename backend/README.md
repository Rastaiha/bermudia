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
│   │   ├── auth.go           # Authentication endpoints
│   │   ├── events.go         # Event handling
│   │   ├── island.go         # Island-related endpoints
│   │   ├── player.go         # Player management
│   │   ├── territory.go      # Territory endpoints
│   │   └── response.go       # Response utilities
│   └── hub/                  # WebSocket hub
│       └── hub.go
└── internal/                 # Internal packages
    ├── config/               # Configuration management
    ├── domain/               # Domain models and interfaces
    ├── mock/                 # Mock data for testing
    ├── repository/           # Data access layer
    └── service/              # Business logic layer
```

## 🚀 Getting Started

### Prerequisites

- Go 1.24 or higher
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

   Create a `.env` file in the backend directory:

   ```env
   # Server Configuration
   PORT=8080
   HOST=0.0.0.0
   
   # Database Configuration
   DATABASE_TYPE=postgres  # or "sqlite"
   DATABASE_URL=postgresql://user:password@localhost:5432/bermudia
   
   # For SQLite (alternative to PostgreSQL)
   # DATABASE_TYPE=sqlite
   # DATABASE_PATH=./data/bermudia.db
   
   # JWT Configuration
   JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
   JWT_EXPIRATION=24h
   
   # CORS Configuration
   CORS_ORIGIN=http://localhost:5173
   CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
   
   # Bot Configuration (Telegram/Bale)
   BOT_TOKEN=your-bot-token-here
   BOT_TYPE=telegram  # or "bale"
   
   # Game Configuration
   STARTING_COINS=100
   STARTING_FUEL=50
   MASTER_KEY_CHANCE=0.3
   TREASURE_VALUE_THRESHOLD=80
   
   # External Services
   GOFINO_API_URL=https://api.gofino.example.com
   GOFINO_API_KEY=your-gofino-api-key
   JITSI_DOMAIN=meet.jit.si
   
   # Development
   DEBUG=true
   LOG_LEVEL=info
   ```

4. **Initialize the database**

   The application will automatically create tables on first run. For manual setup:

   ```bash
   # PostgreSQL
   psql -U user -d bermudia -f scripts/init.sql
   
   # SQLite
   sqlite3 ./data/bermudia.db < scripts/init.sql
   ```

5. **Run the application**

   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8080`

### Development Mode

For development with auto-reload:

```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Run with air
air
```

## 🐳 Docker Deployment

### Build Docker Image

```bash
docker build -t bermudia-backend .
```

### Run with Docker

```bash
docker run -p 8080:8080 \
  -e DATABASE_URL=postgresql://user:password@host:5432/bermudia \
  -e JWT_SECRET=your-secret-key \
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
      - DATABASE_URL=postgresql://postgres:password@db:5432/bermudia
      - JWT_SECRET=your-secret-key
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

### Authentication

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "player1",
  "password": "password123"
}
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "player1",
    "email": "player1@example.com"
  }
}
```

#### Register

```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "newplayer",
  "email": "newplayer@example.com",
  "password": "password123"
}
```

### Player Management

#### Get Player Info

```http
GET /api/player
Authorization: Bearer <token>
```

Response:

```json
{
  "id": 1,
  "username": "player1",
  "coins": 500,
  "fuel": 75,
  "keys": {
    "blue": 5,
    "red": 3,
    "golden": 1,
    "master": 0
  },
  "current_territory": "territory1",
  "current_island": "island_start"
}
```

#### Update Player Inventory

```http
PUT /api/player/inventory
Authorization: Bearer <token>
Content-Type: application/json

{
  "coins": 100,
  "fuel": -10,
  "keys": {
    "blue": 1
  }
}
```

### Territories and Islands

#### List Territories

```http
GET /api/territories
Authorization: Bearer <token>
```

#### Get Territory Details

```http
GET /api/territories/:id
Authorization: Bearer <token>
```

#### Get Island Details

```http
GET /api/islands/:id
Authorization: Bearer <token>
```

#### Get Island Challenge

```http
GET /api/islands/:id/challenge
Authorization: Bearer <token>
```

#### Submit Challenge Answer

```http
POST /api/islands/:id/challenge
Authorization: Bearer <token>
Content-Type: application/json

{
  "answer": "42",
  "question_id": "q123"
}
```

### Treasure System

#### List Player Treasures

```http
GET /api/player/treasures
Authorization: Bearer <token>
```

#### Unlock Treasure

```http
POST /api/treasures/:id/unlock
Authorization: Bearer <token>
Content-Type: application/json

{
  "use_master_key": false
}
```

### Market System

#### Get Market Offers

```http
GET /api/market
Authorization: Bearer <token>
```

#### Create Trade Offer

```http
POST /api/market/offers
Authorization: Bearer <token>
Content-Type: application/json

{
  "offer": {
    "coins": 100,
    "fuel": 20
  },
  "request": {
    "blue_key": 2
  }
}
```

#### Accept Trade Offer

```http
POST /api/market/offers/:id/accept
Authorization: Bearer <token>
```

#### Cancel Trade Offer

```http
DELETE /api/market/offers/:id
Authorization: Bearer <token>
```

### Inbox System

#### Get Messages

```http
GET /api/inbox
Authorization: Bearer <token>
```

#### Mark as Read

```http
PUT /api/inbox/:id/read
Authorization: Bearer <token>
```

### WebSocket Connection

Connect to WebSocket for real-time updates:

```text
WS /api/ws?token=<jwt-token>
```

WebSocket message types:

- `market_update` - New market offers or trades
- `inbox_message` - New inbox messages
- `player_update` - Player state changes
- `event_notification` - Game events

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | No |
| `HOST` | Server host | `0.0.0.0` | No |
| `DATABASE_TYPE` | Database type (postgres/sqlite) | `postgres` | Yes |
| `DATABASE_URL` | PostgreSQL connection string | - | If using PostgreSQL |
| `DATABASE_PATH` | SQLite database file path | `./data/bermudia.db` | If using SQLite |
| `JWT_SECRET` | JWT signing secret | - | Yes |
| `JWT_EXPIRATION` | Token expiration time | `24h` | No |
| `CORS_ORIGIN` | Primary CORS origin | `*` | No |
| `BOT_TOKEN` | Bot token for admin features | - | No |
| `DEBUG` | Enable debug mode | `false` | No |

### Database Configuration

#### PostgreSQL Setup

```sql
CREATE DATABASE bermudia;
CREATE USER bermudia_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE bermudia TO bermudia_user;
```

#### SQLite Setup

SQLite databases are created automatically. Ensure the data directory exists:

```bash
mkdir -p data
```

## 🧪 Testing

### Run All Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test -cover ./...
```

### Run Tests for Specific Package

```bash
go test ./internal/service/...
```

### Generate Coverage Report

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

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
- **internal/mock**: Mock implementations for testing
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

- All passwords are hashed using bcrypt
- JWT tokens for authentication
- CORS protection enabled
- Input validation on all endpoints
- SQL injection prevention through parameterized queries

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

### JWT Token Issues

Ensure `JWT_SECRET` is properly set and matches between backend and any JWT-generating services.

## 📚 Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Chi Router](https://github.com/go-chi/chi)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [JWT Introduction](https://jwt.io/introduction)

## 🤝 Contributing

Contributions are welcome! Please ensure:

- Code passes all tests
- New features include tests
- Code is properly formatted (`go fmt`)
- Commit messages are descriptive

## 📄 License

This project is open source. See LICENSE file for details.

---

For questions or support, contact the development team at Rasta.
