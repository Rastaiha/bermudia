# Bermudia Frontend

The Bermudia frontend is a modern, responsive Vue.js 3 application that provides an immersive gamified learning experience. Built with Vite, Tailwind CSS, and featuring real-time updates through WebSocket connections.

## 🧩 Composables

### useAudioPlayer

Manages audio playback:

```javascript
import { useAudioPlayer } from '@/composables/useAudioPlayer'

const {
  currentTrack,
  isPlaying,
  volume,
  play,
  pause,
  toggle,
  setVolume,
  next,
  previous
} = useAudioPlayer()
```

### useCountdownToNoon

Provides countdown timer functionality:

```javascript
import { useCountdownToNoon } from '@/composables/useCountdownToNoon'

const { timeUntilNoon, formatted } = useCountdownToNoon()
// formatted: "2h 15m 30s"
```

### useNow

Reactive current time:

```javascript
import { useNow } from '@/composables/useNow'

const { now } = useNow()
// Updates every second
```

## 🗺️ Routing

Routes are defined in `src/router/index.js`:

| Route | Component | Description |
|-------|-----------|-------------|
| `/` | Login.vue | Authentication page |
| `/territory/:id` | Territory.vue | Territory map view |
| `/territory/:territoryId/island/:islandId` | TerritoryIsland.vue | Island detail view |

### Navigation Guards

Protected routes require authentication:

```javascript
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  if (to.path !== '/' && !token) {
    next('/')
  } else {
    next()
  }
})
```

## 📱 State Management

### UI State Service

Centralized UI state management in `src/services/uiState.js`:

```javascript
import { uiState } from '@/services/uiState'

// Show/hide components
uiState.showMarket = true
uiState.showBackpack = false

// Modal state
uiState.modalData = { type: 'treasure', data: treasureInfo }
```

### Local Storage

Key data persisted in localStorage:

- `token`: JWT authentication token
- `userId`: Current user ID
- `audioSettings`: Volume and mute preferences
- `uiPreferences`: UI customization settings

## 🎮 Game Flow

### 1. Authentication

```vue
<script setup>
import { ref } from 'vue'
import api from '@/services/api'

const username = ref('')
const password = ref('')

const login = async () => {
  const response = await api.post('/auth/login', {
    username: username.value,
    password: password.value
  })
  
  localStorage.setItem('token', response.token)
  router.push('/territory/1')
}
</script>
```

### 2. Territory Navigation

```vue
<script setup>
import { ref, onMounted } from 'vue'
import api from '@/services/api'

const territories = ref([])
const currentTerritory = ref(null)

onMounted(async () => {
  territories.value = await api.get('/territories')
  currentTerritory.value = territories.value[0]
})

const navigateToIsland = (islandId) => {
  router.push(`/territory/${currentTerritory.value.id}/island/${islandId}`)
}
</script>
```

### 3. Island Challenges

```vue
<script setup>
import { ref } from 'vue'
import api from '@/services/api'
import { showNotification } from '@/services/notificationService'

const challenge = ref(null)
const answer = ref('')

const submitAnswer = async () => {
  try {
    const result = await api.post(`/islands/${islandId}/challenge`, {
      answer: answer.value,
      question_id: challenge.value.id
    })
    
    if (result.correct) {
      showNotification({
        type: 'success',
        message: 'Correct answer! +' + result.reward.coins + ' coins'
      })
    } else {
      showNotification({
        type: 'error',
        message: 'Incorrect answer. Try again!'
      })
    }
  } catch (error) {
    showNotification({
      type: 'error',
      message: 'Error submitting answer'
    })
  }
}
</script>
```

### 4. Market Trading

```vue
<script setup>
import { ref, onMounted } from 'vue'
import { connectWebSocket } from '@/services/marketWebsocket'

const offers = ref([])
const ws = ref(null)

onMounted(() => {
  // Connect to market websocket
  ws.value = connectWebSocket()
  
  ws.value.on('market_update', (newOffers) => {
    offers.value = newOffers
  })
})

const acceptOffer = async (offerId) => {
  await api.post(`/market/offers/${offerId}/accept`)
  showNotification({
    type: 'success',
    message: 'Trade completed successfully!'
  })
}
</script>
```

## 🧪 Testing

### Run Tests

```bash
npm run test
# or
yarn test
```

### Component Testing

```bash
npm run test:unit
# or
yarn test:unit
```

### E2E Testing

```bash
npm run test:e2e
# or
yarn test:e2e
```

## 🔧 Development Tools

### Vite Configuration

Custom Vite configuration in `vite.config.js`:

```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
```

### ESLint Configuration

Code linting configured in `eslint.config.js`:

