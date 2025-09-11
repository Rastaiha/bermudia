import { ref, onMounted, onUnmounted } from 'vue';

export function useNow(interval = 1000) {
    const now = ref(Date.now());
    let timer = null;

    onMounted(() => {
        timer = setInterval(() => {
            now.value = Date.now();
        }, interval);
    });

    onUnmounted(() => {
        clearInterval(timer);
    });

    return { now };
}
