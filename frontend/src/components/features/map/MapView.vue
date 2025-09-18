<template>
    <div
        v-if="hoveredShip"
        class="absolute bg-[#000000A0] text-white p-4 rounded-full text-sm pointer-events-none z-2000"
        :style="{
            left: hoveredShip.x + 'px',
            top: hoveredShip.y - 30 + 'px',
            transform: 'translateX(-50%)',
        }"
    >
        {{ hoveredShip.name }}
    </div>
    <svg
        v-if="islands.length > 0"
        ref="svgRef"
        :viewBox="dynamicViewBox"
        preserveAspectRatio="xMidYMid meet"
        class="w-full h-full block cursor-grab active:cursor-grabbing"
    >
        <g class="edges">
            <path
                v-for="edge in wavyEdges"
                :key="`${edge.from}-${edge.to}`"
                :d="edge.pathD"
                class="fill-none stroke-white"
                stroke-width="0.005"
                stroke-dasharray="0.01, 0.005"
            />
        </g>

        <g class="islands">
            <g
                v-for="island in islands"
                :key="island.id"
                class="cursor-pointer group"
                @pointerdown.stop="handlePointerDown(island)"
                @pointerup.stop="handlePointerUp(island)"
            >
                <ellipse
                    :cx="island.x"
                    :cy="island.y"
                    :rx="island.width / 2"
                    :ry="island.height / 2"
                    fill="transparent"
                />
                <image
                    :href="island.iconAsset"
                    :x="island.x - island.width / 2"
                    :y="island.y - island.height / 2"
                    :width="island.width"
                    :height="island.height"
                    style="transform-box: fill-box"
                    class="transition-transform duration-200 ease-in-out origin-center group-hover:scale-105 pointer-events-none"
                />
            </g>
        </g>

        <g
            v-if="player && player.atTerritory == territoryId"
            class="ship-container"
            :transform="shipTransform"
            :style="{ transition: shipTransition }"
        >
            <image
                :href="shipSrc(username)"
                :width="BOAT_WIDTH"
                :height="BOAT_HEIGHT"
                class="animate-boat"
            />
        </g>

        <g
            v-for="(territory, user) in users"
            :key="user"
            class="user-ship"
            :transform="
                'translate(' +
                territory.position.x +
                ' ' +
                territory.position.y +
                ')'
            "
            :style="{ transition: shipTransition }"
            @mouseenter="updateHoveredShipPosition(user, $event)"
            @mouseleave="hoveredShip = null"
        >
            <image
                :href="shipSrc(user)"
                :width="BOAT_WIDTH * 0.6"
                :height="BOAT_HEIGHT * 0.6"
                class="cursor-pointer"
            />
        </g>
    </svg>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, watch, nextTick } from 'vue';
import panzoom from 'panzoom';
import { getPlayersLocation } from '@/services/api';

const props = defineProps({
    islands: { type: Array, required: true },
    edges: { type: Array, required: true },
    player: { type: Object },
    dynamicViewBox: { type: String, required: true },
    territoryId: { type: String, required: true },
    username: { type: String, required: true },
});

const emit = defineEmits(['nodeClick', 'mapTransformed']);

const users = ref({});
const hoveredShip = ref(null);
const svgRef = ref(null);
let panzoomInstance = null;
let potentialClickNode = null;
let lastTapTime = 0;

const shipTransform = ref('');
const shipTransition = ref('none');
const previousAtIsland = ref(null);

const BOAT_WIDTH = 0.08;
const BOAT_HEIGHT = 0.11;

const getShipPosition = atIsland => {
    const island = props.islands.find(island => island.id === atIsland);
    if (!island) return { x: 0, y: 0 };
    const x = island.x;
    const y = island.y - island.height / 2;
    return { x, y };
};

const getShipPositionRandom = atIsland => {
    const island = props.islands.find(island => island.id === atIsland);
    if (!island) return { x: 0, y: 0 };
    const theta = Math.random() * 360;
    const r = Math.max(island.width) * 1.2;
    const x = island.x - island.width / 4 + (Math.cos(theta) * r) / 2;
    const y = island.y - island.height / 4 + (Math.sin(theta) * r) / 2; //todo improve the randommizing function
    return { x, y };
};

const updateHoveredShipPosition = (user, event) => {
    const bbox = event.target.getBoundingClientRect();
    hoveredShip.value = {
        name: user,
        x: bbox.left + bbox.width / 2,
        y: bbox.top,
    };
};

