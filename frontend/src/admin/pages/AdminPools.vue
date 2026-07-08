<script setup>
import { ref, computed, onMounted } from 'vue';
import { useToast } from 'vue-toastification';
import {
    getTerritories,
    getTerritoryIslandBindings,
    setTerritoryIslandBindings,
    getPools,
    getBook,
} from '../services/api.js';
import AdminModal from '../components/AdminModal.vue';
import IslandContentEditor from '../components/IslandContentEditor.vue';
import '../styles/admin.css';

const toast = useToast();

const tab = ref('bindings'); // 'bindings' | 'books'

// ---- Territory bindings ----
const territories = ref([]);
const selectedTerritory = ref('');
const bindings = ref(null); // { territoryId, islands:[{id,pooled}], poolSettings }
const bindingsLoading = ref(false);
const savingBindings = ref(false);

// island name lookup for nicer labels
const islandNames = ref({});

const loadTerritories = async () => {
    try {
        territories.value = (await getTerritories()) || [];
        const names = {};
        for (const t of territories.value)
            for (const i of t.islands || []) names[i.id] = i.name || i.id;
        islandNames.value = names;
        if (territories.value.length && !selectedTerritory.value) {
            selectTerritory(territories.value[0].id);
        }
    } catch (e) {
        toast.error(`Failed to load territories: ${e.message}`);
    }
};

const selectTerritory = async id => {
    selectedTerritory.value = id;
    bindingsLoading.value = true;
    bindings.value = null;
    try {
        const b = await getTerritoryIslandBindings(id);
        // Merge empty + pooled into one editable list with a `pooled` flag.
        const islands = [
            ...(b.emptyIslands || []).map(x => ({ id: x, pooled: false })),
            ...(b.pooledIslands || []).map(x => ({ id: x, pooled: true })),
        ];
        bindings.value = {
            territoryId: id,
            islands,
            poolSettings: {
                easy: b.poolSettings?.easy || 0,
                medium: b.poolSettings?.medium || 0,
                hard: b.poolSettings?.hard || 0,
            },
        };
    } catch (e) {
        toast.error(`Failed to load bindings: ${e.message}`);
    } finally {
        bindingsLoading.value = false;
    }
};

const pooledCount = computed(
    () => bindings.value?.islands.filter(i => i.pooled).length || 0
);
const settingsTotal = computed(() => {
    const s = bindings.value?.poolSettings;
    return s ? s.easy + s.medium + s.hard : 0;
});
const countsMatch = computed(() => pooledCount.value === settingsTotal.value);

const togglePooled = island => {
    island.pooled = !island.pooled;
};

const saveBindings = async () => {
    if (savingBindings.value || !bindings.value) return;
    if (!countsMatch.value) {
        toast.warning(
            `Pooled islands (${pooledCount.value}) must equal pool settings total (${settingsTotal.value}).`
        );
        return;
    }
    savingBindings.value = true;
    try {
        const payload = {
            territoryId: bindings.value.territoryId,
            emptyIslands: bindings.value.islands
                .filter(i => !i.pooled)
                .map(i => i.id),
            pooledIslands: bindings.value.islands
                .filter(i => i.pooled)
                .map(i => i.id),
            poolSettings: bindings.value.poolSettings,
        };
        await setTerritoryIslandBindings(payload);
        toast.success('Bindings saved');
    } catch (e) {
        toast.error(`Save failed: ${e.message}`);
    } finally {
        savingBindings.value = false;
    }
};

// ---- Pool books ----
const pools = ref([]);
const poolsLoading = ref(false);
const bookPreviews = ref({}); // bookId -> { components, treasures }
const editor = ref(null); // { poolId, bookId } | null

const loadPools = async () => {
    poolsLoading.value = true;
    try {
        pools.value = (await getPools()) || [];
        // fetch a light preview (component counts) for each book
        const previews = {};
        await Promise.all(
            pools.value.flatMap(p =>
                (p.books || []).map(async bid => {
                    try {
                        const b = await getBook(bid);
                        previews[bid] = {
                            components: (b.components || []).length,
                            treasures: (b.treasures || []).length,
                        };
                    } catch {
                        previews[bid] = null;
                    }
                })
            )
        );
        bookPreviews.value = previews;
    } catch (e) {
        toast.error(`Failed to load pools: ${e.message}`);
    } finally {
        poolsLoading.value = false;
    }
};

