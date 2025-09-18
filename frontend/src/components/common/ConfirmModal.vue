<template>
    <VueFinalModal
        class="flex justify-center items-center"
        overlay-class="bg-black/50"
        content-class="flex flex-col max-w-md mx-4 p-6 rounded-xl shadow-lg space-y-4 transition-colors"
        content-transition="vfm-slide-up"
        overlay-transition="vfm-fade"
    >
        <div
            class="p-4 rounded-xl"
            :class="[
                'flex flex-col space-y-4',
                dark
                    ? 'bg-slate-900 text-white border-slate-700'
                    : 'bg-white text-gray-800',
            ]"
        >
            <div class="flex items-center justify-between">
                <h1 class="text-xl font-semibold">{{ title }}</h1>
                <button
                    class="p-1 rounded-full hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
                    @click="handleClose"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-6 w-6"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        :class="dark ? 'text-white' : 'text-gray-600'"
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

            <div>
                <slot name="content">
                    <p class="text-gray-600 dark:text-gray-300">
                        اینجا می‌توانید محتوای مودال را قرار دهید.
                    </p>
                </slot>
            </div>

            <div class="flex justify-end gap-4 mt-4">
                <button
                    class="px-4 py-2 text-sm font-medium rounded-lg transition-colors"
                    :class="
                        dark
                            ? 'text-gray-200 bg-gray-700 hover:bg-gray-600'
                            : 'text-gray-700 bg-gray-100 hover:bg-gray-200'
                    "
                    @click="emit('cancel')"
                >
                    انصراف
                </button>
                <button
                    class="px-4 py-2 text-sm font-medium rounded-lg transition-colors"
                    :class="
                        dark
                            ? 'text-white bg-blue-600 hover:bg-blue-700'
                            : 'text-white bg-blue-600 hover:bg-blue-700'
                    "
                    @click="emit('confirm')"
                >
                    تایید
                </button>
            </div>
        </div>
    </VueFinalModal>
</template>

<script setup>
import { VueFinalModal } from 'vue-final-modal';

defineProps({
    title: {
        type: String,
        default: 'تایید عملیات',
    },
    dark: {
        type: Boolean,
        default: false,
    },
});

const emit = defineEmits(['close', 'confirm', 'cancel']);

const handleClose = () => {
    emit('cancel');
};
</script>
