<template>
    <div class="w-full space-y-2">
        <div class="flex justify-between items-center text-sm">
            <span>قیمت هر واحد:</span>
            <span class="font-semibold">{{ refuel.coinCostPerUnit }}</span>
        </div>
        <div class="flex justify-between items-center text-sm">
            <span>حداکثر واحد:</span>
            <span class="font-semibold">{{ refuel.maxAvailableAmount }}</span>
        </div>
        <input
            ref="fuelInput"
            v-model.number="fuelCount"
            type="number"
            :max="
                refuel
                    ? refuel.maxAvailableAmount
                    : player.fuelCap - player.fuel
            "
            class="w-full mt-1 rounded-lg border border-[#07458bb5] text-center bg-transparent py-1.5"
            @pointerdown.stop="focusFuelInput"
            @dblclick.stop
        />
        <button
            class="btn-hover w-full p-2 rounded-lg bg-[#07458bb5] text-white"
            @pointerdown.stop="buyFuel"
        >
            {{ fuelPriceText }}
        </button>
    </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue';
import { glossary } from '@/services/glossary.js';

const props = defineProps({
    refuel: Object,
    player: Object,
});

const emit = defineEmits(['buyFuel']);

const fuelInput = ref(null);
const fuelCount = ref(0);

const fuelPriceText = computed(() => {
    if (!props.refuel || !fuelCount.value || fuelCount.value <= 0)
        return glossary.buyFuel;
    const totalCost = props.refuel.coinCostPerUnit * fuelCount.value;
    return `${glossary.buy} (${totalCost} ${glossary.coin})`;
});

const focusFuelInput = () => {
    nextTick(() => fuelInput.value?.focus());
};

const buyFuel = () => {
    if (fuelCount.value > 0) emit('buyFuel', fuelCount.value);
};

watch(fuelCount, newValue => {
    if (!props.refuel) return;
    let correctedValue = newValue;
    if (newValue > props.refuel.maxAvailableAmount) {
        correctedValue = props.refuel.maxAvailableAmount;
    } else if (newValue < 0) {
        correctedValue = 0;
    }
    if (correctedValue !== fuelCount.value) {
        fuelCount.value = correctedValue;
    }
});
</script>
