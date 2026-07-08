<template>
    <VueFinalModal
        class="flex justify-center items-center"
        content-class="panel w-full md:w-1/3 max-h-[80vh]"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div class="panel-header">
            <h1 class="panel-title">کوله پشتی</h1>
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

        <div class="panel-body">
            <div v-for="item in INVENTORY_ITEMS" :key="item" class="panel-row">
                <div class="flex items-center">
                    <img
                        :src="COST_ITEMS_INFO[item].icon"
                        class="w-12 h-12 ml-3"
                        :alt="item + 'icon'"
                    />
                    <span class="text-base font-medium text-white">{{
                        COST_ITEMS_INFO[item].name
                    }}</span>
                </div>
                <span class="text-lg font-bold text-accent"
                    >x {{ player[item] }}</span
                >
            </div>
        </div>
    </VueFinalModal>
</template>

<script setup>
import { VueFinalModal } from 'vue-final-modal';
import { COST_ITEMS_INFO, INVENTORY_ITEMS } from '@/services/cost';

defineProps({
    player: Object,
});

const emit = defineEmits(['close']);

function handleClose() {
    emit('close');
}
</script>
