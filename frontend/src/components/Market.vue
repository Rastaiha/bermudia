<template>
    <VueFinalModal
        class="flex justify-center items-center"
        content-class="flex flex-col w-full h-full md:w-1/2 mx-4 p-6 
                       bg-[#5C3A21] border-4 border-[#3E2A17] 
                       rounded-xl shadow-xl space-y-4 text-amber-200"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div
            class="flex items-center justify-between border-b-2 border-[#3E2A17] pb-2 mb-4"
        >
            <h1 class="text-xl font-semibold">بازارچه بنزوئیلا</h1>
            <button
                class="p-1 rounded-full hover:bg-[#3E2A17]"
                @click="handleClose"
            >
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

        <div class="flex justify-around">
            <button @pointerdown="isOffersYours = false">
                معاملات درخواستی
            </button>
            <button @pointerdown="isOffersYours = true">معاملات شما</button>
        </div>

        <div
            v-if="!isOffersYours"
            class="w-full flex flex-col items-center justify-between items-end space-y-2"
        >
            <div class="flex flex-wrap justify-around gap-y-4 pb-2">
                <div
                    v-for="(offer, index) in otherOffers"
                    :key="index"
                    class="relative w-48 h-32 flex flex-col justify-between items-center pb-1 pt-1 bg-gradient-to-b from-yellow-600 via-yellow-700 to-yellow-800 border-l-2 border-yellow-900 rounded-sm shadow-md"
                >
                    <div
                        class="flex flex-row justify-between items-center w-full"
                    >
                        <div class="w-2/5">
                            <CostlyButton
                                :on-click="() => {}"
                                :cost="offer.offered"
                                label="داد"
                                enabled
                                :loading="false"
                            >
                            </CostlyButton>
                        </div>
                        <div class="w-2/5">
                            <CostlyButton
                                :on-click="() => {}"
                                :cost="offer.requested"
                                label="ستد"
                                enabled
                                :loading="false"
                            >
                            </CostlyButton>
                        </div>
                    </div>
                    <div>{{ offer.createdAt }}</div>
                    <button
                        v-if="offer.acceptable"
                        @pointerdown="acceptTradeOffer(offer.id)"
                    >
                        انجام معامله
                    </button>
                    <div v-else></div>
                </div>
            </div>
        </div>

        <div
            v-else
            class="w-full flex flex-col items-center justify-between space-y-2"
        >
            <div class="flex flex-wrap justify-around gap-y-4 pb-2">
                <div
                    v-for="(offer, index) in myOffers"
                    :key="index"
                    class="relative w-48 h-32 flex flex-col justify-between items-center pb-1 pt-1 bg-gradient-to-b from-yellow-600 via-yellow-700 to-yellow-800 border-l-2 border-yellow-900 rounded-sm shadow-md"
                >
                    <div
                        class="flex flex-row justify-between items-center w-full"
                    >
                        <div class="w-2/5">
                            <CostlyButton
                                :on-click="() => {}"
                                :cost="offer.offered"
                                label="داد"
                                enabled
                                :loading="false"
                            >
                            </CostlyButton>
                        </div>
                        <div class="w-2/5">
                            <CostlyButton
                                :on-click="() => {}"
                                :cost="offer.requested"
                                label="ستد"
                                enabled
                                :loading="false"
                            >
                            </CostlyButton>
                        </div>
                    </div>
                    <div>{{ offer.createdAt }}</div>
                    <button @pointerdown="deleteTradeOffer(offer.id)">
                        حذف معامله
                    </button>
                </div>
            </div>
            <button
                class="transition-transform duration-200 hover:scale-110 pointer-events-auto"
                title="معامله جدید"
                @pointerdown="openTrade"
            >
                معامله جدید
            </button>
        </div>
    </VueFinalModal>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { VueFinalModal, useModal } from 'vue-final-modal';
import {
    acceptTradeOffer,
    deleteTradeOffer,
    getTradeOffers,
} from '../services/api';
import Trade from './Trade.vue';
import CostlyButton from './CostlyButton.vue';

const props = defineProps({
    player: Object,
    username: String,
});

const myOffers = ref([]);
const otherOffers = ref([]);
const tradables = ref([]);
const pagesLimit = ref(30);
const pageNumber = ref(0);
const isOffersYours = ref(false);
const emit = defineEmits(['close']);

const { open: openTrade, close: closeTrade } = useModal({
    component: Trade,
    attrs: {
        player: props.player,
        username: props.username,
        tradables: tradables,
        onClose() {
            closeTrade();
        },
    },
});

function handleClose() {
    emit('close');
}

onMounted(async () => {
    try {
        myOffers.value = await getTradeOffers(
            pageNumber.value,
            pagesLimit.value,
            'me'
        );
        otherOffers.value = await getTradeOffers(
            pageNumber.value,
            pagesLimit.value,
            'others'
        );
        tradables.value = ['coin', 'redKey', 'blueKey', 'goldenKey'];
    } catch (err) {
        console.error('Failed to load trade offers:', err);
    }
});
</script>
