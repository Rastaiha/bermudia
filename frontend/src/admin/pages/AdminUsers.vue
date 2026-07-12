<script setup>
import { ref, onMounted, computed } from 'vue';
import { useToast } from 'vue-toastification';
import {
    getUsers,
    createUser,
    getTerritories,
    getPlayerState,
    editPlayerState,
} from '../services/api.js';
import AdminModal from '../components/AdminModal.vue';
import '../styles/admin.css';

const toast = useToast();

const users = ref([]);
const loading = ref(true);
const territories = ref([]);
const search = ref('');

const showCreate = ref(false);
const creating = ref(false);
const form = ref(newForm());
const createdUser = ref(null);

function newForm() {
    return {
        name: '',
        username: '',
        password: '',
        startingTerritory: '',
        meetLink: '',
    };
}

// --- Manage player state ---
const NUMERIC_FIELDS = [
    { key: 'coin', label: 'Coins' },
    { key: 'fuel', label: 'Fuel' },
    { key: 'fuelCap', label: 'Fuel capacity' },
    { key: 'blueKey', label: 'Blue keys' },
    { key: 'redKey', label: 'Red keys' },
    { key: 'goldenKey', label: 'Golden keys' },
    { key: 'masterKey', label: 'Master keys' },
];

const showManage = ref(false);
const manageUser = ref(null);
const manageLoading = ref(false);
const manageSaving = ref(false);
const stateForm = ref(null);

// Islands available for the currently-selected territory in the state form.
const manageIslands = computed(() => {
    if (!stateForm.value) return [];
    const t = territories.value.find(t => t.id === stateForm.value.atTerritory);
    return t?.islands || [];
});

const openManage = async user => {
    if (!territories.value.length) await loadTerritories();
    manageUser.value = user;
    stateForm.value = null;
    showManage.value = true;
    manageLoading.value = true;
    try {
        const p = await getPlayerState(user.id);
        stateForm.value = {
            atTerritory: p.atTerritory,
            atIsland: p.atIsland,
            anchored: !!p.anchored,
            coin: p.coin ?? 0,
            fuel: p.fuel ?? 0,
            fuelCap: p.fuelCap ?? 0,
            blueKey: p.blueKey ?? 0,
            redKey: p.redKey ?? 0,
            goldenKey: p.goldenKey ?? 0,
            masterKey: p.masterKey ?? 0,
        };
    } catch (e) {
        toast.error(`Failed to load player: ${e.message}`);
        showManage.value = false;
    } finally {
        manageLoading.value = false;
    }
};

// When the territory changes, keep the island valid: if the current island
// isn't in the new territory, fall back to its first island.
const onTerritoryChange = () => {
    const islands = manageIslands.value;
    if (!islands.some(i => i.id === stateForm.value.atIsland)) {
        stateForm.value.atIsland = islands[0]?.id || '';
    }
};

const closeManage = () => {
    showManage.value = false;
    manageUser.value = null;
    stateForm.value = null;
};

const submitManage = async () => {
    if (manageSaving.value || !stateForm.value) return;
    const f = stateForm.value;
    if (!f.atTerritory || !f.atIsland) {
        toast.warning('Location (territory and island) is required');
        return;
    }
    for (const { key, label } of NUMERIC_FIELDS) {
        if (!Number.isInteger(f[key]) || f[key] < 0) {
            toast.warning(`${label} must be a non-negative whole number`);
            return;
        }
    }
    if (f.fuel > f.fuelCap) {
        toast.warning('Fuel cannot exceed fuel capacity');
        return;
    }
    manageSaving.value = true;
    try {
        await editPlayerState(manageUser.value.id, {
            atTerritory: f.atTerritory,
            atIsland: f.atIsland,
            anchored: f.anchored,
            coin: f.coin,
            fuel: f.fuel,
            fuelCap: f.fuelCap,
            blueKey: f.blueKey,
            redKey: f.redKey,
            goldenKey: f.goldenKey,
            masterKey: f.masterKey,
        });
        toast.success(`Updated ${manageUser.value.username}'s state`);
        closeManage();
    } catch (e) {
        toast.error(`Failed to update player: ${e.message}`);
    } finally {
        manageSaving.value = false;
    }
};

const filteredUsers = computed(() => {
    const q = search.value.trim().toLowerCase();
    if (!q) return users.value;
    return users.value.filter(
        u =>
            (u.username || '').toLowerCase().includes(q) ||
            (u.name || '').toLowerCase().includes(q)
    );
});

const loadUsers = async () => {
    loading.value = true;
    try {
        users.value = (await getUsers()) || [];
    } catch (e) {
        toast.error(`Failed to load users: ${e.message}`);
    } finally {
        loading.value = false;
    }
};

