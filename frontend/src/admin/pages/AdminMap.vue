<script setup>
import { ref, computed, onMounted } from 'vue';
import { useToast } from 'vue-toastification';
import { getTerritories, setTerritory } from '../services/api.js';
import AssetPicker from '../components/AssetPicker.vue';
import AdminModal from '../components/AdminModal.vue';
import IslandContentEditor from '../components/IslandContentEditor.vue';
import '../styles/admin.css';

const toast = useToast();

const territories = ref([]);
const selectedId = ref('');
const loading = ref(true);
const saving = ref(false);

// Working copy of the currently edited territory (deep-cloned from server),
// plus a snapshot of the last-saved state for dirty tracking.
const draft = ref(null);
const savedSnapshot = ref('');
const selectedIslandId = ref(null);

// Asset picker state
const picker = ref(null); // { kind, target } | null

const svgRef = ref(null);

const clone = obj => JSON.parse(JSON.stringify(obj));

// ---- Dirty tracking ----
const isDirty = computed(
    () => draft.value && JSON.stringify(draft.value) !== savedSnapshot.value
);

const markSaved = () => {
    savedSnapshot.value = JSON.stringify(draft.value);
};

// ---- Loading ----
const loadTerritories = async () => {
    loading.value = true;
    try {
        territories.value = (await getTerritories()) || [];
        if (territories.value.length && !selectedId.value) {
            selectId(territories.value[0].id);
        }
    } catch (e) {
        toast.error(`Failed to load territories: ${e.message}`);
    } finally {
        loading.value = false;
    }
};

// Ensure optional arrays/maps exist so the editor can mutate them safely.
const normalizeDraft = t => {
    t.islands ||= [];
    t.edges ||= [];
    t.refuelIslands ||= [];
    t.terminalIslands ||= [];
    t.islandPrerequisites ||= {};
    return t;
};

const selectId = id => {
    const t = territories.value.find(x => x.id === id);
    if (!t) return;
    if (isDirty.value && !confirmDiscard()) return;
    selectedId.value = id;
    draft.value = normalizeDraft(clone(t));
    markSaved();
    selectedIslandId.value = null;
};

const confirmDiscard = () =>
    window.confirm('You have unsaved changes. Discard them?');

// ---- Derived ----
const selectedIsland = computed(() =>
    draft.value?.islands.find(i => i.id === selectedIslandId.value)
);

const islandCenter = id => {
    const i = draft.value.islands.find(x => x.id === id);
    return i ? { x: i.x, y: i.y } : null;
};

const isRefuel = id => draft.value.refuelIslands.some(r => r.id === id);
const isTerminal = id => draft.value.terminalIslands.some(t => t.id === id);

// Display helper: keep coordinate inputs readable (3 decimals) while editing.
const round3 = v => Math.round((Number(v) || 0) * 1000) / 1000;

// ---- Island selection / drag (normalized 0..1 coordinates) ----
let dragState = null;

// Map a client point to normalized [0,1] coords. The SVG uses
// viewBox "0 0 1 1" with preserveAspectRatio "xMidYMid meet" (to mirror the
// game's crop-to-fill "stick" background), so the unit square is drawn
// centered as the largest square that fits the element. We undo that
// letterboxing to recover the normalized coordinate.
const clientToNorm = e => {
    const rect = svgRef.value.getBoundingClientRect();
    const side = Math.min(rect.width, rect.height);
    const offsetX = (rect.width - side) / 2;
    const offsetY = (rect.height - side) / 2;
    return {
        x: (e.clientX - rect.left - offsetX) / side,
        y: (e.clientY - rect.top - offsetY) / side,
    };
};

const onIslandPointerDown = (island, e) => {
    selectedIslandId.value = island.id;
    if (edgeMode.value) return;
    e.target.setPointerCapture?.(e.pointerId);
    const start = clientToNorm(e);
    dragState = {
        id: island.id,
        dx: island.x - start.x,
        dy: island.y - start.y,
        moved: false,
    };
};

