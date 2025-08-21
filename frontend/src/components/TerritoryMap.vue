<template>
  <div class="territory-container" :style="{ backgroundImage: `url(${backgroundImage})` }">
    <div v-if="hoveredNode" class="info-box" :style="infoBoxStyle">
      {{ hoveredNode.name }}
    </div>

    <svg
      ref="svgRef"
      v-if="nodes.length > 0"
      :viewBox="dynamicViewBox"
      preserveAspectRatio="xMidYMid meet"
      @mousemove="updateMousePosition"
      class="map-svg"
    >
      <g class="edges">
        <path
          v-for="edge in wavyEdges"
          :key="`${edge.from__node_id}-${edge.to_node_id}`"
          :d="edge.pathD"
          class="edge-path"
        />
      </g>

      <g class="nodes">
        <g v-for="node in nodes" :key="node.id" @click="navigateToIsland(node.id)" class="node-link">
          <image
            :href="node.iconPath"
            :x="node.imageX"
            :y="node.imageY"
            :width="node.width"
            :height="node.height"
            class="node-image"
            @mouseover="showInfoBox(node)"
            @mouseleave="hideInfoBox"
          />
        </g>
      </g>
      <g v-if="player && props.territoryId == player.atTerritory" class="ship">
        <image 
          href="/images/ships/ship1.svg" 
          width="0.16" 
          height="0.22" 
          :x="player.atIsland.imageX - 0.08" 
          :y="player.atIsland.imageY - 0.08"
        />
      </g>
    </svg>
    <div v-else class="loading-message">
      {{ loadingMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick } from 'vue';
import { useRouter } from 'vue-router'; // Import the router
import { getPlayer, getToken } from "@/services/api";
import panzoom from 'panzoom';

// --- Define reactive state ---
const svgRef = ref(null);
const nodes = ref([]);
const edges = ref([]);
const player = ref(null);
const backgroundImage = ref('');
const hoveredNode = ref(null);
const mousePosition = ref({ x: 0, y: 0 });
const dynamicViewBox = ref('0 0 1 1');
const loadingMessage = ref('Loading map data...');
const router = useRouter(); // Initialize the router

let panzoomInstance = null;

// --- Define component props ---
const props = defineProps({
  territoryId: {
    type: String, // ID from URL is always a string
    required: true,
  },
});

// --- THE CHANGE IS HERE: Added the navigation function ---
const navigateToIsland = (islandId) => {
  router.push(`/territory/${props.territoryId}/${islandId}`);
};

// --- Computed property to generate wavy paths ---
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

    panzoomInstance = panzoom(svgRef.value, {
        maxZoom: 4,
        minZoom: 1,
    });
}

// --- Fetch and process data from the REAL API ---
const fetchTerritoryData = async (id) => {
  const BASE_URL = 'http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1';
  // THE FIX IS HERE: We now use the 'id' directly from the props.
  const apiUrl = `${BASE_URL}/territories/${id}`;
  
  try {
    loadingMessage.value = 'Fetching data from server...';
    console.log(`Initiating fetch request to REAL API endpoint: ${apiUrl}`);
    const response = await fetch(apiUrl);
    const data = await response.json();

    if (!response.ok) { // Check HTTP status first
        throw new Error(data.error || `Territory not found (${response.status})`);
    }
    if (!data.ok || !data.result) {
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

    await nextTick();
    initializePanzoom();
    await fetchPlayer();

  } catch (error) {
    console.error('Failed to load or process territory data:', error);
    loadingMessage.value = `Error: ${error.message}`;
  }
};

const fetchPlayer = async () => {
  if (!getToken()) return;
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
} 

// --- Helper Functions ---
const getNodeById = (id) => nodes.value.find(node => node.id === id);
const updateMousePosition = (event) => { mousePosition.value = { x: event.clientX, y: event.clientY }; };
const showInfoBox = (node) => { hoveredNode.value = node; };
const hideInfoBox = () => { hoveredNode.value = null; };
const infoBoxStyle = computed(() => ({
  position: 'fixed',
  top: `${mousePosition.value.y + 20}px`,
  left: `${mousePosition.value.x}px`,
  transform: 'translateX(-50%)',
}));

// --- Lifecycle Hook ---
onMounted(() => {
  fetchTerritoryData(props.territoryId);
});
onUnmounted(() => {
  if (panzoomInstance) {
    panzoomInstance.dispose();
  }
});
</script>

<style scoped>
/* Define all colors as CSS variables for easy management */
.territory-container {
  --color-edge: #ffffff;
  --color-info-box-bg: rgba(0, 0, 0, 0.8);
  --color-info-box-text: white;
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
  overflow: hidden; /* Crucial */
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

/* THE CHANGE IS HERE: Added styles for the new group tag */
.node-link {
  cursor: pointer;
}

.node-image, .ship image {
  transition: transform 0.2s ease-in-out;
  transform-box: fill-box;
}

/* THE CHANGE IS HERE: Hover effect is now applied to the group */
.node-link:hover .node-image {
  transform-origin: center center;
  transform: scale(1.1);
}

.info-box {
  background-color: var(--color-info-box-bg);
  color: var(--color-info-box-text);
  padding: 8px 12px;
  border-radius: 6px;
  font-family: sans-serif;
  font-size: 14px;
  pointer-events: none;
  z-index: 100;
  white-space: nowrap;
}

.loading-message {
  font-size: 1.5rem;
  color: var(--color-loading-text);
  font-family: sans-serif;
  background-color: var(--color-info-box-bg);
  padding: 1rem 2rem;
  border-radius: 0.5rem;
}


/*TODO remove this part; temporary fix for svgs*/
image.node-image{
  width: 50%;
  height: 70%;
  transform-origin: 0px 0px !important;
  transform: scale(0.3) !important;
}
.ship image{
  width: 50%;
  height: 70%;
  transform-origin: 0px 0px !important;
}

/* SHIP ANIMATION */
.ship image {
  filter: drop-shadow(0 5px 20px rgba(0,0,0,0.5));
  animation: 10s linear boat-animation infinite;
}

@keyframes boat-animation {
  0% {
    transform: translate(0,0) rotate(10deg) scale(0.3);
  }
  35% {
    transform: translate(0.02px, 0.01px) rotate(-10deg) scale(0.3)
  }
  70% {
    transform: translate(-0.02px, 0.01px) rotate(3deg) scale(0.3)
  }
  100% {
    transform: translate(0,0) rotate(10deg) scale(0.3);
  }
}
</style>