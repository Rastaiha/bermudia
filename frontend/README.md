# Bermudia Frontend

The Bermudia frontend is a modern, responsive Vue.js 3 application that provides an immersive gamified learning experience. Built with Vite, Tailwind CSS, and featuring real-time updates through WebSocket connections.

## 🧩 Composables

### useAudioPlayer

Drives a shared `<audio>` element based on the `set-audio-state` event and mute setting (see the "Audio System" section below for details):

```javascript
import { useAudioPlayer } from '@/composables/useAudioPlayer';

const audioPlayer = ref(null); // ref to an <audio> element
const { handleSongEnd } = useAudioPlayer(audioPlayer);
```

### useCountdownToNoon

Provides countdown timer functionality:

```javascript
import { useCountdownToNoon } from '@/composables/useCountdownToNoon';

const { timeUntilNoon, formatted } = useCountdownToNoon();
// formatted: "2h 15m 30s"
```

### useNow

Reactive current time:

```javascript
import { useNow } from '@/composables/useNow';

const { now } = useNow();
// Updates every second
```

## 🗺️ Routing

Routes are defined in `src/router/index.js`:

| Route                      | Component             | Description                          |
| -------------------------- | --------------------- | ------------------------------------ |
| `/`                        | redirects to `/login` | Root redirect                        |
| `/login`                   | Login.vue             | Authentication page                  |
| `/territory/:id`           | Territory.vue         | Territory map view                   |
| `/territory/:id/:islandId` | TerritoryIsland.vue   | Island detail view                   |
| `/admin/login`             | admin/AdminLogin.vue  | Admin authentication (separate auth) |
| `/admin`                   | admin/AdminSettings   | Admin dashboard — General Settings   |
| `/admin/users`             | admin/AdminUsers      | Admin — Users                        |
| `/admin/map`               | admin/AdminMap        | Admin — Map Editor                   |
| `/admin/islands`           | admin/AdminIslands    | Admin — Island Content               |
| `/admin/pools`             | admin/AdminPools      | Admin — Pools                        |

