<template>
  <div
    class="w-full h-screen flex justify-center items-center p-4 box-border bg-cover bg-center bg-no-repeat overflow-hidden bg-[#0c2036]"
    :style="{ backgroundImage: `url(${backgroundImage})` }">

    <div v-if="hoveredNode"
      class="bg-[rgb(121,200,237,0.8)] text-[#310f0f] px-5 py-3 rounded-md font-vazir text-base z-[10000] whitespace-nowrap flex flex-col justify-center items-center pointer-events-auto"
      :style="infoBoxStyle" @mouseover="isHoveringBox = true" @mouseleave="unHoverBox">
      <div>{{ hoveredNode.name }}</div>
      <div>
        <div v-if="hoveredNode.id == player.atIsland.id && fuelStations.some(station => station.id === hoveredNode.id)"
          class="whitespace-nowrap flex flex-col justify-center items-center">
          <div v-if="refuel"> قیمت هر واحد: {{ refuel.coinCostPerUnit }} </div>
          <div v-if="refuel"> حداکثر واحد قابل اخذ: {{ refuel.maxAvailableAmount }} </div>
          <input type="number" ref="fuelInput" :max="refuel ? refuel.maxAvailableAmount : player.fuelCap - player.fuel"
            v-model.number="fuelCount" @pointerdown="focusFuelInput"
            class="w-16 rounded-lg border border-[#07458bb5] text-center" />
          <button @pointerdown="buyFuelFromIsland" class="p-1.5 rounded-lg bg-[#07458bb5] mt-2.5">{{ fuelPriceText
            }}</button>
        </div>
        <button v-else-if="hoveredNode.id == player.atIsland.id" @pointerdown="navigateToIsland(player.atIsland.id)"
          class="p-1.5 rounded-lg bg-[#07458bb5] mt-2.5 disabled:contrast-50">
          ورود به جزیره
        </button>
        <button v-else :disabled="!edges.some(edge =>
          (edge.from_node_id === player.atIsland.id && edge.to_node_id === hoveredNode.id) ||
          (edge.to_node_id === player.atIsland.id && edge.from_node_id === hoveredNode.id)
        )" @pointerdown="travelToIsland(hoveredNode.id)"
          class="p-1.5 rounded-lg bg-[#07458bb5] mt-2.5 disabled:contrast-50">
          سفر به جزیره
          <span v-if="fuelCost">{{ fuelCost }}</span>
        </button>
      </div>
    </div>

    <svg ref="svgRef" v-if="nodes.length > 0" :viewBox="dynamicViewBox" preserveAspectRatio="xMidYMid meet"
      class="w-full h-full block cursor-grab active:cursor-grabbing">
      <g class="edges">
        <path v-for="edge in wavyEdges" :key="`${edge.from_node_id}-${edge.to_node_id}`" :d="edge.pathD"
          class="fill-none stroke-white" stroke-width="0.005" stroke-dasharray="0.01, 0.005" />
      </g>

      <g class="nodes">
        <g v-for="node in nodes" :key="node.id" class="cursor-pointer group">
          <ellipse :cx="node.x" :cy="node.y" :rx="node.width / 2" :ry="node.height / 2" fill="transparent"
            @pointerdown="showInfoBox(node)" @mouseover="isHoveringNode = true" @mouseleave="unhoverNode" />
          <image :href="node.iconPath" :x="node.imageX" :y="node.imageY" :width="node.width" :height="node.height"
            class="transition-transform duration-200 ease-in-out origin-center group-hover:scale-105 pointer-events-none" />
        </g>
      </g>

      <g v-if="player && props.territoryId == player.atTerritory" class="ship">
        <image href="/images/ships/ship1.png" width="0.16" height="0.22" :x="player.atIsland.imageX - 0.08"
          :y="player.atIsland.imageY - 0.08" class="drop-shadow-lg animate-boat" />
      </g>
    </svg>
    <div v-else class="text-2xl text-gray-300 font-sans bg-[rgb(121,200,237,0.8)] px-8 py-4 rounded-lg">
      {{ loadingMessage }}
    </div>

    <div class="fixed top-1/2 right-5 -translate-y-1/2 w-10 flex flex-col items-center z-[200] font-vazir"
      v-if="player">
      <div class="mb-2 text-sm text-white">
        <img src="/images/icons/drop.svg" alt="Drop Icon" width="24" height="24">
      </div>
      <div class="w-5 h-[150px] bg-blue-900/95 rounded-md overflow-hidden relative flex flex-col-reverse p-0.5">
        <div class="w-full bg-black rounded-t-md relative overflow-hidden transition-height duration-300 ease-in"
          :style="{ height: `${(player.fuel / player.fuelCap) * 100}%` }"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick, watch } from 'vue';