const onPointerMove = e => {
    if (!dragState) return;
    const p = clientToNorm(e);
    const island = draft.value.islands.find(i => i.id === dragState.id);
    if (!island) return;
    dragState.moved = true;
    island.x = round3(Math.min(1, Math.max(0, p.x + dragState.dx)));
    island.y = round3(Math.min(1, Math.max(0, p.y + dragState.dy)));
};

const onPointerUp = () => {
    dragState = null;
};

// ---- Island CRUD ----
let newIslandCounter = 0;
const addIsland = () => {
    const id = `island_new_${Date.now()}_${newIslandCounter++}`;
    const island = {
        id,
        name: 'New Island',
        x: 0.5,
        y: 0.5,
        width: 0.08,
        height: 0.11,
        iconAsset: '/images/islands/educational/1.png',
    };
    draft.value.islands.push(island);
    selectedIslandId.value = id;
    if (!draft.value.startIsland) draft.value.startIsland = id;
};

const removeSelectedIsland = () => {
    const id = selectedIslandId.value;
    if (!id) return;
    draft.value.islands = draft.value.islands.filter(i => i.id !== id);
    draft.value.edges = draft.value.edges.filter(
        e => e.from !== id && e.to !== id
    );
    draft.value.refuelIslands = draft.value.refuelIslands.filter(
        r => r.id !== id
    );
    draft.value.terminalIslands = draft.value.terminalIslands.filter(
        t => t.id !== id
    );
    delete draft.value.islandPrerequisites[id];
    for (const k of Object.keys(draft.value.islandPrerequisites)) {
        draft.value.islandPrerequisites[k] = draft.value.islandPrerequisites[
            k
        ].filter(p => p !== id);
        if (draft.value.islandPrerequisites[k].length === 0)
            delete draft.value.islandPrerequisites[k];
    }
    if (draft.value.startIsland === id) draft.value.startIsland = '';
    selectedIslandId.value = null;
};

const toggleRefuel = () => {
    const id = selectedIslandId.value;
    if (!id) return;
    if (isRefuel(id)) {
        draft.value.refuelIslands = draft.value.refuelIslands.filter(
            r => r.id !== id
        );
    } else {
        draft.value.refuelIslands.push({ id });
    }
};

const toggleTerminal = () => {
    const id = selectedIslandId.value;
    if (!id) return;
    if (isTerminal(id)) {
        draft.value.terminalIslands = draft.value.terminalIslands.filter(
            t => t.id !== id
        );
    } else {
        draft.value.terminalIslands.push({ id });
    }
};

const setAsStart = () => {
    if (selectedIslandId.value)
        draft.value.startIsland = selectedIslandId.value;
};

// ---- Edges ----
const edgeMode = ref(false);
const edgeFrom = ref(null);
const hoveredEdge = ref(null);

const toggleEdgeMode = () => {
    edgeMode.value = !edgeMode.value;
    edgeFrom.value = edgeMode.value ? selectedIslandId.value : null;
};

const onIslandClickForEdge = island => {
    if (!edgeMode.value) return;
    if (!edgeFrom.value) {
        edgeFrom.value = island.id;
        return;
    }
    if (edgeFrom.value === island.id) return;
    const a = edgeFrom.value;
    const b = island.id;
    const exists = draft.value.edges.some(
        e => (e.from === a && e.to === b) || (e.from === b && e.to === a)
    );
    if (!exists) draft.value.edges.push({ from: a, to: b });
    edgeFrom.value = island.id; // chain edges
};

const removeEdge = edge => {
    draft.value.edges = draft.value.edges.filter(
        e => !(e.from === edge.from && e.to === edge.to)
    );
    hoveredEdge.value = null;
};

const edgeKey = edge => `${edge.from}-${edge.to}`;

// ---- Prerequisites ----
const prereqsOf = id => draft.value.islandPrerequisites[id] || [];
const togglePrereq = (id, prereqId) => {
    const list = draft.value.islandPrerequisites[id] || [];
    if (list.includes(prereqId)) {
        draft.value.islandPrerequisites[id] = list.filter(p => p !== prereqId);
        if (draft.value.islandPrerequisites[id].length === 0)
            delete draft.value.islandPrerequisites[id];
    } else {
        draft.value.islandPrerequisites[id] = [...list, prereqId];
    }
};

