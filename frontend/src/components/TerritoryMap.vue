<template>
  <div class="territory-container" :style="{ backgroundImage: `url(${backgroundImage})` }">
    <div v-if="hoveredNode" class="info-box" :style="infoBoxStyle">
      {{ hoveredNode.name }}
    </div>

    <div class="svg-wrapper" v-if="nodes.length > 0">
      <svg
        ref="svgRef"
        :width="mapWidth"
        :height="mapHeight"
        class="map-svg"
        @mousemove="updateMousePosition"
      >
        <image
          :href="backgroundImage"
          x="0"
          y="0"
          :width="mapWidth"
          :height="mapHeight"
        />
        
        <g class="edges">
          <path
            v-for="edge in wavyEdges"
            :key="`${edge.from_node_id}-${edge.to_node_id}`"
            :d="edge.pathD"
            class="edge-path"
          />
        </g>

        <g class="nodes">
          <image
            v-for="node in nodes"
            :key="node.id"
            :href="node.iconPath"
            :x="node.imageX"
            :y="node.imageY"
            :width="node.pixelWidth"
            :height="node.pixelHeight"
            class="node-image"
            @mouseover="showInfoBox(node)"
            @mouseleave="hideInfoBox"
          />
        </g>
      </svg>
    </div>
    
    <div v-else class="loading-message">
      {{ loadingMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick } from 'vue';
import panzoom from 'panzoom';

// Define a fixed coordinate system for our world
const mapWidth = 1500;
const mapHeight = 1500;

// --- Define reactive state ---
const containerRef = ref(null);
const svgRef = ref(null);
const nodes = ref([]);
const edges = ref([]);
const backgroundImage = ref('');
const hoveredNode = ref(null);
const mousePosition = ref({ x: 0, y: 0 });
const loadingMessage = ref('Loading map data...');

let panzoomInstance = null;

// --- Define component props ---
const props = defineProps({
  territoryId: {
    type: [Number, String],
    required: true,
  },
});

// --- Computed property to generate wavy paths ---
const wavyEdges = computed(() => {
  return edges.value.map(edge => {
    const startNode = getNodeById(edge.from_node_id);
    const endNode = getNodeById(edge.to_node_id);
    if (!startNode || !endNode) return { ...edge, pathD: '' };
    const x1 = startNode.pixelX, y1 = startNode.pixelY, x2 = endNode.pixelX, y2 = endNode.pixelY;
    const dx = x2 - x1, dy = y2 - y1;
    const lineLength = Math.sqrt(dx * dx + dy * dy);
    const amplitude = lineLength * 0.1;
    const perpDx = -dy / lineLength, perpDy = dx / lineLength;
    const cp1x = x1 + dx * 0.25 + perpDx * amplitude, cp1y = y1 + dy * 0.25 + perpDy * amplitude;
    const cp2x = x1 + dx * 0.75 - perpDx * amplitude, cp2y = y1 + dy * 0.75 - perpDy * amplitude;
    const pathD = `M ${x1} ${y1} C ${cp1x} ${cp1y}, ${cp2x} ${cp2y}, ${x2} ${y2}`;
    return { ...edge, pathD };
  });
});

const initializePanzoom = () => {
    if (!svgRef.value || !containerRef.value || panzoomInstance) return;
    const containerWidth = containerRef.value.clientWidth;
    const containerHeight = containerRef.value.clientHeight;
    const scaleX = containerWidth / mapWidth;
    const scaleY = containerHeight / mapHeight;
    const minZoom = Math.max(scaleX, scaleY);

    panzoomInstance = panzoom(svgRef.value, {
        maxZoom: 3,
        minZoom: minZoom,
        initialZoom: minZoom,
    });
    const allX = nodes.value.map(n => n.pixelX);
    const allY = nodes.value.map(n => n.pixelY);
    const contentCenterX = (Math.min(...allX) + Math.max(...allX)) / 2;
    const contentCenterY = (Math.min(...allY) + Math.max(...allY)) / 2;
    panzoomInstance.moveTo(
        (containerWidth / 2) - (contentCenterX * minZoom),
        (containerHeight / 2) - (contentCenterY * minZoom)
    );
};

// --- Fetch and process data from the REAL API ---
const fetchTerritoryData = async (id) => {
  const BASE_URL = 'http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1';
  // Corrected the endpoint from "territories" to "territory"
  const apiUrl = `${BASE_URL}/territories/${id}`;
  
  try {
    console.log(`Initiating fetch request to REAL API endpoint: ${apiUrl}`);
    const response = await fetch(apiUrl);
    const data = await response.json();

    if (!data.ok || !data.result) {
      throw new Error(data.error || 'Invalid API response format');
    }

    const rawData = data.result; // Extract data from the 'result' key
    backgroundImage.value = `/images/${rawData.backgroundAsset}`;
    
    const transformedNodes = rawData.islands.map(island => {
      const pixelWidth = island.width * mapWidth;
      const pixelHeight = island.height * mapHeight;
      const pixelX = island.x * mapWidth;
      const pixelY = island.y * mapHeight;
      return {
        id: island.id, name: island.name,
        pixelX: pixelX, pixelY: pixelY,
        pixelWidth: pixelWidth, pixelHeight: pixelHeight,
        iconPath: `/images/islands/${island.iconAsset}`,
        imageX: pixelX - (pixelWidth / 2),
        imageY: pixelY - (pixelHeight / 2),
      };
    });
    const transformedEdges = rawData.edges.map(edge => ({
      from_node_id: edge.from, to_node_id: edge.to,
    }));
    nodes.value = transformedNodes;
    edges.value = transformedEdges;
    
    await nextTick();
    initializePanzoom();
  } catch (error) {
    console.error('Failed to load territory data from API:', error);
    loadingMessage.value = `Error loading map: ${error.message}`;
  }
};

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

// --- Lifecycle Hooks ---
onMounted(() => {
  fetchTerritoryData(props.territoryId);
});
onUnmounted(() => {
  if (panzoomInstance) panzoomInstance.dispose();
});
</script>

<style scoped>
/* Styles remain the same */
.territory-container {
  --color-edge: #ffffff;
  --color-info-box-bg: rgba(0, 0, 0, 0.8);
  --color-info-box-text: white;
  --color-loading-text: #ddd;
  --color-bg-fallback: #0c2036;

  width: 100vw;
  height: 100vh;
  position: relative;
  background-color: var(--color-bg-fallback);
  overflow: hidden;
}
.svg-wrapper {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
}
.map-svg {
  flex-shrink: 0;
  cursor: grab;
}
.map-svg:active {
    cursor: grabbing;
}
.edge-path {
  fill: none;
  stroke: var(--color-edge);
  stroke-width: 2;
  stroke-dasharray: 6, 4;
  pointer-events: none;
}
.node-image {
  transition: transform 0.2s ease-in-out;
  transform-box: fill-box;
}
.node-image:hover {
  transform-origin: center center;
  transform: scale(1.1);
}
.info-box {
  position: absolute;
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
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 1.5rem;
  color: var(--color-loading-text);
  font-family: sans-serif;
  background-color: var(--color-info-box-bg);
  padding: 1rem 2rem;
  border-radius: 0.5rem;
}
</style>