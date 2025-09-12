import { onMounted, onUnmounted } from 'vue';
import { getToken } from '@/services/api/index.js';
import { API_ENDPOINTS } from '@/services/api/config.js';

export function useInboxWebSocket(syncOffset, messages) {
    let socket = null;
    let reconnectTimeoutId = null;
    let reconnectAttempts = 0;
    const maxReconnectAttempts = 10;
    const baseReconnectDelay = 1000;

    const disconnect = () => {
        if (reconnectTimeoutId) clearTimeout(reconnectTimeoutId);
        if (socket) {
            socket.onclose = null;
            socket.onmessage = null;
            socket.onerror = null;
            socket.onopen = null;
            socket.close(1000, 'Connection closed intentionally by client');
            socket = null;
        }
    };

    const connect = () => {
        if (socket) return;

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
            try {
                const data = JSON.parse(event.data);

                if (data.ok === false) {
                    console.error(
                        'Received error in inbox websocket event:',
                        data.error
                    );
                    return;
                }

                const inboxEvent = data.result;

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

        socket.onclose = event => {
            socket = null;
            if (event.code !== 1000) {
                scheduleReconnect();
            }
        };

        socket.onerror = error => {
            console.error('Inbox WebSocket error:', error);
        };
    };

    const scheduleReconnect = () => {
        if (reconnectAttempts >= maxReconnectAttempts) {
            return;
        }
        const delay = Math.min(
            baseReconnectDelay * Math.pow(2, reconnectAttempts),
            30000
        );
        reconnectAttempts++;
        reconnectTimeoutId = setTimeout(connect, delay);
    };

    onMounted(() => {
        if (!socket) {
            connect();
        }
    });

    onUnmounted(() => {
        disconnect();
    });
}