// ---- Asset picking ----
const openIconPicker = () => {
    if (!selectedIsland.value) return;
    picker.value = { kind: 'icon', target: 'island' };
};
const openBackgroundPicker = () => {
    picker.value = { kind: 'background', target: 'territory' };
};
const onAssetSelect = path => {
    if (picker.value?.target === 'island' && selectedIsland.value) {
        selectedIsland.value.iconAsset = path;
    } else if (picker.value?.target === 'territory') {
        draft.value.backgroundAsset = path;
    }
    picker.value = null;
};

// ---- Island content editor ----
const contentEditorFor = ref(null);
// True when the island exists in the last-saved snapshot (so the backend
// knows about it and getIslandHeader will succeed).
const islandExistsOnServer = id => {
    if (!savedSnapshot.value) return false;
    try {
        return JSON.parse(savedSnapshot.value).islands.some(i => i.id === id);
    } catch {
        return false;
    }
};
const openContentEditor = () => {
    const island = selectedIsland.value;
    if (!island) return;
    if (!islandExistsOnServer(island.id)) {
        toast.info('Save the map first, then edit this island’s content.');
        return;
    }
    contentEditorFor.value = { id: island.id, name: island.name };
};

// ---- Save / discard ----
const discardChanges = () => {
    if (!isDirty.value) return;
    if (!confirmDiscard()) return;
    const t = territories.value.find(x => x.id === selectedId.value);
    draft.value = normalizeDraft(clone(t));
    markSaved();
    selectedIslandId.value = null;
    toast.info('Reverted to last saved version');
};

const save = async () => {
    if (saving.value || !draft.value) return;
    saving.value = true;
    try {
        await setTerritory(draft.value);
        const idx = territories.value.findIndex(t => t.id === draft.value.id);
        if (idx >= 0) territories.value[idx] = clone(draft.value);
        markSaved();
        toast.success('Territory saved');
    } catch (e) {
        toast.error(`Save failed: ${e.message}`);
    } finally {
        saving.value = false;
    }
};

// edge path helper (straight line between centers, normalized)
const edgePath = edge => {
    const a = islandCenter(edge.from);
    const b = islandCenter(edge.to);
    if (!a || !b) return '';
    return `M ${a.x} ${a.y} L ${b.x} ${b.y}`;
};

onMounted(loadTerritories);
</script>

