<template>
    <button
        class="h-12 w-12 rounded-lg p-2 transition-colors duration-200 hover:bg-slate-700/90 flex justify-center items-center"
        title="قطع/پخش صدا"
        @pointerdown.stop="toggleMute"
    >
        <img
            :src="iconSrc"
            alt="Mute/Unmute Icon"
            class="h-full w-full object-contain"
        />
    </button>
</template>

<script setup>
import { computed } from 'vue';
import { audioSettings } from '@/services/audio/settings.js';

const iconSrc = computed(() => {
    return audioSettings.isMuted
        ? '/images/icons/mute.png'
        : '/images/icons/unmute.png';
});

const toggleMute = () => {
    audioSettings.toggleMute();
};
</script>
