<template>
    <div
        class="w-full h-screen flex justify-center items-center p-4 box-border overflow-hidden relative"
    >
        <StarryNight />
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
                :dynamic-view-box="dynamicViewBox"
                :territory-id="territoryId"
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
import {
    getPlayer,
    getMe,
    getToken,
    getTerritory,
} from '@/services/api/index.js';
import { usePlayerWebSocket } from '@/services/websocket.js';
import { useInboxWebSocket } from '@/services/inboxWebsocket.js';
import eventBus from '@/services/eventBus.js';

import MapView from '@/components/features/map/MapView.vue';
import IslandInfoBox from '@/components/features/map/IslandInfoBox.vue';
import PlayerInfo from '@/components/layout/PlayerInfo.vue';
import LoadingBar from '@/components/common/LoadingBar.vue';
import Toolbar from '@/components/layout/Toolbar.vue';
import StarryNight from '@/components/common/StarryNight.vue';
import UserProfile from '@/components/layout/UserProfile.vue';
import UnreadMessageIndicator from '@/components/features/notification/UnreadMessageIndicator.vue';

const route = useRoute();
const router = useRouter();
const mapViewComponentRef = ref(null);
const infoBoxRef = ref(null);
const transformCounter = ref(0);

const territoryId = ref(null);
const islands = ref([]);
const refuelIslands = ref([]);
const terminalIslands = ref([]);
const edges = ref([]);
const territoryName = ref('');
const player = ref(null);
const username = ref('...');
const backgroundImage = ref('');
const selectedIsland = ref(null);
const dynamicViewBox = ref('0 0 1 1');
const isLoading = ref(true);
const loadingProgress = ref(0);
const infoBoxStyle = ref({ display: 'none' });

const fetchTerritoryData = async id => {
    return getTerritory(id);
};

const fetchPlayerAndUserData = async () => {
    if (!getToken()) {
        router.push({ name: 'Login' });
        throw new Error('User not authenticated');
    }
    const [playerData, meData] = await Promise.all([getPlayer(), getMe()]);
    return { playerData, meData };
};

const setupTerritoryData = territoryData => {
    backgroundImage.value = territoryData.backgroundAsset;
    territoryName.value = territoryData.name;
    islands.value = territoryData.islands;
    edges.value = territoryData.edges;
    refuelIslands.value = territoryData.refuelIslands;
    terminalIslands.value = territoryData.terminalIslands;
    dynamicViewBox.value = calculateViewBox(territoryData.islands);
};

const setupPlayerAndUserData = (playerAndUserData, currentTerritoryId) => {
    if (!playerAndUserData) return;
    const { playerData, meData } = playerAndUserData;

    if (playerData.atTerritory.toString() !== currentTerritoryId.toString()) {
        router.push({
            name: 'Territory',
            params: { id: playerData.atTerritory },
        });
        throw new Error('Redirecting to correct territory');
    }

    username.value = meData.name;
    player.value = playerData;
};

const loadPageData = async id => {
    if (!id) return;
    isLoading.value = true;
    loadingProgress.value = 0;
    territoryId.value = id;
    hideInfoBox();

    let progressInterval = null;

    try {
        progressInterval = setInterval(() => {
            if (loadingProgress.value < 90) {
                loadingProgress.value += 5;
            }
        }, 100);

        const [territoryData, playerAndUserData] = await Promise.all([
            fetchTerritoryData(id),
            fetchPlayerAndUserData(),
        ]);

        loadingProgress.value = 100;

        setupTerritoryData(territoryData);
        setupPlayerAndUserData(playerAndUserData, id);
    } catch (error) {
        console.error('Failed to load page data:', error.message);
        if (
            !error.message.includes('authenticated') &&
            !error.message.includes('Redirecting')
        ) {
            router.push({ name: 'Login' });
        }
    } finally {
        clearInterval(progressInterval);
        setTimeout(() => {
            isLoading.value = false;
            nextTick(() => {
                mapViewComponentRef.value?.zoomToPlayer();
            });
        }, 500);
    }
};

const handleClickOutside = event => {
    if (!selectedIsland.value || infoBoxRef.value?.$el.contains(event.target)) {
        return;
    }
    hideInfoBox();
};

const calculateInfoBoxStyle = () => {
    const svgElement = mapViewComponentRef.value?.svgRef;
    if (!selectedIsland.value || !svgElement) {
        infoBoxStyle.value = { display: 'none' };
        return;
    }

    const island = selectedIsland.value;
    const pt = svgElement.createSVGPoint();
    pt.x = island.x;
    pt.y = island.y;
    const screenPoint = pt.matrixTransform(svgElement.getScreenCTM());

    infoBoxStyle.value = {
        position: 'fixed',
        top: `${screenPoint.y}px`,
        left: `${screenPoint.x}px`,
        transform: 'translate(-50%, -100%) translateY(-20px)',
    };
};

watch([selectedIsland, transformCounter], calculateInfoBoxStyle, {
    flush: 'post',
});

const updateInfoBoxPosition = () => {
    transformCounter.value++;
};

const calculateViewBox = (islands, padding = 0.1) => {
    if (!islands || islands.length === 0) return '0 0 1 1';
    const bounds = islands.reduce(
        (acc, island) => ({
            minX: Math.min(acc.minX, island.x - island.width / 2),
            maxX: Math.max(acc.maxX, island.x + island.width / 2),
            minY: Math.min(acc.minY, island.y - island.height / 2),
            maxY: Math.max(acc.maxY, island.y + island.height / 2),
        }),
        { minX: Infinity, maxX: -Infinity, minY: Infinity, maxY: -Infinity }
    );
    const { minX, minY, maxX, maxY } = bounds;
    return `${minX - padding} ${minY - padding} ${maxX - minX + padding * 2} ${
        maxY - minY + padding * 2
    }`;
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

usePlayerWebSocket(player, territoryId, router);
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
