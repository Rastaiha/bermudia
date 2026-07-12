# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository layout

This is a monorepo with two independently deployed apps:

- `backend/` — Go API server, WebSocket hub, and admin Telegram/Bale bot
- `frontend/` — Vue 3 SPA (Vite + Tailwind CSS 4)

The `README.md` files at the repo root, `backend/`, and `frontend/` are reliable references for setup, routes, and config. If anything in them ever looks inconsistent with the source, trust the source and update the READMEs to match.

## Commands

### Backend (from `backend/`)

```bash
go run main.go          # run the server
go build ./...           # compile
go vet ./...             # static checks
gofmt -l .                # list unformatted files (there is no test suite in this repo)
```

There are no `*_test.go` files in the backend currently — don't assume `go test ./...` will exercise anything.

### Frontend (from `frontend/`)

```bash
npm install
npm run dev              # Vite dev server, http://localhost:5173
npm run build             # production build to dist/
npm run lint              # eslint --fix
npm run lint:check        # eslint, no autofix
npm run format             # prettier --write
npm run format:check
```

There is no test script in `package.json`. Husky runs `lint-staged` (eslint + prettier) on pre-commit.

## Configuration (backend)

Config is **not** loaded from a `.env` file. It's loaded via `koanf` from environment variables prefixed `BERMUDIA__`, with `__` as the nesting separator, mapped onto `internal/config/Config` (see `internal/config/config.go` and `load.go`). Example: `BERMUDIA__POSTGRES__ENABLE=true` sets `Config.Postgres.Enable`. Defaults live in `defaultConfig()`.

Key fields: `DevMode`, `Postgres.*` (falls back to embedded SQLite if `Postgres.Enable` is false), `TokenSigningKey` (base64, used for JWT signing, HS512, tokens valid 16h from `iat`), `BotToken` (Telegram/Bale bot — uses `go-telegram/bot`, currently pointed at Bale's server URL `https://tapi.bale.ai` since Bale's bot API is Telegram-Bot-API-compatible), `CreateMock` (when true + `DevMode`, seeds mock game content on boot from `internal/mock/data`), `CorrectionGroupsStr` (parsed into a `territory -> chatID` map for the correction bot), `AdminUsername`/`AdminPassword` (separate admin auth, unrelated to player JWT auth). Note: the listen port (`8080`) is hardcoded in `handler.go`, not configurable. CORS is hardcoded to allow all origins.

## Backend architecture

Layering follows `internal/domain` → `internal/repository` → `internal/service` → `api/handler`, wired up in `main.go`:

