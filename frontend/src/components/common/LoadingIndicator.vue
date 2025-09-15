<template>
    <div
        class="text-2xl text-gray-300 font-main bg-[rgb(121,200,237,0.8)] px-8 py-4 rounded-lg"
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
