<template>
  <div class="territory-container" :style="{ backgroundImage: `url(${backgroundImage})` }">
    <div v-if="hoveredNode" class="info-box" :style="infoBoxStyle">
      {{ hoveredNode.name }}
    </div>

    <svg
      v-if="nodes.length > 0"
      :viewBox="dynamicViewBox"
      preserveAspectRatio="xMidYMid meet"
      @mousemove="updateMousePosition"
      class="map-svg"
    >
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
          :width="node.width"
          :height="node.height"
          class="node-image"
          @mouseover="showInfoBox(node)"
          @mouseleave="hideInfoBox"
        />
      </g>
    </svg>
    <div v-else class="loading-message">
      Loading map data...
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';

// --- Define reactive state ---
const nodes = ref([]);
const edges = ref([]);
const backgroundImage = ref('');
const hoveredNode = ref(null);
const mousePosition = ref({ x: 0, y: 0 });
const dynamicViewBox = ref('0 0 1 1');

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


// --- Fetch and process data ---
const fetchTerritoryData = async (id) => {
  const apiUrl = `/api/v1/client/territory/${id}`;
  try {
    const response = await fetch(apiUrl);
    if (!response.ok) throw new Error(`API responded with status: ${response.status}`);
    const rawData = await response.json();
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
  } catch (error) {
    console.error('Failed to load or process territory data:', error);
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

// --- Lifecycle Hook ---
onMounted(() => {
  fetchTerritoryData(props.territoryId);
});
</script>

<style scoped>
/* Define all colors as CSS variables for easy management */
.territory-container {
  --color-edge: #ffffff;
  --color-info-box-bg: rgba(0, 0, 0, 0.8);
  --color-info-box-text: white;
  --color-loading-text: #555;

  width: 100%;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 1rem;
  box-sizing: border-box;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.map-svg {
  width: 100%;
  height: 100%;
  display: block;
  overflow: visible;
}

.edge-path {
  fill: none;
  stroke: var(--color-edge);
  stroke-width: 0.005;
  stroke-dasharray: 0.01, 0.005;
}

.node-image {
  cursor: pointer;
  transition: transform 0.2s ease-in-out;
  transform-box: fill-box;
}

.node-image:hover {
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
}
</style>