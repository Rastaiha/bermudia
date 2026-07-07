<template>
    <div
        class="w-full h-screen flex justify-center items-center p-4 box-border overflow-hidden relative"
    >
        <StarryNight v-if="starryBackground" />
        <div
            v-if="fixedBackground && backgroundImage"
            class="fixed inset-0 -z-10 bg-center bg-no-repeat"
            :style="{
                backgroundImage: `url(${backgroundImage})`,
                backgroundSize: 'contain',
            }"
        ></div>
        <div class="fixed inset-0 bg-[#0c2036] -z-20"></div>
        <div
            v-if="isLoading"
            class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
        >
            <LoadingBar :progress="loadingProgress" />
        </div>

        <template v-else-if="player">
            <MapView
                ref="mapViewComponentRef"
                :islands="islands"
                :edges="edges"
                :player="player"
                :username="username"
                :dynamic-view-box="dynamicViewBox"
                :territory-id="territoryId"
                :background-image="backgroundImage"
                @node-click="showInfoBox"
                @map-transformed="updateInfoBoxPosition"
            />

            <Toolbar :player="player" />

            <UnreadMessageIndicator />

            <UserProfile :username="username" />

            <PlayerInfo :player="player" />

            <Transition name="popup-fade">
                <IslandInfoBox
                    v-if="selectedIsland"
                    ref="infoBoxRef"
                    :key="selectedIsland.id"
                    :selected-island="selectedIsland"
                    :player="player"
                    :info-box-style="infoBoxStyle"
                    :refuel-islands="refuelIslands"
                    :terminal-islands="terminalIslands"
                    :territory-id="territoryId"
                />
            </Transition>

            <div
                class="fixed top-4 left-4 text-xl font-bold text-white drop-shadow-[0_2px_4px_rgba(0,0,0,0.7)] pointer-events-none md:top-6 md:left-1/2 md:-translate-x-1/2 md:text-3xl"
            >
                {{ territoryName }}
            </div>
        </template>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { usePlayerWebSocket } from '@/services/websocket.js';
import { useInboxWebSocket } from '@/services/inboxWebsocket.js';
import { useTerritoryData } from '@/composables/useTerritoryData.js';
import { useInfoBoxPosition } from '@/composables/useInfoBoxPosition.js';
import eventBus from '@/services/eventBus.js';
import { APP_CONFIG, BACKGROUND_MODES } from '@/config/appConfig.js';

import MapView from '@/components/features/map/MapView.vue';
import IslandInfoBox from '@/components/features/map/IslandInfoBox.vue';
import PlayerInfo from '@/components/layout/PlayerInfo.vue';
import LoadingBar from '@/components/common/LoadingBar.vue';
import Toolbar from '@/components/layout/Toolbar.vue';
import StarryNight from '@/components/common/StarryNight.vue';
import UserProfile from '@/components/layout/UserProfile.vue';
import UnreadMessageIndicator from '@/components/features/notification/UnreadMessageIndicator.vue';

const backgroundMode = APP_CONFIG.BACKGROUND_MODE;
const fixedBackground = backgroundMode === BACKGROUND_MODES.FIXED;
const starryBackground = backgroundMode === BACKGROUND_MODES.STARRY;

const route = useRoute();
const router = useRouter();
const mapViewComponentRef = ref(null);
const infoBoxRef = ref(null);
const selectedIsland = ref(null);

const {
    territoryId,
    islands,
    refuelIslands,
    terminalIslands,
    edges,
    territoryName,
    player,
    username,
    backgroundImage,
    dynamicViewBox,
    isLoading,
    loadingProgress,
    loadPageData,
} = useTerritoryData(router, {
    onBeforeLoad: () => hideInfoBox(),
    onLoaded: () => nextTick(() => mapViewComponentRef.value?.zoomToPlayer()),
});

const { infoBoxStyle, updateInfoBoxPosition } = useInfoBoxPosition(
    mapViewComponentRef,
    selectedIsland
);

const handleClickOutside = event => {
    if (!selectedIsland.value || infoBoxRef.value?.$el.contains(event.target)) {
        return;
    }
    hideInfoBox();
};

const showInfoBox = island => {
    if (selectedIsland.value && selectedIsland.value.id === island.id) {
        hideInfoBox();
    } else {
        selectedIsland.value = island;
    }
};

const hideInfoBox = () => {
    selectedIsland.value = null;
};

onMounted(() => {
    eventBus.emit('set-audio-state', 'play');
    loadPageData(route.params.id);
    document.addEventListener('pointerdown', handleClickOutside);
});

onUnmounted(() => {
    eventBus.emit('set-audio-state', 'stop');
    document.removeEventListener('pointerdown', handleClickOutside);
});

watch(
    () => route.params.id,
    newId => {
        if (newId && newId !== territoryId.value) {
            loadPageData(newId);
        }
    }
);

watch(
    () => player.value?.atIsland,
    (newIsland, oldIsland) => {
        if (newIsland && oldIsland && newIsland !== oldIsland) {
            hideInfoBox();
        }
    },
    { deep: true }
);

usePlayerWebSocket(player, territoryId, route, router);
useInboxWebSocket();
</script>

<style scoped>
.popup-fade-enter-active,
.popup-fade-leave-active {
    transition:
        opacity 0.3s ease,
        transform 0.3s ease;
}

.popup-fade-enter-from,
.popup-fade-leave-to {
    opacity: 0;
    transform: translateY(10px);
}
</style>
