<template>
    <VueFinalModal
        class="flex justify-center items-center py-5"
        content-class="flex flex-col h-full w-11/12 md:w-4/5 lg:w-3/5 max-w-4xl bg-slate-800 border border-slate-600 rounded-2xl shadow-2xl text-slate-100"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div
            class="flex-shrink-0 flex items-center justify-between p-6 pb-4 border-b border-slate-600"
        >
            <h1 class="text-2xl font-bold text-amber-400">
                {{ glossary.benzuelaMarketplace }}
            </h1>
            <button
                class="p-2 rounded-full text-slate-400 hover:bg-slate-700 hover:text-white transition-colors"
                aria-label="Close modal"
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

        <div
            class="flex-shrink-0 flex justify-center items-center p-4 border-b border-slate-700 bg-slate-800/50"
        >
            <div class="flex p-1 space-x-1 reverse bg-slate-900/70 rounded-xl">
                <button
                    :class="tabClass(!isOffersYours)"
                    @click="isOffersYours = false"
                >
                    درخواست‌های دیگران
                </button>
                <button
                    :class="tabClass(isOffersYours)"
                    @click="isOffersYours = true"
                >
                    درخواست‌های شما
                </button>
            </div>
        </div>

        <div
            class="flex-1 overflow-y-auto relative"
            @scroll="handleOtherScroll"
        >
            <div
                v-if="!isOffersYours"
                ref="otherOffersContainer"
                class="p-4 md:p-6"
            >
                <div
                    v-if="otherOffers.length > 0"
                    class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
                >
                    <TradeOfferCard
                        v-for="offer in otherOffers"
                        :key="offer.id"
                        :offer="offer"
                        :is-mine="false"
                        @accept="acceptTrade"
                    />
                </div>
                <div
                    v-else-if="!otherPageIsLoading"
                    class="text-center text-slate-400 py-8"
                >
                    معامله‌ای برای نمایش وجود ندارد.
                </div>
            </div>

            <div v-else class="flex flex-col h-full">
                <div class="p-4 flex-shrink-0">
                    <button
                        class="w-full md:w-auto md:px-6 flex-shrink-0 mx-auto flex justify-center items-center p-3 rounded-lg bg-amber-500 text-slate-900 font-bold text-lg transition-all duration-300 hover:bg-amber-400 hover:shadow-lg hover:scale-105"
                        @click="createTrade"
                    >
                        + ساخت معامله جدید
                    </button>
                </div>
                <div
                    ref="myOffersContainer"
                    class="flex-1 overflow-y-auto p-4 md:p-6 pt-0"
                    @scroll="handleMyScroll"
                >
                    <div
                        v-if="myOffers.length > 0"
                        class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
                    >
                        <TradeOfferCard
                            v-for="offer in myOffers"
                            :key="offer.id"
                            :offer="offer"
                            is-mine
                            @delete="deleteTrade"
                        />
                    </div>
                    <div
                        v-else-if="!myPageIsLoading"
                        class="text-center text-slate-400 py-8"
                    >
                        شما هنوز معامله‌ای ثبت نکرده‌اید.
                    </div>
                </div>
            </div>
        </div>
    </VueFinalModal>
</template>

<script setup>
import { watch, ref } from 'vue';
import { VueFinalModal, useModal } from 'vue-final-modal';
import {
    acceptTradeOffer,
    deleteTradeOffer,
    getTradeOffers,
    makeTradeOfferCheck,
} from '@/services/api/index.js';
import { useToast } from 'vue-toastification';
import { useMarketWebSocket } from '@/services/marketWebsocket.js';
import Trade from '@/components/features/market/Trade.vue';
import TradeOfferCard from '@/components/features/market/TradeOfferCard.vue';
import { glossary } from '@/services/glossary.js';

const props = defineProps({
    player: Object,
});

