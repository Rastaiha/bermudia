<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useToast } from 'vue-toastification';
import {
    getGameState,
    setGameState,
    broadcast,
    getConnections,
} from '../services/api.js';

const toast = useToast();

const isPaused = ref(false);
const pauseLoading = ref(false);
const stateLoaded = ref(false);

const message = ref('');
const broadcastLoading = ref(false);

const connections = ref(null);
let pollTimer = null;

const loadGameState = async () => {
    try {
        const res = await getGameState();
        isPaused.value = !!res.isPaused;
        stateLoaded.value = true;
    } catch (e) {
        toast.error(`Failed to load game state: ${e.message}`);
    }
};

const togglePause = async () => {
    if (pauseLoading.value) return;
    pauseLoading.value = true;
    const next = !isPaused.value;
    try {
        const res = await setGameState(next);
        isPaused.value = !!res.isPaused;
        toast.success(isPaused.value ? 'Game paused' : 'Game resumed');
    } catch (e) {
        toast.error(`Failed to update: ${e.message}`);
    } finally {
        pauseLoading.value = false;
    }
};

const sendBroadcast = async () => {
    const text = message.value.trim();
    if (!text) {
        toast.warning('Message is empty');
        return;
    }
    if (broadcastLoading.value) return;
    broadcastLoading.value = true;
    try {
        const res = await broadcast(text);
        toast.success(`Sent to ${res.sentTo} players`);
        message.value = '';
    } catch (e) {
        toast.error(`Broadcast failed: ${e.message}`);
    } finally {
        broadcastLoading.value = false;
    }
};

const loadConnections = async () => {
    try {
        connections.value = await getConnections();
    } catch {
        // keep last known values; avoid toast spam on poll
    }
};

onMounted(() => {
    loadGameState();
    loadConnections();
    pollTimer = setInterval(loadConnections, 10000);
});

onUnmounted(() => {
    if (pollTimer) clearInterval(pollTimer);
});
</script>

<template>
    <div class="page">
        <header class="page-head">
            <h1>General Settings</h1>
            <p>Control game state, broadcast messages, and monitor activity.</p>
        </header>

        <div class="grid">
            <!-- Game state -->
            <section class="card">
                <h2>Game State</h2>
                <div class="state-row">
                    <span
                        class="badge"
                        :class="isPaused ? 'badge-paused' : 'badge-live'"
                    >
                        {{ isPaused ? '⏸ Paused' : '▶ Running' }}
                    </span>
                    <button
                        class="btn"
                        :class="isPaused ? 'btn-green' : 'btn-amber'"
                        :disabled="pauseLoading || !stateLoaded"
                        @click="togglePause"
                    >
                        {{
                            pauseLoading
                                ? 'Working…'
                                : isPaused
                                  ? 'Resume game'
                                  : 'Pause game'
                        }}
                    </button>
                </div>
                <p class="hint">
                    While paused, players cannot travel, refuel, trade, answer,
                    or perform any state-changing action.
                </p>
            </section>

            <!-- Connections -->
            <section class="card">
                <h2>Live Connections</h2>
                <div v-if="connections" class="conn-grid">
                    <div class="conn-cell">
                        <span class="conn-num">{{ connections.players }}</span>
                        <span class="conn-label">Players</span>
                    </div>
                    <div class="conn-cell">
                        <span class="conn-num">{{ connections.market }}</span>
                        <span class="conn-label">Market</span>
                    </div>
                    <div class="conn-cell">
                        <span class="conn-num">{{ connections.inbox }}</span>
                        <span class="conn-label">Inbox</span>
                    </div>
                </div>
                <p v-else class="hint">Loading…</p>
                <p class="hint">Auto-refreshes every 10 seconds.</p>
            </section>
        </div>

        <!-- Broadcast -->
        <section class="card">
            <h2>Broadcast Message</h2>
            <p class="hint">Sends an announcement to every player's inbox.</p>
            <textarea
                v-model="message"
                class="textarea"
                rows="4"
                placeholder="Type your announcement…"
            ></textarea>
            <div class="broadcast-actions">
                <button
                    class="btn btn-blue"
                    :disabled="broadcastLoading"
                    @click="sendBroadcast"
                >
                    {{ broadcastLoading ? 'Sending…' : 'Send to all players' }}
                </button>
            </div>
        </section>
    </div>
</template>

<style scoped>
.page {
    max-width: 900px;
}
.page-head {
    margin-bottom: 24px;
}
.page-head h1 {
    font-size: 24px;
    font-weight: 700;
    color: #f8fafc;
}
.page-head p {
    color: #94a3b8;
    font-size: 14px;
    margin-top: 4px;
}

.grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 18px;
    margin-bottom: 18px;
}
@media (max-width: 720px) {
    .grid {
        grid-template-columns: 1fr;
    }
}

.card {
    background: #111827;
    border: 1px solid #1f2937;
    border-radius: 14px;
    padding: 20px 22px;
}
.card h2 {
    font-size: 15px;
    font-weight: 600;
    color: #f1f5f9;
    margin-bottom: 14px;
}

.state-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    flex-wrap: wrap;
}

.badge {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 6px 12px;
    border-radius: 999px;
    font-size: 13px;
    font-weight: 600;
}
.badge-live {
    background: #064e3b;
    color: #6ee7b7;
}
.badge-paused {
    background: #78350f;
    color: #fcd34d;
}

.btn {
    border: none;
    border-radius: 10px;
    padding: 10px 16px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    color: #fff;
    transition:
        background 0.15s,
        opacity 0.15s;
}
.btn:disabled {
    opacity: 0.6;
    cursor: default;
}
.btn-green {
    background: #059669;
}
.btn-green:hover:not(:disabled) {
    background: #047857;
}
.btn-amber {
    background: #d97706;
}
.btn-amber:hover:not(:disabled) {
    background: #b45309;
}
.btn-blue {
    background: #2563eb;
}
.btn-blue:hover:not(:disabled) {
    background: #1d4ed8;
}

.hint {
    color: #94a3b8;
    font-size: 12.5px;
    margin-top: 10px;
    line-height: 1.5;
}

.conn-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 10px;
}
.conn-cell {
    background: #0b1220;
    border: 1px solid #1f2937;
    border-radius: 10px;
    padding: 14px 8px;
    text-align: center;
}
.conn-num {
    display: block;
    font-size: 26px;
    font-weight: 700;
    color: #60a5fa;
}
.conn-label {
    display: block;
    font-size: 12px;
    color: #94a3b8;
    margin-top: 2px;
}

.textarea {
    width: 100%;
    background: #0b1220;
    border: 1px solid #334155;
    border-radius: 10px;
    padding: 12px;
    color: #f8fafc;
    font-size: 14px;
    resize: vertical;
    outline: none;
    font-family: inherit;
}
.textarea:focus {
    border-color: #2563eb;
}
.broadcast-actions {
    margin-top: 12px;
    display: flex;
    justify-content: flex-end;
}
</style>
