<template>
    <div
        :style="{
            width: `calc(var(--spacing) * ${boxWidth})`,
            ...infoBoxStyle,
        }"
        class="bg-surface/90 text-content border border-border p-4 rounded-xl shadow-2xl font-main text-base z-[10000] flex flex-col items-center pointer-events-auto"
        @pointerdown.stop
    >
        <h3 class="text-lg font-bold text-center shrink-0 text-accent">
            {{ title }}
        </h3>

        <div
            class="w-full grid transition-[grid-template-rows] duration-300 ease-smooth-expand"
            :class="loading ? 'grid-rows-[0fr]' : 'grid-rows-[1fr]'"
        >
            <div class="overflow-hidden">
                <div v-if="!loading" class="w-full mt-3 space-y-3">
                    <slot></slot>
                    <div v-if="buttonText || cost" class="text-sm">
                        <CostlyButton
                            :on-click="() => $emit('action')"
                            :cost="cost"
                            :label="buttonText"
                            :loading="loading"
                            :enabled="buttonEnabled && !errorText"
                        >
                        </CostlyButton>
                    </div>
                    <div
                        v-if="errorText"
                        class="text-center text-sm text-danger-content font-semibold bg-danger-soft border border-danger/40 p-2 rounded-md"
                    >
                        {{ errorText }}
                    </div>
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
    buttonText: String,
    buttonEnabled: Boolean,
    errorText: String,
    cost: Object,
    loading: Boolean,
    boxWidth: Number,
});

defineEmits(['action']);
</script>
