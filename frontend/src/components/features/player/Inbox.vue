<template>
    <VueFinalModal
        class="flex justify-center items-center"
        content-class="flex flex-col w-full md:w-[450px] max-h-[80vh] mx-4 
                       bg-gray-800 border-4 border-gray-700 
                       rounded-xl shadow-xl"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div class="flex-shrink-0 p-5 border-b-2 border-gray-700">
            <div class="flex items-center justify-between">
                <h1 class="text-xl font-semibold text-gray-200">صندوق پیام</h1>
                <button
                    class="p-1 rounded-full hover:bg-gray-700 transition-colors"
                    @click="handleClose"
                >
                    <XMarkIcon class="h-6 w-6 text-gray-300" />
                </button>
            </div>
            <p v-if="messages.length > 0" class="text-sm text-gray-400 mt-1">
                {{ messages.length }} پیام جدید
            </p>
        </div>

        <div
            ref="scrollContainer"
            class="scrollable-content flex-grow p-5 space-y-4 overflow-y-auto"
        >
            <div v-if="isLoading" class="text-center py-10">
                <p class="text-gray-400">در حال بارگذاری پیام‌ها...</p>
            </div>

            <div v-else-if="messages.length === 0" class="text-center py-10">
                <p class="text-gray-400">صندوق پیام شما خالی است.</p>
            </div>

            <NotificationItem
                v-for="message in messages"
                :key="message.createdAt"
                :message="message"
            />

            <div v-if="canLoadMore" class="text-center">
                <button
                    class="text-blue-400 hover:text-blue-300 text-sm"
                    @click="loadMoreMessages"
                >
                    بارگذاری بیشتر
                </button>
            </div>
        </div>
    </VueFinalModal>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { VueFinalModal } from 'vue-final-modal';
import { XMarkIcon } from '@heroicons/vue/24/outline';
import { getInboxMessages } from '@/services/api/index.js';
import eventBus from '@/services/eventBus.js';
import NotificationItem from './NotificationItem.vue';

const emit = defineEmits(['close']);

const messages = ref([]);
const syncOffset = ref(null);
const isLoading = ref(true);
const canLoadMore = ref(true);
const scrollContainer = ref(null);

const fetchMessages = async (offset = null) => {
    if (!canLoadMore.value && offset) return;
    isLoading.value = true;
    try {
        const limit = 15;
        const result = await getInboxMessages(offset, limit);
        if (result.length < limit) {
            canLoadMore.value = false;
        }
        const newMessages = offset ? [...messages.value, ...result] : result;
        messages.value = newMessages.sort((a, b) => b.createdAt - a.createdAt);
    } catch (error) {
        console.error('Failed to fetch inbox messages:', error);
    } finally {
        isLoading.value = false;
    }
};

const handleNewMessage = newMessage => {
    messages.value.unshift(newMessage);
};

const handleSync = syncEvent => {
    syncOffset.value = syncEvent.offset;
    fetchMessages(syncOffset.value);
};

const loadMoreMessages = () => {
    if (messages.value.length > 0) {
        const lastMessage = messages.value[messages.value.length - 1];
        fetchMessages(lastMessage.createdAt);
    }
};

function handleClose() {
    emit('close');
}

onMounted(() => {
    eventBus.on('inbox-sync', handleSync);
    eventBus.on('inbox-new-message', handleNewMessage);

    if (!syncOffset.value) {
        fetchMessages();
    }
});

onUnmounted(() => {
    eventBus.off('inbox-sync', handleSync);
    eventBus.off('inbox-new-message', handleNewMessage);
});
</script>

<style scoped>
.scrollable-content::-webkit-scrollbar {
    width: 8px;
}
.scrollable-content::-webkit-scrollbar-track {
    background: transparent;
}
.scrollable-content::-webkit-scrollbar-thumb {
    background-color: #4b5563;
    border-radius: 4px;
}
.scrollable-content::-webkit-scrollbar-thumb:hover {
    background-color: #6b7280;
}
</style>