<template>
    <div class="admin-scope" dir="ltr">
        <header class="page-head">
            <div>
                <h1>Map Editor</h1>
                <p>
                    Arrange islands, connections, and roles for each territory.
                </p>
            </div>
            <div class="head-actions">
                <select
                    v-if="territories.length"
                    class="terr-select"
                    :value="selectedId"
                    @change="selectId($event.target.value)"
                >
                    <option v-for="t in territories" :key="t.id" :value="t.id">
                        {{ t.name || t.id }}
                    </option>
                </select>
                <template v-if="isDirty">
                    <span class="dirty-dot" title="Unsaved changes"></span>
                    <button class="btn btn-ghost" @click="discardChanges">
                        Discard
                    </button>
                    <button class="btn" :disabled="saving" @click="save">
                        {{ saving ? 'Saving…' : '💾 Save changes' }}
                    </button>
                </template>
                <span v-else class="saved-note">All changes saved</span>
            </div>
        </header>

        <p v-if="loading" class="hint">Loading…</p>

        <div v-else-if="draft" class="editor">
            <!-- Canvas -->
            <div class="canvas-wrap card">
                <div class="canvas-toolbar">
                    <button class="chip" @click="addIsland">
                        ➕ Add island
                    </button>
                    <button
                        class="chip"
                        :class="{ 'chip-on': edgeMode }"
                        @click="toggleEdgeMode"
                    >
                        {{
                            edgeMode
                                ? '✓ Drawing… (click two islands)'
                                : '↔ Draw edges'
                        }}
                    </button>
                    <button class="chip" @click="openBackgroundPicker">
                        🖼️ Background
                    </button>
                    <span class="canvas-hint">
                        Drag islands to reposition. Click an edge to delete it.
                    </span>
                </div>

                <div class="canvas-frame">
                    <svg
                        ref="svgRef"
                        viewBox="0 0 1 1"
                        preserveAspectRatio="xMidYMid meet"
                        class="map-svg"
                        @pointermove="onPointerMove"
                        @pointerup="onPointerUp"
                        @pointerleave="onPointerUp"
                    >
                        <image
                            v-if="draft.backgroundAsset"
                            :href="draft.backgroundAsset"
                            x="0"
                            y="0"
                            width="1"
                            height="1"
                            preserveAspectRatio="xMidYMid slice"
                            class="map-bg"
                        />
                        <rect
                            v-else
                            x="0"
                            y="0"
                            width="1"
                            height="1"
                            fill="#0b1220"
                        />

                        <!-- edges (with wide invisible hit-area for delete) -->
                        <g class="edges">
                            <g v-for="edge in draft.edges" :key="edgeKey(edge)">
                                <path
                                    :d="edgePath(edge)"
                                    class="edge"
                                    :class="{
                                        'edge-hover':
                                            hoveredEdge === edgeKey(edge),
                                    }"
                                />
                                <path
                                    :d="edgePath(edge)"
                                    class="edge-hit"
                                    @pointerenter="hoveredEdge = edgeKey(edge)"
                                    @pointerleave="hoveredEdge = null"
                                    @click.stop="removeEdge(edge)"
                                />
                            </g>
                        </g>

                        <!-- islands -->
                        <g class="islands">
                            <g
                                v-for="island in draft.islands"
                                :key="island.id"
                                @pointerdown.stop="
                                    onIslandPointerDown(island, $event)
                                "
                                @click.stop="onIslandClickForEdge(island)"
                            >
                                <ellipse
                                    :cx="island.x"
                                    :cy="island.y"
                                    :rx="island.width / 2 + 0.006"
                                    :ry="island.height / 2 + 0.006"
                                    fill="none"
                                    :class="[
                                        'ring',
                                        {
                                            'ring-selected':
                                                island.id === selectedIslandId,
                                            'ring-start':
                                                island.id === draft.startIsland,
                                        },
                                    ]"
                                />
                                <image
                                    :href="island.iconAsset"
                                    :x="island.x - island.width / 2"
                                    :y="island.y - island.height / 2"
                                    :width="island.width"
                                    :height="island.height"
                                    class="island-icon"
                                    :class="{ grabbable: !edgeMode }"
                                />
                                <text
                                    v-if="isRefuel(island.id)"
                                    :x="island.x - island.width / 4"
                                    :y="island.y - island.height / 2 - 0.006"
                                    class="badge-text"
                                >
                                    ⛽
                                </text>
                                <text
                                    v-if="isTerminal(island.id)"
                                    :x="island.x + island.width / 4"
                                    :y="island.y - island.height / 2 - 0.006"
                                    class="badge-text"
                                >
                                    🛣️
                                </text>
                            </g>
                        </g>
                    </svg>
                </div>
            </div>

            <!-- Side panel -->
            <aside class="side card">
                <div v-if="!selectedIsland" class="side-empty">
                    <p class="hint">
                        Select an island to edit it, or add a new one.
                    </p>
                    <div class="terr-meta">
                        <div>
                            <b>Territory:</b> {{ draft.name || draft.id }}
                        </div>
                        <div><b>Islands:</b> {{ draft.islands.length }}</div>
                        <div><b>Edges:</b> {{ draft.edges.length }}</div>
                        <div><b>Start:</b> {{ draft.startIsland || '—' }}</div>
                    </div>
                </div>

                <div v-else class="island-panel">
                    <div class="panel-title">
                        <span>Island</span>
                        <code>{{ selectedIsland.id }}</code>
                    </div>

                    <label class="field">
                        <span>Name</span>
                        <input v-model="selectedIsland.name" type="text" />
                    </label>

                    <div class="row2">
                        <label class="field">
                            <span>X</span>
                            <input
                                :value="selectedIsland.x"
                                type="number"
                                step="0.01"
                                min="0"
                                max="1"
                                @input="
                                    selectedIsland.x = round3(
                                        $event.target.value
                                    )
                                "
                            />
                        </label>
                        <label class="field">
                            <span>Y</span>
                            <input
                                :value="selectedIsland.y"
                                type="number"
                                step="0.01"
                                min="0"
                                max="1"
                                @input="
                                    selectedIsland.y = round3(
                                        $event.target.value
                                    )
                                "
                            />
                        </label>
                    </div>
                    <div class="row2">
                        <label class="field">
                            <span>Width</span>
                            <input
                                :value="selectedIsland.width"
                                type="number"
                                step="0.01"
                                min="0"
                                max="1"
                                @input="
                                    selectedIsland.width = round3(
                                        $event.target.value
                                    )
                                "
                            />
                        </label>
                        <label class="field">
                            <span>Height</span>
                            <input
                                :value="selectedIsland.height"
                                type="number"
                                step="0.01"
                                min="0"
                                max="1"
                                @input="
                                    selectedIsland.height = round3(
                                        $event.target.value
                                    )
                                "
                            />
                        </label>
                    </div>

                    <div class="icon-row">
                        <img
                            :src="selectedIsland.iconAsset"
                            class="icon-preview"
                            alt=""
                        />
                        <button class="btn btn-ghost" @click="openIconPicker">
                            Change icon
                        </button>
                    </div>

                    <button
                        class="btn btn-ghost content-btn"
                        @click="openContentEditor"
                    >
                        📝 Edit island content
                    </button>

                    <div class="role-buttons">
                        <button
                            class="chip"
                            :class="{
                                'chip-on':
                                    selectedIsland.id === draft.startIsland,
                            }"
                            @click="setAsStart"
                        >
                            ⭐ Start
                        </button>
                        <button
                            class="chip"
                            :class="{ 'chip-on': isRefuel(selectedIsland.id) }"
                            @click="toggleRefuel"
                        >
                            ⛽ Refuel
                        </button>
                        <button
                            class="chip"
                            :class="{
                                'chip-on': isTerminal(selectedIsland.id),
                            }"
                            @click="toggleTerminal"
                        >
                            🛣️ Terminal
                        </button>
                    </div>

                    <div class="prereq">
                        <span class="prereq-label">Prerequisites</span>
                        <p class="hint">
                            Islands that must be completed before this one
                            unlocks.
                        </p>
                        <div class="prereq-list">
                            <label
                                v-for="other in draft.islands.filter(
                                    i => i.id !== selectedIsland.id
                                )"
                                :key="other.id"
                                class="prereq-item"
                            >
                                <input
                                    type="checkbox"
                                    :checked="
                                        prereqsOf(selectedIsland.id).includes(
                                            other.id
                                        )
                                    "
                                    @change="
                                        togglePrereq(
                                            selectedIsland.id,
                                            other.id
                                        )
                                    "
                                />
                                <span>{{ other.name || other.id }}</span>
                            </label>
                        </div>
                    </div>

                    <button
                        class="btn btn-danger"
                        @click="removeSelectedIsland"
                    >
                        🗑 Delete island
                    </button>
                </div>
            </aside>
        </div>

        <p v-else class="hint">No territories found.</p>

        <AssetPicker
            v-if="picker"
            :kind="picker.kind"
            :current="
                picker.target === 'island'
                    ? selectedIsland?.iconAsset
                    : draft?.backgroundAsset
            "
            :title="
                picker.kind === 'background'
                    ? 'Choose a background'
                    : 'Choose an island icon'
            "
            @select="onAssetSelect"
            @close="picker = null"
        />

        <AdminModal
            v-if="contentEditorFor"
            wide
            :title="`Edit content — ${contentEditorFor.name || contentEditorFor.id}`"
            @close="contentEditorFor = null"
        >
            <IslandContentEditor :island-id="contentEditorFor.id" />
        </AdminModal>
    </div>
