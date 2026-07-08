<template>
    <VueFinalModal
        class="flex justify-center items-center"
        content-class="panel w-full md:w-1/2 max-h-[80vh]"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div class="panel-header">
            <h1 class="panel-title">
                {{ glossary.casino }}
            </h1>
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

        <div
            class="panel-body !gap-6 justify-center items-center min-h-[200px]"
        >
            <div v-if="isLoading" class="text-content-muted text-lg">
                در حال بارگذاری اطلاعات بورس...
            </div>

            <div v-else-if="loadError" class="w-full text-center">
                <h2 class="text-3xl font-bold text-danger-content">
                    بورس در دسترس نیست
                </h2>
                <p class="text-content-muted mt-4 text-lg">
                    در حال حاضر امکان اتصال به بورس وجود ندارد. لطفاً بعداً
                    دوباره تلاش کنید.
                </p>
                <button class="panel-btn mt-6" @click="doInvestCheck">
                    تلاش دوباره
                </button>
            </div>

            <div v-else-if="checkResult" class="w-full text-center">
                <div v-if="countdown" class="text-accent-strong text-lg mb-2">
                    زمان باقی‌مانده: {{ countdown }}
                </div>
                <div
                    v-if="checkResult.feasible"
                    class="w-full max-w-md mx-auto"
                >
                    <p
                        class="text-s text-content-muted mb-2 text-justify preserve-lines"
                    >
                        {{ checkResult.session.text }}
                    </p>
                    <div class="space-y-4">
                        <div>
                            <label
                                class="block text-sm text-content-muted mb-2 flex items-center justify-center gap-1"
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
                                class="w-full p-3 text-lg text-center text-content bg-surface-alt rounded-lg border-2 border-border focus:border-accent focus:ring-0 outline-none transition-colors"
                                :placeholder="`مبلغ به ${glossary.coin}`"
                            />
                            <p
                                class="text-xs text-content-subtle mt-2 flex items-center justify-center gap-1"
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
                            class="panel-btn w-full"
                            @click="handleInvest"
                        >
                            <span v-if="isInvesting">در حال ثبت...</span>
                            <span v-else>سرمایه‌گذاری کن</span>
                        </button>
                    </div>
                </div>
                <div v-else>
                    <div v-if="checkResult.investments?.length > 0">
                        <p class="text-lg text-content-muted">
                            شما در این دوره سرمایه‌گذاری کرده‌اید.
                        </p>
                        <div
                            v-for="(coins, index) in checkResult.investments"
                            :key="index"
                            class="mt-4 text-3xl font-bold text-accent-strong bg-black/20 py-4 rounded-lg flex items-center justify-center gap-2"
                        >
                            <span>{{ coins.coin }}</span>
                            <img
                                :src="COST_ITEMS_INFO.coin.icon"
                                class="w-8 h-8"
                            />
                        </div>
                    </div>
                    <div v-else>
                        <h2 class="text-3xl font-bold text-danger-content">
                            بورس بسته است
                        </h2>
                        <p class="text-content-muted mt-4 text-lg">
                            {{ checkResult.reason }}
                        </p>
                    </div>
                </div>
            </div>
        </div>
    </VueFinalModal>
</template>

<script setup>
import { VueFinalModal } from 'vue-final-modal';
import { glossary } from '@/services/glossary.js';
import { investCheck, invest } from '@/services/api/index.js';
import { onMounted, onUnmounted, ref, watch } from 'vue';
import { useToast } from 'vue-toastification';
import { COST_ITEMS_INFO } from '@/services/cost.js';
import { logger } from '@/services/logger.js';

const emit = defineEmits(['close']);

const isLoading = ref(true);
const isInvesting = ref(false);
const checkResult = ref(null);
const loadError = ref(false);
const investAmount = ref(0);
const toast = useToast();
const countdown = ref('');
let countdownInterval = null;

watch(investAmount, newValue => {
    if (checkResult.value && newValue > checkResult.value.maxCoin) {
        investAmount.value = checkResult.value.maxCoin;
    }
    if (newValue < 0) {
        investAmount.value = 0;
    }
});

function updateCountdown() {
    if (!checkResult.value?.session?.endAt) {
        countdown.value = '';
        return;
    }

    const now = new Date().getTime();
    const endTime = new Date(
        parseInt(checkResult.value.session.endAt)
    ).getTime();
    const distance = endTime - now;

    if (distance < 0) {
        clearInterval(countdownInterval);
        countdown.value = 'پایان یافته';
        doInvestCheck();
        return;
    }

    const hours = Math.floor(
        (distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)
    );
    const minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
    const seconds = Math.floor((distance % (1000 * 60)) / 1000);

    countdown.value = `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
}

watch(
    checkResult,
    newResult => {
        clearInterval(countdownInterval);
        if (newResult?.session?.endAt) {
            updateCountdown();
            countdownInterval = setInterval(updateCountdown, 1000);
        }
    },
    { immediate: true }
);

async function doInvestCheck() {
    isLoading.value = true;
    loadError.value = false;
    try {
        const result = await investCheck();
        checkResult.value = result;
    } catch (error) {
        loadError.value = true;
        toast.error(error.message || 'خطا در دریافت اطلاعات بورس');
        logger.error('Error calling investCheck:', error);
    } finally {
        isLoading.value = false;
    }
}

async function handleInvest() {
    if (isInvesting.value) return;
    isInvesting.value = true;
    try {
        await invest(checkResult.value.session.id, investAmount.value);
        toast.success('سرمایه‌گذاری شما با موفقیت ثبت شد.');
        await doInvestCheck();
    } catch (error) {
        toast.error(error.message || 'خطا در سرمایه‌گذاری');
        logger.error('Error calling invest:', error);
    } finally {
        isInvesting.value = false;
    }
}

function handleClose() {
    emit('close');
}

onMounted(doInvestCheck);

onUnmounted(() => {
    clearInterval(countdownInterval);
});
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
