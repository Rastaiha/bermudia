<script setup>
import { ref, computed } from 'vue';
import AdminModal from './AdminModal.vue';
import {
    islandIconGroups,
    islandTypes,
    territoryBackgrounds,
} from '../services/assets.js';

const props = defineProps({
    // 'icon' -> island icons grouped by type; 'background' -> territory bgs
    kind: { type: String, default: 'icon' },
    current: { type: String, default: '' },
    title: { type: String, default: 'Choose an image' },
});

const emit = defineEmits(['select', 'close']);

const activeType = ref(islandTypes[0] || '');

const groups = computed(() => {
    if (props.kind === 'background') {
        return { territory: territoryBackgrounds };
    }
    return islandIconGroups;
});

const tabs = computed(() => (props.kind === 'background' ? [] : islandTypes));

const shownAssets = computed(() => {
    if (props.kind === 'background') return territoryBackgrounds;
    return groups.value[activeType.value] || [];
});

const pick = path => emit('select', path);
</script>

<template>
    <AdminModal :title="title" class="picker-modal" @close="emit('close')">
        <div class="picker">
            <div v-if="tabs.length" class="picker-tabs">
                <button
                    v-for="t in tabs"
                    :key="t"
                    class="picker-tab"
                    :class="{ 'is-active': t === activeType }"
                    @click="activeType = t"
                >
                    {{ t }}
                </button>
            </div>

            <div class="picker-grid">
                <button
                    v-for="asset in shownAssets"
                    :key="asset"
                    class="picker-cell"
                    :class="{ 'is-current': asset === current }"
                    :title="asset"
                    @click="pick(asset)"
                >
                    <img :src="asset" alt="" loading="lazy" />
                </button>
                <p v-if="shownAssets.length === 0" class="picker-empty">
                    No images found in this category.
                </p>
            </div>
        </div>
    </AdminModal>
</template>

<style scoped>
.picker {
    width: min(70vw, 640px);
    max-width: 100%;
}
.picker-tabs {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin-bottom: 14px;
}
.picker-tab {
    background: #0b1220;
    border: 1px solid #334155;
    color: #cbd5e1;
    border-radius: 999px;
    padding: 5px 12px;
    font-size: 12.5px;
    cursor: pointer;
    text-transform: capitalize;
}
.picker-tab:hover {
    background: #1f2937;
}
.picker-tab.is-active {
    background: #2563eb;
    border-color: #2563eb;
    color: #fff;
}

.picker-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(72px, 1fr));
    gap: 10px;
    max-height: 52vh;
    overflow-y: auto;
    padding: 4px;
}
.picker-cell {
    aspect-ratio: 1;
    background: #0b1220;
    border: 2px solid #1f2937;
    border-radius: 12px;
    padding: 6px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition:
        border-color 0.12s,
        transform 0.12s;
}
.picker-cell:hover {
    border-color: #2563eb;
    transform: translateY(-2px);
}
.picker-cell.is-current {
    border-color: #22c55e;
}
.picker-cell img {
    max-width: 100%;
    max-height: 100%;
    object-fit: contain;
}
.picker-empty {
    grid-column: 1 / -1;
    color: #94a3b8;
    font-size: 13px;
    text-align: center;
    padding: 24px 0;
}
</style>
