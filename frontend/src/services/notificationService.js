import { reactive, computed } from 'vue';

const RECEIVED_MESSAGES_KEY = 'received_messages';
const SEEN_MESSAGES_KEY = 'seen_messages';

const state = reactive({
    receivedMessages: [],
    seenMessages: [],
});

const getFromLocalStorage = key => {
    try {
        const data = localStorage.getItem(key);
        return data ? JSON.parse(data) : [];
    } catch (e) {
        console.error(`Error reading ${key} from LS:`, e);
        return [];
    }
};

const saveToLocalStorage = (key, data) => {
    try {
        localStorage.setItem(key, JSON.stringify(data));
    } catch (e) {
        console.error('Error saving to LS', e);
    }
};

function init() {
    state.receivedMessages = getFromLocalStorage(RECEIVED_MESSAGES_KEY);
    state.seenMessages = getFromLocalStorage(SEEN_MESSAGES_KEY);
}

function setReceivedMessages(messages) {
    const messageInfos = messages
        .filter(msg => msg && msg.createdAt && msg.content)
        .map(msg => ({
            createdAt: msg.createdAt,
            reason: Object.keys(msg.content)[0],
        }));
    state.receivedMessages = messageInfos;
    saveToLocalStorage(RECEIVED_MESSAGES_KEY, state.receivedMessages);
}

function markAllAsSeen() {
    state.seenMessages = [...state.receivedMessages];
    saveToLocalStorage(SEEN_MESSAGES_KEY, state.seenMessages);
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
