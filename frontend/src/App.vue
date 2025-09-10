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
import { ref, onMounted, onUnmounted, watch } from 'vue';
import eventBus from './services/eventBus';
import { getCurrentTrack } from './services/radioService';
import { audioSettings } from './services/audioSettings';

const audioPlayer = ref(null);
const isMusicPageActive = ref(false);
let isPlaybackBlockedByBrowser = false;
let fadeInterval = null;
const FADE_DURATION = 1500;

const attemptToUnlockAudio = () => {
    if (
        isMusicPageActive.value &&
        !audioSettings.isMuted &&
        isPlaybackBlockedByBrowser
    ) {
        playMusic();
    }
};

const playMusic = async () => {
    clearInterval(fadeInterval);

    try {
        const trackInfo = getCurrentTrack();
        const fullUrl = window.location.origin + trackInfo.url;

        if (audioPlayer.value.src !== fullUrl) {
            audioPlayer.value.src = fullUrl;
        }

        audioPlayer.value.currentTime = trackInfo.currentTime;
        audioPlayer.value.volume = 0;
        await audioPlayer.value.play();

        isPlaybackBlockedByBrowser = false;
        window.removeEventListener('pointerdown', attemptToUnlockAudio);

        const stepTime = 50;
        const volumeStep = 1 / (FADE_DURATION / stepTime);

        fadeInterval = setInterval(() => {
            const newVolume = audioPlayer.value.volume + volumeStep;
            if (newVolume >= 1) {
                audioPlayer.value.volume = 1;
                clearInterval(fadeInterval);
            } else {
                audioPlayer.value.volume = newVolume;
            }
        }, stepTime);
    } catch {
        isPlaybackBlockedByBrowser = true;
        window.addEventListener('pointerdown', attemptToUnlockAudio, {
            once: true,
        });
    }
};

const stopMusic = () => {
    clearInterval(fadeInterval);
    window.removeEventListener('pointerdown', attemptToUnlockAudio);

    if (!audioPlayer.value || audioPlayer.value.paused) {
        return;
    }

    const stepTime = 50;
    const currentVolume = audioPlayer.value.volume;
    const volumeStep = currentVolume / (FADE_DURATION / stepTime);

    fadeInterval = setInterval(() => {
        const newVolume = audioPlayer.value.volume - volumeStep;
        if (newVolume <= 0) {
            audioPlayer.value.volume = 0;
            clearInterval(fadeInterval);
            audioPlayer.value.pause();
            audioPlayer.value.src = '';
        } else {
            audioPlayer.value.volume = newVolume;
        }
    }, stepTime);
};

const handleSongEnd = () => {
    if (isMusicPageActive.value && !audioSettings.isMuted) {
        playMusic();
    }
};

watch(
    [isMusicPageActive, () => audioSettings.isMuted],
    ([pageActive, isMuted]) => {
        if (pageActive && !isMuted) {
            playMusic();
        } else {
            stopMusic();
        }
    },
    { immediate: true }
);

onMounted(() => {
    eventBus.on('set-audio-state', state => {
        isMusicPageActive.value = state === 'play';
    });
});

onUnmounted(() => {
    eventBus.off('set-audio-state');
    window.removeEventListener('pointerdown', attemptToUnlockAudio);
});
</script>
