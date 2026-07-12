// The admin API lives at `/admin/*` on the same host as the player API,
// which is served under `/api/v1`. Derive the admin base by stripping the
// `/api/v1` suffix from the configured player API base URL.
const API_BASE_URL =
    import.meta.env.VITE_API_BASE_URL ||
    'https://bermudia-api-internal.darkube.ir/api/v1';

const ADMIN_BASE_URL =
    import.meta.env.VITE_ADMIN_BASE_URL ||
    API_BASE_URL.replace(/\/api\/v1\/?$/, '') + '/admin';

export const ADMIN_ENDPOINTS = {
    login: `${ADMIN_BASE_URL}/login`,

    // General settings
    gameState: `${ADMIN_BASE_URL}/game_state`,
    broadcast: `${ADMIN_BASE_URL}/broadcast`,
    connections: `${ADMIN_BASE_URL}/connections`,

    // Users
    users: `${ADMIN_BASE_URL}/users`,
    player: userId => `${ADMIN_BASE_URL}/players/${userId}`,

    // Territories & map
    territories: `${ADMIN_BASE_URL}/territories`,
    territoryIslandBindings: id =>
        `${ADMIN_BASE_URL}/territories/${id}/island_bindings`,

    // Islands & books
    island: id => `${ADMIN_BASE_URL}/islands/${id}`,
    islandBook: id => `${ADMIN_BASE_URL}/islands/${id}/book`,
    book: id => `${ADMIN_BASE_URL}/books/${id}`,
    pools: `${ADMIN_BASE_URL}/pools`,
    poolBooks: id => `${ADMIN_BASE_URL}/pools/${id}/books`,
};
