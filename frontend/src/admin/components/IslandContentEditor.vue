<script setup>
import { ref, watch } from 'vue';
import { useToast } from 'vue-toastification';
import { getIslandHeader, getBook, setIslandBook } from '../services/api.js';

const props = defineProps({
    islandId: { type: String, required: true },
});
const emit = defineEmits(['saved']);

const toast = useToast();

const REWARD_SOURCES = [
    '',
    'edu1',
    'edu2',
    'edu3',
    'edu4',
    'edu5',
    'edu6',
    'final',
];
const INPUT_TYPES = ['text', 'number', 'file'];

const loading = ref(true);
const saving = ref(false);
const header = ref(null);
const fromPool = ref(false);
const bookId = ref('');
const components = ref([]); // [{ kind:'iframe'|'question', ... }]
const treasures = ref([]); // [{ id }]

const load = async () => {
    loading.value = true;
    components.value = [];
    treasures.value = [];
    try {
        const h = await getIslandHeader(props.islandId);
        header.value = h;
        fromPool.value = !!h.fromPool;
        bookId.value = h.bookId || '';

        if (bookId.value) {
            const book = await getBook(bookId.value);
            components.value = (book.components || []).map(toEditable);
            treasures.value = (book.treasures || []).map(t => ({
                id: t.id || '',
            }));
        }
    } catch (e) {
        toast.error(`Failed to load island content: ${e.message}`);
    } finally {
        loading.value = false;
    }
};

const toEditable = c => {
    if (c.iframe) {
        return { kind: 'iframe', url: c.iframe.url || '' };
    }
    const q = c.question || {};
    return {
        kind: 'question',
        id: q.id || '',
        text: q.text || '',
        inputType: q.inputType || 'text',
        inputAccept: (q.inputAccept || []).join(', '),
        knowledgeAmount: q.knowledgeAmount ?? 0,
        rewardSource: q.rewardSource || '',
        correctionHintMessage: q.correctionHintMessage || '',
    };
};

// ---- Component list mutations ----
const addIframe = () => components.value.push({ kind: 'iframe', url: '' });
const addQuestion = () =>
    components.value.push({
        kind: 'question',
        id: '',
        text: '',
        inputType: 'text',
        inputAccept: '',
        knowledgeAmount: 0,
        rewardSource: '',
        correctionHintMessage: '',
    });
const removeComponent = i => components.value.splice(i, 1);
const move = (i, dir) => {
    const j = i + dir;
    if (j < 0 || j >= components.value.length) return;
    const arr = components.value;
    [arr[i], arr[j]] = [arr[j], arr[i]];
};

const addTreasure = () => treasures.value.push({ id: '' });
const removeTreasure = i => treasures.value.splice(i, 1);

// ---- Save ----
const buildPayload = () => {
    const payloadComponents = components.value.map(c => {
        if (c.kind === 'iframe') {
            return { iframe: { url: c.url.trim() } };
        }
        const q = {
            id: c.id || undefined,
            text: c.text,
            inputType: c.inputType,
            knowledgeAmount: Number(c.knowledgeAmount) || 0,
            rewardSource: c.rewardSource || undefined,
            correctionHintMessage: c.correctionHintMessage || undefined,
        };
        if (c.inputType === 'file') {
            q.inputAccept = c.inputAccept
                .split(',')
                .map(s => s.trim())
                .filter(Boolean);
        }
        return { question: q };
    });
    return {
        bookId: bookId.value || undefined,
        components: payloadComponents,
        treasures: treasures.value.map(t => ({ id: t.id || undefined })),
    };
};

const validate = () => {
    for (const c of components.value) {
        if (c.kind === 'iframe' && !c.url.trim()) {
            return 'An iframe component is missing its URL.';
        }
        if (c.kind === 'question') {
            if (!c.text.trim()) return 'A question is missing its text.';
            if (c.inputType === 'file' && !c.inputAccept.trim())
                return 'A file question needs at least one accepted type (e.g. image/png).';
        }
    }
    return null;
};

