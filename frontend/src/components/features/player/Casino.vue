<template>
    <VueFinalModal
        class="flex justify-center items-center"
        content-class="flex flex-col w-full md:w-1/2 mx-4 p-6 
                   bg-[#5C2A21] border-4 border-[#3E1A17] 
                   rounded-xl shadow-xl space-y-4"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div
            class="flex items-center justify-between border-b-2 border-[#3E1A17] pb-2 mb-4"
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

        <div class="w-full flex flex-col justify-between items-end space-y-2">
            <div
                class="w-full h-full flex items-center justify-center text-[3.75vmin]"
            >
                <div
                    class="relative w-[18.75vmin] h-[13.75vmin] m-5 shrink-0 origin-[center_0.625vmin] animate-[init_1000ms_forwards,init-sign-move_1000ms_ease-out_1000ms,sign-move_3000ms_2000ms_infinite]"
                >
                    <div
                        class="absolute bottom-0 w-full h-[8.75vmin] flex items-center justify-center font-montserrat font-bold text-white bg-[#ff5625] rounded-[1vmin] shadow-[0_-0.015em_0_#992100]"
                    >
                        Closed
                    </div>

                    <div
                        class="absolute left-1/2 w-[7.5vmin] h-[7.5vmin] border-[0.225vmin] border-[#893d00] border-b-0 border-r-0 origin-top-left rotate-45"
                    ></div>

                    <div
                        class="absolute top-0 left-1/2 -ml-[0.625vmin] w-[1.25vmin] h-[1.25vmin] rounded-full bg-[#9f9f9f]"
                    ></div>
                    <div
                        class="absolute top-[5.375vmin] left-[3.25vmin] w-[1.25vmin] h-[1.25vmin] rounded-full bg-[#e63300]"
                    ></div>
                    <div
                        class="absolute top-[5.375vmin] right-[3.25vmin] w-[1.25vmin] h-[1.25vmin] rounded-full bg-[#e63300]"
                    ></div>
                </div>
            </div>
        </div>
    </VueFinalModal>
</template>

<script setup>
import { VueFinalModal } from 'vue-final-modal';
import { glossary } from '@/services/glossary.js';

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
