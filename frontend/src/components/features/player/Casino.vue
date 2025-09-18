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
            class="w-full flex flex-col justify-center items-center space-y-6 py-8"
        >
            <div v-if="isPastNoon" class="text-center">
                <h2 class="text-3xl font-bold text-amber-100">
                    بازگشایی به زودی...
                </h2>
                <p class="text-amber-200 mt-2">
                    بورس تا دقایقی دیگر باز می‌شود
                </p>
            </div>
            <div v-else class="text-center">
                <h2 class="text-2xl font-bold text-amber-200 mb-6">
                    زمان باقی‌مانده تا بازگشایی
                </h2>
                <div
                    class="flex flex-row-reverse items-center justify-center gap-4"
                >
                    <div class="flex flex-col items-center">
                        <div
                            class="text-6xl p-4 bg-black/20 rounded-lg shadow-inner text-amber-200"
                        >
                            {{ hours }}
                        </div>
                        <div class="text-sm text-amber-200 mt-2">ساعت</div>
                    </div>
                    <div class="text-6xl text-amber-200 pb-8">:</div>
                    <div class="flex flex-col items-center">
                        <div
                            class="text-6xl p-4 bg-black/20 rounded-lg shadow-inner text-amber-200"
                        >
                            {{ minutes }}
                        </div>
                        <div class="text-sm text-amber-200 mt-2">دقیقه</div>
                    </div>
                    <div class="text-6xl text-amber-200 pb-8">:</div>
                    <div class="flex flex-col items-center">
                        <div
                            class="text-6xl p-4 bg-black/20 rounded-lg shadow-inner text-amber-200"
                        >
                            {{ seconds }}
                        </div>
                        <div class="text-sm text-amber-200 mt-2">ثانیه</div>
                    </div>
                </div>
            </div>
        </div>
    </VueFinalModal>
</template>

<script setup>
import { VueFinalModal } from 'vue-final-modal';
import { glossary } from '@/services/glossary.js';
import { useCountdownToNoon } from '@/composables/useCountdownToNoon.js';

const { hours, minutes, seconds, isPastNoon } = useCountdownToNoon();

const emit = defineEmits(['close']);

function handleClose() {
    emit('close');
}
</script>

<style>
@import url('https://fonts.googleapis.com/css?family=Montserrat:400,400i,700');

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
