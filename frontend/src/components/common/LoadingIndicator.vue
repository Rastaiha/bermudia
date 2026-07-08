<template>
    <div
        class="text-2xl text-content font-main bg-surface/90 border border-border shadow-2xl px-8 py-4 rounded-lg"
    >
        {{ message }}
    </div>
</template>

<script setup>
defineProps({
    message: {
        type: String,
        default: 'Loading...',
    },
});
</script>
