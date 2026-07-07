import { ref, onMounted, onUnmounted } from 'vue';
import { getToken, getInboxMessages } from '@/services/api/index.js';
import { API_ENDPOINTS } from '@/services/api/config.js';
import { notificationService } from '@/services/notificationService.js';
import { uiState } from '@/services/uiState.js';
import { createReconnectingSocket } from '@/services/ws/reconnectingSocket.js';

export const messages = ref([]);
export const syncOffset = ref(null);

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

function handleMessage(data) {
    if (data.ok === false) {
        console.error('[InboxWS] Received error event:', data.error);
        return;
    }
    if (data.sync) {
        syncOffset.value = data.sync.offset;
        return;
    }
    const newMessage = data.newMessage || data;
    if (newMessage.content && newMessage.createdAt) {
        messages.value.unshift(newMessage);
        triggerNotificationCheck();
    }
}

const { connect, disconnect } = createReconnectingSocket({
    buildUrl: () => {
        const token = getToken();
        return token ? `${API_ENDPOINTS.inboxEvents}?token=${token}` : null;
    },
    onMessage: handleMessage,
    label: 'InboxWS',
});

export function useInboxWebSocket() {
    onMounted(() => {
        connect();
    });

    onUnmounted(() => {
        disconnect();
    });
}
