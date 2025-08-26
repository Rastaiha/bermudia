<template>
    <div v-if="player"
        class="fixed top-1/2 right-5 -translate-y-1/2 w-10 flex flex-col items-center z-[200] font-vazir">
        <div class="mb-2 text-sm text-white">
            <img src="/images/icons/drop.svg" alt="Drop Icon" width="24" height="24" />
        </div>
        <div class="w-5 h-[150px] bg-blue-900/95 rounded-md overflow-hidden relative flex flex-col-reverse p-0.5">
            <div class="w-full bg-cyan-400 rounded-t-md relative overflow-hidden transition-height duration-300 ease-in"
                :style="{ height: `${fuelPercentage}%` }"></div>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
    player: {
        type: Object,
        required: true,
    },
});

const fuelPercentage = computed(() => {
    if (!props.player || props.player.fuelCap === 0) {
        return 0;
    }
    return (props.player.fuel / props.player.fuelCap) * 100;
});
</script>

<style scoped>
.transition-height {
    transition: height 0.3s ease-in-out;
}
</style>