const loadTerritories = async () => {
    try {
        territories.value = (await getTerritories()) || [];
    } catch (e) {
        toast.error(`Failed to load territories: ${e.message}`);
    }
};

const openCreate = () => {
    form.value = newForm();
    if (territories.value.length) {
        form.value.startingTerritory = territories.value[0].id;
    }
    createdUser.value = null;
    showCreate.value = true;
};

const submitCreate = async () => {
    if (creating.value) return;
    const f = form.value;
    if (!f.username.trim()) {
        toast.warning('Username is required');
        return;
    }
    if (!f.startingTerritory) {
        toast.warning('Starting territory is required');
        return;
    }
    creating.value = true;
    try {
        const payload = {
            name: f.name.trim(),
            username: f.username.trim(),
            startingTerritory: f.startingTerritory,
            meetLink: f.meetLink.trim(),
        };
        // Only send password if the admin typed one; otherwise backend
        // generates and returns one.
        if (f.password) payload.password = f.password;

        const result = await createUser(payload);
        createdUser.value = result;
        toast.success(`User "${result.username}" created`);
        await loadUsers();
    } catch (e) {
        toast.error(`Failed to create user: ${e.message}`);
    } finally {
        creating.value = false;
    }
};

const closeCreate = () => {
    showCreate.value = false;
    createdUser.value = null;
};

const copyPassword = async () => {
    if (!createdUser.value?.password) return;
    try {
        await navigator.clipboard.writeText(createdUser.value.password);
        toast.success('Password copied');
    } catch {
        toast.info('Copy failed — select the password manually');
    }
};

onMounted(() => {
    loadUsers();
    loadTerritories();
});
</script>

<template>
    <div class="admin-scope" dir="ltr">
        <header class="page-head">
            <div>
                <h1>Users</h1>
                <p>{{ users.length }} registered players.</p>
            </div>
            <button class="btn" @click="openCreate">+ Add user</button>
        </header>

        <div class="toolbar">
            <input
                v-model="search"
                class="search"
                type="text"
                placeholder="Search by username or name…"
            />
        </div>

        <div class="table-wrap">
            <table>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Username</th>
                        <th>Meet link</th>
                        <th class="col-actions">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="loading">
                        <td colspan="4" class="empty">Loading…</td>
                    </tr>
                    <tr v-else-if="filteredUsers.length === 0">
                        <td colspan="4" class="empty">No users found.</td>
                    </tr>
                    <tr v-for="u in filteredUsers" :key="u.username">
                        <td>{{ u.name || '—' }}</td>
                        <td>
                            <span class="mono">{{ u.username }}</span>
                        </td>
                        <td>
                            <a
                                v-if="u.meetLink"
                                :href="u.meetLink"
                                target="_blank"
                                rel="noopener"
                                class="link"
                                >{{ u.meetLink }}</a
                            >
                            <span v-else>—</span>
                        </td>
                        <td class="col-actions">
                            <button
                                class="btn-ghost btn btn-sm"
                                :disabled="!u.id"
                                :title="
                                    u.id
                                        ? 'Edit player state'
                                        : 'This user has no player id'
                                "
                                @click="openManage(u)"
                            >
                                Manage state
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- Create user modal -->
        <AdminModal v-if="showCreate" title="Add user" @close="closeCreate">
            <!-- Success view: show credentials -->
            <div v-if="createdUser" class="created">
                <p class="hint">
                    User created. Share these credentials with the player — the
                    password is shown only once.
                </p>
                <div class="cred">
                    <div class="cred-row">
                        <span class="cred-label">Username</span>
                        <span class="mono">{{ createdUser.username }}</span>
                    </div>
                    <div class="cred-row">
                        <span class="cred-label">Password</span>
                        <span class="mono">{{ createdUser.password }}</span>
                        <button class="copy" @click="copyPassword">Copy</button>
                    </div>
                </div>
                <div class="modal-actions">
                    <button class="btn-ghost btn" @click="openCreate">
                        Add another
                    </button>
                    <button class="btn" @click="closeCreate">Done</button>
                </div>
            </div>

            <!-- Form view -->
            <form v-else class="form" @submit.prevent="submitCreate">
                <label class="field">
                    <span>Username *</span>
                    <input v-model="form.username" type="text" required />
                </label>
                <label class="field">
                    <span>Name</span>
                    <input v-model="form.name" type="text" />
                </label>
                <label class="field">
                    <span>Starting territory *</span>
                    <select v-model="form.startingTerritory" required>
                        <option value="" disabled>Select…</option>
                        <option
                            v-for="t in territories"
                            :key="t.id"
                            :value="t.id"
                        >
                            {{ t.name || t.id }} ({{ t.id }})
                        </option>
                    </select>
                </label>
                <label class="field">
                    <span>Password</span>
                    <input
                        v-model="form.password"
                        type="text"
                        placeholder="Leave blank to auto-generate"
                    />
                </label>
                <label class="field">
                    <span>Meet link</span>
                    <input v-model="form.meetLink" type="text" />
                </label>

                <div class="modal-actions">
                    <button
                        type="button"
                        class="btn-ghost btn"
                        @click="closeCreate"
                    >
                        Cancel
                    </button>
                    <button type="submit" class="btn" :disabled="creating">
                        {{ creating ? 'Creating…' : 'Create user' }}
                    </button>
                </div>
            </form>
        </AdminModal>

        <!-- Manage player state modal -->
        <AdminModal
            v-if="showManage"
            :title="
                manageUser
                    ? `Player state — ${manageUser.username}`
                    : 'Player state'
            "
            @close="closeManage"
        >
            <div v-if="manageLoading" class="empty">Loading player…</div>

            <form
                v-else-if="stateForm"
                class="form"
                @submit.prevent="submitManage"
            >
                <p class="hint">
                    Changes are applied immediately and pushed to the player if
                    they are online.
                </p>

                <div class="field-grid">
                    <label class="field">
                        <span>Territory</span>
                        <select
                            v-model="stateForm.atTerritory"
                            @change="onTerritoryChange"
                        >
                            <option
                                v-for="t in territories"
                                :key="t.id"
                                :value="t.id"
                            >
                                {{ t.name || t.id }} ({{ t.id }})
                            </option>
                        </select>
                    </label>
                    <label class="field">
                        <span>Island</span>
                        <select v-model="stateForm.atIsland">
                            <option
                                v-for="i in manageIslands"
                                :key="i.id"
                                :value="i.id"
                            >
                                {{ i.name || i.id }}
                            </option>
                        </select>
                    </label>
                </div>

                <label class="field field-inline">
                    <input v-model="stateForm.anchored" type="checkbox" />
                    <span>Anchored</span>
                </label>

                <div class="field-grid">
                    <label
                        v-for="f in NUMERIC_FIELDS"
                        :key="f.key"
                        class="field"
                    >
                        <span>{{ f.label }}</span>
                        <input
                            v-model.number="stateForm[f.key]"
                            type="number"
                            min="0"
                            step="1"
                        />
                    </label>
                </div>

                <div class="modal-actions">
                    <button
                        type="button"
                        class="btn-ghost btn"
                        @click="closeManage"
                    >
                        Cancel
                    </button>
                    <button type="submit" class="btn" :disabled="manageSaving">
                        {{ manageSaving ? 'Saving…' : 'Save state' }}
                    </button>
                </div>
            </form>
        </AdminModal>
    </div>
