<template>
    <VueFinalModal
        class="flex justify-center items-center"
        content-class="flex flex-col w-full md:w-1/2 mx-4 p-6 
                   bg-[#205647] border-4 border-[#508677] 
                   rounded-xl shadow-xl space-y-4"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div
            class="flex items-center justify-between border-b-2 border-[#508677] pb-2 mb-4"
        >
            <h1 class="text-xl font-semibold text-amber-200">
                {{ glossary.casino }}
            </h1>
            <button
                class="p-1 rounded-full hover:bg-[#3E1A17]"
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

        <div
            class="w-full flex flex-col justify-center items-center space-y-6 min-h-[200px]"
        >
            <!-- <div v-if="isLoading" class="text-amber-200 text-lg">
                در حال بارگذاری اطلاعات بورس...
            </div>

            <div v-else-if="checkResult" class="w-full text-center">
                <div v-if="countdown" class="text-amber-200 text-lg mb-2">
                    زمان باقی‌مانده: {{ countdown }}
                </div>
                <div
                    v-if="checkResult.feasible"
                    class="w-full max-w-md mx-auto"
                >
                    <p
                        class="text-s text-amber-100 mb-2 text-justify preserve-lines"
                    >
                        {{ checkResult.session.text }}
                    </p>
                    <div class="space-y-4">
                        <div>
                            <label
                                class="block text-sm text-amber-200 mb-2 flex items-center justify-center gap-1"
                            >
                                <span>مبلغ سرمایه‌گذاری</span>
                                <img
                                    :src="COST_ITEMS_INFO.coin.icon"
                                    class="w-8 h-8"
                                />
                            </label>
                            <input
                                v-model.number="investAmount"
                                type="number"
                                :max="checkResult.maxCoin"
                                min="0"
                                class="w-full p-3 text-lg text-center text-gray-100 bg-slate-800/70 rounded-lg border-2 border-slate-600 focus:border-cyan-500 focus:ring-0 outline-none transition-colors"
                                :placeholder="`مبلغ به ${glossary.coin}`"
                            />
                            <p
                                class="text-xs text-gray-400 mt-2 flex items-center justify-center gap-1"
                            >
                                <span>حداکثر:</span>
                                <span>{{
                                    checkResult.maxCoin.toLocaleString()
                                }}</span>
                                <img
                                    :src="COST_ITEMS_INFO.coin.icon"
                                    class="w-8 h-8"
                                />
                            </p>
                        </div>
                        <button
                            :disabled="
                                isInvesting ||
                                investAmount <= 0 ||
                                investAmount > checkResult.maxCoin
                            "
                            class="w-full px-6 py-3 text-lg font-semibold text-white bg-green-600 rounded-lg transition-all disabled:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-green-500"
                            @click="handleInvest"
                        >
                            <span v-if="isInvesting">در حال ثبت...</span>
                            <span v-else>سرمایه‌گذاری کن</span>
                        </button>
                    </div>
                </div>
                <div v-else>
                    <div v-if="checkResult.investments.length > 0">
                        <p class="text-lg text-gray-200">
                            شما در این دوره سرمایه‌گذاری کرده‌اید.
                        </p>
                        <div
                            v-for="(coins, index) in checkResult.investments"
                            :key="index"
                            class="mt-4 text-3xl font-bold text-amber-300 bg-black/20 py-4 rounded-lg flex items-center justify-center gap-2"
                        >
                            <span>{{ coins.coin }}</span>
                            <img
                                :src="COST_ITEMS_INFO.coin.icon"
                                class="w-8 h-8"
                            />
                        </div>
                    </div>
                    <div v-else>
                        <h2 class="text-3xl font-bold text-red-400">
                            بورس بسته است
                        </h2>
                        <p class="text-amber-200 mt-4 text-lg">
                            {{ checkResult.reason }}
                        </p>
                    </div>
                </div>
            </div> -->
            <div
                class="w-full flex flex-col justify-center items-center space-y-4"
            >
                <div v-if="!isPastNoon" class="text-center">
                    <p class="text-lg text-gray-300">
                        زمان باقی‌مانده تا بازگشایی:
                    </p>
                    <div
                        class="text-4xl font-bold text-amber-400 tracking-widest"
                    >
                        {{ hours }}:{{ minutes }}:{{ seconds }}
                    </div>
                </div>
                <div v-else class="text-center">
                    <h3 class="text-3xl font-bold text-green-400">
                        درحال بازگشایی
                    </h3>
                    <p class="text-gray-300 mt-2">
                        بازگشایی بورس تا دقایقی دیگر...
                    </p>
                </div>
            </div>
        </div>
    </VueFinalModal>
</template>

<script setup>
import { VueFinalModal } from 'vue-final-modal';
import { glossary } from '@/services/glossary.js';
// import { investCheck, invest } from '@/services/api/index.js';
// import { onMounted, onUnmounted, ref, watch } from 'vue';
// import { useToast } from 'vue-toastification';
// import { COST_ITEMS_INFO } from '@/services/cost.js';
import { useCountdownToNoon } from '@/composables/useCountdownToNoon.js';

const emit = defineEmits(['close']);
const { hours, minutes, seconds, isPastNoon } = useCountdownToNoon();

// const isLoading = ref(true);
// const isInvesting = ref(false);
// const checkResult = ref(null);
// const investAmount = ref(0);
// const toast = useToast();
// const countdown = ref('');
// let countdownInterval = null;

// watch(investAmount, newValue => {
//     if (checkResult.value && newValue > checkResult.value.maxCoin) {
//         investAmount.value = checkResult.value.maxCoin;
//     }
//     if (newValue < 0) {
//         investAmount.value = 0;
//     }
// });

// function updateCountdown() {
//     if (!checkResult.value?.session?.endAt) {
//         countdown.value = '';
//         return;
//     }

//     const now = new Date().getTime();
//     const endTime = new Date(
//         parseInt(checkResult.value.session.endAt)
//     ).getTime();
//     const distance = endTime - now;

//     if (distance < 0) {
//         clearInterval(countdownInterval);
//         countdown.value = 'پایان یافته';
//         doInvestCheck();
//         return;
//     }

//     const hours = Math.floor(
//         (distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)
//     );
//     const minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
//     const seconds = Math.floor((distance % (1000 * 60)) / 1000);

//     countdown.value = `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
// }

// watch(
//     checkResult,
//     newResult => {
//         clearInterval(countdownInterval);
//         if (newResult?.session?.endAt) {
//             updateCountdown();
//             countdownInterval = setInterval(updateCountdown, 1000);
//         }
//     },
//     { immediate: true }
// );

// async function doInvestCheck() {
//     isLoading.value = true;
//     try {
//         const result = await investCheck();
//         checkResult.value = result;
//     } catch (error) {
//         toast.error(error.message || 'خطا در دریافت اطلاعات بورس');
//         console.error('Error calling investCheck:', error);
//     } finally {
//         isLoading.value = false;
//     }
// }

// async function handleInvest() {
//     if (isInvesting.value) return;
//     isInvesting.value = true;
//     try {
//         await invest(checkResult.value.session.id, investAmount.value);
//         toast.success('سرمایه‌گذاری شما با موفقیت ثبت شد.');
//         await doInvestCheck(); // Refresh the state
//     } catch (error) {
//         toast.error(error.message || 'خطا در سرمایه‌گذاری');
//         console.error('Error calling invest:', error);
//     } finally {
//         isInvesting.value = false;
//     }
// }

function handleClose() {
    emit('close');
}

// onMounted(doInvestCheck);

// onUnmounted(() => {
//     clearInterval(countdownInterval);
// });
</script>

<style>
.preserve-lines {
    white-space: pre-wrap;
}

@keyframes init {
    0% {
        transform: scale(0);
    }
    40% {
        transform: scale(1.1);
    }
    60% {
        transform: scale(0.9);
    }
    80% {
        transform: scale(1.05);
    }
    100% {
        transform: scale(1);
    }
}

@keyframes init-sign-move {
    100% {
        transform: rotateZ(3deg);
    }
}

@keyframes sign-move {
    0% {
        transform: rotateZ(3deg);
    }
    50% {
        transform: rotateZ(-3deg);
    }
    100% {
        transform: rotateZ(3deg);
    }
}
</style>
