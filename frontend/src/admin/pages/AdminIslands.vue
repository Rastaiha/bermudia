<script setup>
import { ref, computed, onMounted } from 'vue';
import { useToast } from 'vue-toastification';
import { getTerritories } from '../services/api.js';
import AdminModal from '../components/AdminModal.vue';
import IslandContentEditor from '../components/IslandContentEditor.vue';
import '../styles/admin.css';

const toast = useToast();

const territories = ref([]);
const loading = ref(true);
const search = ref('');
const territoryFilter = ref('');

const editing = ref(null); // { id, name, territory } | null

// Flatten islands across all territories into table rows.
const rows = computed(() => {
    const out = [];
    for (const t of territories.value) {
        for (const isl of t.islands || []) {
            out.push({
                id: isl.id,
                name: isl.name || '',
                territoryId: t.id,
                territoryName: t.name || t.id,
            });
        }
    }
    return out;
});

const filteredRows = computed(() => {
    const q = search.value.trim().toLowerCase();
    return rows.value.filter(r => {
        if (territoryFilter.value && r.territoryId !== territoryFilter.value)
            return false;
        if (!q) return true;
        return (
            r.id.toLowerCase().includes(q) || r.name.toLowerCase().includes(q)
        );
    });
});

const load = async () => {
    loading.value = true;
    try {
        territories.value = (await getTerritories()) || [];
    } catch (e) {
        toast.error(`Failed to load islands: ${e.message}`);
    } finally {
        loading.value = false;
    }
};

const openEditor = row => {
    editing.value = row;
};
const closeEditor = () => {
    editing.value = null;
};

onMounted(load);
</script>

<template>
    <div class="admin-scope" dir="ltr">
        <header class="page-head">
            <div>
                <h1>Island Content</h1>
                <p>
                    {{ rows.length }} islands across
                    {{ territories.length }} territories.
                </p>
            </div>
        </header>

        <div class="toolbar">
            <input
                v-model="search"
                class="search"
                type="text"
                placeholder="Search by island id or name…"
            />
            <select v-model="territoryFilter" class="terr-filter">
                <option value="">All territories</option>
                <option v-for="t in territories" :key="t.id" :value="t.id">
                    {{ t.name || t.id }}
                </option>
            </select>
        </div>

        <div class="table-wrap">
            <table>
                <thead>
                    <tr>
                        <th>Island</th>
                        <th>ID</th>
                        <th>Territory</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="loading">
                        <td colspan="4" class="empty">Loading…</td>
                    </tr>
                    <tr v-else-if="filteredRows.length === 0">
                        <td colspan="4" class="empty">No islands found.</td>
                    </tr>
                    <tr v-for="r in filteredRows" :key="r.id">
                        <td>{{ r.name || '—' }}</td>
                        <td>
                            <span class="mono">{{ r.id }}</span>
                        </td>
                        <td>{{ r.territoryName }}</td>
                        <td class="right">
                            <button class="btn btn-sm" @click="openEditor(r)">
                                Edit content
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <AdminModal
            v-if="editing"
            wide
            :title="`Edit content — ${editing.name || editing.id}`"
            @close="closeEditor"
        >
            <IslandContentEditor :island-id="editing.id" />
        </AdminModal>
    </div>
</template>

<style scoped>
.toolbar {
    display: flex;
    gap: 10px;
    margin-bottom: 14px;
    flex-wrap: wrap;
}
.search,
.terr-filter {
    background: #0b1220;
    border: 1px solid #334155;
    border-radius: 10px;
    padding: 9px 12px;
    color: #f8fafc;
    font-size: 14px;
    outline: none;
}
.search {
    flex: 1;
    min-width: 200px;
}
.mono {
    font-family: ui-monospace, monospace;
    font-size: 12.5px;
    color: #cbd5e1;
}
.empty {
    text-align: center;
    color: #64748b;
    padding: 24px 0;
}
.right {
    text-align: right;
}
.btn-sm {
    padding: 6px 12px;
    font-size: 13px;
}
</style>
