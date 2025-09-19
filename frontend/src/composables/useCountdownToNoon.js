import { ref, onMounted, onUnmounted, computed } from 'vue';

export function useCountdownToNoon() {
    const now = ref(new Date());
    const hours = ref(0);
    const minutes = ref(0);
    const seconds = ref(0);
    const isPastNoon = ref(false);

    let intervalId;

    const updateCountdown = () => {
        now.value = new Date();
        const target = new Date(now.value);
        target.setHours(12, 30, 0, 0);

        if (now.value >= target) {
            isPastNoon.value = true;
            hours.value = 0;
            minutes.value = 0;
            seconds.value = 0;
            if (intervalId) clearInterval(intervalId);
            return;
        }

        isPastNoon.value = false;
        const diff = target.getTime() - now.value.getTime();

        hours.value = Math.floor(diff / (1000 * 60 * 60));
        minutes.value = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
        seconds.value = Math.floor((diff % (1000 * 60)) / 1000);
    };

    onMounted(() => {
        updateCountdown();
        intervalId = setInterval(updateCountdown, 1000);
    });

    onUnmounted(() => {
        if (intervalId) clearInterval(intervalId);
    });

    const padZero = num => num.toString().padStart(2, '0');

    return {
        hours: computed(() => padZero(hours.value)),
        minutes: computed(() => padZero(minutes.value)),
        seconds: computed(() => padZero(seconds.value)),
        isPastNoon,
    };
}
