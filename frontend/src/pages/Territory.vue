<template>
  <div
    class="w-full h-screen flex justify-center items-center p-4 box-border bg-cover bg-center bg-no-repeat overflow-hidden bg-[#0c2036]"
    :style="{ backgroundImage: `url(${backgroundImage})` }" @pointerdown="hideInfoBox">
    <LoadingIndicator v-if="isLoading" :message="loadingMessage" />

    <template v-else>
      <MapView ref="mapViewComponentRef" :nodes="nodes" :edges="edges" :player="player" :dynamicViewBox="dynamicViewBox"
        :territoryId="territoryId" @nodeClick="showInfoBox" @mapTransformed="updateInfoBoxPosition" />

      <PlayerInfo :player="player" :username="username" v-if="player" />

      <Transition name="popup-fade">
        <IslandInfoBox v-if="hoveredNode" :key="hoveredNode.id" :hoveredNode="hoveredNode" :player="player"
          :refuel="refuel" :travel="travel" :infoBoxStyle="infoBoxStyle" :isFuelStation="isHoveredNodeFuelStation"
          :isAdjacent="isHoveredNodeAdjacent" :loading="isInfoBoxLoading" @navigateToIsland="navigateToIsland"
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
const nodes = ref([]);
const fuelStations = ref([]);
const edges = ref([]);
const player = ref(null);
const username = ref('...');
const travel = ref(null);
const travelError = ref(null);
const refuel = ref(null);
const backgroundImage = ref('');
const hoveredNode = ref(null);
const dynamicViewBox = ref('0 0 1 1');
const loadingMessage = ref('Loading map data...');
const isLoading = ref(true);

// --- Computed Properties ---
const isHoveredNodeFuelStation = computed(() => {
  if (!hoveredNode.value) return false;
  return fuelStations.value.some(station => station.id === hoveredNode.value.id);
});

const isHoveredNodeAdjacent = computed(() => {
  if (!hoveredNode.value || !player.value) return false;
  return edges.value.some(edge =>
    (edge.from_node_id === player.value.atIsland.id && edge.to_node_id === hoveredNode.value.id) ||
    (edge.to_node_id === player.value.atIsland.id && edge.from_node_id === hoveredNode.value.id)
  );
});

const infoBoxStyle = computed(() => {
  const _ = transformCounter.value;
  const svgElement = mapViewComponentRef.value?.svgRef;
  if (!hoveredNode.value || !svgElement) return { display: 'none' };
  const pt = svgElement.createSVGPoint();
  pt.x = hoveredNode.value.x;
  pt.y = hoveredNode.value.y;
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

const getNodeById = (id) => nodes.value.find(node => node.id === id);

const navigateToIsland = (islandId) => {
  router.push({ name: 'Island', params: { id: territoryId.value, islandId: islandId } });
};

const travelToIsland = async (dest) => {
  try {
    travelError.value = null;
    await travelTo(player.value.atIsland.id, dest);
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
  
  nodes.value = territoryData.islands.map(island => ({
    ...island,
    iconPath: `/images/islands/${island.iconAsset}`,
    imageX: island.x - island.width / 2,
    imageY: island.y - island.height / 2,
  }));
  
  edges.value = territoryData.edges.map(edge => ({
    from_node_id: edge.from,
    to_node_id: edge.to
  }));
  
  fuelStations.value = territoryData.refuelIslands;
};

const setupPlayerAndUserData = (playerAndUserData) => {
  if (!playerAndUserData) return;
  
  const { playerData, meData } = playerAndUserData;
  username.value = meData.username;
  
  const island = getNodeById(playerData.atIsland);
  player.value = {
    ...playerData,
    atIsland: island,
  };
};

const updateTravel = async () => {
  if (!player.value || !hoveredNode.value) return;
  try {
    travel.value = await checkTravel(player.value.atIsland.id, hoveredNode.value.id);
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
const showInfoBox = async (node) => {
  if (!player.value) return;
  if (hoveredNode.value && hoveredNode.value.id === node.id) {
    hideInfoBox();
    return;
  }

  isInfoBoxLoading.value = true;
  hoveredNode.value = node;
  travel.value = null;
  refuel.value = null;
  travelError.value = null;

  try {
    const isCurrent = node.id === player.value.atIsland.id;
    if (isCurrent) {
      if (isHoveredNodeFuelStation.value) {
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
  hoveredNode.value = null;
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
usePlayerWebSocket(player, nodes);

// --- Watcher to update infobox on arrival ---
watch(() => player.value?.atIsland, (newIsland, oldIsland) => {
  if (newIsland && oldIsland && newIsland.id !== oldIsland.id) {
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
