<template>
    <div
        :style="{
            width: `calc(var(--spacing) * ${boxWidth})`,
            ...infoBoxStyle,
        }"
        class="bg-[rgb(121,200,237,0.8)] text-[#310f0f] p-4 rounded-xl font-vazir text-base z-[10000] flex flex-col items-center pointer-events-auto"
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
                    <slot></slot>
                    <div
                        v-if="!errorText && (buttonText || cost)"
                        class="text-sm"
                    >
                        <CostlyButton
                            :on-click="() => $emit('action')"
                            :cost="cost"
                            :label="buttonText"
                            :loading="loading"
                            :enabled="buttonEnabled"
                        >
                        </CostlyButton>
                    </div>
                    <div
                        v-if="errorText"
                        class="text-center text-sm text-red-700 font-semibold bg-red-200 p-2 rounded-md"
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
