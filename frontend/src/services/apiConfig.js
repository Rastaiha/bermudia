// frontend/src/services/apiConfig.js

const API_BASE_URL =
    'http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1';
const WS_BASE_URL =
    'ws://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1';

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

    // WebSocket - Assuming it must also go through the /api/v1 proxy path
    events: `${WS_BASE_URL}/events`,
};
