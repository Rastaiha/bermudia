<template>
    <VueFinalModal
        v-bind="$attrs"
        class="flex justify-center items-center"
        content-class="flex flex-col w-full md:w-[450px] max-h-[80vh] mx-4 
                       bg-gray-800 border-4 border-gray-700 
                       rounded-xl shadow-xl"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
        @update:model-value="val => (uiState.isInboxOpen = val)"
    >
        <div class="flex-shrink-0 p-5 border-b-2 border-gray-700">
            <div class="flex items-center justify-between">
                <h1 class="text-xl font-semibold text-gray-200">صندوق پیام</h1>
                <button
                    class="p-1 rounded-full hover:bg-gray-700 transition-colors"
                    @click="emit('close')"
                >
                    <XMarkIcon class="h-6 w-6 text-gray-300" />
                </button>
            </div>
            <p v-if="messages.length > 0" class="text-sm text-gray-400 mt-1">
                {{ messages.length }} پیام
            </p>
        </div>

        <div
            v-if="isLoading"
            class="flex-grow flex items-center justify-center"
        >
            <p class="text-gray-400">در حال بارگذاری پیام‌ها...</p>
        </div>

        <div
            v-else-if="messages.length === 0"
            class="flex-grow flex items-center justify-center"
        >
            <p class="text-gray-400">صندوق پیام شما خالی است.</p>
        </div>

        <div
            v-else
            ref="scrollContainer"
            class="scrollable-content flex-grow p-5 space-y-4 overflow-y-auto"
            @scroll="handleScroll"
        >
            <NotificationItem
                v-for="message in messages"
                :key="message.id"
                :message="message"
            />
            <div v-if="isFetchingMore" class="text-center py-2">
                <p class="text-gray-400 text-sm">در حال بارگذاری بیشتر...</p>
            </div>
        </div>
    </VueFinalModal>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { VueFinalModal } from 'vue-final-modal';
import { XMarkIcon } from '@heroicons/vue/24/outline';
import { getInboxMessages } from '@/services/api/index.js';
import { messages } from '@/services/inboxWebsocket.js';
import NotificationItem from './NotificationItem.vue';
import { notificationService } from '@/services/notificationService';
import { uiState } from '@/services/uiState.js';

const emit = defineEmits(['close']);

const isLoading = ref(true);
const isFetchingMore = ref(false);
const hasMore = ref(true);
const scrollContainer = ref(null);

const fetchMessages = async (offset = null) => {
    if (offset && !hasMore.value) return;

    if (offset) {
        isFetchingMore.value = true;
    } else {
        isLoading.value = true;
    }

    try {
        const result = await getInboxMessages(offset, 15);
        const newMessagesArray = result?.messages || result || [];

        if (newMessagesArray.length < 15) {
            hasMore.value = false;
        }

        if (offset) {
            messages.value.push(...newMessagesArray);
        } else {
            messages.value = newMessagesArray;
            notificationService.setReceivedMessages(messages.value);
        }
    } catch (error) {
        console.error('Failed to fetch inbox messages:', error);
    } finally {
        isLoading.value = false;
        isFetchingMore.value = false;
    }
};

const handleScroll = event => {
    if (isFetchingMore.value || !hasMore.value) return;

    const { scrollTop, scrollHeight, clientHeight } = event.target;
    if (scrollTop + clientHeight >= scrollHeight - 10) {
        const lastMessage = messages.value[messages.value.length - 1];
        if (lastMessage) {
            fetchMessages(lastMessage.createdAt);
        }
    }
};

onMounted(async () => {
    await fetchMessages();
    notificationService.markAllAsSeen();
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
