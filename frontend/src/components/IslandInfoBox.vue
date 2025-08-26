<template>
    <div v-if="hoveredNode"
        class="bg-[rgb(121,200,237,0.8)] text-[#310f0f] px-5 py-3 rounded-md font-vazir text-base z-[10000] whitespace-nowrap flex flex-col justify-center items-center pointer-events-auto"
        :style="infoBoxStyle" @pointerdown.stop>
        <div>{{ hoveredNode.name }}</div>
        <div>
            <div v-if="isCurrentIsland && isFuelStation"
                class="whitespace-nowrap flex flex-col justify-center items-center">
                <div v-if="refuel">قیمت هر واحد: {{ refuel.coinCostPerUnit }}</div>
                <div v-if="refuel">حداکثر واحد قابل اخذ: {{ refuel.maxAvailableAmount }}</div>
                <input type="number" ref="fuelInput"
                    :max="refuel ? refuel.maxAvailableAmount : player.fuelCap - player.fuel" v-model.number="fuelCount"
                    @pointerdown.stop="focusFuelInput" class="w-16 rounded-lg border border-[#07458bb5] text-center" />
                <button @pointerdown.stop="buyFuel" class="p-1.5 rounded-lg bg-[#07458bb5] mt-2.5">
                    {{ fuelPriceText }}
                </button>
            </div>
            <button v-else-if="isCurrentIsland" @pointerdown.stop="$emit('navigateToIsland', player.atIsland.id)"
                class="p-1.5 rounded-lg bg-[#07458bb5] mt-2.5 disabled:contrast-50">
                ورود به جزیره
            </button>
            <button v-else :disabled="!isAdjacent" @pointerdown.stop="$emit('travelToIsland', hoveredNode.id)"
                class="p-1.5 rounded-lg bg-[#07458bb5] mt-2.5 disabled:contrast-50">
                سفر به جزیره
                <span v-if="travel && travel.fuelCost">{{ travel.fuelCost }}</span>
            </button>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue';

const props = defineProps({
    hoveredNode: Object,
    player: Object,
    refuel: Object,
    travel: Object,
    infoBoxStyle: Object,
    isFuelStation: Boolean,
    isAdjacent: Boolean,
});

const emit = defineEmits(['buyFuel', 'navigateToIsland', 'travelToIsland', 'update:fuelCount']);

const fuelInput = ref(null);
const fuelCount = ref(0);

const isCurrentIsland = computed(() => props.hoveredNode && props.player && props.hoveredNode.id === props.player.atIsland.id);

const fuelPriceText = computed(() => {
    if (!props.refuel) return "خرید سوخت";
    return `خرید سوخت ${props.refuel.coinCostPerUnit * fuelCount.value}`;
});

const focusFuelInput = () => {
    nextTick(() => {
        fuelInput.value?.focus();
    });
};

const buyFuel = () => {
    emit('buyFuel', fuelCount.value);
};

watch(fuelCount, (newValue) => {
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