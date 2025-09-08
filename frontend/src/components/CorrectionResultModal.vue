<template>
    <VueFinalModal
        class="flex justify-center items-center"
        content-class="flex flex-col w-full md:w-1/3 mx-4 p-6 
                       rounded-xl shadow-xl space-y-4"
        :content-style="{
            backgroundColor: isCorrect ? '#166534' : '#991B1B',
            borderColor: isCorrect ? '#15803D' : '#B91C1C',
            borderWidth: '4px',
        }"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div
            class="flex items-center justify-between pb-2 mb-4"
            :style="{
                borderBottom: `2px solid ${isCorrect ? '#15803D' : '#B91C1C'}`,
            }"
        >
            <h1 class="text-2xl font-bold" :class="textColor">
                {{ title }}
            </h1>
            <button
                class="p-1 rounded-full"
                :class="hoverBgColor"
                @click="emit('close')"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-6 w-6"
                    :class="textColor"
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

        <div class="text-center">
            <p class="text-lg" :class="textColor">{{ message }}</p>
        </div>

        <div v-if="isCorrect && rewards.length > 0" class="pt-4">
            <h2
                class="text-lg font-semibold text-center mb-3"
                :class="textColor"
            >
                جوایز شما:
            </h2>
            <div class="flex flex-col items-center space-y-3">
                <div
                    v-for="reward in rewards"
                    :key="reward.name"
                    class="flex items-center p-2 rounded-lg w-full"
                    :style="{
                        backgroundColor: isCorrect
                            ? 'rgba(22, 163, 74, 0.5)'
                            : '',
                    }"
                >
                    <img :src="reward.icon" class="w-8 h-8 ml-3" alt="" />
                    <span class="font-medium" :class="textColor"
                        >{{ reward.name }}: +{{ reward.amount }}</span
                    >
                </div>
            </div>
        </div>

        <div class="flex justify-center pt-4">
            <button
                class="px-6 py-2 rounded-lg font-bold text-white transition-transform duration-200 hover:scale-105"
                :class="buttonBgColor"
                @click="emit('close')"
            >
                فهمیدم
            </button>
        </div>
    </VueFinalModal>
</template>

<script setup>
import { computed } from 'vue';
import { VueFinalModal } from 'vue-final-modal';

const props = defineProps({
    isCorrect: {
        type: Boolean,
        required: true,
    },
    rewards: {
        type: Array,
        default: () => [],
    },
});

const emit = defineEmits(['close']);

const title = computed(() =>
    props.isCorrect ? 'پاسخ صحیح بود!' : 'پاسخ غلط بود!'
);
const message = computed(() =>
    props.isCorrect
        ? 'آفرین! با موفقیت به این سوال پاسخ دادی.'
        : 'اشکالی نداره، دوباره تلاش کن!'
);

const textColor = computed(() =>
    props.isCorrect ? 'text-green-100' : 'text-red-100'
);
const hoverBgColor = computed(() =>
    props.isCorrect ? 'hover:bg-green-800' : 'hover:bg-red-800'
);
const buttonBgColor = computed(() =>
    props.isCorrect ? 'bg-green-600' : 'bg-red-600'
);
</script>
