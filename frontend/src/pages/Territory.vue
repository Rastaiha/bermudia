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
        <IslandInfoBox v-if="hoveredIsland" :key="hoveredIsland" :hoveredIsland="hoveredIsland" :hoveredIslandName="getIslandById(hoveredIsland).name"
        :player="player"
          :refuel="refuel" :travel="travel" :infoBoxStyle="infoBoxStyle" :isFuelStation="ishoveredIslandFuelStation"
          :isAdjacent="ishoveredIslandAdjacent" :loading="isInfoBoxLoading" @navigateToIsland="navigateToIsland"
          @travelToIsland="travelToIsland" @buyFuel="buyFuelFromIsland" />
      </Transition>

    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, nextTick, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getPlayer, getMe, getToken, checkTravel, travelTo, refuelCheck, buyFuel, logout, getTerritory } from "@/services/api.js";
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
const travelError = ref(null);
const refuel = ref(null);
const backgroundImage = ref('');
const hoveredIsland = ref(null);
const dynamicViewBox = ref('0 0 1 1');
const loadingMessage = ref('Loading map data...');
const isLoading = ref(true);

// --- Computed Properties ---
const ishoveredIslandFuelStation = computed(() => {
  if (!hoveredIsland.value) return false;
  console.log(refuelIslands.value);
  console.log(hoveredIsland.value);
  return refuelIslands.value.some(station => station.id === hoveredIsland.value);
});

const ishoveredIslandAdjacent = computed(() => {
  if (!hoveredIsland.value || !player.value) return false;
  return edges.value.some(edge =>
    (edge.from === player.value.atIsland && edge.to === hoveredIsland.value) ||
    (edge.to === player.value.atIsland && edge.from === hoveredIsland.value)
  );
});

const infoBoxStyle = computed(() => {
  const _ = transformCounter.value;
  const svgElement = mapViewComponentRef.value?.svgRef;
  if (!hoveredIsland.value || !svgElement) return { display: 'none' };
  const island = getIslandById(hoveredIsland.value);
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

const getIslandById = (id) => islands.value.find(island => island.id === id);

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
  if (!player.value || !hoveredIsland.value) return;
  try {
    travel.value = await checkTravel(player.value.atIsland, hoveredIsland.value.id);
  } catch (err) {
    travelError.value = err.message;
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
  if (hoveredIsland.value && hoveredIsland.value === island.id) {
    hideInfoBox();
    return;
  }

  isInfoBoxLoading.value = true;
  hoveredIsland.value = island.id;
  travel.value = null;
  refuel.value = null;
  travelError.value = null;

  try {
    const isCurrent = island.id === player.value.atIsland;
    if (isCurrent) {
      if (ishoveredIslandFuelStation.value) {
        await updateRefuel();
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
  hoveredIsland.value = null;
  travel.value = null;
  refuel.value = null;
  travelError.value = null;
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