The `/admin/*` routes are contributed by `src/admin/router.js` (spread into
the main route table) and lazy-loaded, so they add nothing to the player
bundle. See the [Admin Panel](#-admin-panel) section below.

### Navigation Guards

Player routes with `meta: { requiresAuth: true }` require a player token.
`/admin/*` routes have their own independent guard and token — the global
`beforeEach` delegates any path starting with `/admin` to `adminGuard`:

```javascript
router.beforeEach((to, from, next) => {
    // Admin routes use their own token, independent of the player token.
    if (to.path.startsWith('/admin')) {
        if (adminGuard(to, from, next)) next();
        return;
    }

    const isLoggedIn = !!getToken();
    if (to.meta.requiresAuth && !isLoggedIn) {
        next({ name: 'Login' });
    } else {
        next();
    }
});
```

## 📱 State Management

### UI State Service

Centralized UI state management in `src/services/uiState.js`:

```javascript
import { uiState } from '@/services/uiState';

// Show/hide components
uiState.showMarket = true;
uiState.showBackpack = false;

// Modal state
uiState.modalData = { type: 'treasure', data: treasureInfo };
```

### Local Storage

Key data persisted client-side:

- `authToken`: player JWT, stored in `localStorage` if "remember me" is checked at login, otherwise in `sessionStorage`
- `adminToken`: admin JWT (separate from the player token), always in `localStorage`; see [Admin Panel](#-admin-panel)
- `adminSidebarCollapsed`: `"1"`/`"0"`, remembers whether the admin sidebar is collapsed

## 🛠️ Admin Panel

The admin panel is a self-contained area under `src/admin/`, isolated from
the player app. It shares the build, deploy, and Tailwind setup but keeps its
own routes, auth, API client, and styles so the two never bleed into each
other.

### Structure

```text
src/admin/
├── router.js                  # adminRoutes + adminGuard (spread into src/router)
├── layout/AdminLayout.vue     # collapsible sidebar shell (nav + logout)
├── pages/
│   ├── AdminLogin.vue         # admin login
│   ├── AdminSettings.vue      # pause/resume, broadcast, live connections
│   ├── AdminUsers.vue         # user table + create-user modal + manage-player-state modal
│   ├── AdminMap.vue           # map editor (islands/edges/roles + new territory)
│   ├── AdminIslands.vue       # island-content table view
│   └── AdminPools.vue         # pool bindings + pool books
├── components/
│   ├── AdminModal.vue         # shared modal (supports a `wide` variant)
│   ├── AssetPicker.vue        # image picker over public/images/**
│   └── IslandContentEditor.vue# reusable book editor (island + pool modes)
├── services/
│   ├── config.js              # ADMIN_ENDPOINTS (derives /admin base from API host)
│   └── api.js                 # admin fetch wrapper + adminToken helpers
└── styles/admin.css           # shared .admin-scope styles
```

### Auth

Admin auth is completely separate from the player flow. It hits
`POST /admin/login` (backed by the backend's `AdminUsername`/`AdminPassword`
config, **not** the player user store) and stores the returned JWT under
`adminToken`. A `401` from any admin request triggers `adminLogout()`. The
admin API tree lives at `/admin/*` (not `/api/v1`); `config.js` derives its
base URL from the player API host by stripping the `/api/v1` suffix, and it is
overridable via `VITE_ADMIN_BASE_URL`.

### Assets

`AssetPicker.vue` never uploads files. It enumerates images that already exist
in `public/images/**` via `import.meta.glob` and emits their public paths
(e.g. `/images/islands/educational/10.png`, `/images/backgrounds/territory/1.jpg`) —
the exact form the game stores in `iconAsset` / `backgroundAsset`. To offer a
new image in the pickers, drop the file into the matching `public/images/`
subdirectory.

### Sections

- **General Settings** — pause/resume the game, broadcast a message to every
  player's inbox, and view live WebSocket connection counts (auto-refreshed).
- **Users** — list users (name/username/meet link), create new ones, and
  **Manage state** per user: a modal that loads the player's current state and
  lets an admin edit location (territory + island), anchored flag, fuel, fuel
  capacity, coins, and the four key counts. Saving posts to
  `POST /admin/players/{userID}`; changes are applied immediately and pushed to
  the player if they are online.
- **Map Editor** — always-editing canvas over the territory background:
  drag/resize islands, pick icons/background, draw/delete edges, set the start
  island, toggle refuel/terminal roles, edit prerequisites, **create brand-new
  territories**, and jump to an island's content editor. A Save/Discard pair
  appears only when there are unsaved changes.
- **Island Content** — searchable table of every island; each opens the
  `IslandContentEditor` to edit its book (articles/iframes, questions,
  treasures).
- **Pools** — mark territory islands as normal vs. pooled and set the
  easy/medium/hard pool counts, plus author the books inside each pool (the
  content editor runs in "pool mode").

### Removal semantics (important)

The backend has no general delete endpoints, so the panel can only remove what
the backend supports, and some removals are partial. See the
[backend README](../backend/README.md#admin-api) for the authoritative
behavior. In short:

- Removing a **question/treasure** from a book truly deletes it (but does not
  clean up players' existing answers to that question).
- Removing an **island** from a territory drops it from the map, but leaves a
  stale row in the backend `islands` table (its id stays reserved to that
  territory and can't be reused by another territory).
- There is **no** way to remove a book from a pool, or delete a whole book or
  territory — those endpoints don't exist.

## 🎮 Game Flow

### 1. Authentication

```vue
<script setup>
import { ref } from 'vue';
import { login, getPlayer } from '@/services/api/index.js';
import { useRouter } from 'vue-router';

const username = ref('');
const password = ref('');
const router = useRouter();

const handleLogin = async () => {
    const result = await login(username.value, password.value);
    // result.token is already persisted to localStorage/sessionStorage by login()

    const playerData = await getPlayer();
    router.push({ name: 'Territory', params: { id: playerData.atTerritory } });
};
</script>
```

### 2. Territory Navigation

```vue
<script setup>
import { ref, onMounted } from 'vue';
import { getTerritory } from '@/services/api/index.js';

const territory = ref(null);

onMounted(async () => {
    territory.value = await getTerritory(territoryId);
});

const navigateToIsland = islandId => {
    router.push({ name: 'Island', params: { id: territoryId, islandId } });
};
</script>
```

### 3. Island Challenges

```vue
<script setup>
import { ref } from 'vue';
import { submitAnswer } from '@/services/api/index.js';
import { useToast } from 'vue-toastification';

const toast = useToast();
const island = ref(null);

const handleSubmit = async formData => {
    try {
        const result = await submitAnswer(island.value.id, formData);
        toast.success('Correct answer!');
    } catch (error) {
        toast.error(error.message || 'Error submitting answer');
    }
};
</script>
```

### 4. Market Trading

```vue
<script setup>
import { ref, onMounted } from 'vue';
import { getTradeOffers, acceptTradeOffer } from '@/services/api/index.js';
import { useMarketWebSocket } from '@/services/marketWebsocket';

const offers = ref([]);

onMounted(async () => {
    offers.value = await getTradeOffers();
    useMarketWebSocket(/* ... */);
});

const acceptOffer = async offerId => {
    await acceptTradeOffer(offerId);
};
</script>
```

## 🧪 Testing

There is currently no automated test suite (no unit or e2e test scripts are configured in `package.json`). Verification is done manually by running the dev server and exercising the app.

## 🔧 Development Tools

### Vite Configuration

Custom Vite configuration in `vite.config.js`:

```javascript
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';
import tailwindcss from '@tailwindcss/vite';
import svgLoader from 'vite-svg-loader';

export default defineConfig({
    resolve: {
        alias: {
            '@': path.resolve(__dirname, './src'),
        },
    },
    plugins: [vue(), tailwindcss(), svgLoader()],
});
```

Note: there is no dev-server API proxy configured — the app talks directly to the absolute API base URL defined in `src/services/api/config.js` (overridable via `VITE_API_BASE_URL`).

### Testing

Unit tests use [Vitest](https://vitest.dev/) (configured in the `test` block of `vite.config.js`, running in a `jsdom` environment). Test files live under the top-level `tests/` directory (mirroring the `src/` layout) as `*.test.js`, and currently target the pure-logic modules — `src/services/{storage,cost,notificationService}.js`, `src/services/api/config.js`, and the `src/composables/useNow.js` / `useCountdownToNoon.js` timers.

```bash
# Run the suite once
npm run test

# Watch mode
npm run test:watch
```

### ESLint Configuration

Code linting configured in `eslint.config.js`:

```bash
# Run linter (auto-fixes issues)
npm run lint

# Check only, no auto-fix
npm run lint:check
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

A pre-commit hook runs `lint-staged`, which lints and formats staged files automatically:

```bash
# .husky/pre-commit
cd frontend
npx lint-staged
```

`lint-staged` runs `eslint --fix` and `prettier --write` on staged `.js`/`.jsx`/`.vue` files, and `prettier --write` on staged `.css`/`.scss`/`.less`/`.html`/`.json`/`.md` files. There is no type checking (the project is plain JavaScript, not TypeScript).

## 📊 Performance Optimization

### Code Splitting

Routes are lazy-loaded for better performance:

```javascript
const routes = [
    {
        path: '/territory/:id',
        component: () => import('@/pages/Territory.vue'),
    },
];
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

### Console Logging

There is no built-in `.env`-driven debug flag; the app does not read any `import.meta.env` variables. Add `console.log` statements directly where needed during development and remove them before committing.

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
import { ref, computed, onMounted } from 'vue';
import { getPlayer } from '@/services/api/index.js';

// 2. Props
const props = defineProps({
    id: String,
    data: Object,
});

// 3. Emits
const emit = defineEmits(['update', 'close']);

// 4. Reactive state
const loading = ref(false);
const player = ref(null);

// 5. Computed properties
const isReady = computed(() => !loading.value && player.value !== null);

// 6. Methods
const fetchData = async () => {
    loading.value = true;
    player.value = await getPlayer();
    loading.value = false;
};

// 7. Lifecycle hooks
onMounted(() => {
    fetchData();
});
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

- JWT token stored in `localStorage` or `sessionStorage` (depending on "remember me")
- Token included in the `Authorization: Bearer <token>` header on authenticated requests
- A `401` response automatically triggers logout (token cleared, redirect to `/login`)
- There is no automatic token refresh

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

### API Base URL

The API and WebSocket base URLs default to the production endpoints hardcoded in `src/services/api/config.js`. To point at a different backend, set `VITE_API_BASE_URL`/`VITE_WS_BASE_URL` in a `.env` file (see `.env.example`) and rebuild:

```bash
VITE_API_BASE_URL=https://bermudia-api-internal.darkube.ir/api/v1
VITE_WS_BASE_URL=wss://bermudia-api-internal.darkube.ir/api/v1
```

### Nginx Configuration

The included `nginx.conf` handles:

- Static file serving with SPA routing fallback (`try_files $uri /index.html`)
- Gzip compression for common text/asset types

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
    "rewrites": [{ "source": "/(.*)", "destination": "/" }]
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
    npm run lint
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

Check that the backend WebSocket server is running and that the `VITE_WS_BASE_URL` (or the default in `src/services/api/config.js`) is correct.

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

3. **Configure the API base URL**

    The app defaults to the production backend. To point at a different backend (e.g. local dev), copy `.env.example` to `.env` and set:

    ```bash
    VITE_API_BASE_URL=http://localhost:8080/api/v1
    VITE_WS_BASE_URL=ws://localhost:8080/api/v1
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
            - '80:80'
        depends_on:
            - backend
```

Note: since `VITE_API_BASE_URL`/`VITE_WS_BASE_URL` (see `src/services/api/config.js`) are baked into the JS bundle at build time, passing runtime environment variables to the container has no effect — the target backend URL must be set before `npm run build` runs (i.e. before/during the Docker image build).

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
<ChallengeBox :challenge="currentChallenge" @submit="handleSubmit" />
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
<Inbox :messages="messages" @mark-read="handleMarkRead" />
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

`useAudioPlayer` is a composable that drives a single shared `<audio>` element (passed in as a ref) based on an `eventBus` event (`set-audio-state`) and the current mute setting; it does not expose `play`/`pause`/`setVolume` methods directly:

```javascript
import { useAudioPlayer } from '@/composables/useAudioPlayer';

const audioPlayer = ref(null); // ref to an <audio> element
const { handleSongEnd } = useAudioPlayer(audioPlayer);
```

Playback is toggled elsewhere by emitting on the shared event bus (`eventBus.emit('set-audio-state', 'play' | 'pause')`), and mute state lives in `src/services/audio/settings.js` (`audioSettings.isMuted`).

### Audio Configuration

The track list lives in `src/services/audio/playlist.js`:

```javascript
export const playlist = [
    {
        title: '1',
        url: '/audio/1.mp3',
        duration: 141,
    },
    // ...more tracks
];
```

Track selection/progression logic lives in `src/services/audio/radioService.js` (`getCurrentTrack()`), not in the composable itself.

## 🎨 Styling

### Tailwind CSS

The project uses Tailwind CSS v4, loaded via the `@tailwindcss/vite` plugin (see `vite.config.js`) and imported in `src/styles/main.css` with `@import 'tailwindcss';`. A `tailwind.config.js` is still used for theme extensions (ESM syntax, not CommonJS):

```javascript
/** @type {import('tailwindcss').Config} */
export default {
    content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
    theme: {
        extend: {
            fontFamily: {
                vazir: ['Vazirmatn', 'sans-serif'],
            },
            transitionTimingFunction: {
                'smooth-expand': 'cubic-bezier(0.25, 0.1, 0.25, 1.0)',
            },
            keyframes: {
                'boat-animation': {
                    '0%, 100%': {
                        transform: 'translate(0, 0) rotate(10deg) scale(0.3)',
                    },
                    '35%': {
                        transform:
                            'translate(0.02px, 0.01px) rotate(-10deg) scale(0.3)',
                    },
                    '70%': {
                        transform:
                            'translate(-0.02px, 0.01px) rotate(3deg) scale(0.3)',
                    },
                },
            },
            animation: {
                boat: 'boat-animation 10s linear infinite',
            },
        },
    },
    plugins: [],
};
```

Custom fonts (`@font-face` for Pelak and Vazirmatn) are declared directly in `src/styles/main.css`, not in the Tailwind config.

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

Located in `src/services/api/` (`index.js` for request functions, `config.js` for the base URL and endpoint URLs). It is a thin `fetch` wrapper exposing one named async function per endpoint, not a generic `get`/`post` client:

```javascript
import { getPlayer, submitAnswer } from '@/services/api/index.js';

const player = await getPlayer();
const result = await submitAnswer('123', formData);
```

### WebSocket Service

Real-time communication through WebSockets, split by concern:

- `src/services/websocket.js` exposes `usePlayerWebSocket(player, territoryId, route, router)` for the main player/territory event stream
- `src/services/marketWebsocket.js` exposes `useMarketWebSocket(...)` for market updates
- `src/services/inboxWebsocket.js` handles inbox message updates

```javascript
import { usePlayerWebSocket } from '@/services/websocket';

usePlayerWebSocket(player, territoryId, route, router);
```

### Event Bus

A shared `mitt` instance for component communication, exported as the default export:

```javascript
import eventBus from '@/services/eventBus';

// Emit event
eventBus.emit('set-audio-state', 'play');

// Listen to event
eventBus.on('set-audio-state', state => {
    console.log('Audio state:', state);
});
```

### Notification Service

`src/services/notificationService.js` tracks read/unread inbox message state (backed by `localStorage`), not toast notifications:

```javascript
import { notificationService } from '@/services/notificationService';

notificationService.setReceivedMessages(messages);
notificationService.markAllAsSeen();
// notificationService.hasUnreadMessages is a computed boolean
```

Toast notifications (success/error/warning messages) are shown directly with `vue-toastification`'s `useToast()` composable, used ad hoc in components — there is no dedicated wrapper service for it.
