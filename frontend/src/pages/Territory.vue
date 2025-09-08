<template>
    <div
        class="w-full h-screen flex justify-center items-center p-4 box-border bg-cover bg-center bg-no-repeat overflow-hidden bg-[#0c2036]"
        :style="{ backgroundImage: `url(${backgroundImage})` }"
    >
        <LoadingIndicator v-if="isLoading" :message="loadingMessage" />

        <template v-else-if="player">
            <div
                class="fixed md:top-8 top-0 text-2xl font-bold text-[#f5deb3] drop-shadow-[2px_4px_6px_white]"
            >
                {{ territoryName }}
            </div>
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

            <PlayerInfo :player="player" :username="username" />

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
        </template>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getPlayer, getMe, getToken, getTerritory } from '@/services/api.js';
import { usePlayerWebSocket } from '@/components/service/WebSocket.js';

import MapView from '@/components/MapView.vue';
import IslandInfoBox from '@/components/IslandInfoBox.vue';
import PlayerInfo from '@/components/PlayerInfo.vue';
import LoadingIndicator from '@/components/LoadingIndicator.vue';

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
const loadingMessage = ref('در حال بارگذاری نقشه...');
const isLoading = ref(true);
const infoBoxStyle = ref({ display: 'none' });

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

const loadPageData = async id => {
    if (!id) return;
    isLoading.value = true;
    hideInfoBox();

    try {
        if (!getToken()) {
            router.push({ name: 'Login' });
            return;
        }

        loadingMessage.value = 'در حال دریافت اطلاعات کاربری...';
        const [playerData, meData] = await Promise.all([getPlayer(), getMe()]);

        if (playerData.atTerritory.toString() !== id.toString()) {
            router.push({
                name: 'Territory',
                params: { id: playerData.atTerritory },
            });
            return;
        }

        loadingMessage.value = 'در حال بارگذاری نقشه...';
        const territoryData = await getTerritory(id);

        player.value = playerData;
        username.value = meData.username;
        territoryId.value = id;

        backgroundImage.value = `/images/${territoryData.backgroundAsset}`;
        territoryName.value = territoryData.name;
        islands.value = territoryData.islands;
        edges.value = territoryData.edges;
        refuelIslands.value = territoryData.refuelIslands;
        terminalIslands.value = territoryData.terminalIslands;
        dynamicViewBox.value = calculateViewBox(territoryData.islands);
    } catch (error) {
        console.error('Failed to load territory data:', error);
        router.push({ name: 'Login' });
    } finally {
        isLoading.value = false;
    }
};

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
    return `${minX - padding} ${minY - padding} ${maxX - minX + padding * 2} ${maxY - minY + padding * 2}`;
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
    loadPageData(route.params.id);
    document.addEventListener('pointerdown', handleClickOutside);
});

onUnmounted(() => {
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