const openNewBook = poolId => {
    editor.value = { poolId, bookId: '' };
};
const openBook = (poolId, bookId) => {
    editor.value = { poolId, bookId };
};
const onBookSaved = () => {
    // refresh pool lists/previews so a newly created book shows up
    loadPools();
};

const switchTab = t => {
    tab.value = t;
    if (t === 'books' && !pools.value.length) loadPools();
};

onMounted(loadTerritories);
</script>

<template>
    <div class="admin-scope" dir="ltr">
        <header class="page-head">
            <div>
                <h1>Pools</h1>
                <p>
                    Manage which islands draw from shared question pools, and
                    the books inside each difficulty pool.
                </p>
            </div>
        </header>

        <div class="tabs">
            <button
                class="tab"
                :class="{ 'tab-on': tab === 'bindings' }"
                @click="switchTab('bindings')"
            >
                Territory bindings
            </button>
            <button
                class="tab"
                :class="{ 'tab-on': tab === 'books' }"
                @click="switchTab('books')"
            >
                Pool books
            </button>
        </div>

        <!-- ===== Bindings tab ===== -->
        <section v-if="tab === 'bindings'">
            <div class="binding-head">
                <select
                    v-if="territories.length"
                    class="terr-select"
                    :value="selectedTerritory"
                    @change="selectTerritory($event.target.value)"
                >
                    <option v-for="t in territories" :key="t.id" :value="t.id">
                        {{ t.name || t.id }}
                    </option>
                </select>
            </div>

            <p v-if="bindingsLoading" class="hint">Loading…</p>

            <div v-else-if="bindings" class="binding-grid">
                <!-- pool settings -->
                <div class="card settings-card">
                    <h2>Pool settings</h2>
                    <p class="hint">
                        How many pooled islands pull from each difficulty. The
                        total must equal the number of islands marked
                        <b>Pooled</b>.
                    </p>
                    <div class="settings-rows">
                        <label class="field">
                            <span>Easy</span>
                            <input
                                v-model.number="bindings.poolSettings.easy"
                                type="number"
                                min="0"
                            />
                        </label>
                        <label class="field">
                            <span>Medium</span>
                            <input
                                v-model.number="bindings.poolSettings.medium"
                                type="number"
                                min="0"
                            />
                        </label>
                        <label class="field">
                            <span>Hard</span>
                            <input
                                v-model.number="bindings.poolSettings.hard"
                                type="number"
                                min="0"
                            />
                        </label>
                    </div>
                    <div class="counts" :class="countsMatch ? 'ok' : 'bad'">
                        Pooled islands: {{ pooledCount }} / settings total:
                        {{ settingsTotal }}
                        <span v-if="countsMatch">✓</span>
                        <span v-else>— must match to save</span>
                    </div>
                    <button
                        class="btn"
                        :disabled="savingBindings || !countsMatch"
                        @click="saveBindings"
                    >
                        {{ savingBindings ? 'Saving…' : '💾 Save bindings' }}
                    </button>
                </div>

                <!-- island list -->
                <div class="card islands-card">
                    <h2>Islands</h2>
                    <p class="hint">
                        Only islands without a directly-authored book appear
                        here. Toggle each between normal (empty) and pooled.
                    </p>
                    <div
                        v-if="!bindings.islands.length"
                        class="hint empty-line"
                    >
                        No empty/pooled islands in this territory.
                    </div>
                    <div
                        v-for="isl in bindings.islands"
                        :key="isl.id"
                        class="island-row"
                    >
                        <div class="island-id">
                            <span>{{ islandNames[isl.id] || isl.id }}</span>
                            <code>{{ isl.id }}</code>
                        </div>
                        <button
                            class="toggle"
                            :class="isl.pooled ? 'pooled' : 'normal'"
                            @click="togglePooled(isl)"
                        >
                            {{ isl.pooled ? '🎲 Pooled' : '○ Normal' }}
                        </button>
                    </div>
                </div>
            </div>
        </section>

        <!-- ===== Pool books tab ===== -->
        <section v-else>
            <p v-if="poolsLoading" class="hint">Loading…</p>
            <div v-else class="pools-grid">
                <div v-for="p in pools" :key="p.id" class="card pool-card">
                    <div class="pool-card-head">
                        <h2 class="pool-name">{{ p.id }}</h2>
                        <span class="pool-count"
                            >{{ (p.books || []).length }} books</span
                        >
                    </div>
                    <div
                        v-for="bid in p.books"
                        :key="bid"
                        class="book-row"
                        @click="openBook(p.id, bid)"
                    >
                        <code class="book-id">{{ bid }}</code>
                        <span v-if="bookPreviews[bid]" class="book-meta">
                            {{ bookPreviews[bid].components }} items ·
                            {{ bookPreviews[bid].treasures }} treasures
                        </span>
                    </div>
                    <div v-if="!(p.books || []).length" class="hint empty-line">
                        No books yet.
                    </div>
                    <button
                        class="btn btn-ghost add-book"
                        @click="openNewBook(p.id)"
                    >
                        + Add book
                    </button>
                </div>
            </div>
        </section>

        <!-- Pool book editor -->
        <AdminModal
            v-if="editor"
            wide
            :title="`${editor.bookId ? 'Edit' : 'New'} ${editor.poolId} pool book`"
            @close="editor = null"
        >
            <IslandContentEditor
                :pool-id="editor.poolId"
                :initial-book-id="editor.bookId"
                @saved="onBookSaved"
            />
        </AdminModal>
    </div>
