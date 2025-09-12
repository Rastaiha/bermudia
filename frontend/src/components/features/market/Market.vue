<template>
    <VueFinalModal
        class="flex justify-center items-center"
        content-class="flex flex-col w-full h-full md:w-5/7 mx-4 p-6 
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
            <button
                class="p-1 rounded-[5px]"
                :class="
                    !isOffersYours
                        ? 'bg-[#fee685] text-[#5c3a21]'
                        : 'border-[5px] border-solid border-[#fee685]'
                "
                @pointerdown="isOffersYours = false"
            >
                معاملات درخواستی
            </button>
            <button
                class="p-1 rounded-[5px]"
                :class="
                    isOffersYours
                        ? 'bg-[#fee685] text-[#5c3a21]'
                        : 'border-[5px] border-solid border-[#fee685]'
                "
                @pointerdown="isOffersYours = true"
            >
                معاملات شما
            </button>
        </div>

        <div
            v-if="!isOffersYours"
            ref="otherOffersContainer"
            class="w-full flex flex-col items-center justify-between space-y-2 overflow-y-auto max-h-[fit-content]"
            @scroll="handleOtherScroll"
        >
            <div
                v-if="otherOffers.length > 0"
                class="w-full flex flex-wrap justify-around gap-y-4 pb-2"
            >
                <div
                    v-for="(offer, index) in otherOffers"
                    :key="index"
                    class="relative w-64 h-32 flex flex-col justify-between items-center p-2 bg-gradient-to-b from-yellow-600 via-yellow-700 to-yellow-800 border-l-2 border-yellow-900 rounded-sm shadow-md"
                >
                    <div
                        class="flex flex-row justify-between items-center w-full"
                    >
                        <div class="w-3/7">
                            <CostlyButton
                                :on-click="() => {}"
                                :cost="offer.offered"
                                label="داد"
                                enabled
                                :loading="false"
                            />
                        </div>
                        <div class="w-3/7">
                            <CostlyButton
                                :on-click="() => {}"
                                :cost="offer.requested"
                                label="ستد"
                                enabled
                                :loading="false"
                                background-color="#480202"
                            />
                        </div>
                    </div>
                    <div>{{ timeCommenter(offer.created_at) }}</div>
                    <button
                        v-if="offer.acceptable"
                        class="transition-transform duration-200 hover:scale-110 pointer-events-auto p-1 rounded-[5px] bg-[#fee685] text-[#5c3a21]"
                        @pointerdown="acceptTrade(offer)"
                    >
                        انجام معامله
                    </button>
                </div>
            </div>
            <div v-else class="text-center w-full">معامله‌ای یافت نشد.</div>
        </div>

        <div
            v-else
            class="w-full h-full flex flex-col items-center justify-between space-y-2"
        >
            <button
                class="transition-transform duration-200 hover:scale-110 pointer-events-auto p-1 rounded-[5px] bg-[#fee685] text-[#5c3a21]"
                title="معامله جدید"
                @click="createTrade"
            >
                معامله جدید
            </button>
            <div
                ref="myOffersContainer"
                class="w-full flex flex-col items-center justify-between space-y-2 overflow-y-auto max-h-[fit-content]"
                @scroll="handleMyScroll"
            >
                <div
                    v-if="myOffers.length > 0"
                    class="w-full flex flex-wrap justify-around gap-y-4 pb-2"
                >
                    <div
                        v-for="(offer, index) in myOffers"
                        :key="index"
                        class="relative w-64 h-32 flex flex-col justify-between items-center p-2 bg-gradient-to-b from-yellow-600 via-yellow-700 to-yellow-800 border-l-2 border-yellow-900 rounded-sm shadow-md"
                    >
                        <div
                            class="flex flex-row justify-between items-center w-full"
                        >
                            <div class="w-3/7">
                                <CostlyButton
                                    :on-click="() => {}"
                                    :cost="offer.offered"
                                    label="داد"
                                    enabled
                                    :loading="false"
                                    background-color="#480202"
                                >
                                </CostlyButton>
                            </div>
                            <div class="w-3/7">
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
                        <div>{{ timeCommenter(offer.created_at) }}</div>
                        <button
                            class="transition-transform duration-200 hover:scale-110 pointer-events-auto p-1 rounded-[5px] bg-[#fee685] text-[#5c3a21]"
                            @pointerdown="deleteTrade(offer)"
                        >
                            حذف معامله
                        </button>
                    </div>
                </div>
                <div v-else class="text-center w-full">معامله‌ای یافت نشد.</div>
            </div>
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
} from '@/services/api/index.js';
import { useToast } from 'vue-toastification';
import { useNow } from '@/composables/useNow.js';
import Trade from '@/components/features/market/Trade.vue';
import CostlyButton from '@/components/common/CostlyButton.vue';
import { makeTradeOfferCheck } from '../../../services/api';