</template>

<style scoped>
.toolbar {
    margin-bottom: 14px;
}
.search {
    width: 100%;
    max-width: 340px;
    background: #0b1220;
    border: 1px solid #334155;
    border-radius: 10px;
    padding: 10px 12px;
    color: #f8fafc;
    font-size: 14px;
    outline: none;
}
.search:focus {
    border-color: #2563eb;
}

.mono {
    font-family: ui-monospace, 'SF Mono', Menlo, monospace;
    font-size: 13px;
}
.link {
    color: #60a5fa;
    text-decoration: none;
    max-width: 320px;
    display: inline-block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    vertical-align: bottom;
}
.link:hover {
    text-decoration: underline;
}
.empty {
    text-align: center;
    color: #64748b;
    padding: 28px 16px;
}

.form {
    display: flex;
    flex-direction: column;
    gap: 14px;
}
.field-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 14px;
}
.field-inline {
    flex-direction: row;
    align-items: center;
    gap: 8px;
}
.field-inline span {
    margin: 0;
}
.field-inline input {
    width: auto;
}
.col-actions {
    text-align: right;
    white-space: nowrap;
}
.btn-sm {
    padding: 5px 12px;
    font-size: 12.5px;
}
.modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 8px;
}

.created {
    display: flex;
    flex-direction: column;
    gap: 16px;
}
.cred {
    background: #0b1220;
    border: 1px solid #1f2937;
    border-radius: 12px;
    padding: 14px 16px;
    display: flex;
    flex-direction: column;
    gap: 10px;
}
.cred-row {
    display: flex;
    align-items: center;
    gap: 10px;
}
.cred-label {
    width: 84px;
    color: #94a3b8;
    font-size: 12.5px;
}
.copy {
    margin-left: auto;
    background: #1f2937;
    border: none;
    color: #cbd5e1;
    border-radius: 8px;
    padding: 5px 10px;
    font-size: 12px;
    cursor: pointer;
}
.copy:hover {
    background: #374151;
}
</style>
