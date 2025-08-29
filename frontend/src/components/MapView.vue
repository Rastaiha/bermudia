<template>
    <svg ref="svgRef" v-if="islands.length > 0" :viewBox="dynamicViewBox" preserveAspectRatio="xMidYMid meet"
        class="w-full h-full block cursor-grab active:cursor-grabbing">
        <g class="edges">
            <path v-for="edge in wavyEdges" :key="`${edge.from}-${edge.to}`" :d="edge.pathD"
                class="fill-none stroke-white" stroke-width="0.005" stroke-dasharray="0.01, 0.005" />
        </g>

        <g class="islands">
            <g v-for="island in islands" :key="island.id" class="cursor-pointer group"
                @pointerdown.stop="handlePointerDown(island)" @pointerup.stop="handlePointerUp(island)">
                <ellipse :cx="island.x" :cy="island.y" :rx="island.width / 2" :ry="island.height / 2" fill="transparent" />
                <image :href="`/images/islands/${island.iconAsset}`"
                    :x="island.x - island.width / 2"
                    :y="island.y - island.height / 2"
                    :width="island.width" :height="island.height"
                    style="transform-box: fill-box"
                    class="transition-transform duration-200 ease-in-out origin-center group-hover:scale-105 pointer-events-none" />
            </g>
        </g>

        <g v-if="player && player.atTerritory == territoryId" class="ship-container" :transform="shipTransform"
            :style="{ transition: shipTransition }">
            <image href="/images/ships/ship1.png" :width="BOAT_WIDTH" :height="BOAT_HEIGHT"
                class="drop-shadow-lg animate-boat" />
        </g>
    </svg>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, watch, nextTick } from 'vue';
import panzoom from 'panzoom';

const props = defineProps({
    islands: { type: Array, required: true },
    edges: { type: Array, required: true },
    player: { type: Object },
    dynamicViewBox: { type: String, required: true },
    territoryId: { type: String, required: true },
});

const emit = defineEmits(['nodeClick', 'mapTransformed']);

const svgRef = ref(null);
let panzoomInstance = null;
let potentialClickNode = null;

// --- Ship animation state ---
const shipTransform = ref('');
const shipTransition = ref('none');
const previousIsland = ref(null);


const BOAT_WIDTH = 0.16;
const BOAT_HEIGHT = 0.22;

const getShipPosition = (island) => {
  const x = island.x; // center X
  const y = island.y - island.height / 2; // top edge
  return { x, y };
};

// --- Watch for player's island changes to animate the ship ---
watch(() => props.player?.atIsland, (newIsland) => {
    if (!newIsland) return;

    const startIsland = previousIsland.value;
    const endIsland = newIsland;
    if (!startIsland) {
        const { x, y } = getShipPosition(endIsland);
        shipTransition.value = 'none';
        shipTransform.value = `translate(${x} ${y})`;
    } else if (startIsland.id !== endIsland.id) {
        const { x: startX, y: startY } = getShipPosition(startIsland);
        const { x: endX, y: endY } = getShipPosition(endIsland);

        shipTransition.value = 'none';
        shipTransform.value = `translate(${startX} ${startY})`;

        nextTick(() => {
            shipTransition.value = 'transform 2s ease-in-out';
            shipTransform.value = `translate(${endX} ${endY})`;
        });
    }
    previousIsland.value = endIsland;
}, { deep: true, immediate: true });


const handlePointerDown = (island) => {
    potentialClickNode = island;
};

const handlePointerUp = (island) => {
    if (potentialClickNode && potentialClickNode.id === island.id) {
        emit('nodeClick', island);
    }
    potentialClickNode = null;
};

const getIslandById = (id) => props.islands.find((island) => island.id === id);

const wavyEdges = computed(() => {
    return props.edges.map((edge) => {
        const startNode = getIslandById(edge.from);
        const endNode = getIslandById(edge.to);

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
        friction: 1,
        smoothScroll: false
    });

    panzoomInstance.on('panstart', () => {
        potentialClickNode = null;
    });

    panzoomInstance.on('transform', () => {
        emit('mapTransformed');
    });
};

onMounted(() => {
    initializePanzoom();
});

onUnmounted(() => {
    if (panzoomInstance) panzoomInstance.dispose();
});

defineExpose({
    svgRef,
});
</script>

<style scoped>
@keyframes boat-bobbing {
  0%   { transform: translateY(0); }
  50%  { transform: translateY(-3%); }
  100% { transform: translateY(0); }
}

.animate-boat {
    animation: boat-bobbing 4s ease-in-out infinite;
    transform-origin: center;
}
</style>
