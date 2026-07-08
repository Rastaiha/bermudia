<script setup>
import { ref, onMounted, computed } from 'vue';
import { useToast } from 'vue-toastification';
import { getUsers, createUser, getTerritories } from '../services/api.js';
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
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="loading">
                        <td colspan="3" class="empty">Loading…</td>
                    </tr>
                    <tr v-else-if="filteredUsers.length === 0">
                        <td colspan="3" class="empty">No users found.</td>
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
