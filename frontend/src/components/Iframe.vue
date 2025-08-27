<template>
    <div ref="containerRef"
        class="relative bg-slate-900/80 border border-sky-500/40 rounded-3xl shadow-lg overflow-hidden transition-all duration-300 ease-linear"
        :class="{ 'fullscreen-active': isFullscreen }">
        <button
            class="absolute top-4 left-4 z-10 bg-black/40 border border-white/20 rounded-full w-11 h-11 flex items-center justify-center cursor-pointer transition-all duration-200 hover:bg-black/70 hover:scale-110"
            @click="toggleFullScreen"
            @mouseover="$emit('showTooltip', isFullscreen ? 'خروج از تمام صفحه' : 'تمام صفحه')"
            @mouseleave="$emit('hideTooltip')">
            <svg v-if="!isFullscreen" xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px"
                fill="#FFFFFF">
                <path d="M0 0h24v24H0V0z" fill="none" />
                <path d="M7 14H5v5h5v-2H7v-3zm-2-4h2V7h3V5H5v5zm12 7h-3v2h5v-5h-2v3zM14 5v2h3v3h2V5h-5z" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px"
                fill="#FFFFFF">
                <path d="M0 0h24v24H0V0z" fill="none" />
                <path d="M5 16h3v3h2v-5H5v2zm3-8H5v2h5V5H8v3zm6 11h2v-3h3v-2h-5v5zm2-11V5h-2v5h5V8h-3z" />
            </svg>
        </button>
        <iframe :src="url" frameborder="0" allowfullscreen allow="fullscreen"
            class="w-full h-[550px] border-none block"></iframe>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';

defineProps({
    url: {
        type: String,
        required: true,
    },
});

defineEmits(['showTooltip', 'hideTooltip']);

const isFullscreen = ref(false);
const containerRef = ref(null);

const toggleFullScreen = () => {
    if (isFullscreen.value) {
        if (document.exitFullscreen) {
            document.exitFullscreen();
        }
    } else {
        const el = containerRef.value;
        if (el && el.requestFullscreen) {
            el.requestFullscreen();
        }
    }
};

const onFullscreenChange = () => {
    isFullscreen.value = !!document.fullscreenElement;
};

onMounted(() => {
    document.addEventListener('fullscreenchange', onFullscreenChange);
});

onUnmounted(() => {
    document.removeEventListener('fullscreenchange', onFullscreenChange);
});
</script>

<style scoped>
/* Fullscreen Overrides */
.fullscreen-active {
    width: 100vw !important;
    height: 100vh !important;
    max-width: 100% !important;
    margin: 0 !important;
    padding: 0 !important;
    border-radius: 0 !important;
    position: fixed;
    top: 0;
    left: 0;
    z-index: 9999;
    background-color: #000;
}

.fullscreen-active iframe {
    height: 100% !important;
}
</style>