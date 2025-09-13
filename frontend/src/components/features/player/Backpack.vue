<template>
    <VueFinalModal
        class="flex justify-center items-center"
        content-class="flex flex-col w-full md:w-1/3 mx-4 p-6 
                       bg-[#1E3A8A] border-4 border-[#1E40AF] 
                       rounded-xl shadow-xl space-y-4"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div
            class="flex items-center justify-between border-b-2 border-[#1E40AF] pb-2 mb-4"
        >
            <h1 class="text-xl font-semibold text-blue-200">کوله پشتی</h1>
            <button
                class="p-1 rounded-full hover:bg-[#1E40AF]"
                @click="handleClose"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-6 w-6 text-blue-200"
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

        <div class="w-full flex flex-col space-y-4">
            <div
                v-for="item in INVENTORY_ITEMS"
                :key="item"
                class="flex items-center justify-between p-3 bg-blue-900 bg-opacity-50 rounded-lg"
            >
                <div class="flex items-center">
                    <img
                        :src="COST_ITEMS_INFO[item].icon"
                        class="w-8 h-8 ml-3"
                        :alt="item + 'icon'"
                    />
                    <span class="text-base font-medium text-white">{{
                        COST_ITEMS_INFO[item].name
                    }}</span>
                </div>
                <span class="text-lg font-bold text-blue-200"
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
