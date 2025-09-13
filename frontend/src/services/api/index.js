import { API_ENDPOINTS } from '@/services/api/config.js';

// --- Helper Functions ---
const getAuthHeaders = () => {
    const token = getToken();
    return token ? { Authorization: `Bearer ${token}` } : {};
};

const handleResponse = async response => {
    const data = await response.json();
    if (!response.ok || data.ok === false) {
        if (response.status === 401) {
            logout();
        }
        throw new Error(data.error || 'An unknown error occurred');
    }
    return data.result;
};

// --- Token Management ---
export const setToken = token => {
    if (token) {
        localStorage.setItem('authToken', token);
    } else {
        localStorage.removeItem('authToken');
    }
};

export const getToken = () => localStorage.getItem('authToken');

export const logout = () => {
    setToken(null);
    window.location.pathname = '/login';
};

// --- API Functions ---
export const login = async (username, password) => {
    const response = await fetch(API_ENDPOINTS.login, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
    });
    const result = await handleResponse(response);
    if (result.token) {
        setToken(result.token);
    }
    return result;
};

export const getMe = async () => {
    const response = await fetch(API_ENDPOINTS.getMe, {
        method: 'GET',
        headers: getAuthHeaders(),
    });
    return handleResponse(response);
};

export const getPlayer = async () => {
    const response = await fetch(API_ENDPOINTS.getPlayer, {
        method: 'GET',
        headers: getAuthHeaders(),
    });
    return handleResponse(response);
};

export const getTerritory = async id => {
    const response = await fetch(API_ENDPOINTS.getTerritory(id), {
        headers: getAuthHeaders(),
    });
    return handleResponse(response);
};

export const travelCheck = async (from, dest) => {
    const response = await fetch(API_ENDPOINTS.travelCheck, {
        method: 'POST',
        headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
        body: JSON.stringify({ fromIsland: from, toIsland: dest }),
    });
    return handleResponse(response);
};

export const anchorCheck = async island => {
    const response = await fetch(API_ENDPOINTS.anchorCheck, {
        method: 'POST',
        headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
        body: JSON.stringify({ island: island }),
    });
    return handleResponse(response);
};

export const travelTo = async (from, dest) => {
    const response = await fetch(API_ENDPOINTS.travelTo, {
        method: 'POST',
        headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
        body: JSON.stringify({ fromIsland: from, toIsland: dest }),
    });
    return handleResponse(response);
};

export const dropAnchorAtIsland = async currentIsland => {
    const response = await fetch(API_ENDPOINTS.dropAnchor, {
        method: 'POST',
        headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
        body: JSON.stringify({ island: currentIsland }),
    });
    return handleResponse(response);
};

export const getIsland = async id => {
    const response = await fetch(API_ENDPOINTS.getIsland(id), {
        headers: getAuthHeaders(),
    });
    return handleResponse(response);
};

export const submitAnswer = async (id, formData) => {
    const response = await fetch(API_ENDPOINTS.submitAnswer(id), {
        method: 'POST',
        headers: getAuthHeaders(),
        body: formData,
    });
    return handleResponse(response);
};

export const refuelCheck = async () => {
    const response = await fetch(API_ENDPOINTS.refuelCheck, {
        method: 'POST',
        headers: getAuthHeaders(),
    });
    return handleResponse(response);
};

export const buyFuel = async amount => {
    const response = await fetch(API_ENDPOINTS.buyFuel, {
        method: 'POST',
        headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
        body: JSON.stringify({ amount }),
    });
    return handleResponse(response);
};

export const migrateCheck = async () => {
    const response = await fetch(API_ENDPOINTS.migrateCheck, {
        method: 'POST',
        headers: getAuthHeaders(),
    });
    return handleResponse(response);
};

export const migrateTo = async territory => {
    const response = await fetch(API_ENDPOINTS.migrate, {
        method: 'POST',
        headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
        body: JSON.stringify({ toTerritory: territory }),
    });
    return handleResponse(response);
};

export const treasureCheck = async treasureId => {
    const response = await fetch(API_ENDPOINTS.treasureCheck, {
        method: 'POST',
        headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
        body: JSON.stringify({ treasureID: treasureId }),
    });
    return handleResponse(response);
};

export const treasureUnlock = async treasureId => {
    const response = await fetch(API_ENDPOINTS.treasureUnlock, {
        method: 'POST',
        headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
        body: JSON.stringify({ treasureID: treasureId }),
    });
    return handleResponse(response);
};

export const makeTradeOfferCheck = async () => {
    const response = await fetch(API_ENDPOINTS.makeOfferCheck, {
        method: 'POST',
        headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
    });
    return handleResponse(response);
};

export const makeTradeOffer = async (offered, requested) => {
    const response = await fetch(API_ENDPOINTS.makeOffer, {
        method: 'POST',
        headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
        body: JSON.stringify({ offered: offered, requested: requested }),
    });
    return handleResponse(response);
};

export const acceptTradeOffer = async offerID => {
    const response = await fetch(API_ENDPOINTS.acceptOffer, {
        method: 'POST',
        headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
        body: JSON.stringify({ offerID: offerID }),
    });
    return handleResponse(response);
};

export const deleteTradeOffer = async offerID => {
    const response = await fetch(API_ENDPOINTS.deleteOffer, {
        method: 'POST',
        headers: { ...getAuthHeaders(), 'Content-Type': 'application/json' },
        body: JSON.stringify({ offerID: offerID }),
    });
    return handleResponse(response);
};

export const getTradeOffers = async (offset = 0, limit = 5, by = null) => {
    const response = await fetch(
        `${API_ENDPOINTS.getOffers(offset, limit, by)}`,
        {
            method: 'GET',
            headers: getAuthHeaders(),
        }
    );
    return handleResponse(response);
};

export const getInboxMessages = async (offset = null, limit = 15) => {
    const response = await fetch(
        API_ENDPOINTS.getInboxMessages(offset, limit),
        {
            method: 'GET',
            headers: getAuthHeaders(),
        }
    );
    return handleResponse(response);
};
