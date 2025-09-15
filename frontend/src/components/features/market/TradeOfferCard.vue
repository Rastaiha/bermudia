<template>
    <div
        class="flex justify-between flex-col bg-slate-900/70 rounded-xl border border-slate-700 shadow-lg p-4 space-y-4 transition-transform hover:scale-105"
    >
        <div class="flex flex-row justify-around gap-2.5 items-center">
            <span class="text-sm font-bold text-green-400 mb-1">می‌دهد</span>
            <span></span>
            <span class="text-sm font-bold text-red-400 mb-1">می‌خواهد</span>
        </div>
        <div class="flex items-center justify-center text-center">
            <div class="flex-1 flex flex-col gap-2.5 items-center">
                <div
                    v-if="offer.offered.items.length > 0"
                    class="flex flex-col items-center gap-2"
                >
                    <div
                        v-for="item in offer.offered.items"
                        :key="item.type"
                        class="flex items-center gap-2"
                        :title="COST_ITEMS_INFO[item.type].name"
                    >
                        <img
                            :src="COST_ITEMS_INFO[item.type].icon"
                            class="w-10 h-10"
                        />
                        <span class="font-semibold text-base text-slate-200">{{
                            item.amount
                        }}</span>
                    </div>
                </div>
                <div v-else class="text-slate-500">-</div>
            </div>

            <div class="flex-shrink-0">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke-width="1.5"
                    stroke="currentColor"
                    class="w-7 h-7 text-slate-400 mx-3"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M7.5 21 3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5"
                    />
                </svg>
            </div>

            <div class="flex-1 flex flex-col gap-2.5 items-center">
                <div
                    v-if="offer.requested.items.length > 0"
                    class="flex flex-col items-center gap-2"
                >
                    <div
                        v-for="item in offer.requested.items"
                        :key="item.type"
                        class="flex items-center gap-2"
                        :title="COST_ITEMS_INFO[item.type].name"
                    >
                        <img
                            :src="COST_ITEMS_INFO[item.type].icon"
                            class="w-10 h-10"
                        />
                        <span class="font-semibold text-base text-slate-200">{{
                            item.amount
                        }}</span>
                    </div>
                </div>
                <div v-else class="text-slate-500">-</div>
            </div>
        </div>

        <div
            class="pt-3 border-t border-slate-700 flex items-center justify-end"
        >
            <button
                v-if="!isMine && offer.acceptable"
                class="px-4 py-1.5 rounded-md text-white text-sm font-bold transition-colors bg-green-500 hover:bg-green-400"
                @click="$emit('accept', offer.id)"
            >
                انجام معامله
            </button>
            <button
                v-if="isMine"
                class="px-4 py-1.5 rounded-md text-white text-sm font-bold transition-colors bg-red-500 hover:bg-red-400"
                @click="$emit('delete', offer.id)"
            >
                حذف
            </button>
        </div>
    </div>
</template>

<script setup>
import { COST_ITEMS_INFO } from '@/services/cost';

defineProps({
    offer: Object,
    isMine: Boolean,
});

defineEmits(['accept', 'delete']);
</script>

<style scoped>
/* No styles needed here*/
</style>