const props = defineProps({
    player: Object,
    username: String,
});

const myOffers = ref([]);
const otherOffers = ref([]);
const tradables = ref([]);
const myPagesLimit = ref(12);
const myPageNumber = ref(0);
const myPageIsLoading = ref(true);
const myPageIsLoadedAll = ref(false);
const otherPagesLimit = ref(12);
const otherPageNumber = ref(0);
const otherPageIsLoading = ref(true);
const otherPageIsLoadedAll = ref(false);
const isOffersYours = ref(false);
const myOffersContainer = ref(null);
const otherOffersContainer = ref(null);
const toast = useToast();
const { now } = useNow(60000);
const emit = defineEmits(['close']);

const createTrade = async () => {
    try {
        const makeOfferCheck = await makeTradeOfferCheck();
        if (makeOfferCheck.feasible) {
            tradables.value.splice(0, tradables.value.length);
            makeOfferCheck.tradableItems.items.array.forEach(element => {
                tradables.value.push(element.type);
            });
            openTrade();
        } else {
            toast.error(makeOfferCheck.reason || 'مشکل در ثبت معامله');
        }
    } catch (err) {
        toast.error(err.message || 'مشکل در ثبت معامله');
    }
};

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

async function loadMoreMyOffers() {
    myPageIsLoading.value = true;
    const newOffers = await getTradeOffers(
        myPageNumber.value + 1,
        myPagesLimit.value,
        'me'
    );
    if (newOffers.length) {
        myPageNumber.value += 1;
        myOffers.value.push(...newOffers);
        myPageIsLoading.value = false;
    } else {
        myPageIsLoadedAll.value = true;
    }
}

async function loadMoreOtherOffers() {
    otherPageIsLoading.value = true;
    const newOffers = await getTradeOffers(
        otherPageNumber.value + 1,
        otherPagesLimit.value,
        'others'
    );
    if (newOffers.length) {
        otherPageNumber.value += 1;
        otherOffers.value.push(...newOffers);
        otherPageIsLoading.value = false;
    } else {
        otherPageIsLoadedAll.value = true;
    }
}

function handleMyScroll() {
    if (myPageIsLoadedAll.value || myPageIsLoading.value) return;
    const container = myOffersContainer.value;
    if (!container) return;
    if (
        container.scrollTop + container.clientHeight >=
        container.scrollHeight - 50
    ) {
        loadMoreMyOffers();
    }
}

function handleOtherScroll() {
    if (otherPageIsLoadedAll.value || otherPageIsLoading.value) return;
    const container = otherOffersContainer.value;
    if (!container) return;
    if (
        container.scrollTop + container.clientHeight >=
        container.scrollHeight - 50
    ) {
        loadMoreOtherOffers();
    }
}
const acceptTrade = offer => {
    try {
        acceptTradeOffer(offer.id);
        toast.success('معامله جوش خورد.');
    } catch (err) {
        toast.error(err.message || 'در حین تایید معامله خطایی رخ داد');
    }
};

const deleteTrade = offer => {
    try {
        deleteTradeOffer(offer.id);
        toast.success('معامله حذف شد.');
    } catch (err) {
        toast.error(err.message || 'در حین حذف معامله خطایی رخ داد');
    }
};

const timeCommenter = time => {
    let diff = now.value - time;
    diff /= 1000;
    if (diff < 60) return 'ثانیه‌هایی پیش';
    diff = Math.floor(diff / 60);
    if (diff < 60) return diff + ' دقیقه پیش';
    return Math.floor(diff / 60) + ' ساعت پیش';
};

onMounted(async () => {
    try {
        myOffers.value = await getTradeOffers(
            myPageNumber.value,
            myPagesLimit.value,
            'me'
        );
        myPageIsLoading.value = false;
        otherOffers.value = await getTradeOffers(
            otherPageNumber.value,
            otherPagesLimit.value,
            'others'
        );
        otherPageIsLoading.value = false;
    } catch (err) {
        toast.error(err.message || 'در حین دریافت معاملات خطایی رخ داد');
    }
});
</script>
<style>
::-webkit-scrollbar {
    width: 10px;
    height: 10px;
}

::-webkit-scrollbar-track {
    background: #3e2a17;
    border-radius: 25px;
}

::-webkit-scrollbar-thumb {
    background: #fee685;
    border-radius: 25px;
    border: 2px solid #3e2a17;
}

::-webkit-scrollbar-thumb:hover {
    background: #fcd34d;
}

* {
    scrollbar-color: #fee685 #3e2a17;
}
</style>