- **`internal/domain`** — entities/value types (Player, Island, Territory, Treasure, Market, Inbox, Invest, Question, User, Cost, Reward...) and a small typed `Error` (with a `Reason()` used to map to HTTP status codes).
- **`internal/repository`** — SQL access (Postgres or SQLite, same schema/queries — see `db.go` for connection setup). One repository file per aggregate.
- **`internal/service`** — business logic. Services take repos as constructor args and are wired together in `main.go`; note the cross-service callback wiring, e.g. `islandService.OnNewPortableIsland(playerService.HandleNewPortableIsland)` and the `h.playerService.On*` event hooks registered in `handler.Start()` — this is how player/trade/inbox state changes get pushed out over WebSocket.
- **`api/handler`** — chi router + HTTP handlers, one file per resource (`auth.go`, `island.go`, `player.go`, `territory.go`, `admin.go`, `events.go`). `response.go` has the shared `sendResult`/`handleError` helpers that translate `domain.Error` reasons into HTTP status codes.
- **`api/hub`** — a generic WebSocket connection hub (`hub.go`) reused for three separate channels: player events, trade/market events, inbox events. Each gets its own `Hub` instance in `handler.Handler`.
- **`adminbot`** — Telegram/Bale bot (via `go-telegram/bot`, pointed at Bale's API server since its bot API is Telegram-compatible) used for challenge correction workflows and admin notifications; it holds references into the handler and several services.
- **`internal/mock`** — seeds a full game world (territories, islands, books, pool settings, users) from JSON/zip fixtures under `internal/mock/data` when `DevMode && CreateMock` — this is the easiest way to get a locally runnable game without hand-authoring content.

### Actual API surface

Routes are mounted at `/api/v1/*` (not `/api/*`), defined in `api/handler/handler.go`:

- `POST /api/v1/login` and WS endpoints (`/events`, `/trade/events`, `/inbox/events`) are unauthenticated at the route level (auth happens via token in the connection).
- An authenticated-but-not-paused group (`authMiddleware`) covers reads and "check" endpoints (`*_check` routes let the frontend validate an action's cost/outcome before committing it).
- An authenticated-and-paused group (`authMiddleware` + `pauseCheckMiddleware`) covers all mutating actions (`/travel`, `/refuel`, `/anchor`, `/migrate`, `/unlock_treasure`, `/trade/*`, `/invest`, `/answer/{inputID}`) — the pause check presumably gates actions during maintenance/events.
- A completely separate `/admin/*` tree with its own login and `authMiddleware` (`adminHandler`, backed by `AdminUsername`/`AdminPassword` config, not JWT player auth) for content management (territories, island bindings, books, pools, users).
- `GET /health` for liveness checks.

The `*_check` / non-`_check` pairing (e.g. `travel_check` + `travel`) is a recurring pattern — check the paired handler when modifying either.

## Frontend architecture

- **`src/pages`** — route-level views (`Login.vue`, `Territory.vue`, `TerritoryIsland.vue`), wired in `src/router/index.js`. Root `/` redirects to `/login`; protected routes are marked with `meta: { requiresAuth: true }` and gated in a `router.beforeEach` guard.
- **`src/components/{common,features,layout}`** — common (generic UI), features (domain-specific: market, treasures, inbox, challenges), layout (toolbar, player info, etc.).
- **`src/services/api`** — a named-export `fetch` wrapper (`login`, `getPlayer`, `getTerritory`, `submitAnswer`, `getTradeOffers`, `acceptTradeOffer`, ...), not an axios-style `api.get`/`api.post` client. Base URLs default to the production endpoints hardcoded in `src/services/api/config.js`, overridable via `VITE_API_BASE_URL`/`VITE_WS_BASE_URL` (see `.env.example`) — changing the backend target requires a rebuild since Vite inlines env vars at build time.
- **`src/services`** (rest) — `websocket.js` + `marketWebsocket.js` + `inboxWebsocket.js` expose `usePlayerWebSocket`/`useMarketWebSocket`/etc. composable-style connectors (three separate WS connections mirroring the backend's three hubs), `uiState.js` (centralized reactive UI state), `eventBus.js` (default-exported `mitt` instance), `notificationService.js` (tracks read/unread inbox messages — not a toast wrapper; toasts are driven directly via `vue-toastification`'s `useToast()`), `cost.js`, `glossary.js`.
- Auth token key is `authToken`, stored in `localStorage` if "remember me" was checked at login, otherwise `sessionStorage`. A 401 response triggers logout; there is no token-refresh mechanism.
- **`src/composables`** — `useAudioPlayer` (drives an `<audio>` ref + eventBus, not a full play/pause/next/previous API), `useCountdownToNoon`, `useNow`.
- Path alias `@` → `src/` (configured in `vite.config.js`). No dev-server proxy for `/api` — the app always calls the absolute base URL.
- Styling: Tailwind CSS 4 via `@tailwindcss/vite` plugin (ESM `export default` config, not classic `module.exports`/PostCSS), plus `tailwindcss-rtl` (the app supports Persian/RTL — custom Pelak and Vazirmatn fonts are loaded via `@font-face` in `src/styles/main.css`).
- SVGs are imported as Vue components via `vite-svg-loader`.
- **`src/admin`** — the admin panel, deliberately isolated from the player app (own routes, auth, API client, and styles under this one directory). `admin/router.js` exports `adminRoutes` + `adminGuard`, spread into `src/router/index.js`; the global `beforeEach` delegates any `/admin/*` path to `adminGuard`. Admin routes are lazy-loaded so they don't touch the player bundle. Auth is separate from the player flow: it calls `POST /admin/login` (backed by backend `AdminUsername`/`AdminPassword`, not the player user store) and stores the JWT under `adminToken` (always `localStorage`, distinct from the player `authToken`). The admin API lives at `/admin/*`, not `/api/v1` — `admin/services/config.js` derives its base by stripping `/api/v1` from the API host (overridable via `VITE_ADMIN_BASE_URL`). Sections: General Settings (pause/broadcast/connections), Users (view/create, plus per-user "Manage state" to edit a player's location/anchored/fuel/coin/keys via `GET`/`POST /admin/players/{userID}`), Map Editor (island/edge/role editing + create new territories), Island Content, Pools. `AssetPicker.vue` picks from images already in `public/images/**` via `import.meta.glob` (no upload) and emits their public paths — add a new selectable image by dropping it into the right `public/images/` subdir. `IslandContentEditor.vue` is reused in both island mode and pool mode.

## Backend removal semantics

The backend has **no general delete endpoints**; deletions are limited and some are partial. Removing a question/treasure from a book truly deletes it (via `POST /admin/islands/{id}/book` diffing) but does not clean up players' existing `answers` for it. Removing an island from a territory drops it from the map JSON but leaves a stale row in the `islands` table — its id stays reserved to that territory (reusable there, rejected elsewhere). There is no way to remove a book from a pool, or delete a whole book or territory. See `backend/README.md` → "Removal / deletion semantics" for detail.

## Game domain concepts

Useful vocabulary when reading code or handlers: players navigate a map of **territories**, each containing **islands** (educational, challenge, refuel, terminal, final). Islands present **questions**; correct answers can unlock **treasures** (gated by Blue/Red/Golden/Master keys) and yield **rewards**. Players hold **coins**/**fuel**/keys/**books** as inventory, trade them in a real-time **market** (offer/accept/cancel), and receive updates through the **inbox**. The admin bot handles manual **correction** of free-form answers for the correction/education workflow, routed to per-territory chat groups (`CorrectionGroups`).
