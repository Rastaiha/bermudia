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
                    <div v-if="!errorText" class="text-sm">
                        <slot></slot>
                        <CostlyButton
                            :on-click="() => $emit('action')"
                            :cost="cost"
                            :label="buttonText"
                            :loading="loading"
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
    errorText: String,
    cost: Object,
    loading: Boolean,
});

defineEmits(['action']);
</script>