```bash
# Run linter
npm run lint

# Fix auto-fixable issues
npm run lint:fix
```

### Prettier Configuration

Code formatting with Prettier:

```bash
# Format code
npm run format

# Check formatting
npm run format:check
```

### Husky Git Hooks

Pre-commit hooks ensure code quality:

```bash
# Hooks run automatically on git commit
# - Linting
# - Formatting
# - Type checking
```

## 📊 Performance Optimization

### Code Splitting

Routes are lazy-loaded for better performance:

```javascript
const routes = [
  {
    path: '/territory/:id',
    component: () => import('@/pages/Territory.vue')
  }
]
```

### Image Optimization

- Use WebP format when possible
- Lazy load images with `loading="lazy"`
- Optimize image sizes for different screen resolutions

### Bundle Size

Monitor bundle size:

```bash
npm run build -- --mode production
# Check dist/ folder size
```

## 🐛 Debugging

### Vue Devtools

Install Vue Devtools browser extension for debugging:

- Component inspection
- Vuex state (if used)
- Event tracking
- Performance profiling

### Debug Mode

Enable debug mode in `.env`:

```env
VITE_DEBUG_MODE=true
```

### Console Logging

Service logs can be enabled:

```javascript
// In any service file
const DEBUG = import.meta.env.VITE_DEBUG_MODE === 'true'

if (DEBUG) {
  console.log('API Request:', endpoint, data)
}
```

## 🌐 Browser Support

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

For older browser support, consider adding polyfills.

## 📝 Code Style Guide

### Vue Component Structure

```vue
<script setup>
// 1. Imports
import { ref, computed, onMounted } from 'vue'
import api from '@/services/api'

// 2. Props
const props = defineProps({
  id: String,
  data: Object
})

// 3. Emits
const emit = defineEmits(['update', 'close'])

// 4. Reactive state
const loading = ref(false)
const items = ref([])

// 5. Computed properties
const filteredItems = computed(() => {
  return items.value.filter(item => item.active)
})

// 6. Methods
const fetchData = async () => {
  loading.value = true
  items.value = await api.get('/items')
  loading.value = false
}

// 7. Lifecycle hooks
onMounted(() => {
  fetchData()
})
</script>

<template>
  <!-- Template content -->
</template>

<style scoped>
/* Component-specific styles */
</style>
```

### Naming Conventions

- **Components**: PascalCase (e.g., `PlayerInfo.vue`)
- **Composables**: camelCase with 'use' prefix (e.g., `useAudioPlayer.js`)
- **Services**: camelCase (e.g., `notificationService.js`)
- **Constants**: UPPER_SNAKE_CASE (e.g., `MAX_FUEL_CAPACITY`)

## 🔐 Security Best Practices

### Authentication

- JWT tokens stored in localStorage
- Tokens included in Authorization header
- Automatic token refresh before expiration
- Logout clears all stored tokens

### XSS Prevention

- Vue's template escaping by default
- Sanitize user input
- Use `v-html` only with trusted content

### CORS

Backend CORS configuration must allow frontend origin:

```javascript
// Backend configuration
CORS_ALLOWED_ORIGINS=http://localhost:5173,https://bermudia.example.com
```

## 🚀 Deployment

### Production Build

```bash
npm run build
```

Output in `dist/` directory.

### Environment Variables

Set production environment variables:

```env
VITE_API_BASE_URL=https://api.bermudia.example.com/api
VITE_WS_BASE_URL=wss://api.bermudia.example.com/api/ws
```

### Nginx Configuration

The included `nginx.conf` handles:

- Static file serving
- Gzip compression
- Caching headers
- SPA routing fallback

### Deploy to Static Hosting

#### Netlify

```bash
# netlify.toml
[build]
  command = "npm run build"
  publish = "dist"

[[redirects]]
  from = "/*"
  to = "/index.html"
  status = 200
```

#### Vercel

```json
{
  "buildCommand": "npm run build",
  "outputDirectory": "dist",
  "rewrites": [
    { "source": "/(.*)", "destination": "/" }
  ]
}
```

#### AWS S3 + CloudFront

1. Build the project
2. Upload `dist/` to S3 bucket
3. Configure CloudFront distribution
4. Set up custom domain

## 📚 Additional Resources