</template>

<style scoped>
.head-actions {
    display: flex;
    gap: 8px;
    align-items: center;
    flex-wrap: wrap;
}
.terr-select {
    background: #0b1220;
    border: 1px solid #334155;
    border-radius: 10px;
    padding: 9px 12px;
    color: #f8fafc;
    font-size: 14px;
    outline: none;
}
.dirty-dot {
    width: 9px;
    height: 9px;
    border-radius: 50%;
    background: #eab308;
    box-shadow: 0 0 0 3px #eab30833;
}
.saved-note {
    color: #64748b;
    font-size: 12.5px;
}
.btn-danger {
    background: #dc2626;
    width: 100%;
    margin-top: 12px;
}
.btn-danger:hover:not(:disabled) {
    background: #b91c1c;
}

.editor {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 320px;
    gap: 18px;
    align-items: start;
}
@media (max-width: 900px) {
    .editor {
        grid-template-columns: 1fr;
    }
}

.canvas-wrap {
    padding: 12px;
    min-width: 0;
}
.canvas-toolbar {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
    align-items: center;
    margin-bottom: 10px;
}
.canvas-hint {
    color: #64748b;
    font-size: 12px;
}
.canvas-frame {
    display: flex;
    justify-content: center;
    align-items: center;
    background: #0b1220;
    border-radius: 10px;
    overflow: hidden;
    max-height: 78vh;
}

