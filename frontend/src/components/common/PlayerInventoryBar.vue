<template>
    <div
        class="mb-3 transition-all duration-300"
        :class="{
            'drop-shadow-[0_0_8px_rgba(239,68,68,0.7)]':
                barData.englishName === 'Coin' && barData.value < 0,
        }"
    >
        <div
            class="flex justify-between px-1 mb-1 text-xs text-gray-300 drop-shadow-md"
        >
            <span class="hidden sm:block">{{ barData.name }}</span>
            <span
                :style="{
                    visibility: barData.total == -1 ? 'hidden' : 'visible',
                }"
                class="hidden sm:block"
                >ظرفیت: {{ barData.total }}</span
            >
        </div>
        <div
            class="relative flex items-center h-6 rounded-md shadow-inner"
            :class="[
                barData.englishName === 'Coin' && barData.value < 0
                    ? 'bg-red-800/50'
                    : 'bg-black/30',
            ]"
            :style="{
                width: barData.width ? barData.width * 100 + '%' : '100%',
            }"
        >
            <div
                v-if="barData.required"
                class="border-l-6 h-8 absolute z-20"
                :class="requiredLineClass"
                :style="{
                    right:
                        ((barData.width ? barData.width : 1) *
                            barData.required *
                            99) /
                            barData.total +
                        '%',
                }"
            ></div>
            <div
                class="absolute inset-0 flex items-center justify-end w-full h-full z-11 gap-2 flex-row-reverse pr-2"
            >
                <div
                    class="h-4 text-white text-xs font-bold drop-shadow-md text-right"
                >
                    {{ barData.value }}
                </div>
                <img
                    :src="barData.icon"
                    :alt="barData.englishName + ' Icon'"
                    class="w-8 h-8"
                />
            </div>
            <div
                class="absolute top-0 right-0 h-full rounded-md transition-[width] duration-500 ease-in-out bar-shadow"
                :style="{
                    width: barPercentage + '%',
                    backgroundImage:
                        'linear-gradient(to left, ' +
                        barData.gradientFrom +
                        ', ' +
                        barData.gradientTo +
                        ')',
                    filter:
                        'drop-shadow(0 0px 3px ' + barData.shadowColor + ')',
                }"
            ></div>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
    barData: Object,
});

const barPercentage = computed(() => {
    const value = props.barData.value;
    const total = props.barData.total;

    if (total === 0 || total === -1) return 0;

    const percentage = (value / total) * 100;
    return Math.max(0, percentage);
});

const requiredLineClass = computed(() => {
    if (props.barData.required <= props.barData.value) {
        return 'border-green-500';
    } else {
        return 'border-red-500 border-dotted';
    }
});
</script>