watch(
    () => props.player?.atIsland,
    newAtIsland => {
        if (!newAtIsland) return;

        const startIsland = previousAtIsland.value;
        const endIsland = newAtIsland;
        if (!startIsland) {
            const { x, y } = getShipPosition(endIsland);
            shipTransition.value = 'none';
            shipTransform.value = `translate(${x} ${y})`;
        } else {
            const { x: startX, y: startY } = getShipPosition(startIsland);
            const { x: endX, y: endY } = getShipPosition(endIsland);

            shipTransition.value = 'none';
            shipTransform.value = `translate(${startX} ${startY})`;

            nextTick(() => {
                shipTransition.value = 'transform 2s ease-in-out';
                shipTransform.value = `translate(${endX} ${endY})`;
            });
        }
        previousAtIsland.value = endIsland;
    },
    { deep: true, immediate: true }
);

const handlePointerDown = island => {
    potentialClickNode = island;
};

const handlePointerUp = island => {
    if (potentialClickNode && potentialClickNode.id === island.id) {
        emit('nodeClick', island);
    }
    potentialClickNode = null;
};

const getIslandById = id => props.islands.find(island => island.id === id);

const wavyEdges = computed(() => {
    return props.edges.map(edge => {
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
        smoothScroll: false,
        zoomDoubleClickSpeed: 0,
        zoomOutDoubleClickSpeed: 0,
    });

    svgRef.value.addEventListener(
        'dblclick',
        e => {
            e.preventDefault();
            e.stopPropagation();
        },
        { capture: true }
    );

    svgRef.value.addEventListener(
        'touchend',
        e => {
            const currentTime = new Date().getTime();
            const tapLength = currentTime - lastTapTime;
            if (tapLength < 500 && tapLength > 0) {
                e.preventDefault();
                e.stopPropagation();
            }
            lastTapTime = currentTime;
        },
        { passive: false }
    );

    panzoomInstance.on('panstart', () => {
        potentialClickNode = null;
    });

    panzoomInstance.on('transform', () => {
        const transform = panzoomInstance.getTransform();
        emit('mapTransformed', {
            x: transform.x,
            y: transform.y,
            scale: transform.scale,
        });
    });
};

const zoomToPlayer = () => {
    if (!props.player || !props.player.atIsland || !panzoomInstance) return;
    const svg = svgRef.value;
    const ctm = svg.getScreenCTM();
    if (!ctm) return;

    const { x: topLeftX, y: topLeftY } = getShipPosition(props.player.atIsland);
    const shipCenterX_svg = topLeftX + BOAT_WIDTH / 2;
    const shipCenterY_svg = topLeftY + BOAT_HEIGHT / 2;

    const point = svg.createSVGPoint();
    point.x = shipCenterX_svg;
    point.y = shipCenterY_svg;

    const shipCenterOnScreen = point.matrixTransform(ctm);

    const clientW = svg.clientWidth;
    const clientH = svg.clientHeight;
    const screenCenterX = clientW / 2;
    const screenCenterY = clientH / 2;

    const currentTransform = panzoomInstance.getTransform();
    const currentPanX = currentTransform.x;
    const currentPanY = currentTransform.y;

    const deltaX = screenCenterX - shipCenterOnScreen.x;
    const deltaY = screenCenterY - shipCenterOnScreen.y;

    const finalPanX = currentPanX + deltaX;
    const finalPanY = currentPanY + deltaY;

    const targetScale = 2;

    panzoomInstance.smoothZoom(screenCenterX, screenCenterY, targetScale);
    panzoomInstance.smoothMoveTo(finalPanX, finalPanY);
};

const fetchOtherPlayers = async () => {
    const otherPlayers = await getPlayersLocation(props.territoryId);
    const result = {};

    for (let i = 0; i < otherPlayers.length; i++) {
        if (i == 10) break;
        let island = otherPlayers[i];
        island.players.forEach(user => {
            result[user.name] = {
                island: island.islandId,
                position: getShipPositionRandom(island.islandId),
                delay: Math.random() * -2,
            };
        });
    }
    users.value = result;
};

const shipSrc = name => {
    let sum = 0;
    for (let i = 0; i < name.length; i++) {
        sum += name.charCodeAt(i);
    }
    return '/images/ships/' + ((sum % 11) + 1) + '.png';
};

onMounted(() => {
    initializePanzoom();
    fetchOtherPlayers();
});

onUnmounted(() => {
    if (panzoomInstance) panzoomInstance.dispose();
});

defineExpose({
    svgRef,
    zoomToPlayer,
});
</script>

<style scoped>
@keyframes boat-bobbing {
    0% {
        transform: translateY(0);
    }
    50% {
        transform: translateY(-3%);
    }
    100% {
        transform: translateY(0);
    }
}

.animate-boat {
    animation: boat-bobbing 4s ease-in-out infinite;
    transform-origin: center;
}
</style>
