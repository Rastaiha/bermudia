import eventBus from './eventBus';
import { getToken } from './api/index.js';
import { API_ENDPOINTS } from './api/config.js';

let socket = null;
let reconnectTimeoutId = null;
let reconnectAttempts = 0;
const maxReconnectAttempts = 10;
const baseReconnectDelay = 1000;

function disconnect() {
    if (reconnectTimeoutId) clearTimeout(reconnectTimeoutId);
    if (socket) {
        socket.onclose = null;
        socket.close(1000, 'Connection closed by client');
        socket = null;
    }
}

function scheduleReconnect() {
    if (reconnectAttempts >= maxReconnectAttempts) {
        console.error(
            'Max Inbox WebSocket reconnection attempts reached. Giving up.'
        );
        return;
    }
    const delay = Math.min(
        baseReconnectDelay * Math.pow(2, reconnectAttempts),
        30000
    );
    reconnectAttempts++;

    reconnectTimeoutId = setTimeout(() => connect(), delay);
}

function connect() {
    if (socket) {
        return;
    }

    const token = getToken();
    if (!token) {
        console.error('Inbox WebSocket: No auth token found.');
        return;
    }

    socket = new WebSocket(`${API_ENDPOINTS.inboxEvents}?token=${token}`);

    socket.onopen = () => {
        reconnectAttempts = 0;
    };

    socket.onmessage = event => {
        try {
            const data = JSON.parse(event.data);

            if (data.ok === false) {
                console.error(
                    'Received error in inbox websocket event:',
                    data.error
                );
                return;
            }

            const inboxEvent = data.result || data;

            if (inboxEvent.sync) {
                eventBus.emit('inbox-sync', inboxEvent.sync);
            } else if (inboxEvent.newMessage) {
                eventBus.emit('inbox-new-message', inboxEvent.newMessage);
            }
        } catch (error) {
            console.error('Error parsing Inbox WebSocket message:', error);
        }
    };

    socket.onclose = event => {
        socket = null;
        if (event.code !== 1000) {
            scheduleReconnect();
        }
    };

    socket.onerror = error => {
        console.error('Inbox WebSocket error:', error);
    };
}

export default {
    connect,
    disconnect,
};
