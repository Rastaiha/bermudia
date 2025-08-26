<template>
  <div
    class="w-full h-screen flex justify-center items-center p-4 box-border bg-cover bg-center bg-no-repeat overflow-hidden bg-[#0c2036]"
    :style="{ backgroundImage: `url(${backgroundImage})` }" @pointerdown="hideInfoBox">
    <LoadingIndicator v-if="isLoading" :message="loadingMessage" />

    <template v-else>
      <MapView ref="mapViewComponentRef" :nodes="nodes" :edges="edges" :player="player" :dynamicViewBox="dynamicViewBox"
        :territoryId="territoryId" @nodeClick="showInfoBox" />
      <PlayerInfo :player="player" />
      <IslandInfoBox :hoveredNode="hoveredNode" :player="player" :refuel="refuel" :travel="travel"
        :infoBoxStyle="infoBoxStyle" :isFuelStation="isHoveredNodeFuelStation" :isAdjacent="isHoveredNodeAdjacent"
        @navigateToIsland="navigateToIsland" @travelToIsland="travelToIsland" @buyFuel="buyFuelFromIsland" />
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getPlayer, getToken, checkTravel, travelTo, refuelCheck, buyFuel, logout } from "@/services/api.js";
import { usePlayerWebSocket } from '@/components/service/WebSocket.js';

import MapView from '@/components/MapView.vue';
import IslandInfoBox from '@/components/IslandInfoBox.vue';
import PlayerInfo from '@/components/PlayerInfo.vue';
import LoadingIndicator from '@/components/LoadingIndicator.vue';

// --- State ---
const route = useRoute();
const router = useRouter();
const mapViewComponentRef = ref(null);

const territoryId = ref(route.params.id);
const nodes = ref([]);
const fuelStations = ref([]);
const edges = ref([]);
const player = ref(null);
const travel = ref(null);
const refuel = ref(null);
const backgroundImage = ref('');
const hoveredNode = ref(null);
const dynamicViewBox = ref('0 0 1 1');
const loadingMessage = ref('Loading map data...');
const isLoading = ref(true);

// --- Constants ---
const BASE_URL = 'http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1';

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
  const svgElement = mapViewComponentRef.value?.svgRef;
  if (!hoveredNode.value || !svgElement) {
    return { display: 'none' };
  }

  const pt = svgElement.createSVGPoint();
  pt.x = hoveredNode.value.x;
  pt.y = hoveredNode.value.y;

  const screenPoint = pt.matrixTransform(svgElement.getScreenCTM());
  return {
    position: 'fixed',
    top: `${screenPoint.y - 100}px`,
    left: `${screenPoint.x}px`,
    transform: 'translateX(-50%)',
    display: 'flex',
  };
});

// --- Methods ---
const getNodeById = (id) => nodes.value.find(node => node.id === id);

const navigateToIsland = (islandId) => {
  router.push({ name: 'Island', params: { id: territoryId.value, islandId: islandId } });
};

const travelToIsland = (dest) => {
  travelTo(player.value.atIsland.id, dest);
};

const buyFuelFromIsland = (fuelAmount) => {
  buyFuel(fuelAmount);
};

// --- API Calls & Data Fetching ---
const fetchTerritoryData = async (id) => {
  isLoading.value = true;
  // ... (این بخش بدون تغییر باقی می‌ماند)
  loadingMessage.value = 'Fetching data from server...';
  try {
    const response = await fetch(`${BASE_URL}/territories/${id}`);
    const data = await response.json();
    if (!response.ok || !data.ok || !data.result) throw new Error(data.error || 'Invalid API response');

    loadingMessage.value = 'Processing data...';
    const rawData = data.result;
    backgroundImage.value = `/images/${rawData.backgroundAsset}`;

    const padding = 0.1;
    let minX = Infinity, maxX = -Infinity, minY = Infinity, maxY = -Infinity;
    rawData.islands.forEach(island => {
      minX = Math.min(minX, island.x - island.width / 2);
      maxX = Math.max(maxX, island.x + island.width / 2);
      minY = Math.min(minY, island.y - island.height / 2);
      maxY = Math.max(maxY, island.y + island.height / 2);
    });
    dynamicViewBox.value = `${minX - padding} ${minY - padding} ${maxX - minX + padding * 2} ${maxY - minY + padding * 2}`;

    nodes.value = rawData.islands.map(island => ({
      ...island,
      iconPath: `/images/islands/${island.iconAsset}`,
      imageX: island.x - island.width / 2,
      imageY: island.y - island.height / 2,
    }));
    edges.value = rawData.edges.map(edge => ({ from_node_id: edge.from, to_node_id: edge.to }));
    fuelStations.value = rawData.refuelIslands;

    await fetchPlayer();
  } catch (error) {
    console.error('Failed to load territory data:', error);
    loadingMessage.value = `Error: ${error.message}`;
  } finally {
    isLoading.value = false;
  }
};

const fetchPlayer = async () => {
  // ... (این بخش بدون تغییر باقی می‌ماند)
  if (!getToken()) {
    logout();
    router.push({ name: 'Login' });
    return;
  }
  try {
    const playerData = await getPlayer();
    const island = getNodeById(playerData.atIsland);
    player.value = {
      atTerritory: playerData.atTerritory,
      atIsland: island,
      fuel: playerData.fuel,
      fuelCap: playerData.fuelCap
    };
  } catch (err) {
    console.error("Failed to get player data:", err);
  }
};

const updateTravel = async () => {
  if (!player.value || !hoveredNode.value) return;
  try {
    travel.value = await checkTravel(player.value.atIsland.id, hoveredNode.value.id);
  } catch (err) {
    console.error("Failed to get travel data:", err);
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

  // اگر روی همان جزیره باز کلیک شد، آن را ببند
  if (hoveredNode.value && hoveredNode.value.id === node.id) {
    hideInfoBox();
    return;
  }

  hoveredNode.value = node;

  travel.value = null;
  refuel.value = null;

  if (node.id === player.value.atIsland.id) {
    if (isHoveredNodeFuelStation.value) {
      updateRefuel();
    }
  } else {
    updateTravel();
  }
};

const hideInfoBox = () => {
  hoveredNode.value = null;
  travel.value = null;
  refuel.value = null;
};

// --- Lifecycle Hooks ---
onMounted(() => {
  fetchTerritoryData(territoryId.value);
});

// --- WebSocket ---
usePlayerWebSocket(player, nodes);
</script>