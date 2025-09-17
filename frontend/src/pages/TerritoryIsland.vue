<template>
    <div
        class="w-full min-h-screen p-8 box-border flex justify-center items-center bg-cover bg-fixed bg-center font-main text-gray-200"
        :style="{ backgroundImage: `url(${backgroundImage})` }"
        @mousemove="updateMousePosition"
    >
        <FloatingUI
            :tooltip-text="tooltipText"
            :mouse-position="mousePosition"
        />
        <BackButton
            v-if="isLoaded"
            :territory-id="id"
            @show-tooltip="showTooltip"
            @hide-tooltip="hideTooltip"
        />

        <div v-if="isLoaded" class="w-full max-w-4xl flex flex-col gap-10">
            <template
                v-for="componentData in components"
                :key="componentData.id"
            >
                <Iframe
                    v-if="componentData.iframe"
                    :url="componentData.iframe.url"
                    @show-tooltip="showTooltip"
                    @hide-tooltip="hideTooltip"
                />
                <ChallengeBox
                    v-else-if="componentData.input"
                    :challenge="componentData.input"
                    @submit="handleChallengeSubmit"
                    @help-requested="handleHelpRequested"
                />
            </template>
            <template
                v-for="(treasureData, index) in treasures"
                :key="treasureData.id"
            >
                <Treasure
                    v-model="treasures[index]"
                    :treasure-data="treasureData"
                    :player="player"
                    @treasureOpened="openRewardModal"
                />
            </template>
        </div>

        <div
            v-else
            class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-gray-200 text-xl"
        >
            {{ loadingMessage }}
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, toRef } from 'vue';
import { useToast } from 'vue-toastification';
import { useRouter } from 'vue-router';
import {
    getIsland,
    submitAnswer,
    getPlayer,
    requestHelp,
} from '@/services/api';
import { usePlayerWebSocket } from '@/services/websocket.js';
import { useInboxWebSocket } from '@/services/inboxWebsocket.js';
import { useModal } from 'vue-final-modal';
import eventBus from '@/services/eventBus.js';

import Iframe from '@/components/common/Iframe.vue';
import ChallengeBox from '@/components/features/island/ChallengeBox.vue';
import BackButton from '@/components/features/island/BackButton.vue';
import Treasure from '@/components/features/island/Treasure.vue';
import FloatingUI from '@/components/common/FloatingUI.vue';
import TreasureRewardModal from '@/components/features/island/TreasureRewardModal.vue';

const props = defineProps({
    id: { type: String, required: true },
    islandId: { type: String, required: true },
});

const isLoaded = ref(false);
const backgroundImage = ref('');
const tooltipText = ref('');
const mousePosition = ref({ x: 0, y: 0 });
const loadingMessage = ref('درحال بارگذاری اطلاعات سیاره...');
const components = ref([]);
const treasures = ref([]);
const player = ref(null);
const toast = useToast();
const router = useRouter();
const territoryId = toRef(props, 'id');

usePlayerWebSocket(player, territoryId, router);
useInboxWebSocket();

const fetchIslandData = async id => {
    try {
        const [islandData, playerData] = await Promise.all([
            getIsland(id),
            getPlayer(),
        ]);

        player.value = playerData;
        backgroundImage.value = `/images/backgrounds/island/background.jpg`;
        components.value = islandData.components;
        treasures.value = islandData.treasures;
        isLoaded.value = true;
    } catch (error) {
        console.error('Failed to load island data from API:', error);
        loadingMessage.value = `خطا در بارگذاری سیاره: ${error.message}`;
        toast.error(error.message || 'خطا در دریافت اطلاعات سیاره');
    }
};

const openRewardModal = rewards => {
    const { open, close } = useModal({
        component: TreasureRewardModal,
        attrs: {
            rewards: rewards,
            onClose: () => close(),
        },
    });
    setTimeout(() => open(), 0);
};

const handleChallengeSubmit = async ({ inputId, data }) => {
    if (data === '' || data === null || data === undefined) {
        toast.error(`چیزی برای ارسال وارد نشده‌است!`);
        return;
    }
    const formData = new FormData();
    formData.append('data', data);
    try {
        await submitAnswer(inputId, formData);
        toast.success(`پاسخ شما با موفقیت ثبت شد.`);
        await fetchIslandData(props.islandId);
    } catch (error) {
        console.error('Error submitting answer:', error);
        toast.error(error.message || 'خطا در ارسال پاسخ!');
    }
};

const handleHelpRequested = async challenge => {
    try {
        const response = await requestHelp(challenge.id);
        if (challenge.submissionState.hasRequestedHelp) {
            window.open(response.meetLink, '_blank');
        } else {
            await fetchIslandData(props.islandId);
        }
    } catch (err) {
        toast.error(err.message);
    }
};

const updateMousePosition = event => {
    mousePosition.value = { x: event.clientX, y: event.clientY };
};

const showTooltip = text => {
    tooltipText.value = text;
};
const hideTooltip = () => {
    tooltipText.value = '';
};

onMounted(() => {
    eventBus.emit('set-audio-state', 'stop');
    fetchIslandData(props.islandId);
});
</script>

<style scoped>
.w-full {
    width: 100%;
}
</style>