- [Vue.js Documentation](https://vuejs.org/)
- [Vite Documentation](https://vitejs.dev/)
- [Tailwind CSS Documentation](https://tailwindcss.com/)
- [Vue Router Documentation](https://router.vuejs.org/)

## 🤝 Contributing

### Development Workflow

1. Create a feature branch

   ```bash
   git checkout -b feature/new-feature
   ```

2. Make changes and test locally

3. Run linter and formatter

   ```bash
   npm run lint:fix
   npm run format
   ```

4. Commit with descriptive message

   ```bash
   git commit -m "feat: add new island navigation feature"
   ```

5. Push and create pull request

### Component Development Guidelines

- Keep components small and focused
- Use composition API for new components
- Write props validation
- Document complex components
- Include examples in comments

## 📄 License

This project is open source. See LICENSE file for details.

## 🆘 Troubleshooting

### Development Server Won't Start

```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

### WebSocket Connection Failed

Check that backend WebSocket server is running and `VITE_WS_BASE_URL` is correct.

### Build Errors

```bash
# Clear Vite cache
rm -rf node_modules/.vite
npm run build
```

### Font Not Loading

Ensure fonts are in `src/assets/fonts/` and referenced in `src/styles/main.css`.

---

For questions or support, contact the development team at Rasta.

## 🎨 Features

- **Interactive Map System**: Explore territories and islands with visual navigation
- **Real-time Updates**: WebSocket integration for live market and inbox updates
- **Responsive Design**: Fully responsive UI built with Tailwind CSS
- **Rich Media**: Audio player with playlist management and background music
- **Component-based Architecture**: Modular Vue 3 components with Composition API
- **Smooth Animations**: Engaging UI transitions and effects
- **Starry Night Background**: Immersive visual experience
- **Custom Persian Fonts**: Pelak and Vazirmatn font families

## 🏗️ Architecture

```text
frontend/
├── public/                    # Static assets
│   ├── audio/                # Audio files
│   ├── images/               # Game images
│   │   ├── backgrounds/      # Territory and island backgrounds
│   │   ├── icons/           # UI icons
│   │   ├── islands/         # Island-specific images
│   │   ├── profiles/        # Player profile images
│   │   ├── ships/           # Ship/vehicle images
│   │   └── territories/     # Territory maps
│   └── map-generator.html   # Map generation tool
├── src/
│   ├── assets/              # Build-time assets
│   │   └── fonts/          # Custom fonts
│   ├── components/          # Vue components
│   │   ├── common/         # Reusable UI components
│   │   ├── features/       # Feature-specific components
│   │   └── layout/         # Layout components
│   ├── composables/         # Vue composables
│   ├── pages/              # Page components
│   ├── router/             # Vue Router configuration
│   ├── services/           # API and WebSocket services
│   └── styles/             # Global styles
│   ├── App.vue             # Root component
│   └── main.js             # Application entry point
├── nginx.conf              # Nginx configuration for production
├── Dockerfile              # Docker configuration
├── vite.config.js          # Vite configuration
├── tailwind.config.js      # Tailwind CSS configuration
└── package.json            # Dependencies
```

## 🚀 Getting Started

### Prerequisites

- Node.js 18+ or higher
- npm 9+ or yarn 1.22+
- Modern web browser (Chrome, Firefox, Safari, Edge)

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/Rastaiha/bermudia.git
   cd bermudia/frontend
   ```

2. **Install dependencies**

   ```bash
   npm install
   # or
   yarn install
   ```

3. **Set up environment variables**

   Create a `.env` file in the frontend directory:

   ```env
   # API Configuration
   VITE_API_BASE_URL=http://localhost:8080/api
   VITE_WS_BASE_URL=ws://localhost:8080/api/ws
   
   # Feature Flags
   VITE_ENABLE_AUDIO=true
   VITE_ENABLE_ANIMATIONS=true
   
   # Debug
   VITE_DEBUG_MODE=false
   ```

4. **Run development server**

   ```bash
   npm run dev
   # or
   yarn dev
   ```

   The application will be available at `http://localhost:5173`

### Build for Production

```bash
npm run build
# or
yarn build
```

The built files will be in the `dist/` directory.

### Preview Production Build

```bash
npm run preview
# or
yarn preview
```

## 🐳 Docker Deployment

### Build Docker Image

```bash
docker build -t bermudia-frontend .
```

### Run with Docker

```bash
docker run -p 80:80 bermudia-frontend
```

### Docker Compose

```yaml
version: '3.8'

services:
  frontend:
    build: .
    ports:
      - "80:80"
    environment:
      - VITE_API_BASE_URL=http://backend:8080/api
    depends_on:
      - backend
```

## 📦 Component Library

### Common Components

#### ConfirmModal.vue

Modal dialog for user confirmations.

```vue
<ConfirmModal
  :show="showConfirm"
  title="Confirm Action"
  message="Are you sure?"
  @confirm="handleConfirm"
  @cancel="handleCancel"
/>
```

#### CostlyButton.vue

Button component that displays resource costs.

```vue
<CostlyButton
  :cost="{ coins: 100, fuel: 10 }"
  :disabled="!canAfford"
  @click="handlePurchase"
>
  Purchase Item
</CostlyButton>
```

#### LoadingIndicator.vue

Loading spinner component.

```vue
<LoadingIndicator :loading="isLoading" />
```

#### InfoBox.vue

Information display box with customizable content.

```vue
<InfoBox title="Island Info" :data="islandData" />
```

#### PlayerInventoryBar.vue

Display player's current resources.

```vue
<PlayerInventoryBar :player="playerData" />
```

### Feature Components

#### Map System

**MapView.vue**: Main map component displaying territories and islands.

```vue
<MapView
  :territories="territories"
  :current-location="currentIsland"
  @island-click="handleIslandClick"
/>
```

**IslandInfoBox.vue**: Display island information when selected.

**RefuelBox.vue**: Refueling interface for islands with fuel stations.

#### Island Challenges

**ChallengeBox.vue**: Display and answer educational challenges.

```vue
<ChallengeBox
  :challenge="currentChallenge"
  @submit="handleSubmit"
/>
```

**Treasure.vue**: Treasure collection interface.

**TreasureRewardModal.vue**: Display rewards after opening treasures.

#### Market System

**Market.vue**: Main market component.

**Trade.vue**: Create and manage trade offers.

**TradeOfferCard.vue**: Display individual trade offers.

```vue
<TradeOfferCard
  :offer="tradeOffer"
  @accept="handleAccept"
  @cancel="handleCancel"
/>
```

#### Player Features

**Backpack.vue**: Inventory management interface.

**Bookshelf.vue**: View collected educational content.

**Brain.vue**: Knowledge and achievements display.

**Casino.vue**: Mini-game interface (if enabled).

**Inbox.vue**: Message inbox with real-time updates.

```vue
<Inbox
  :messages="messages"
  @mark-read="handleMarkRead"
/>
```

**NotificationItem.vue**: Individual notification display.

### Layout Components

**LoginTemplate.vue**: Authentication page layout.

**PlayerInfo.vue**: Player information sidebar.

**Toolbar.vue**: Main navigation toolbar.

**UserProfile.vue**: User profile management.

## 🎵 Audio System

The frontend includes a sophisticated audio system:

### Features

- Background music playlist
- Sound effects for actions
- Volume control
- Mute toggle
- Automatic track progression

### Usage

```javascript
import { useAudioPlayer } from '@/composables/useAudioPlayer'

const { play, pause, setVolume } = useAudioPlayer()

// Play specific track
play('background-1')

// Adjust volume (0-1)
setVolume(0.5)
```

### Audio Configuration

Edit `src/services/audio/playlist.js`:

```javascript
export const playlist = [
  {
    id: 'background-1',
    title: 'Ocean Theme',
    src: '/audio/ocean-theme.mp3',
    loop: true
  },
  // Add more tracks...
]
```

## 🎨 Styling

### Tailwind CSS

The project uses Tailwind CSS for styling. Configuration in `tailwind.config.js`:

```javascript
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: '#1e40af',
        secondary: '#7c3aed',
        // Add custom colors
      },
      fontFamily: {
        pelak: ['Pelak', 'sans-serif'],
        vazir: ['Vazirmatn', 'sans-serif'],
      }
    },
  },
  plugins: [],
}
```

### Custom Fonts

Persian fonts are included:

- **Pelak**: Black, ExtraBold, Regular
- **Vazirmatn**: Full weight range (Thin to Black)

Usage in components:

```vue
<template>
  <div class="font-vazir">
    <!-- Content with Vazirmatn font -->
  </div>
</template>
```

## 🔌 Services

### API Service

Located in `src/services/api/`:

```javascript
import api from '@/services/api'

// Make authenticated requests
const player = await api.get('/player')
const result = await api.post('/islands/123/challenge', { answer: '42' })
```

### WebSocket Service

Real-time communication through WebSockets:

```javascript
import { connectWebSocket } from '@/services/websocket'

const ws = connectWebSocket()

ws.on('market_update', (data) => {
  console.log('New market offer:', data)
})

ws.on('inbox_message', (message) => {
  console.log('New message:', message)
})
```

### Event Bus

Global event system for component communication:

```javascript
import { eventBus } from '@/services/eventBus'

// Emit event
eventBus.emit('player:updated', playerData)

// Listen to event
eventBus.on('player:updated', (data) => {
  console.log('Player updated:', data)
})
```

### Notification Service

Display toast notifications:

```javascript
import { showNotification } from '@/services/notificationService'

showNotification({
  type: 'success',
  title: 'Success!',
  message: 'Challenge completed'
})
```