const save = async () => {
    if (saving.value) return;
    const err = validate();
    if (err) {
        toast.warning(err);
        return;
    }
    saving.value = true;
    try {
        const result = await setIslandBook(props.islandId, buildPayload());
        // backend returns the canonical book (with generated ids); refresh
        bookId.value = result.bookId || bookId.value;
        components.value = (result.components || []).map(toEditable);
        treasures.value = (result.treasures || []).map(t => ({
            id: t.id || '',
        }));
        toast.success('Island content saved');
        emit('saved', props.islandId);
    } catch (e) {
        toast.error(`Save failed: ${e.message}`);
    } finally {
        saving.value = false;
    }
};

watch(() => props.islandId, load, { immediate: true });
</script>

<template>
    <div class="ice">
        <p v-if="loading" class="hint">Loading…</p>

        <template v-else>
            <div v-if="fromPool" class="pool-warn">
                ⚠️ This island draws its content from a shared pool. Editing a
                book here won't apply while it's pool-bound.
            </div>

            <div class="ice-head">
                <div>
                    <div class="ice-title">
                        {{ header?.name || islandId }}
                    </div>
                    <code class="ice-sub">{{ islandId }}</code>
                </div>
                <code v-if="bookId" class="ice-book">{{ bookId }}</code>
            </div>

            <!-- Components -->
            <div class="section-label">Content ({{ components.length }})</div>
            <div v-if="!components.length" class="hint empty-line">
                No components yet. Add an article or a question below.
            </div>

            <div
                v-for="(c, i) in components"
                :key="i"
                class="comp"
                :class="c.kind"
            >
                <div class="comp-head">
                    <span class="comp-kind">
                        {{ c.kind === 'iframe' ? '🖼️ Article' : '❓ Question' }}
                    </span>
                    <div class="comp-tools">
                        <button
                            class="tool"
                            :disabled="i === 0"
                            title="Move up"
                            @click="move(i, -1)"
                        >
                            ↑
                        </button>
                        <button
                            class="tool"
                            :disabled="i === components.length - 1"
                            title="Move down"
                            @click="move(i, 1)"
                        >
                            ↓
                        </button>
                        <button
                            class="tool tool-danger"
                            title="Remove"
                            @click="removeComponent(i)"
                        >
                            ✕
                        </button>
                    </div>
                </div>

                <!-- iframe -->
                <label v-if="c.kind === 'iframe'" class="field">
                    <span>Article URL</span>
                    <input
                        v-model="c.url"
                        type="text"
                        placeholder="https://…"
                    />
                </label>

                <!-- question -->
                <template v-else>
                    <label class="field">
                        <span>Question text</span>
                        <textarea v-model="c.text" rows="2"></textarea>
                    </label>
                    <div class="row3">
                        <label class="field">
                            <span>Input type</span>
                            <select v-model="c.inputType">
                                <option
                                    v-for="t in INPUT_TYPES"
                                    :key="t"
                                    :value="t"
                                >
                                    {{ t }}
                                </option>
                            </select>
                        </label>
                        <label class="field">
                            <span>Knowledge</span>
                            <input
                                v-model.number="c.knowledgeAmount"
                                type="number"
                                min="0"
                            />
                        </label>
                        <label class="field">
                            <span>Reward</span>
                            <select v-model="c.rewardSource">
                                <option
                                    v-for="r in REWARD_SOURCES"
                                    :key="r"
                                    :value="r"
                                >
                                    {{ r || '(none)' }}
                                </option>
                            </select>
                        </label>
                    </div>
                    <label v-if="c.inputType === 'file'" class="field">
                        <span>Accepted file types (comma-separated MIME)</span>
                        <input
                            v-model="c.inputAccept"
                            type="text"
                            placeholder="image/png, image/jpeg"
                        />
                    </label>
                    <label class="field">
                        <span>Correction hint (optional)</span>
                        <input
                            v-model="c.correctionHintMessage"
                            type="text"
                            placeholder="Shown to correctors"
                        />
                    </label>
                </template>
            </div>

            <div class="add-row">
                <button class="btn btn-ghost" @click="addIframe">
                    + Article
                </button>
                <button class="btn btn-ghost" @click="addQuestion">
                    + Question
                </button>
            </div>

            <!-- Treasures -->
            <div class="section-label">Treasures ({{ treasures.length }})</div>
            <p class="hint">
                Each treasure is a hidden reward unlocked with keys. IDs are
                generated on save; leave blank to create new ones.
            </p>
            <div
                v-for="(t, i) in treasures"
                :key="`t-${i}`"
                class="treasure-row"
            >
                <input
                    v-model="t.id"
                    class="treasure-id"
                    type="text"
                    placeholder="(new treasure)"
                    readonly
                />
                <button
                    class="tool tool-danger"
                    title="Remove"
                    @click="removeTreasure(i)"
                >
                    ✕
                </button>
            </div>
            <div class="add-row">
                <button class="btn btn-ghost" @click="addTreasure">
                    + Treasure
                </button>
            </div>

            <div class="ice-actions">
                <button class="btn" :disabled="saving" @click="save">
                    {{ saving ? 'Saving…' : '💾 Save content' }}
                </button>
            </div>
        </template>
    </div>
