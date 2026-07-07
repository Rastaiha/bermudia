import { reactive, computed } from 'vue';
import { getJSON, setJSON } from '@/services/storage.js';

const RECEIVED_MESSAGES_KEY = 'received_messages';
const SEEN_MESSAGES_KEY = 'seen_messages';

const state = reactive({
    receivedMessages: [],
    seenMessages: [],
});

function init() {
    state.receivedMessages = getJSON(RECEIVED_MESSAGES_KEY);
    state.seenMessages = getJSON(SEEN_MESSAGES_KEY);
}

function setReceivedMessages(messages) {
    const messageInfos = messages
        .filter(msg => msg && msg.createdAt && msg.content)
        .map(msg => ({
            createdAt: msg.createdAt,
            reason: Object.keys(msg.content)[0],
        }));
    state.receivedMessages = messageInfos;
    setJSON(RECEIVED_MESSAGES_KEY, state.receivedMessages);
}

function markAllAsSeen() {
    state.seenMessages = [...state.receivedMessages];
    setJSON(SEEN_MESSAGES_KEY, state.seenMessages);
}

const hasUnreadMessages = computed(() => {
    const receivedIds = new Set(
        state.receivedMessages.map(msg => msg.createdAt)
    );
    const seenIds = new Set(state.seenMessages.map(msg => msg.createdAt));

    for (const id of receivedIds) {
        if (!seenIds.has(id)) {
            return true;
        }
    }
    return false;
});

init();

export const notificationService = {
    setReceivedMessages,
    markAllAsSeen,
    hasUnreadMessages,
};
