import { describe, it, expect } from 'vitest';
import { API_ENDPOINTS } from '@/services/api/config.js';

// Mirror the base-URL resolution in config.js so these tests exercise the
// URL-building logic regardless of whether a developer's .env overrides the
// base URLs. (Vitest loads .env the same way Vite does.)
const API =
    import.meta.env.VITE_API_BASE_URL ||
    'https://bermudia-api-internal.darkube.ir/api/v1';
const WS =
    import.meta.env.VITE_WS_BASE_URL ||
    'wss://bermudia-api-internal.darkube.ir/api/v1';

describe('API_ENDPOINTS', () => {
    it('builds static endpoints from the API base URL', () => {
        expect(API_ENDPOINTS.login).toBe(`${API}/login`);
        expect(API_ENDPOINTS.getPlayer).toBe(`${API}/player`);
        expect(API_ENDPOINTS.travelTo).toBe(`${API}/travel`);
    });

    it('builds WebSocket endpoints from the WS base URL', () => {
        expect(API_ENDPOINTS.events).toBe(`${WS}/events`);
        expect(API_ENDPOINTS.marketEvents).toBe(`${WS}/trade/events`);
        expect(API_ENDPOINTS.inboxEvents).toBe(`${WS}/inbox/events`);
    });

    it('interpolates ids into parameterized endpoints', () => {
        expect(API_ENDPOINTS.getTerritory('t1')).toBe(`${API}/territories/t1`);
        expect(API_ENDPOINTS.getIsland('is-9')).toBe(`${API}/islands/is-9`);
        expect(API_ENDPOINTS.submitAnswer('in-2')).toBe(`${API}/answer/in-2`);
        expect(API_ENDPOINTS.requestHelp('in-2')).toBe(
            `${API}/answer/in-2/help`
        );
    });

    describe('getOffers', () => {
        it('omits the by param when by is null or undefined', () => {
            expect(API_ENDPOINTS.getOffers(0, 20, null)).toBe(
                `${API}/trade/offers?offset=0&limit=20`
            );
            expect(API_ENDPOINTS.getOffers(0, 20)).toBe(
                `${API}/trade/offers?offset=0&limit=20`
            );
        });

        it('includes the by param when provided', () => {
            expect(API_ENDPOINTS.getOffers(10, 5, 'me')).toBe(
                `${API}/trade/offers?offset=10&limit=5&by=me`
            );
        });
    });

    describe('getInboxMessages', () => {
        it('omits the offset param when offset is null or undefined', () => {
            expect(API_ENDPOINTS.getInboxMessages(null, 20)).toBe(
                `${API}/inbox/messages?limit=20`
            );
            expect(API_ENDPOINTS.getInboxMessages(undefined, 20)).toBe(
                `${API}/inbox/messages?limit=20`
            );
        });

        it('includes the offset param when provided', () => {
            expect(API_ENDPOINTS.getInboxMessages(40, 20)).toBe(
                `${API}/inbox/messages?offset=40&limit=20`
            );
        });

        it('includes offset when it is 0 only if passed explicitly as non-null', () => {
            // offset === 0 is not null, so it takes the offset branch.
            expect(API_ENDPOINTS.getInboxMessages(0, 20)).toBe(
                `${API}/inbox/messages?offset=0&limit=20`
            );
        });
    });
});