</template>

<style scoped>
.ice {
    display: flex;
    flex-direction: column;
    gap: 12px;
}
.pool-warn {
    background: #78350f44;
    border: 1px solid #b45309;
    color: #fcd34d;
    border-radius: 10px;
    padding: 10px 12px;
    font-size: 12.5px;
}
.ice-head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 12px;
}
.ice-title {
    font-weight: 700;
    color: #f1f5f9;
    font-size: 15px;
}
.ice-sub,
.ice-book {
    font-size: 11px;
    color: #64748b;
}
.section-label {
    font-size: 12px;
    font-weight: 600;
    color: #94a3b8;
    text-transform: uppercase;
    letter-spacing: 0.03em;
    margin-top: 6px;
}
.empty-line {
    padding: 4px 0;
}

.comp {
    border: 1px solid #1f2937;
    border-radius: 12px;
    padding: 12px;
    background: #0b1220;
    display: flex;
    flex-direction: column;
    gap: 10px;
}
.comp.iframe {
    border-left: 3px solid #38bdf8;
}
.comp.question {
    border-left: 3px solid #a78bfa;
}
.comp-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
}
.comp-kind {
    font-size: 13px;
    font-weight: 600;
    color: #e2e8f0;
}
.comp-tools {
    display: flex;
    gap: 4px;
}
.tool {
    background: #111827;
    border: 1px solid #334155;
    color: #cbd5e1;
    width: 26px;
    height: 26px;
    border-radius: 8px;
    cursor: pointer;
    font-size: 13px;
    line-height: 1;
}
.tool:hover:not(:disabled) {
    background: #1f2937;
}
.tool:disabled {
    opacity: 0.4;
    cursor: default;
}
.tool-danger {
    color: #f87171;
}
.tool-danger:hover {
    background: #7f1d1d33;
}

.row3 {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    gap: 8px;
}
.row3 .field input,
.row3 .field select {
    min-width: 0;
}

.add-row {
    display: flex;
    gap: 8px;
}

.treasure-row {
    display: flex;
    gap: 8px;
    align-items: center;
}
.treasure-id {
    flex: 1;
    min-width: 0;
    background: #0b1220;
    border: 1px solid #334155;
    border-radius: 10px;
    padding: 8px 12px;
    color: #94a3b8;
    font-size: 12px;
    font-family: ui-monospace, monospace;
}

.ice-actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 6px;
    border-top: 1px solid #1f2937;
    padding-top: 12px;
}
</style>