.map-svg {
    width: 100%;
    aspect-ratio: 1;
    max-height: 78vh;
    display: block;
    touch-action: none;
    cursor: crosshair;
}
.map-bg {
    opacity: 0.92;
}

.edge {
    fill: none;
    stroke: #93c5fd;
    stroke-width: 0.004;
    stroke-dasharray: 0.012, 0.008;
    pointer-events: none;
}
.edge-hover {
    stroke: #f87171;
    stroke-width: 0.007;
}
.edge-hit {
    fill: none;
    stroke: transparent;
    stroke-width: 0.03;
    cursor: pointer;
}

.island-icon {
    -webkit-user-select: none;
    user-select: none;
}
.island-icon.grabbable {
    cursor: grab;
}
.ring {
    stroke-width: 0.004;
}
.ring-selected {
    stroke: #22c55e;
}
.ring-start {
    stroke: #eab308;
    stroke-dasharray: 0.01, 0.006;
}
.badge-text {
    font-size: 0.028px;
    text-anchor: middle;
    dominant-baseline: middle;
}

.side {
    position: sticky;
    top: 16px;
    min-width: 0;
}
.side-empty .terr-meta {
    margin-top: 14px;
    display: flex;
    flex-direction: column;
    gap: 6px;
    font-size: 13px;
    color: #cbd5e1;
}

.panel-title {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 14px;
    font-weight: 700;
    color: #f1f5f9;
}
.panel-title code {
    font-size: 11px;
    color: #64748b;
    max-width: 55%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
.island-panel .field {
    margin-bottom: 10px;
}
.row2 {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 10px;
}
/* keep number inputs from overflowing the panel */
.row2 .field input {
    width: 100%;
    min-width: 0;
}
.icon-row {
    display: flex;
    align-items: center;
    gap: 12px;
    margin: 12px 0;
}
.icon-preview {
    width: 44px;
    height: 44px;
    object-fit: contain;
    background: #0b1220;
    border: 1px solid #1f2937;
    border-radius: 10px;
    padding: 4px;
}
.content-btn {
    width: 100%;
    margin-bottom: 4px;
}

.role-buttons {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin: 12px 0;
}
.chip {
    background: #0b1220;
    border: 1px solid #334155;
    color: #cbd5e1;
    border-radius: 999px;
    padding: 6px 12px;
    font-size: 12.5px;
    cursor: pointer;
}
.chip:hover {
    background: #1f2937;
}
.chip-on {
    background: #2563eb;
    border-color: #2563eb;
    color: #fff;
}

.prereq {
    margin-top: 8px;
}
.prereq-label {
    font-size: 12px;
    font-weight: 600;
    color: #94a3b8;
    text-transform: uppercase;
    letter-spacing: 0.03em;
}
.prereq-list {
    margin-top: 8px;
    display: flex;
    flex-direction: column;
    gap: 4px;
    max-height: 160px;
    overflow-y: auto;
}
.prereq-item {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    color: #e2e8f0;
}
</style>
