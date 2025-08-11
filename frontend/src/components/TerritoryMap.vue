<template>
  <div class="territory-container">
    <div
      v-if="hoveredNode"
      class="info-box"
      :style="infoBoxStyle"
    >
      {{ hoveredNode.name }}
    </div>

    <svg
      v-if="nodes.length > 0"
      :viewBox="`0 0 ${width} ${height}`"
      preserveAspectRatio="xMidYMid meet"
      @mousemove="updateMousePosition"
      class="map-svg"
    >
      <g class="edges">
        <line
          v-for="edge in edges"
          :key="`${edge.from_node_id}-${edge.to_node_id}`"
          :x1="getNodeById(edge.from_node_id)?.position.x"
          :y1="getNodeById(edge.from_node_id)?.position.y"
          :x2="getNodeById(edge.to_node_id)?.position.x"
          :y2="getNodeById(edge.to_node_id)?.position.y"
          stroke="#FFF"
          stroke-width="2"
        />
      </g>

      <g class="nodes">
        <circle
          v-for="node in nodes"
          :key="node.id"
          :cx="node.position.x"
          :cy="node.position.y"
          r="25"
          fill="#FFD700"
          stroke="#FFF"
          stroke-width="3"
          class="node-circle"
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

// Define map dimensions
const width = 1500;
const height = 1500;

// Define reactive state
const nodes = ref([]);
const edges = ref([]);
const hoveredNode = ref(null);
const mousePosition = ref({ x: 0, y: 0 });

// Define component props
const props = defineProps({
  territoryId: {
    type: [Number, String],
    required: true,
  },
});

// Fetch and process territory data from the mock API
const fetchTerritoryData = async (id) => {
  const apiUrl = `/api/v1/client/territory/${id}`;
  console.log(`Initiating fetch request to API endpoint: ${apiUrl}`);
  try {
    const response = await fetch(apiUrl);
    if (!response.ok) {
      throw new Error(`API responded with status: ${response.status}`);
    }
    const rawData = await response.json();
    console.log('Successfully fetched raw data. Now transforming data structure...');

    // Transform raw data to match component's expected structure
    const transformedNodes = rawData.islands.map(island => ({
      id: island.id,
      name: island.name,
      position: {
        x: island.x * width,
        y: island.y * height,
      }
    }));

    const transformedEdges = rawData.edges.map(edge => ({
      from_node_id: edge.from,
      to_node_id: edge.to,
    }));

    nodes.value = transformedNodes;
    edges.value = transformedEdges;
    console.log('Data transformation complete. Component is ready to render.');

  } catch (error) {
    console.error('Failed to load or process territory data:', error);
  }
};

// Find a node by its ID
const getNodeById = (id) => {
  return nodes.value.find(node => node.id === id);
};

// Update mouse position for the info box
const updateMousePosition = (event) => {
  mousePosition.value = { x: event.clientX, y: event.clientY };
};

// Handle hover events
const showInfoBox = (node) => {
  hoveredNode.value = node;
};
const hideInfoBox = () => {
  hoveredNode.value = null;
};

// Compute the style for the info box
const infoBoxStyle = computed(() => ({
  position: 'fixed',
  top: `${mousePosition.value.y + 20}px`,
  left: `${mousePosition.value.x}px`,
  transform: 'translateX(-50%)',
}));

// Fetch data when the component is mounted
onMounted(() => {
  fetchTerritoryData(props.territoryId);
});
</script>

<style scoped>
.territory-container {
  width: 100%;
  height: 100vh;
  overflow: auto;
  background-color: #f0f0f0;
  display: flex;
  justify-content: center;
  align-items: center;
}

.map-svg {
  width: 1500px;
  height: 1500px;
  display: block;
  background-color: #a7d1e8;
}

.node-circle {
  cursor: pointer;
  transition: transform 0.2s ease-in-out;
  transform-box: fill-box; /* Fix hover transformation origin */
}

.node-circle:hover {
  transform-origin: center center;
  transform: scale(1.2);
}

.info-box {
  background-color: rgba(0, 0, 0, 0.8);
  color: white;
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
  color: #555;
  font-family: sans-serif;
}
</style>