<template>
    <UnreadMessageIndicator />
    <router-view />
    <ModalsContainer />
    <audio
        ref="audioPlayer"
        style="display: none"
        @ended="handleSongEnd"
    ></audio>
</template>

<script setup>
import { ModalsContainer } from 'vue-final-modal';
import { ref } from 'vue';
import { useAudioPlayer } from '@/composables/useAudioPlayer.js';
import UnreadMessageIndicator from '@/components/features/notification/UnreadMessageIndicator.vue';

const audioPlayer = ref(null);

const { handleSongEnd } = useAudioPlayer(audioPlayer);
</script>
