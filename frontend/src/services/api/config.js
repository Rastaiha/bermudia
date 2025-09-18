// frontend/src/services/apiConfig.js

import { BASE_URLS } from './base_url.js';

const API_BASE_URL = BASE_URLS.API;
const WS_BASE_URL = BASE_URLS.WS;

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
    requestHelp: id => `${API_BASE_URL}/answer/${id}/help`,

    // Market
    makeOfferCheck: `${API_BASE_URL}/trade/make_offer_check`,
    makeOffer: `${API_BASE_URL}/trade/make_offer`,
    acceptOffer: `${API_BASE_URL}/trade/accept_offer`,
    deleteOffer: `${API_BASE_URL}/trade/delete_offer`,
    getOffers: (offset, limit, by) =>
        by == null
            ? `${API_BASE_URL}/trade/offers?offset=${offset}&limit=${limit}`
            : `${API_BASE_URL}/trade/offers?offset=${offset}&limit=${limit}&by=${by}`,

    // Investment:
    investCheck: `${API_BASE_URL}/invest_check`,
    invest: `${API_BASE_URL}/invest`,

    // Inbox
    getInboxMessages: (offset, limit) =>
        offset == null
            ? `${API_BASE_URL}/inbox/messages?limit=${limit}`
            : `${API_BASE_URL}/inbox/messages?offset=${offset}&limit=${limit}`,

    // WebSocket Endpoints
    events: `${WS_BASE_URL}/events`,
    marketEvents: `${WS_BASE_URL}/trade/events`,
    inboxEvents: `${WS_BASE_URL}/inbox/events`,
};
