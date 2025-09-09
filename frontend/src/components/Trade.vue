<template>
    <VueFinalModal
        class="flex justify-center items-center"
        content-class="flex flex-col w-full md:w-1/2 mx-4 p-6 
                       bg-[#5C3A21] border-4 border-[#3E2A17] 
                       rounded-xl shadow-xl space-y-4 text-amber-200"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div
            class="flex items-center justify-between border-b-2 border-[#3E2A17] pb-2 mb-4"
        >
            <h1 class="text-xl font-semibold text-amber-200">
                ساخت معامله جدید
            </h1>
            <button
                class="p-1 rounded-full hover:bg-[#3E2A17]"
                @click="handleClose"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-6 w-6 text-amber-200"
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

        <div class="flex flex-row justify-around items-center">
            <div class="flex flex-col justify-between items-center space-y-2">
                <div>داد</div>
                <div
                    v-for="tradable in tradables"
                    :key="tradable"
                    class="flex flex-row-reverse"
                >
                    <img
                        :src="getIconByType(tradable)"
                        :alt="tradable + ' Icon'"
                        class="w-5 h-5"
                    />
                    <input
                        v-model="offered[tradable]"
                        type="number"
                        class="w-10"
                    />
                </div>
            </div>

            <div class="flex flex-col justify-between items-center space-y-2">
                <div>ستد</div>
                <div
                    v-for="tradable in tradables"
                    :key="tradable"
                    class="flex flex-row-reverse"
                >
                    <img
                        :src="getIconByType(tradable)"
                        :alt="tradable + ' Icon'"
                        class="w-5 h-5"
                    />
                    <input
                        v-model="requested[tradable]"
                        type="number"
                        class="w-10"
                    />
                </div>
            </div>
        </div>

        <button @pointerdown="handleSubmit">ثبت معامله</button>
    </VueFinalModal>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { VueFinalModal } from 'vue-final-modal';
import { makeTradeOffer } from '../services/api';

const props = defineProps({
    player: Object,
    username: String,
    tradables: Array,
});

const offered = ref([]);
const requested = ref([]);
const emit = defineEmits(['close']);

function handleClose() {
    emit('close');
}

const getIconByType = type => {
    switch (type) {
        case 'blueKey':
            return '/images/icons/blueKeys.png';
        case 'goldenKey':
            return '/images/icons/goldenKeys.png';
        case 'redKey':
            return '/images/icons/redKeys.png';
    }
    return '/images/icons/' + type + '.png';
};

function buildPayload(offered, requested) {
    const toItems = obj =>
        Object.entries(obj).map(([type, amount]) => ({
            type,
            amount: Number(amount),
        }));

    return {
        offered: { items: toItems(offered) },
        requested: { items: toItems(requested) },
    };
}

async function handleSubmit() {
    const payload = buildPayload(offered.value, requested.value);
    try {
        await makeTradeOffer(payload.offered, payload.requested);
        handleClose();
    } catch (err) {
        console.error('Failed to make trade offer:', err);
    }
}

onMounted(() => {
    props.tradables.forEach(tradable => {
        offered.value[tradable] = 0;
        requested.value[tradable] = 0;
    });
});
</script>
