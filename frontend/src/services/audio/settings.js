import { reactive } from 'vue';

const STORAGE_KEY = 'bermudia_audio_muted';

export const audioSettings = reactive({
    isMuted: localStorage.getItem(STORAGE_KEY) === 'true',

    toggleMute() {
        this.isMuted = !this.isMuted;
        localStorage.setItem(STORAGE_KEY, this.isMuted);
    },
});
