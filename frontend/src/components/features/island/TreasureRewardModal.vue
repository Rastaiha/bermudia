<template>
    <VueFinalModal
        class="flex justify-center items-center"
        content-class="flex flex-col w-full md:w-1/3 mx-4 p-6 bg-yellow-900/90 border-4 border-amber-500 rounded-xl shadow-xl space-y-4"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div
            class="flex items-center justify-between border-b-2 border-amber-600 pb-2 mb-4"
        >
            <h1 class="text-xl font-semibold text-amber-100">
                پاداش {{ glossary.treasure }}!
            </h1>
            <button
                class="p-1 rounded-full hover:bg-amber-800"
                @click="handleClose"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-6 w-6 text-amber-100"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M6 18L18 6M6 6l12 12"
                    />
                </svg>
            </button>
        </div>

        <div
            class="flex flex-col items-center gap-4 max-h-96 overflow-y-auto p-2"
        >
            <div
                v-for="(reward, index) in rewards.items"
                :key="index"
                class="flex items-center w-full bg-amber-800/50 p-3 rounded-lg border border-amber-700"
            >
                <img
                    :src="COST_ITEMS_INFO[reward.type].icon"
                    :alt="reward.type + 'Icon'"
                    class="w-12 h-12 ml-4 drop-shadow-lg"
                />
                <div class="text-right">
                    <p class="text-lg font-bold text-white">
                        {{ COST_ITEMS_INFO[reward.type].name }}
                    </p>
                    <p class="text-md text-amber-200">
                        تعداد: {{ reward.amount }}
                    </p>
                </div>
            </div>
        </div>
    </VueFinalModal>
</template>

<script setup>
import { VueFinalModal } from 'vue-final-modal';
import { COST_ITEMS_INFO } from '@/services/cost.js';
import { glossary } from '@/services/glossary.js';

defineProps({
    rewards: {
        type: Object,
        default: () => {},
    },
});

const emit = defineEmits(['close']);

function handleClose() {
    emit('close');
}
</script>
