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

// Stable key for messages that lack a server-side id and may share a createdAt.
function keyOf(msg, index) {
    const reason = Object.keys(msg.content)[0];
    return `${msg.createdAt}|${reason}|${JSON.stringify(msg.content)}|${index}`;
}

function setReceivedMessages(messages) {
    const messageInfos = messages
        .filter(msg => msg && msg.createdAt && msg.content)
        .map((msg, index) => ({
            createdAt: msg.createdAt,
            reason: Object.keys(msg.content)[0],
            key: keyOf(msg, index),
        }));
    state.receivedMessages = messageInfos;
    setJSON(RECEIVED_MESSAGES_KEY, state.receivedMessages);
}

function markAllAsSeen() {
    state.seenMessages = [...state.receivedMessages];
    setJSON(SEEN_MESSAGES_KEY, state.seenMessages);
}

const hasUnreadMessages = computed(() => {
    const seenKeys = new Set(state.seenMessages.map(msg => msg.key));
    return state.receivedMessages.some(msg => !seenKeys.has(msg.key));
});

init();

export const notificationService = {
    setReceivedMessages,
    markAllAsSeen,
    hasUnreadMessages,
};
