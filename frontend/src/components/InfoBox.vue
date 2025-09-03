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
                    <div v-if="!error" class="text-sm">
                        <slot></slot>
                        <CostlyButton
                            :on-click="() => $emit('action')"
                            :cost="cost"
                            :label="button"
                            :loading="loading"
                        >
                        </CostlyButton>
                    </div>
                    <div v-else-if="error">{{ error }}</div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import CostlyButton from './CostlyButton.vue';

defineProps({
    infoBoxStyle: Object,
    title: String,
    button: String,
    error: String,
    cost: Object,
    loading: Boolean,
});

defineEmits(['action']);
</script>
