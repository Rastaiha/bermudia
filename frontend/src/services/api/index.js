import { API_ENDPOINTS } from '@/services/api/config.js';

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

const apiGet = async url => {
    const response = await fetch(url, {
        method: 'GET',
        headers: getAuthHeaders(),
    });
    return handleResponse(response);
};

const apiPost = async (url, body) => {
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

const TOKEN_KEY = 'authToken';

export const setToken = (token, remember = true) => {
    localStorage.removeItem(TOKEN_KEY);
    sessionStorage.removeItem(TOKEN_KEY);
    if (token) {
        (remember ? localStorage : sessionStorage).setItem(TOKEN_KEY, token);
    }
};

export const getToken = () =>
    localStorage.getItem(TOKEN_KEY) || sessionStorage.getItem(TOKEN_KEY);

export const logout = () => {
    setToken(null);
    window.location.pathname = '/login';
};

export const login = async (username, password, remember = true) => {
    const result = await apiPost(API_ENDPOINTS.login, { username, password });
    if (result.token) {
        setToken(result.token, remember);
    }
    return result;
};

export const getMe = async () => apiGet(API_ENDPOINTS.getMe);

export const getPlayer = async () => apiGet(API_ENDPOINTS.getPlayer);

export const getTerritory = async id => apiGet(API_ENDPOINTS.getTerritory(id));

export const getPlayersLocation = async id =>
    apiGet(API_ENDPOINTS.getPlayersLocation(id));

export const travelCheck = async (from, dest) =>
    apiPost(API_ENDPOINTS.travelCheck, { fromIsland: from, toIsland: dest });

export const anchorCheck = async island =>
    apiPost(API_ENDPOINTS.anchorCheck, { island: island });

export const travelTo = async (from, dest) =>
    apiPost(API_ENDPOINTS.travelTo, { fromIsland: from, toIsland: dest });

export const dropAnchorAtIsland = async currentIsland =>
    apiPost(API_ENDPOINTS.dropAnchor, { island: currentIsland });

export const getIsland = async id => apiGet(API_ENDPOINTS.getIsland(id));

export const submitAnswer = async (id, formData) => {
    const response = await fetch(API_ENDPOINTS.submitAnswer(id), {
        method: 'POST',
        headers: getAuthHeaders(),
        body: formData,
    });
    return handleResponse(response);
};

export const requestHelp = async inputId =>
    apiGet(API_ENDPOINTS.requestHelp(inputId));

export const refuelCheck = async () => apiPost(API_ENDPOINTS.refuelCheck);

export const buyFuel = async amount =>
    apiPost(API_ENDPOINTS.buyFuel, { amount });

export const migrateCheck = async () => apiPost(API_ENDPOINTS.migrateCheck);

export const migrateTo = async territory =>
    apiPost(API_ENDPOINTS.migrate, { toTerritory: territory });

export const treasureCheck = async treasureId =>
    apiPost(API_ENDPOINTS.treasureCheck, { treasureID: treasureId });

export const treasureUnlock = async (treasureId, chosenCost) =>
    apiPost(API_ENDPOINTS.treasureUnlock, {
        treasureID: treasureId,
        chosenCost: chosenCost,
    });

export const makeTradeOfferCheck = async () =>
    apiPost(API_ENDPOINTS.makeOfferCheck);

export const makeTradeOffer = async (offered, requested) =>
    apiPost(API_ENDPOINTS.makeOffer, {
        offered: offered,
        requested: requested,
    });

export const acceptTradeOffer = async offerID =>
    apiPost(API_ENDPOINTS.acceptOffer, { offerID: offerID });

export const deleteTradeOffer = async offerID =>
    apiPost(API_ENDPOINTS.deleteOffer, { offerID: offerID });

export const getTradeOffers = async (offset = 0, limit = 5, by = null) =>
    apiGet(API_ENDPOINTS.getOffers(offset, limit, by));

export const investCheck = async () => apiPost(API_ENDPOINTS.investCheck);

export const invest = async (sessionID, coin) =>
    apiPost(API_ENDPOINTS.invest, { sessionID, coin });

export const getInboxMessages = async (offset = null, limit = 15) =>
    apiGet(API_ENDPOINTS.getInboxMessages(offset, limit));
