<template>
    <VueFinalModal
        class="flex justify-center items-center"
        content-class="flex flex-col max-w-md mx-4 p-6 bg-white border border-gray-200 rounded-xl shadow-lg space-y-4"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div class="flex items-center justify-between">
            <h1 class="text-xl font-semibold text-gray-800">
                {{ title }}
            </h1>
            <button
                class="p-1 rounded-full hover:bg-gray-200"
                @click="handleClose"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-6 w-6 text-gray-600"
                    fill="none"
                    viewBox="0 0 24"
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

        <div>
            <slot name="content">
                <p class="text-gray-600">
                    اینجا می‌توانید محتوای مودال را قرار دهید.
                </p>
            </slot>
        </div>

        <div class="flex justify-end gap-4 mt-4">
            <button
                class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200"
                @click="emit('cancel')"
            >
                انصراف
            </button>
            <button
                class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700"
                @click="emit('confirm')"
            >
                تایید
            </button>
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
});

const emit = defineEmits(['close', 'confirm', 'cancel']);

const handleClose = () => {
    console.log('Close button clicked!');
    emit('cancel');
};
</script>
