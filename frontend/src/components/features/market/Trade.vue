<template>
    <VueFinalModal
        class="flex justify-center items-center py-5"
        content-class="flex flex-col h-full w-11/12 md:w-3/4 lg:w-2/3 max-w-4xl bg-slate-800 border border-slate-600 rounded-2xl shadow-2xl text-slate-100"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div
            class="flex-shrink-0 flex items-center justify-between p-6 pb-4 border-b border-slate-600"
        >
            <h1 class="text-2xl font-bold text-amber-400">ساخت معامله جدید</h1>
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

        <div class="flex-1 overflow-y-auto p-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-8 items-start">
                <div class="flex flex-col space-y-4">
                    <h2
                        class="text-xl font-semibold text-center text-green-400"
                    >
                        شما می‌دهید
                    </h2>
                    <div
                        v-for="(entry, type) in offered"
                        :key="`offered-${type}`"
                        class="flex items-center gap-3 p-3 bg-slate-700/50 rounded-lg border border-transparent hover:border-slate-500 transition-colors"
                    >
                        <img
                            :src="COST_ITEMS_INFO[type].icon"
                            :alt="type"
                            class="w-10 h-10 flex-shrink-0"
                        />
                        <div class="flex-grow text-right">
                            <span class="font-semibold">{{
                                COST_ITEMS_INFO[type].name
                            }}</span>
                            <span class="text-xs text-slate-400 block"
                                >موجودی: {{ entry.max }}</span
                            >
                        </div>
                        <div class="flex-shrink-0">
                            <div class="flex items-center gap-1">
                                <button
                                    class="w-7 h-7 flex items-center justify-center bg-slate-700 text-amber-400 rounded-full font-bold text-lg transition-colors hover:bg-slate-600"
                                    @click="decrement(entry)"
                                >
                                    -
                                </button>
                                <input
                                    v-model.number="entry.current"
                                    type="number"
                                    class="trade-input w-16 bg-slate-900 border-2 border-slate-600 rounded-md text-center text-lg font-semibold p-2 transition-colors focus:outline-none focus:ring-2 focus:ring-amber-500 focus:border-amber-500"
                                    :name="`offered_${type}`"
                                    placeholder="0"
                                />
                                <button
                                    class="w-7 h-7 flex items-center justify-center bg-slate-700 text-amber-400 rounded-full font-bold text-lg transition-colors hover:bg-slate-600"
                                    @click="increment(entry)"
                                >
                                    +
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="flex flex-col space-y-4">
                    <h2 class="text-xl font-semibold text-center text-red-400">
                        شما می‌خواهید
                    </h2>
                    <div
                        v-for="(entry, type) in requested"
                        :key="`requested-${type}`"
                        class="flex items-center gap-3 p-3 bg-slate-700/50 rounded-lg border border-transparent hover:border-slate-500 transition-colors"
                    >
                        <img
                            :src="COST_ITEMS_INFO[type].icon"
                            :alt="type"
                            class="w-10 h-10 flex-shrink-0"
                        />
                        <span class="font-semibold flex-grow text-right">{{
                            COST_ITEMS_INFO[type].name
                        }}</span>
                        <div class="flex-shrink-0">
                            <div class="flex items-center gap-1">
                                <button
                                    class="w-7 h-7 flex items-center justify-center bg-slate-700 text-amber-400 rounded-full font-bold text-lg transition-colors hover:bg-slate-600"
                                    @click="decrement(entry)"
                                >
                                    -
                                </button>
                                <input
                                    v-model.number="entry.current"
                                    type="number"
                                    class="trade-input w-16 bg-slate-900 border-2 border-slate-600 rounded-md text-center text-lg font-semibold p-2 transition-colors focus:outline-none focus:ring-2 focus:ring-amber-500 focus:border-amber-500"
                                    :name="`requested_${type}`"
                                    placeholder="0"
                                />
                                <button
                                    class="w-7 h-7 flex items-center justify-center bg-slate-700 text-amber-400 rounded-full font-bold text-lg transition-colors hover:bg-slate-600"
                                    @click="increment(entry)"
                                >
                                    +
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="flex-shrink-0 p-6 pt-4 border-t border-slate-600">
            <button
                :disabled="isSubmitDisabled"
                class="w-full md:w-1/2 mx-auto flex justify-center items-center p-3 rounded-lg bg-amber-500 text-slate-900 font-bold text-lg transition-all duration-300 hover:bg-amber-400 hover:shadow-lg hover:scale-105 disabled:bg-slate-600 disabled:text-slate-400 disabled:cursor-not-allowed disabled:scale-100"
                @click="handleSubmit"
            >
                ثبت معامله
            </button>
        </div>
    </VueFinalModal>
</template>

<script setup>
import { onMounted, ref, watch, computed } from 'vue';
import { VueFinalModal } from 'vue-final-modal';
import { useToast } from 'vue-toastification';
import { makeTradeOffer } from '@/services/api/index.js';
import { COST_ITEMS_INFO } from '@/services/cost';

const props = defineProps({
    player: Object,
    tradables: Object,
});

const emit = defineEmits(['close']);
const offered = ref({});
const requested = ref({});
const toast = useToast();

const handleClose = () => emit('close');

const increment = entry => {
    entry.current = clamp(
        (Number(entry.current) || 0) + 1,
        entry.min,
        entry.max
    );
};

const decrement = entry => {
    entry.current = clamp(
        (Number(entry.current) || 0) - 1,
        entry.min,
        entry.max
    );
};

const buildPayload = () => {
    const toItems = obj =>
        Object.entries(obj)
            .map(([type, entry]) => ({
                type,
                amount: Number(entry.current) || 0,
            }))
            .filter(item => item.amount > 0);

    return {
        offered: { items: toItems(offered.value) },
        requested: { items: toItems(requested.value) },
    };
};

const handleSubmit = async () => {
    const payload = buildPayload();
    if (
        payload.offered.items.length === 0 &&
        payload.requested.items.length === 0
    ) {
        toast.error('باید حداقل یک آیتم برای داد و ستد مشخص کنید.');
        return;
    }

    try {
        await makeTradeOffer(payload.offered, payload.requested);
        toast.success('معامله با موفقیت ثبت شد.');
        handleClose();
    } catch (err) {
        toast.error(err.message || 'خطا در ثبت معامله');
    }
};

const isSubmitDisabled = computed(() => {
    const totalOffered = Object.values(offered.value).reduce(
        (sum, entry) => sum + (entry.current || 0),
        0
    );
    const totalRequested = Object.values(requested.value).reduce(
        (sum, entry) => sum + (entry.current || 0),
        0
    );
    return totalOffered === 0 || totalRequested === 0;
});

const clamp = (val, min, max) => Math.max(min, Math.min(val, max));

const watchAndClamp = source => {
    watch(
        source,
        newVal => {
            for (const type in newVal) {
                const entry = newVal[type];
                entry.current = clamp(
                    Number(entry.current) || 0,
                    entry.min,
                    entry.max
                );
            }
        },
        { deep: true }
    );
};

watchAndClamp(offered);
watchAndClamp(requested);

onMounted(() => {
    props.tradables.items.forEach(tradable => {
        const baseEntry = { min: 0, current: 0 };
        offered.value[tradable.type] = { ...baseEntry, max: tradable.amount };
        requested.value[tradable.type] = { ...baseEntry, max: Infinity };
    });
});
</script>

<style scoped>
.trade-input::-webkit-outer-spin-button,
.trade-input::-webkit-inner-spin-button {
    -webkit-appearance: none;
    margin: 0;
}
.trade-input[type='number'] {
    -moz-appearance: textfield;
}
</style>
