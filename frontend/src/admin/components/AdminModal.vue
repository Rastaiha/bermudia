<script setup>
import { onMounted, onUnmounted } from 'vue';

const emit = defineEmits(['close']);

defineProps({
    title: { type: String, default: '' },
    wide: { type: Boolean, default: false },
});

const onKey = e => {
    if (e.key === 'Escape') emit('close');
};

onMounted(() => window.addEventListener('keydown', onKey));
onUnmounted(() => window.removeEventListener('keydown', onKey));
</script>

<template>
    <div class="admin-modal-overlay" @click.self="emit('close')">
        <div class="admin-modal-card" :class="{ wide }" dir="ltr">
            <div class="modal-head">
                <h2>{{ title }}</h2>
                <button class="modal-x" @click="emit('close')">✕</button>
            </div>
            <slot />
        </div>
    </div>
</template>

<style scoped>
.modal-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 18px;
}
.modal-head h2 {
    font-size: 17px;
    font-weight: 700;
    color: #f8fafc;
}
.modal-x {
    background: transparent;
    border: none;
    color: #94a3b8;
    font-size: 16px;
    cursor: pointer;
    padding: 4px 8px;
    border-radius: 8px;
}
.modal-x:hover {
    background: #1f2937;
    color: #f8fafc;
}
</style>
