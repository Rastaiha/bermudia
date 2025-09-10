<template>
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
import { useAudioPlayer } from './services/useAudioPlayer';

const audioPlayer = ref(null);

const { handleSongEnd } = useAudioPlayer(audioPlayer);
</script>
