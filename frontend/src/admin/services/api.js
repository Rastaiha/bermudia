import { ADMIN_ENDPOINTS } from './config.js';

const ADMIN_TOKEN_KEY = 'adminToken';

export const setAdminToken = token => {
    if (token) {
        localStorage.setItem(ADMIN_TOKEN_KEY, token);
    } else {
        localStorage.removeItem(ADMIN_TOKEN_KEY);
    }
};

export const getAdminToken = () => localStorage.getItem(ADMIN_TOKEN_KEY);

export const adminLogout = () => {
    setAdminToken(null);
    window.location.pathname = '/admin/login';
};

const getAuthHeaders = () => {
    const token = getAdminToken();
    return token ? { Authorization: `Bearer ${token}` } : {};
};

const handleResponse = async response => {
    let data = {};
    try {
        data = await response.json();
    } catch {
        // non-JSON body
    }
    if (!response.ok || data.ok === false) {
        if (response.status === 401) {
            adminLogout();
        }
        throw new Error(data.error || `Request failed (${response.status})`);
    }
    return data.result;
};

const adminGet = async url => {
    const response = await fetch(url, {
        method: 'GET',
        headers: getAuthHeaders(),
    });
    return handleResponse(response);
};

const adminPost = async (url, body) => {
    const response = await fetch(url, {
        method: 'POST',
        headers:
            body !== undefined
                ? { ...getAuthHeaders(), 'Content-Type': 'application/json' }
                : getAuthHeaders(),
        body: body !== undefined ? JSON.stringify(body) : undefined,
    });
    return handleResponse(response);
};

// --- Auth ---
export const adminLogin = async (username, password) => {
    const result = await adminPost(ADMIN_ENDPOINTS.login, {
        username,
        password,
    });
    if (result?.token) {
        setAdminToken(result.token);
    }
    return result;
};

// --- General settings ---
export const getGameState = async () => adminGet(ADMIN_ENDPOINTS.gameState);

export const setGameState = async isPaused =>
    adminPost(ADMIN_ENDPOINTS.gameState, { isPaused });

export const broadcast = async message =>
    adminPost(ADMIN_ENDPOINTS.broadcast, { message });

export const getConnections = async () => adminGet(ADMIN_ENDPOINTS.connections);

// --- Users ---
export const getUsers = async () => adminGet(ADMIN_ENDPOINTS.users);

export const createUser = async user => adminPost(ADMIN_ENDPOINTS.users, user);

// --- Territories & map ---
export const getTerritories = async () => adminGet(ADMIN_ENDPOINTS.territories);

export const setTerritory = async territory =>
    adminPost(ADMIN_ENDPOINTS.territories, territory);

// --- Islands & books ---
export const getIslandHeader = async id => adminGet(ADMIN_ENDPOINTS.island(id));

export const getBook = async id => adminGet(ADMIN_ENDPOINTS.book(id));

export const setIslandBook = async (id, input) =>
    adminPost(ADMIN_ENDPOINTS.islandBook(id), input);

// --- Pools & territory island bindings ---
export const getTerritoryIslandBindings = async territoryId =>
    adminGet(ADMIN_ENDPOINTS.territoryIslandBindings(territoryId));

export const setTerritoryIslandBindings = async bindings =>
    adminPost(
        ADMIN_ENDPOINTS.territoryIslandBindings(bindings.territoryId),
        bindings
    );

export const getPools = async () => adminGet(ADMIN_ENDPOINTS.pools);

export const setPoolBook = async (poolId, input) =>
    adminPost(ADMIN_ENDPOINTS.poolBooks(poolId), input);
