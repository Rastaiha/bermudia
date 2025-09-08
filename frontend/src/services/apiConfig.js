// frontend/src/services/apiConfig.js

const API_BASE_URL = 'https://bermudia-api-internal.darkube.app/api/v1';
const WS_BASE_URL = 'wss://bermudia-api-internal.darkube.app/api/v1';

export const API_ENDPOINTS = {
    // Authentication
    login: `${API_BASE_URL}/login`,
    getMe: `${API_BASE_URL}/me`,

    // Player & Travel
    getPlayer: `${API_BASE_URL}/player`,
    travelCheck: `${API_BASE_URL}/travel_check`,
    travelTo: `${API_BASE_URL}/travel`,
    anchorCheck: `${API_BASE_URL}/anchor_check`,
    dropAnchor: `${API_BASE_URL}/anchor`,
    refuelCheck: `${API_BASE_URL}/refuel_check`,
    buyFuel: `${API_BASE_URL}/refuel`,
    migrateCheck: `${API_BASE_URL}/migrate_check`,
    migrate: `${API_BASE_URL}/migrate`,

    // Territory & Island
    getTerritory: id => `${API_BASE_URL}/territories/${id}`,
    getIsland: id => `${API_BASE_URL}/islands/${id}`,
    submitAnswer: id => `${API_BASE_URL}/answer/${id}`,
    treasureCheck: `${API_BASE_URL}/unlock_treasure_check`,
    treasureUnlock: `${API_BASE_URL}/unlock_treasure`,

    // WebSocket - Assuming it must also go through the /api/v1 proxy path
    events: `${WS_BASE_URL}/events`,
};
