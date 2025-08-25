<template>
  <div class="territory-container" :style="{ backgroundImage: `url(${backgroundImage})` }">

    <div v-if="hoveredNode" class="info-box" :style="infoBoxStyle" @mouseover="isHoveringBox = true"
      @mouseleave="unHoverBox">
      <div>{{ hoveredNode.name }}</div>
      <div>
        <div v-if="hoveredNode.id == player.atIsland.id && fuelStations.some(station => station.id === hoveredNode.id)"
          class="refuel">
          <div v-if="refuel"> قیمت هر واحد: {{ refuel.coinCostPerUnit }} </div>
          <div v-if="refuel"> حداکثر واحد قابل اخذ: {{ refuel.maxAvailableAmount }} </div>
          <input type="number" ref="fuelInput" :max="refuel ? refuel.maxAvailableAmount : player.fuelCap - player.fuel"
            v-model.number="fuelCount" @pointerdown="focusFuelInput" />
          <button @pointerdown="buyFuelFromIsland">{{ fuelPriceText }}</button>
        </div>
        <button v-else-if="hoveredNode.id == player.atIsland.id" @pointerdown="navigateToIsland(player.atIsland.id)">
          ورود به جزیره
        </button>
        <button v-else :disabled="!edges.some(edge =>
          (edge.from_node_id === player.atIsland.id && edge.to_node_id === hoveredNode.id) ||
          (edge.to_node_id === player.atIsland.id && edge.from_node_id === hoveredNode.id)
        )" @pointerdown="travelToIsland(hoveredNode.id)">
          سفر به جزیره
          <span v-if="fuelCost">{{ fuelCost }}</span>
        </button>
      </div>
    </div>

    <svg ref="svgRef" v-if="nodes.length > 0" :viewBox="dynamicViewBox" preserveAspectRatio="xMidYMid meet"
      class="map-svg">
      <g class="edges">
        <path v-for="edge in wavyEdges" :key="`${edge.from__node_id}-${edge.to_node_id}`" :d="edge.pathD"
          class="edge-path" />
      </g>

      <g class="nodes">
        <g v-for="node in nodes" :key="node.id" class="node-link">
          <image :href="node.iconPath" :x="node.imageX" :y="node.imageY" :width="node.width" :height="node.height"
            class="node-image" @pointerdown="showInfoBox(node)" @mouseover="isHoveringNode = true"
            @mouseleave="unhoverNode" />
        </g>
      </g>

      <g v-if="player && props.territoryId == player.atTerritory" class="ship">
        <image href="/images/ships/ship1.svg" width="0.16" height="0.22" :x="player.atIsland.imageX - 0.08"
          :y="player.atIsland.imageY - 0.08" />
      </g>
    </svg>
    <div v-else class="loading-message">
      {{ loadingMessage }}
    </div>

    <div class="fuel-bar-container" v-if="player">
      <div class="fuel-bar-label">
        <img src="/images/icons/drop.svg" alt="Drop Icon" width="24" height="24">
      </div>
      <div class="fuel-bar">
        <div class="fuel-bar-fill" :style="{ height: `${(player.fuel / player.fuelCap) * 100}%` }"></div>
        <div class="fuel-wave"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick, watch } from 'vue';
import { useRouter } from 'vue-router';
import { logout, getPlayer, getToken, checkTravel, travelTo, refuelCheck, buyFuel } from "@/services/api";
import { usePlayerWebSocket } from '@/components/service/WebSocket';
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
  router.push(`/territory/${props.territoryId}/${islandId}`);
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
    router.push(`/login`);
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

<style scoped>
.territory-container {
  --color-edge: #ffffff;
  --color-info-box-bg: rgb(121 200 237 / 80%);
  --color-info-box-text: #310f0f;
  --color-loading-text: #ddd;
  --color-bg-fallback: #0c2036;
  width: 100%;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 1rem;
  box-sizing: border-box;
  background-color: var(--color-bg-fallback);
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  overflow: hidden;
}

.map-svg {
  width: 100%;
  height: 100%;
  display: block;
  cursor: grab;
}

.map-svg:active {
  cursor: grabbing;
}

.edge-path {
  fill: none;
  stroke: var(--color-edge);
  stroke-width: 0.005;
  stroke-dasharray: 0.01, 0.005;
}

.node-link {
  cursor: pointer;
}

.node-image,
.ship image {
  transition: transform 0.2s ease-in-out;
  transform-box: fill-box;
}

.node-link:hover .node-image {
  transform-origin: center center;
  transform: scale(1.1);
}

.info-box {
  background-color: var(--color-info-box-bg);
  color: var(--color-info-box-text);
  padding: 13px 20px;
  border-radius: 6px;
  font-family: var(--font-vazir);
  font-size: 16px;
  z-index: 10000;
  white-space: nowrap;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  pointer-events: auto;
}

.refuel {
  white-space: nowrap;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.info-box button {
  padding: 5px;
  border-radius: 10px;
  background: #07458bb5;
  margin: 10px 0 0;
}

.info-box input {
  width: 4rem;
  border-radius: 10px;
  border: 1px solid #07458bb5;
  text-align: center;
}

.info-box button[disabled] {
  filter: contrast(0.5);
}

.loading-message {
  font-size: 1.5rem;
  color: var(--color-loading-text);
  font-family: sans-serif;
  background-color: var(--color-info-box-bg);
  padding: 1rem 2rem;
  border-radius: 0.5rem;
}

image.node-image {
  width: 50%;
  height: 70%;
  transform-origin: 0px 0px !important;
  transform: scale(0.3) !important;
}

.ship image {
  width: 50%;
  height: 70%;
  transform-origin: 0px 0px !important;
}

.ship image {
  filter: drop-shadow(0 5px 20px rgba(0, 0, 0, 0.5));
  animation: 10s linear boat-animation infinite;
}

@keyframes boat-animation {
  0% {
    transform: translate(0, 0) rotate(10deg) scale(0.3);
  }

  35% {
    transform: translate(0.02px, 0.01px) rotate(-10deg) scale(0.3)
  }

  70% {
    transform: translate(-0.02px, 0.01px) rotate(3deg) scale(0.3)
  }

  100% {
    transform: translate(0, 0) rotate(10deg) scale(0.3);
  }
}

.fuel-bar-container {
  position: fixed;
  top: 50%;
  right: 20px;
  transform: translateY(-50%);
  width: 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
  z-index: 200;
  font-family: var(--font-vazir);
}

.fuel-bar-label {
  margin-bottom: 8px;
  font-size: 14px;
  color: #fff;
}

.fuel-bar {
  width: 20px;
  height: 150px;
  background-color: rgb(38 72 174 / 95%);
  border-radius: 6px;
  overflow: hidden;
  position: relative;
  display: flex;
  flex-direction: column-reverse;
  padding: 3px;
}

.fuel-bar-fill {
  width: 100%;
  background-color: #000000;
  border-radius: 6px 6px 0 0;
  position: relative;
  overflow: hidden;
  transition: height 0.3s ease;
}
</style>
