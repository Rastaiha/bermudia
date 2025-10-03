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
- **External Services**: Integration with Gofino and Jitsi for extended functionality

## 📚 Documentation

- [Backend Documentation](./backend/README.md) - API endpoints, services, and backend architecture
- [Frontend Documentation](./frontend/README.md) - Component structure, setup, and development guide
- [API Documentation](./docs/api.md) - Detailed API reference
- [Phase 2 Documentation](./docs/phase2.md) - Feature roadmap and implementation details
- [Phase 3 Documentation](./docs/phase3.md) - Advanced features and future plans

## 🚀 Quick Start

### Prerequisites

- Go 1.24+ (for backend)
- Node.js 18+ and npm/yarn (for frontend)
- PostgreSQL 15+ or SQLite (for database)
- Docker and Docker Compose (optional, for containerized deployment)

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
cp .env.example .env  # Configure your environment variables
go run main.go
```

The backend will start on `http://localhost:8080` by default.

#### 3. Frontend Setup

```bash
cd frontend
npm install
cp .env.example .env  # Configure your environment variables
npm run dev
```

The frontend will start on `http://localhost:5173` by default.

For detailed setup instructions, see the [Backend README](./backend/README.md) and [Frontend README](./frontend/README.md).

## 🐳 Docker Deployment

Build and run with Docker Compose:

```bash
docker-compose up -d
```

This will start both the backend and frontend services with proper networking.

### Individual Service Deployment

**Backend:**

```bash
cd backend
docker build -t bermudia-backend .
docker run -p 8080:8080 bermudia-backend
```

**Frontend:**

```bash
cd frontend
docker build -t bermudia-frontend .
docker run -p 80:80 bermudia-frontend
```

## 🔧 Configuration

### Backend Configuration

Environment variables can be set in `backend/.env`:

- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for JWT token generation
- `PORT` - Server port (default: 8080)
- `CORS_ORIGIN` - Allowed CORS origins
- `BOT_TOKEN` - Telegram/Bale bot token for admin features

### Frontend Configuration

Environment variables can be set in `frontend/.env`:

- `VITE_API_BASE_URL` - Backend API URL
- `VITE_WS_BASE_URL` - WebSocket server URL

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

### Bale/Telegram Bot (Correction Bot)

The admin bot handles:

- Challenge correction and verification
- Player communication
- Administrative notifications

Configure bot token in backend environment variables.

### Gofino Integration

Used for extended gameplay features and social interactions.

### Jitsi Integration

Provides video conferencing capabilities for multiplayer features.

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

### Running Tests

**Backend:**

```bash
cd backend
go test ./...
```

**Frontend:**

```bash
cd frontend
npm run test
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

Key API endpoints:

- `POST /api/auth/login` - User authentication
- `GET /api/territories` - List all territories
- `GET /api/islands/:id` - Get island details
- `POST /api/islands/:id/challenge` - Submit challenge answer
- `GET /api/player` - Get player information
- `GET /api/market` - Get market offers
- `WS /api/ws` - WebSocket connection for real-time updates

For complete API documentation, see [API Documentation](./docs/api.md).

## 👥 Authors

- **Seyed Ali Hosseini** - Core Developer
- **Meysam Bavi** - Core Developer
- **Roham Ghasemi** - Core Developer
- **Fardad Arab** - Core Developer

Developed by the team at [Rasta](https://rastaiha.ir).

## 📄 License

This project is open source. Please check the LICENSE file for more details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 🐛 Bug Reports

If you encounter any bugs or issues, please report them on the [Issues](https://github.com/Rastaiha/bermudia/issues) page.

## 📧 Contact

For questions and support, please contact the development team at Rasta.

---

Made with ❤️ by the Rasta Development Team
