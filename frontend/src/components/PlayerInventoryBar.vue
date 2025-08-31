<template>
    <div class="mb-3">
        <div class="flex justify-between px-1 mb-1 text-xs text-gray-300 drop-shadow-md">
            <!-- <span>{{ barData.name }}</span> -->
            <span></span>
            <span :style="{ visibility: barData.total == -1 ? 'hidden' : 'visible' }">ظرفیت: {{ barData.total }}</span>
        </div>
        <div class="relative flex items-center h-6 rounded-md bg-black/30 shadow-inner"
            :class="barData.total == -1 ? 'w-3/4' : 'w-full'">
            <div class="absolute inset-0 flex items-center justify-end w-full h-full z-11 gap-2 flex-row-reverse pr-2">
                <div class="h-4 text-white text-xs font-bold drop-shadow-md text-right">
                    {{ barData.value }}
                </div>
                <img :src="barData.icon" :alt="barData.englishName + ' Icon'" class="w-12 h-12">
            </div>
            <div class="absolute top-0 right-0 h-full rounded-md transition-[width] duration-500 ease-in-out bar-shadow"
                :style="{
                    width: barPercentage + '%',
                    backgroundImage: 'linear-gradient(to left, ' + barData.gradientFrom + ', ' + barData.gradientTo + ')',
                    filter: 'drop-shadow(0 0px 3px ' + barData.shadowColor + ')'
                }"></div>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
    barData: Object
});


const barPercentage = computed(() => {
    if (!props.barData.value || props.barData.total === 0) return 0;
    if (props.barData.total === -1) return 0;
    return (props.barData.value / props.barData.total) * 100;
});
</script>
