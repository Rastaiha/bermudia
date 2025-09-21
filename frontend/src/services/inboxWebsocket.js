import { ref, onMounted, onUnmounted } from 'vue';
import { getToken, getInboxMessages } from '@/services/api/index.js';
import { API_ENDPOINTS } from '@/services/api/config.js';
import { notificationService } from '@/services/notificationService.js';
import { uiState } from '@/services/uiState.js';

export const messages = ref([]);
export const syncOffset = ref(null);

let socket = null;
let reconnectTimeoutId = null;
let reconnectAttempts = 0;
const maxReconnectAttempts = 10;
const baseReconnectDelay = 1000;

async function triggerNotificationCheck() {
    try {
        const result = await getInboxMessages(null, 20);
        const allMessages = result?.messages || result || [];
        notificationService.setReceivedMessages(allMessages);
        if (uiState.isInboxOpen) {
            notificationService.markAllAsSeen();
        }
    } catch (apiError) {
        console.error(
            '[Inbox WS] Failed to fetch inbox messages after WS event:',
            apiError
        );
    }
}

function handleMessage(event) {
    try {
        const data = JSON.parse(event.data);
        if (data.ok === false) {
            console.error('[InboxWS] Received error event:', data.error);
            return;
        }
        const eventData = data;
        if (eventData.sync) {
            syncOffset.value = eventData.sync.offset;
            return;
        }
        const newMessage = eventData.newMessage || eventData;
        if (newMessage.content && newMessage.createdAt) {
            messages.value.unshift(newMessage);
            triggerNotificationCheck();
        }
    } catch (error) {
        console.error('[InboxWS] Error parsing message:', error);
    }
}

export function useInboxWebSocket() {
    const scheduleReconnect = () => {
        if (reconnectAttempts >= maxReconnectAttempts) {
            console.error(
                '[InboxWS] Max reconnection attempts reached. Giving up.'
            );
            return;
        }
        const delay = Math.min(
            baseReconnectDelay * Math.pow(2, reconnectAttempts),
            30000
        );
        reconnectAttempts++;
        reconnectTimeoutId = setTimeout(connect, delay);
    };

    const disconnect = () => {
        if (reconnectTimeoutId) clearTimeout(reconnectTimeoutId);
        if (socket) {
            socket.onclose = null;
            socket.close(1000, 'Connection closed by client.');
            socket = null;
        }
    };

    const connect = () => {
        if (socket && socket.readyState < 2) {
            return;
        }
        const token = getToken();
        if (!token) {
            return;
        }

        socket = new WebSocket(`${API_ENDPOINTS.inboxEvents}?token=${token}`);

        socket.onopen = () => {
            reconnectAttempts = 0;
        };
        socket.onmessage = handleMessage;
        socket.onclose = event => {
            socket = null;
            if (event.code !== 1000) {
                scheduleReconnect();
            }
        };
        socket.onerror = error => {
            console.error('[InboxWS] WebSocket error:', error);
        };
    };

    onMounted(() => {
        connect();
    });

    onUnmounted(() => {
        disconnect();
    });
}