</template>

<style scoped>
.tabs {
    display: flex;
    gap: 6px;
    margin-bottom: 18px;
    border-bottom: 1px solid #1f2937;
}
.tab {
    background: transparent;
    border: none;
    border-bottom: 2px solid transparent;
    color: #94a3b8;
    padding: 10px 14px;
    font-size: 14px;
    cursor: pointer;
}
.tab:hover {
    color: #e2e8f0;
}
.tab-on {
    color: #fff;
    border-bottom-color: #2563eb;
}

.binding-head {
    margin-bottom: 14px;
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

.binding-grid {
    display: grid;
    grid-template-columns: 300px 1fr;
    gap: 18px;
    align-items: start;
}
@media (max-width: 820px) {
    .binding-grid {
        grid-template-columns: 1fr;
    }
}
.card h2 {
    font-size: 15px;
    font-weight: 600;
    color: #f1f5f9;
    margin-bottom: 8px;
}
.settings-rows {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    gap: 8px;
    margin: 12px 0;
}
.settings-rows .field input {
    min-width: 0;
}
.counts {
    font-size: 13px;
    padding: 8px 10px;
    border-radius: 8px;
    margin-bottom: 12px;
}
.counts.ok {
    background: #064e3b55;
    color: #6ee7b7;
}
.counts.bad {
    background: #7f1d1d55;
    color: #fca5a5;
}

.island-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
    padding: 8px 0;
    border-bottom: 1px solid #1f2937;
}
.island-row:last-child {
    border-bottom: none;
}
.island-id {
    display: flex;
    flex-direction: column;
    min-width: 0;
}
.island-id span {
    color: #e2e8f0;
    font-size: 14px;
}
.island-id code {
    color: #64748b;
    font-size: 11px;
}
.toggle {
    border: 1px solid #334155;
    background: #0b1220;
    color: #cbd5e1;
    border-radius: 999px;
    padding: 6px 14px;
    font-size: 13px;
    cursor: pointer;
    flex-shrink: 0;
}
.toggle.pooled {
    background: #2563eb;
    border-color: #2563eb;
    color: #fff;
}

.pools-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
    gap: 18px;
}
.pool-card-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 12px;
}
.pool-name {
    text-transform: capitalize;
    margin: 0;
}
.pool-count {
    font-size: 12px;
    color: #94a3b8;
}
.book-row {
    display: flex;
    flex-direction: column;
    gap: 2px;
    padding: 10px;
    border: 1px solid #1f2937;
    border-radius: 10px;
    background: #0b1220;
    margin-bottom: 8px;
    cursor: pointer;
}
.book-row:hover {
    border-color: #2563eb;
}
.book-id {
    font-size: 12px;
    color: #cbd5e1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
.book-meta {
    font-size: 11.5px;
    color: #64748b;
}
.add-book {
    width: 100%;
    margin-top: 4px;
}
.empty-line {
    padding: 6px 0;
}
</style>
