<template>
    <Transition name="popup-fade">
        <div v-if="hoveredIsland" :style="infoBoxStyle"
            class="bg-[rgb(121,200,237,0.8)] text-[#310f0f] p-4 rounded-xl font-vazir text-base z-[10000] flex flex-col items-center pointer-events-auto w-60"
            @pointerdown.stop>

            <h3 class="text-lg font-bold text-center shrink-0">{{ hoveredIslandName }}</h3>

            <div class="w-full grid transition-[grid-template-rows] duration-300 ease-smooth-expand"
                :class="loading ? 'grid-rows-[0fr]' : 'grid-rows-[1fr]'">

                <div class="overflow-hidden">
                    <div v-if="!loading" class="w-full mt-3 space-y-3">

                        <div v-if="isCurrentIsland && isRefuelIsland && refuel" class="w-full space-y-2">
                            <div class="flex justify-between items-center text-sm">
                                <span>قیمت هر واحد:</span>
                                <span class="font-semibold">{{ refuel.coinCostPerUnit }}</span>
                            </div>
                            <div class="flex justify-between items-center text-sm">
                                <span>حداکثر واحد:</span>
                                <span class="font-semibold">{{ refuel.maxAvailableAmount }}</span>
                            </div>
                            <input type="number" ref="fuelInput"
                                :max="refuel ? refuel.maxAvailableAmount : player.fuelCap - player.fuel"
                                v-model.number="fuelCount" @pointerdown.stop="focusFuelInput" @dblclick.stop
                                class="w-full mt-1 rounded-lg border border-[#07458bb5] text-center bg-transparent py-1.5" />
                            <button @pointerdown.stop="buyFuel"
                                class="btn-hover w-full p-2 rounded-lg bg-[#07458bb5] text-white">
                                {{ fuelPriceText }}
                            </button>
                        </div>

                        <div v-else-if="isCurrentIsland && !isRefuelIsland" class="w-full space-y-3">
                            <p class="text-center text-sm text-gray-800">شما در این جزیره قرار دارید.</p>
                            <button @pointerdown.stop="$emit('navigateToIsland', player.atIsland)"
                                class="btn-hover w-full p-2 rounded-lg bg-sky-600 text-white">
                                ورود به جزیره
                            </button>
                        </div>

                        <div v-else class="w-full space-y-3">
                            <div v-if="isAdjacent && travel" class="flex justify-between items-center text-sm">
                                <span class="text-gray-800">هزینه سفر:</span>
                                <div class="flex items-center gap-x-1">
                                    <span class="text-gray-900 font-bold">{{ travel.fuelCost }}</span>
                                    <img src="/images/icons/fuel.png" alt="Fuel Icon" class="w-5 h-5" />
                                </div>
                            </div>
                            <button
                                :disabled="loading || !isAdjacent || !!travelError || (travel && player.fuel < travel.fuelCost)"
                                @pointerdown.stop="$emit('travelToIsland', hoveredIsland)"
                                class="btn-hover w-full p-2 rounded-lg bg-green-600 text-white disabled:opacity-50 disabled:cursor-not-allowed text-xs">
                                <span v-if="!isAdjacent">مسیر مستقیمی وجود ندارد</span>
                                <span v-else-if="travel && player.fuel < travel.fuelCost">سوخت کافی نیست</span>
                                <span v-else>سفر به این جزیره</span>
                            </button>
                            <p v-if="travelError"
                                class="text-center text-sm text-red-700 font-semibold bg-red-200 p-2 rounded-md">
                                {{ travelError }}
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </Transition>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue';

const props = defineProps({
    hoveredIsland: String,
    hoveredIslandName: String,
    player: Object,
    refuel: Object,
    travel: Object,
    infoBoxStyle: Object,
    isRefuelIsland: Boolean,
    isAdjacent: Boolean,
    loading: Boolean,
    travelError: String,
});

const emit = defineEmits(['buyFuel', 'navigateToIsland', 'travelToIsland']);

const fuelInput = ref(null);
const fuelCount = ref(0);

const isCurrentIsland = computed(() => props.hoveredIsland && props.player && props.hoveredIsland === props.player.atIsland);

const fuelPriceText = computed(() => {
    if (!props.refuel || !fuelCount.value || fuelCount.value <= 0) return "خرید سوخت";
    const totalCost = props.refuel.coinCostPerUnit * fuelCount.value;
    return `خرید (${totalCost} سکه)`;
});

const focusFuelInput = () => {
    nextTick(() => {
        fuelInput.value?.focus();
    });
};

const buyFuel = () => {
    if (fuelCount.value > 0) {
        emit('buyFuel', fuelCount.value);
    }
};

watch(() => props.hoveredIsland, () => {
    fuelCount.value = 0;
});

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

<style scoped>
.popup-fade-enter-active,
.popup-fade-leave-active {
    transition: opacity 0.3s ease, transform 0.3s ease;
}

.popup-fade-enter-from,
.popup-fade-leave-to {
    opacity: 0;
    transform: translateY(10px);
}

.btn-hover {
    transition: transform 0.2s ease, filter 0.2s ease;
}

.btn-hover:hover:not(:disabled) {
    transform: scale(1.05);
    filter: brightness(1.1);
}
</style>