import { useRouter } from 'vue-router';
import { getPlayer, getToken, checkTravel, travelTo, refuelCheck, buyFuel, logout } from "@/services/api.js";
import { usePlayerWebSocket } from '@/components/service/WebSocket.js';
import panzoom from 'panzoom';

const BASE_URL = 'http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1';
const svgRef = ref(null);
const fuelInput = ref(null);
const nodes = ref([]);
const fuelStations = ref([]);
const edges = ref([]);
const player = ref(null);
const travel = ref(null);
const fuelCount = ref(0);
const refuel = ref(null);
const backgroundImage = ref('');
const hoveredNode = ref(null);
const dynamicViewBox = ref('0 0 1 1');
const loadingMessage = ref('Loading map data...');
const isHoveringBox = ref(false);
const isHoveringNode = ref(false);
const router = useRouter();
let panzoomInstance = null;

const props = defineProps({
  territoryId: {
    type: String,
    required: true,
  },
});

const fuelCost = computed(() => travel.value?.fuelCost ?? null);

const navigateToIsland = (islandId) => {
  router.push({ name: 'Island', params: { id: props.territoryId, islandId: islandId } });
};

const travelToIsland = (dest) => {
  travelTo(player.value.atIsland.id, dest);
};

const buyFuelFromIsland = () => {
  buyFuel(fuelCount.value);
};

const wavyEdges = computed(() => {
  return edges.value.map(edge => {
    const startNode = getNodeById(edge.from_node_id);
    const endNode = getNodeById(edge.to_node_id);

    if (!startNode || !endNode) {
      return { ...edge, pathD: '' };
    }

    const x1 = startNode.x;
    const y1 = startNode.y;
    const x2 = endNode.x;
    const y2 = endNode.y;

    const dx = x2 - x1;
    const dy = y2 - y1;
    const lineLength = Math.sqrt(dx * dx + dy * dy);

    const amplitude = lineLength * 0.2;
    const perpDx = -dy / lineLength;
    const perpDy = dx / lineLength;

    const cp1x = x1 + dx * 0.25 + perpDx * amplitude;
    const cp1y = y1 + dy * 0.25 + perpDy * amplitude;
    const cp2x = x1 + dx * 0.75 - perpDx * amplitude;
    const cp2y = y1 + dy * 0.75 - perpDy * amplitude;

    const pathD = `M ${x1} ${y1} C ${cp1x} ${cp1y}, ${cp2x} ${cp2y}, ${x2} ${y2}`;

    return { ...edge, pathD };
  });
});

const initializePanzoom = () => {
  if (!svgRef.value || panzoomInstance) return;
  panzoomInstance = panzoom(svgRef.value, { maxZoom: 4, minZoom: 1 });

  svgRef.value.addEventListener('pointerdown', (e) => {
    if (e.target.closest('.info-box')) {
      e.stopPropagation();
    }
  });
};

