<template>
  <div
    class="w-full h-screen flex justify-center items-center p-4 box-border bg-cover bg-center bg-no-repeat overflow-hidden bg-[#0c2036]"
    :style="{ backgroundImage: `url(${backgroundImage})` }" @pointerdown="hideInfoBox">
    <LoadingIndicator v-if="isLoading" :message="loadingMessage" />

    <template v-else>
      <MapView ref="mapViewComponentRef" :islands="islands" :edges="edges" :player="player" :dynamicViewBox="dynamicViewBox"
        :territoryId="territoryId" @nodeClick="showInfoBox" @mapTransformed="updateInfoBoxPosition" />

      <PlayerInfo :player="player" :username="username" v-if="player" />

      <Transition name="popup-fade">
        <IslandInfoBox v-if="selectedIsland" :key="selectedIsland" :selectedIsland="selectedIsland" :player="player"
          :refuel="refuel" :travel="travel" :anchor="anchor" :infoBoxStyle="infoBoxStyle" :isRefuelIsland="isSelectedIslandRefuelIsland"
          :isAdjacent="isSelectedIslandAdjacent" :loading="isInfoBoxLoading" @navigateToIsland="navigateToIsland"
          @travelToIsland="travelToIsland" @buyFuel="buyFuelFromIsland" @dropAnchor="dropAnchor"/>
      </Transition>

    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, nextTick, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getPlayer, getMe, getToken, travelCheck, travelTo, refuelCheck, buyFuel, anchorCheck, dropAnchorAtIsland, logout, getTerritory } from "@/services/api.js";
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
const isInfoBoxLoading = ref(false);

const territoryId = ref(route.params.id);
const islands = ref([]);
const refuelIslands = ref([]);
const edges = ref([]);
const player = ref(null);
const username = ref('...');
const travel = ref(null);
const anchor = ref(null);
const travelError = ref(null);
const anchorError = ref(null);
const refuel = ref(null);
const backgroundImage = ref('');
const selectedIsland = ref(null);
const dynamicViewBox = ref('0 0 1 1');
const loadingMessage = ref('Loading map data...');
const isLoading = ref(true);

// --- Computed Properties ---
const isSelectedIslandRefuelIsland = computed(() => {
  if (!selectedIsland.value) return false;
  return refuelIslands.value.some(island => island.id === selectedIsland.value.id);
});

const isSelectedIslandAdjacent = computed(() => {
  if (!selectedIsland.value || !player.value) return false;
  return edges.value.some(edge =>
    (edge.from === player.value.atIsland && edge.to === selectedIsland.value.id) ||
    (edge.to === player.value.atIsland && edge.from === selectedIsland.value.id)
  );
});

const infoBoxStyle = computed(() => {
  const _ = transformCounter.value;
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

const navigateToIsland = (islandId) => {
  router.push({ name: 'Island', params: { id: territoryId.value, islandId: islandId } });
};

const travelToIsland = async (dest) => {
  try {
    travelError.value = null;
    await travelTo(player.value.atIsland, dest);
  } catch (error) {
    travelError.value = error.message;
  }
};

const dropAnchor = async () => {
  try {
    anchorError.value = null;
    await dropAnchorAtIsland(player.value.atIsland);
  } catch (error) {
    anchorError.value = error.message;
  }
};

const buyFuelFromIsland = (fuelAmount) => {
  buyFuel(fuelAmount);
};

const calculateViewBox = (islands, padding = 0.1) => {
  const bounds = islands.reduce((acc, island) => ({
    minX: Math.min(acc.minX, island.x - island.width / 2),
    maxX: Math.max(acc.maxX, island.x + island.width / 2),
    minY: Math.min(acc.minY, island.y - island.height / 2),
    maxY: Math.max(acc.maxY, island.y + island.height / 2),
  }), { minX: Infinity, maxX: -Infinity, minY: Infinity, maxY: -Infinity });
  const { minX, minY, maxX, maxY } = bounds;
  return `${minX - padding} ${minY - padding} ${maxX - minX + padding * 2} ${maxY - minY + padding * 2}`;
};

// --- API Calls & Data Fetching ---
const fetchTerritoryData = async (id) => {
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
    console.error("Failed to get player/user data:", err);
    throw err;
  }
};

// --- Setup Functions (Processing and state setting) ---
const setupTerritoryData = (territoryData) => {
  backgroundImage.value = `/images/${territoryData.backgroundAsset}`;
  dynamicViewBox.value = calculateViewBox(territoryData.islands);
  islands.value = territoryData.islands;
  edges.value = territoryData.edges;
  refuelIslands.value = territoryData.refuelIslands;
};

const setupPlayerAndUserData = (playerAndUserData) => {
  if (!playerAndUserData) return;
  
  const { playerData, meData } = playerAndUserData;
  username.value = meData.username;
  
  player.value = playerData;
};

const updateTravel = async () => {
  if (!player.value || !selectedIsland.value) return;
  try {
    travel.value = await travelCheck(player.value.atIsland, selectedIsland.value.id);
  } catch (err) {
    travelError.value = err.message;
  }
};

const updateAnchor = async () => {
  if (!player.value || !selectedIsland.value) return;
  try {
    anchor.value = await anchorCheck(player.value.atIsland);
  } catch (err) {
    anchorError.value = err.message;
  }
};

const updateRefuel = async () => {
  try {
    refuel.value = await refuelCheck();
  } catch (err) {
    console.error("Failed to get refuel data:", err);
  }
};

// --- Event Handlers from Child Components ---
const showInfoBox = async (island) => {
  if (!player.value) return;
  if (selectedIsland.value && selectedIsland.value.id === island.id) {
    hideInfoBox();
    return;
  }

  isInfoBoxLoading.value = true;
  selectedIsland.value = island;
  travel.value = null;
  refuel.value = null;
  travelError.value = null;
  anchorError.value = null;

  try {
    const isCurrent = island.id === player.value.atIsland;
    if (isCurrent) {
      if (isSelectedIslandRefuelIsland.value) {
        await updateRefuel();
      } else if (!player.value.anchored) {
        await updateAnchor();
      } else {
        await nextTick();
      }
    } else {
      await updateTravel();
    }
  } finally {
    isInfoBoxLoading.value = false;
  }
};

const hideInfoBox = () => {
  selectedIsland.value = null;
  travel.value = null;
  refuel.value = null;
  travelError.value = null;
  anchorError.value = null;
  isInfoBoxLoading.value = false;
};

// --- Lifecycle Hooks ---
onMounted(async () => {
  isLoading.value = true;
  try {
    loadingMessage.value = 'Fetching data...';
    const [territoryData, playerAndUserData] = await Promise.all([
      fetchTerritoryData(territoryId.value),
      fetchPlayerAndUserData()
    ]);
    
    loadingMessage.value = 'Setting up data...';
    
    setupTerritoryData(territoryData);
    setupPlayerAndUserData(playerAndUserData);
    
  } finally {
    isLoading.value = false;
  }
});
// --- WebSocket ---
usePlayerWebSocket(player);

// --- Watcher to update infobox on arrival ---
watch(() => player.value?.atIsland, (newIsland, oldIsland) => {
  if (newIsland && oldIsland && newIsland !== oldIsland) {
    hideInfoBox();
  }
}, { deep: true });
</script>

<style scoped>
.popup-fade-enter-active,
.popup-fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.popup-fade-enter-from,
.popup-fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
