<template>
    <VueFinalModal
        class="flex justify-center items-center"
        overlay-class="bg-black/50"
        content-class="panel max-w-md w-full !p-0"
        content-transition="vfm-slide-up"
        overlay-transition="vfm-fade"
    >
        <div class="panel-header">
            <h1 class="panel-title !text-xl">{{ title }}</h1>
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

        <div class="p-6 space-y-4">
            <div>
                <slot name="content">
                    <p class="text-content-muted">
                        اینجا می‌توانید محتوای مودال را قرار دهید.
                    </p>
                </slot>
            </div>

            <div class="flex justify-end gap-4 mt-2">
                <button
                    class="px-4 py-2 text-sm font-medium rounded-lg transition-colors text-content-muted bg-surface-raised hover:bg-border"
                    @click="emit('cancel')"
                >
                    انصراف
                </button>
                <button
                    class="panel-btn !px-5 !py-2 !text-sm"
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
    // Retained for backward compatibility with existing callers; the modal
    // now always uses the unified slate/amber panel theme.
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
