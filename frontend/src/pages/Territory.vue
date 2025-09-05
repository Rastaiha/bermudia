<template>
    <div
        class="w-full h-screen flex justify-center items-center p-4 box-border bg-cover bg-center bg-no-repeat overflow-hidden bg-[#0c2036]"
        :style="{ backgroundImage: `url(${backgroundImage})` }"
        @pointerdown="hideInfoBox"
    >
        <LoadingIndicator v-if="isLoading" :message="loadingMessage" />

        <template v-else>
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

            <PlayerInfo v-if="player" :player="player" :username="username" />

            <Transition name="popup-fade">
                <IslandInfoBox
                    v-if="selectedIsland"
                    :key="selectedIsland"
                    :selected-island="selectedIsland"
                    :player="player"
                    :info-box-style="infoBoxStyle"
                    :refuel-islands="refuelIslands"
                    :terminal-islands="terminalIslands"
                    :territory-id="territoryId"
                />
            </Transition>
        </template>
    </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import {
    getPlayer,
    getMe,
    getToken,
    logout,
    getTerritory,
} from '@/services/api.js';
import { usePlayerWebSocket } from '@/components/service/WebSocket.js';

import MapView from '@/components/MapView.vue';
import IslandInfoBox from '@/components/IslandInfoBox.vue';
import PlayerInfo from '@/components/PlayerInfo.vue';
import LoadingIndicator from '@/components/LoadingIndicator.vue';

// --- State ---
const route = useRoute();
const router = useRouter();
const mapViewComponentRef = ref(null);
const transformCounter = ref(0);

const territoryId = ref(route.params.id);
const islands = ref([]);
const refuelIslands = ref([]);
const terminalIslands = ref([]);
const edges = ref([]);
const player = ref(null);
const username = ref('...');
const backgroundImage = ref('');
const selectedIsland = ref(null);
const dynamicViewBox = ref('0 0 1 1');
const loadingMessage = ref('Loading map data...');
const isLoading = ref(true);

// --- Computed Properties ---
const infoBoxStyle = computed(() => {
    const svgElement = mapViewComponentRef.value?.svgRef;
    if (!selectedIsland.value || !svgElement) return { display: 'none' };
    const island = selectedIsland.value;
    const pt = svgElement.createSVGPoint();
    pt.x = island.x;
    pt.y = island.y;
    const screenPoint = pt.matrixTransform(svgElement.getScreenCTM());
    return {
        position: 'fixed',
        top: `${screenPoint.y}px`,
        left: `${screenPoint.x}px`,
        transform: 'translate(-50%, -100%) translateY(-20px)',
    };
});

// --- Methods ---
const updateInfoBoxPosition = () => {
    transformCounter.value++;
};

const calculateViewBox = (islands, padding = 0.1) => {
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
    return `${minX - padding} ${minY - padding} ${maxX - minX + padding * 2} ${maxY - minY + padding * 2}`;
};

// --- API Calls & Data Fetching ---
const fetchTerritoryData = async id => {
    return await getTerritory(id);
};

const fetchPlayerAndUserData = async () => {
    if (!getToken()) {
        logout();
        router.push({ name: 'Login' });
        return null;
    }
    try {
        const [playerData, meData] = await Promise.all([getPlayer(), getMe()]);
        return { playerData, meData };
    } catch (err) {
        console.error('Failed to get player/user data:', err);
        throw err;
    }
};

// --- Setup Functions (Processing and state setting) ---
const setupTerritoryData = territoryData => {
    backgroundImage.value = `/images/${territoryData.backgroundAsset}`;
    dynamicViewBox.value = calculateViewBox(territoryData.islands);
    islands.value = territoryData.islands;
    edges.value = territoryData.edges;
    refuelIslands.value = territoryData.refuelIslands;
    terminalIslands.value = territoryData.terminalIslands;
};

const setupPlayerAndUserData = playerAndUserData => {
    if (!playerAndUserData) return;

    const { playerData, meData } = playerAndUserData;
    username.value = meData.username;

    player.value = playerData;
    if (player.value.atTerritory != territoryId.value) {
        router.push({
            name: 'Territory',
            params: { id: player.value.atTerritory },
        });
    }
};

// --- Event Handlers from Child Components ---
const showInfoBox = async island => {
    if (!player.value) return;
    if (selectedIsland.value && selectedIsland.value.id === island.id) {
        hideInfoBox();
        return;
    }
    selectedIsland.value = island;
};

const hideInfoBox = () => {
    selectedIsland.value = null;
};

// --- Lifecycle Hooks ---
onMounted(async () => {
    isLoading.value = true;
    try {
        loadingMessage.value = 'Fetching data...';
        const [territoryData, playerAndUserData] = await Promise.all([
            fetchTerritoryData(territoryId.value),
            fetchPlayerAndUserData(),
        ]);

        loadingMessage.value = 'Setting up data...';

        setupTerritoryData(territoryData);
        setupPlayerAndUserData(playerAndUserData);
    } finally {
        isLoading.value = false;
    }
});
// --- WebSocket ---
usePlayerWebSocket(player, territoryId, router);

// --- Watcher to update infobox on arrival ---
watch(
    () => player.value?.atIsland,
    (newIsland, oldIsland) => {
        if (newIsland && oldIsland && newIsland !== oldIsland) {
            hideInfoBox();
        }
    },
    { deep: true }
);

watch(
    () => route.params.id,
    async newTerritoryId => {
        territoryId.value = newTerritoryId;
        isLoading.value = true;
        try {
            const territoryData = await fetchTerritoryData(newTerritoryId);
            setupTerritoryData(territoryData);

            const playerAndUserData = await fetchPlayerAndUserData();
            setupPlayerAndUserData(playerAndUserData);
        } finally {
            isLoading.value = false;
        }
    },
    { immediate: false }
);
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
