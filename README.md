# Bermudia

Bermudia is an open-source gamified learning platform that combines education with adventure. Players explore a virtual archipelago, solve challenges, collect treasures, and engage in a dynamic trading system while learning new concepts through an immersive gaming experience.

## 🎮 Features

- **Island Exploration**: Navigate through multiple territories with unique islands containing educational challenges
- **Treasure System**: Unlock treasures using different types of keys (Blue, Red, Golden, Master)
- **Dynamic Market**: Real-time trading system with websocket-based market operations
- **Challenge System**: Multiple difficulty levels (Easy, Medium, Hard) with educational content
- **Player Progression**: Track achievements, collect resources, and manage inventory
- **Real-time Notifications**: Inbox system for player communications and updates
- **Audio Experience**: Immersive background music and sound effects
- **Admin Bot Integration**: Telegram/Bale bot for correction and administration

## 🏗️ Architecture

Bermudia follows a modern full-stack architecture:

- **Backend**: Go (Golang) with Chi router, WebSocket support, and PostgreSQL/SQLite
- **Frontend**: Vue.js 3 with Vite, Tailwind CSS, and responsive design
- **Real-time Communication**: WebSocket-based events for market, inbox, and game state
- **External Services**: Telegram/Bale bot API for admin/correction workflows; per-user configurable meet links for challenge help

## 📚 Documentation

- [Backend Documentation](./backend/README.md) - API endpoints, services, and backend architecture
- [Frontend Documentation](./frontend/README.md) - Component structure, setup, and development guide
- [API Documentation](./docs/api.md) - Detailed API reference
- [Phase 2 Documentation](./docs/phase2.md) - Feature roadmap and implementation details
- [Phase 3 Documentation](./docs/phase3.md) - Advanced features and future plans

## 🚀 Quick Start

### Prerequisites

- Go 1.25+ (for backend)
- Node.js 18+ and npm (for frontend)
- PostgreSQL 15+ or SQLite (for database)
- Docker (optional, for containerized deployment; there is no docker-compose.yml in this repo)

### Installation

#### 1. Clone the Repository

```bash
git clone https://github.com/Rastaiha/bermudia.git
cd bermudia
```

#### 2. Backend Setup

```bash
cd backend
go mod download
export BERMUDIA__TOKEN_SIGNING_KEY=$(echo -n "your-secret-key" | base64)
go run main.go
```

Configuration is via environment variables prefixed `BERMUDIA__` (not a `.env` file) — see the [Backend README](./backend/README.md#configuration) for the full list. The backend listens on the hardcoded port `8080`.

#### 3. Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

The frontend will start on `http://localhost:5173` by default. There is no frontend `.env` mechanism — the API/WebSocket base URLs are hardcoded in `frontend/src/services/api/base_url.js`.

For detailed setup instructions, see the [Backend README](./backend/README.md) and [Frontend README](./frontend/README.md).

## 🐳 Docker Deployment

There is no `docker-compose.yml` in this repository. Each service has its own `Dockerfile` and must be built/run individually.

**Backend:**

```bash
cd backend
docker build -t bermudia-backend .
docker run -p 8080:8080 -e BERMUDIA__TOKEN_SIGNING_KEY=your-base64-secret-key bermudia-backend
```

**Frontend:**

```bash
cd frontend
docker build -t bermudia-frontend .
docker run -p 80:80 bermudia-frontend
```

## 🔧 Configuration

### Backend Configuration

The backend has no `.env` file. Configuration is loaded from environment variables prefixed `BERMUDIA__` (via koanf), e.g. `BERMUDIA__POSTGRES__ENABLE`, `BERMUDIA__TOKEN_SIGNING_KEY`, `BERMUDIA__BOT_TOKEN`. CORS is currently hardcoded to allow all origins. See the [Backend README](./backend/README.md#environment-variables) for the full variable table.

### Frontend Configuration

The frontend has no `.env` file either. The API and WebSocket base URLs are hardcoded in `frontend/src/services/api/base_url.js` and require a rebuild to change.

## 🎯 Game Mechanics

### Resources

Players collect and manage various resources:

- **Coins**: Primary currency for trading and purchases
- **Fuel**: Required for island navigation
- **Keys**: Blue, Red, Golden, and Master keys for unlocking treasures
- **Books**: Educational content collected from islands

### Territories and Islands

The game world is organized into territories, each containing multiple islands:

- **Educational Islands**: Contain learning challenges
- **Challenge Islands**: Test player knowledge
- **Refuel Stations**: Replenish fuel supplies
- **Terminal Islands**: Story progression points
- **Final Islands**: End-game content

### Trading System

Real-time marketplace where players can:

- Create trade offers
- Accept offers from other players
- Exchange resources dynamically
- View market history

## 🔌 External Dependencies

### Telegram/Bale Bot (Correction Bot)

The admin bot uses the `go-telegram/bot` library, which implements the Telegram Bot API; since Bale's bot API is practically the same, the bot is currently pointed at the Bale server (`https://tapi.bale.ai`), but it could equally be pointed at Telegram's API. It handles:

- Challenge correction and verification
- Player communication
- Administrative notifications

Configure the bot token via `BERMUDIA__BOT_TOKEN`.

### Per-user meet link

Each user can have a `meetLink` (e.g. a video call link) configured; when a player requests help answering a challenge, the backend returns that user's `meetLink` (see `internal/service/island.go`). There is no Gofino or Jitsi integration in the codebase.

## 🛠️ Development

### Project Structure

```text
bermudia/
├── backend/           # Go backend service
│   ├── api/          # API handlers and WebSocket hub
│   ├── internal/     # Core business logic
│   │   ├── domain/   # Domain models
│   │   ├── repository/ # Data access layer
│   │   └── service/  # Business logic
│   └── adminbot/     # Bot integration
├── frontend/         # Vue.js frontend
│   ├── src/
│   │   ├── components/ # Vue components
│   │   ├── pages/    # Page components
│   │   ├── services/ # API and WebSocket services
│   │   └── router/   # Vue Router configuration
│   └── public/       # Static assets
└── docs/            # Documentation
```

### Code Formatting

**Backend:**

```bash
go fmt ./...
```

**Frontend:**

```bash
npm run format
```

## 📖 API Overview

Key API endpoints (all under `/api/v1` unless noted):

- `POST /api/v1/login` - User authentication (no self-service registration; users are created via the admin API)
- `GET /api/v1/territories/{territoryID}` - Get territory details
- `GET /api/v1/islands/{islandID}` - Get island details
- `POST /api/v1/answer/{inputID}` - Submit a challenge answer (multipart form data)
- `GET /api/v1/player` - Get player information
- `GET /api/v1/trade/offers` - Get market trade offers
- `GET /api/v1/events`, `/api/v1/trade/events`, `/api/v1/inbox/events` - WebSocket connections for real-time updates
- `/admin/*` - Separate, independently authenticated admin API

For complete endpoint documentation, see the [Backend README](./backend/README.md#api-documentation).

## 🐛 Bug Reports

If you encounter any bugs or issues, please report them on the [Issues](https://github.com/Rastaiha/bermudia/issues) page.

## 📧 Contact

For questions and support, please contact the development team at Rasta.

---

Made with ❤️ by the [Rasta](https://rastaiha.ir) Development Team
