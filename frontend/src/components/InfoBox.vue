<template>
    <div
        :style="infoBoxStyle"
        class="bg-[rgb(121,200,237,0.8)] text-[#310f0f] p-4 rounded-xl font-vazir text-base z-[10000] flex flex-col items-center pointer-events-auto w-60"
        @pointerdown.stop
    >
        <h3 class="text-lg font-bold text-center shrink-0">
            {{ title }}
        </h3>

        <div
            class="w-full grid transition-[grid-template-rows] duration-300 ease-smooth-expand"
            :class="loading ? 'grid-rows-[0fr]' : 'grid-rows-[1fr]'"
        >
            <div class="overflow-hidden">
                <div v-if="!loading" class="w-full mt-3 space-y-3">
                    <div v-if="button && !error" class="text-sm">
                        <slot></slot>
                        <button
                            :disabled="loading"
                            class="btn-hover w-full p-2 rounded-lg bg-green-600 text-white disabled:opacity-50 disabled:cursor-not-allowed text-xs"
                            @pointerdown.stop="$emit('action')"
                        >
                            {{ button }}
                        </button>
                    </div>
                    <div v-else>
                        {{
                            anchorError ? anchorError : 'خطا در دریافت اطلاعات'
                        }}
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
defineProps({
    title: String,
    button: String,
    error: String,
    loading: Boolean,
});

defineEmits(['action']);
</script>

<style scoped>
.btn-hover {
    transition:
        transform 0.2s ease,
        filter 0.2s ease;
}

.btn-hover:hover:not(:disabled) {
    transform: scale(1.05);
    filter: brightness(1.1);
}
</style>
