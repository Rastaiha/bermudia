import { logger } from '@/services/logger.js';

const MAX_RECONNECT_ATTEMPTS = 10;
const BASE_RECONNECT_DELAY = 1000;
const MAX_RECONNECT_DELAY = 30000;

export function createReconnectingSocket({
    buildUrl,
    onOpen,
    onMessage,
    onError,
    label = 'WebSocket',
}) {
    let socket = null;
    let reconnectTimeoutId = null;
    let reconnectAttempts = 0;

    const disconnect = () => {
        if (reconnectTimeoutId) clearTimeout(reconnectTimeoutId);

        if (socket) {
            socket.onclose = null;
            socket.onmessage = null;
            socket.onerror = null;
            socket.onopen = null;
            socket.close(1000, 'Connection closed intentionally by client');
            socket = null;
            logger.log(`${label} connection cleanly disconnected.`);
        }
    };

    const scheduleReconnect = () => {
        if (reconnectAttempts >= MAX_RECONNECT_ATTEMPTS) {
            logger.error(
                `Max ${label} reconnection attempts reached. Giving up.`
            );
            return;
        }
        const delay = Math.min(
            BASE_RECONNECT_DELAY * Math.pow(2, reconnectAttempts),
            MAX_RECONNECT_DELAY
        );
        reconnectAttempts++;

        logger.log(
            `Scheduling ${label} reconnection attempt ${reconnectAttempts} in ${delay}ms`
        );
        reconnectTimeoutId = setTimeout(connect, delay);
    };

    const connect = () => {
        if (socket) return;

        const url = buildUrl();
        if (!url) {
            logger.error(`${label}: No auth token found, connection aborted.`);
            return;
        }

        logger.log(`Attempting to connect ${label}...`);
        socket = new WebSocket(url);

        socket.onopen = () => {
            logger.log(`${label} connection established.`);
            reconnectAttempts = 0;
            onOpen?.();
        };

        socket.onmessage = event => {
            try {
                onMessage(JSON.parse(event.data));
            } catch (error) {
                logger.error(`Error parsing ${label} message:`, error);
            }
        };

        socket.onclose = event => {
            logger.log(`${label} connection closed. Code:`, event.code);
            socket = null;

            if (event.code !== 1000) {
                scheduleReconnect();
            }
        };

        socket.onerror = error => {
            logger.error(`${label} error:`, error);
            onError?.(error);
        };
    };

    return { connect, disconnect };
}
