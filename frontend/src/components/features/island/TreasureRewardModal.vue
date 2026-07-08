<template>
    <VueFinalModal
        class="flex justify-center items-center"
        content-class="panel w-full md:w-1/3 max-h-[80vh] !border-accent/60"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div class="panel-header">
            <h1 class="panel-title !text-xl">پاداش {{ glossary.treasure }}!</h1>
            <button class="panel-close" @click="handleClose">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-6 w-6"
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

        <div class="panel-body items-center">
            <div
                v-for="(reward, index) in rewards.items"
                :key="index"
                class="flex items-center w-full bg-accent/10 p-3 rounded-lg border border-accent/30"
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
                    <p class="text-md text-warning-content">
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