const myOffers = ref([]);
const otherOffers = ref([]);
const tradables = ref({});
const myPagesLimit = ref(12);
const mySyncTrade = ref('');
const myPageIsLoading = ref(true);
const myPageIsLoadedAll = ref(false);
const otherPagesLimit = ref(100);
const otherSyncTrade = ref('');
const otherPageIsLoading = ref(true);
const otherPageIsLoadedAll = ref(false);
const isOffersYours = ref(false);
const myOffersContainer = ref(null);
const otherOffersContainer = ref(null);
const toast = useToast();
const emit = defineEmits(['close']);

const createTrade = async () => {
    try {
        const makeOfferCheck = await makeTradeOfferCheck();
        if (makeOfferCheck.feasible) {
            tradables.value = makeOfferCheck.tradableItems;
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
        tradables: tradables,
        onClose() {
            closeTrade();
        },
    },
});

function handleClose() {
    emit('close');
}

const tabClass = isActive => {
    return isActive
        ? 'w-full bg-amber-500 text-slate-900 rounded-lg px-4 py-2 text-sm font-bold transition-colors'
        : 'w-full text-slate-300 hover:bg-slate-700/50 rounded-lg px-4 py-2 text-sm font-bold transition-colors';
};

async function loadMoreMyOffers() {
    myPageIsLoading.value = true;
    const newOffers = await getTradeOffers(
        mySyncTrade.value,
        myPagesLimit.value,
        'me'
    );
    if (newOffers.length) {
        mySyncTrade.value = newOffers[newOffers.length - 1].created_at;
        myOffers.value.push(...newOffers);
        myPageIsLoading.value = false;
    } else {
        myPageIsLoadedAll.value = true;
    }
}

async function loadMoreOtherOffers() {
    otherPageIsLoading.value = true;
    const newOffers = await getTradeOffers(
        otherSyncTrade.value,
        otherPagesLimit.value,
        'others'
    );
    if (newOffers.length) {
        otherSyncTrade.value = newOffers[newOffers.length - 1].created_at;
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
const acceptTrade = async id => {
    try {
        await acceptTradeOffer(id);
        toast.success('معامله جوش خورد.');
    } catch (err) {
        toast.error(err.message || 'در حین تایید معامله خطایی رخ داد');
    }
};

const deleteTrade = async id => {
    try {
        await deleteTradeOffer(id);
        toast.success('معامله حذف شد.');
    } catch (err) {
        toast.error(err.message || 'در حین حذف معامله خطایی رخ داد');
    }
};

const loadMyOffers = watch(mySyncTrade, async newVal => {
    if (newVal.length > 0) {
        try {
            myOffers.value = await getTradeOffers(
                mySyncTrade.value,
                myPagesLimit.value,
                'me'
            );
            if (myOffers.value.length)
                mySyncTrade.value =
                    myOffers.value[myOffers.value.length - 1].created_at;
            myPageIsLoading.value = false;
        } catch (err) {
            toast.error(err.message || 'در حین دریافت معاملات خطایی رخ داد');
        }
    }
    loadMyOffers();
});

const loadOtherOffers = watch(otherSyncTrade, async newVal => {
    if (newVal.length > 0) {
        try {
            otherOffers.value = await getTradeOffers(
                otherSyncTrade.value,
                otherPagesLimit.value,
                'others'
            );
            if (otherOffers.value.length)
                otherSyncTrade.value =
                    otherOffers.value[otherOffers.value.length - 1].created_at;
            otherPageIsLoading.value = false;
        } catch (err) {
            toast.error(err.message || 'در حین دریافت معاملات خطایی رخ داد');
        }
    }
    loadOtherOffers();
});

useMarketWebSocket(mySyncTrade, otherSyncTrade, myOffers, otherOffers);
</script>
<style>
::-webkit-scrollbar {
    width: 10px;
    height: 10px;
}
::-webkit-scrollbar-track {
    background: #1e293b;
    border-radius: 25px;
}
::-webkit-scrollbar-thumb {
    background: #f59e0b;
    border-radius: 25px;
    border: 2px solid #1e293b;
}
::-webkit-scrollbar-thumb:hover {
    background: #fcd34d;
}
* {
    scrollbar-color: #f59e0b #1e293b;
}
.reverse {
    direction: rtl;
}
</style>