const fetchTerritoryData = async (id) => {
  const apiUrl = `${BASE_URL}/territories/${id}`;

  try {
    loadingMessage.value = 'Fetching data from server...';
    const response = await fetch(apiUrl);
    const data = await response.json();

    if (!response.ok || !data.ok || !data.result) {
      throw new Error(data.error || 'Invalid API response format');
    }

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
    const viewBoxX = minX - padding;
    const viewBoxY = minY - padding;
    const viewBoxWidth = (maxX - minX) + (padding * 2);
    const viewBoxHeight = (maxY - minY) + (padding * 2);
    dynamicViewBox.value = `${viewBoxX} ${viewBoxY} ${viewBoxWidth} ${viewBoxHeight}`;

    const transformedNodes = rawData.islands.map(island => ({
      id: island.id, name: island.name,
      x: island.x, y: island.y,
      width: island.width, height: island.height,
      iconPath: `/images/islands/${island.iconAsset}`,
      imageX: island.x - island.width / 2,
      imageY: island.y - island.height / 2,
    }));
    const transformedEdges = rawData.edges.map(edge => ({
      from_node_id: edge.from,
      to_node_id: edge.to,
    }));
    nodes.value = transformedNodes;
    edges.value = transformedEdges;
    fuelStations.value = rawData.refuelIslands;

    await nextTick();
    initializePanzoom();
    await fetchPlayer();

  } catch (error) {
    console.error('Failed to load or process territory data:', error);
    loadingMessage.value = `Error: ${error.message}`;
  }
};

const fetchPlayer = async () => {
  if (!getToken()) {
    logout();
    router.push({ name: 'Login' });
    return;
  }
  try {
    const playerData = await getPlayer();
    const island = nodes.value.find(node => node.id === playerData.atIsland);
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

const fuelPriceText = computed(() => {
  if (!refuel.value) return "خرید سوخت";
  return `خرید سوخت ${refuel.value.coinCostPerUnit * fuelCount.value}`;
});

const getNodeById = (id) => nodes.value.find(node => node.id === id);

const showInfoBox = async (node) => {
  isHoveringNode.value = true;
  if (!player) return;
  hoveredNode.value = node;
  if (hoveredNode.value.id == player.value.atIsland.id) {
    if (fuelStations.value.some(station => station.id === hoveredNode.value.id)) {
      updateRefuel();
    }
  } else {
    updateTravel();
  }
};

const focusFuelInput = () => {
  requestAnimationFrame(() => {
    fuelInput.value?.focus();
  });
};

const updateTravel = async () => {
  try {
    const travelData = await checkTravel(player.value.atIsland.id, hoveredNode.value.id);
    travel.value = {
      feasible: travelData.feasible,
      fuelCost: travelData.fuelCost,
      reason: travelData.feasible ? "" : travelData.reason,
    }
  } catch (err) {
    console.error("Failed to get travel data:", err);
    loadingMessage.value = "Error: " + err;
  }
};

const updateRefuel = async () => {
  try {
    const refuelData = await refuelCheck();
    refuel.value = {
      maxReason: refuelData.maxReason,
      coinCostPerUnit: refuelData.coinCostPerUnit,
      maxAvailableAmount: refuelData.maxAvailableAmount,
    }
  } catch (err) {
    console.error("Failed to get refuel data:", err);
    loadingMessage.value = "Error: " + err;
  }
};

const hideInfoBox = () => {
  hoveredNode.value = null;
  travel.value = null;
  refuel.value = null;
};

const unhoverNode = () => {
  isHoveringNode.value = false;
  setTimeout(() => {
    if (!isHoveringBox.value) hideInfoBox();
  }, 1000);
};

const unHoverBox = () => {
  isHoveringBox.value = false;
  setTimeout(() => {
    if (!isHoveringNode.value) hideInfoBox();
  }, 1000);
};

const infoBoxStyle = computed(() => {
  if (!hoveredNode.value || !svgRef.value) return {};
  const svg = svgRef.value;
  const pt = svg.createSVGPoint();
  pt.x = hoveredNode.value.x;
  pt.y = hoveredNode.value.y;
  const screenPoint = pt.matrixTransform(svg.getScreenCTM());
  return {
    position: 'fixed',
    top: `${screenPoint.y - 100}px`,
    left: `${screenPoint.x}px`,
    transform: 'translateX(-50%)',
  };
});

onMounted(() => { fetchTerritoryData(props.territoryId); });
onUnmounted(() => { if (panzoomInstance) panzoomInstance.dispose(); });

watch(fuelCount, (newValue) => {
  if (!refuel.value) return;
  if (newValue > refuel.value.maxAvailableAmount) fuelCount.value = refuel.value.maxAvailableAmount;
  else if (newValue < 0) fuelCount.value = 0;
});

usePlayerWebSocket(player, nodes);
</script>