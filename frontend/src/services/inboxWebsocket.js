import { onMounted, onUnmounted } from 'vue';
import { getToken } from '@/services/api/index.js';
import { API_ENDPOINTS } from '@/services/api/config.js';

let socket = null;
let reconnectTimeoutId = null;
let reconnectAttempts = 0;
const maxReconnectAttempts = 10;
const baseReconnectDelay = 1000;

const messageListeners = new Set();

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

function disconnect() {
    if (reconnectTimeoutId) clearTimeout(reconnectTimeoutId);
    if (socket) {
        socket.onclose = null;
        socket.onmessage = null;
        socket.onerror = null;
        socket.onopen = null;
        socket.close(1000, 'Connection closed: no active listeners.');
        socket = null;
        messageListeners.clear();
    }
}

function connect() {
    if (socket && socket.readyState < 2) {
        return;
    }

    const token = getToken();
    if (!token) {
        console.error(
            'Inbox WebSocket: No auth token found, connection aborted.'
        );
        return;
    }

    socket = new WebSocket(`${API_ENDPOINTS.inboxEvents}?token=${token}`);

    socket.onopen = () => {
        reconnectAttempts = 0;
    };

    socket.onmessage = event => {
        messageListeners.forEach(handler => handler(event));
    };

    socket.onclose = event => {
        socket = null;
        if (event.code !== 1000) {
            scheduleReconnect();
        }
    };

    socket.onerror = error => {
        console.error('Inbox WebSocket error:', error);
        if (socket) {
            socket.close();
        }
    };
}

export function useInboxWebSocket(syncOffset, messages) {
    const handleMessage = event => {
        try {
            const data = JSON.parse(event.data);

            if (data.ok === false) {
                console.error(
                    'Received error in inbox websocket event:',
                    data.error
                );
                return;
            }

            const inboxEvent = data;

            if (inboxEvent.sync) {
                syncOffset.value = inboxEvent.sync.offset;
            }
            if (inboxEvent.newMessage) {
                messages.value.unshift(inboxEvent.newMessage);
                syncOffset.value = inboxEvent.newMessage.createdAt;
            }
        } catch (error) {
            console.error('Error parsing Inbox WebSocket message:', error);
        }
    };

    onMounted(() => {
        connect();
        messageListeners.add(handleMessage);
    });

    onUnmounted(() => {
        messageListeners.delete(handleMessage);
        if (messageListeners.size === 0) {
            disconnect();
        }
    });
}
