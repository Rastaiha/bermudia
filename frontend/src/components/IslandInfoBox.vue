<template>
    <Transition name="popup-fade">
        <div v-if="selectedIsland" :style="infoBoxStyle"
            class="bg-[rgb(121,200,237,0.8)] text-[#310f0f] p-4 rounded-xl font-vazir text-base z-[10000] flex flex-col items-center pointer-events-auto w-60"
            @pointerdown.stop>
            <h3 class="text-lg font-bold text-center shrink-0">
                {{ selectedIsland.name }}
            </h3>

            <div class="w-full grid transition-[grid-template-rows] duration-300 ease-smooth-expand"
                :class="loading ? 'grid-rows-[0fr]' : 'grid-rows-[1fr]'">
                <div class="overflow-hidden">
                    <div v-if="!loading" class="w-full mt-3 space-y-3">
                        <RefuelIslandInfoBox v-if="isCurrentIsland && isRefuelIsland && refuel" :refuel="refuel"
                            :player="player" @buyFuel="$emit('buyFuel', $event)" />
                        <div v-else-if="isCurrentIsland && !isRefuelIsland" class="w-full space-y-3">
                            <p class="text-center text-sm text-gray-800">
                                شما در این جزیره قرار دارید.
                            </p>
                            <button @pointerdown.stop="
                                $emit('navigateToIsland', player.atIsland)
                                " class="btn-hover w-full p-2 rounded-lg bg-sky-600 text-white">
                                ورود به جزیره
                            </button>
                        </div>
                        <div v-else class="w-full space-y-3">
                            <div v-if="isAdjacent && travel" class="text-sm">
                                <div class="text-gray-800 mb-1">هزینه سفر:</div>
                                <div v-for="(costItem, index) in travel.travelCost.items" :key="index"
                                    class="flex justify-between items-center flex-row-reverse">
                                    <div class="flex items-center gap-x-1">
                                        <span class="text-gray-900 font-bold">{{ costItem.amount }}</span>
                                        <img :src="getIconByType(costItem.type)" :alt="costItem.type + ' Icon'"
                                            class="w-5 h-5" />
                                    </div>
                                </div>
                            </div>
                            <button :disabled="loading ||
                                !isAdjacent ||
                                !!travelError ||
                                (travel && player.fuel < travel.fuelCost)
                                " @pointerdown.stop="
                                    $emit('travelToIsland', selectedIsland.id)
                                    "
                                class="btn-hover w-full p-2 rounded-lg bg-green-600 text-white disabled:opacity-50 disabled:cursor-not-allowed text-xs">
                                <span v-if="!isAdjacent">مسیر مستقیمی وجود ندارد</span>
                                <span v-else-if="
                                    travel && player.fuel < travel.fuelCost
                                ">سوخت کافی نیست</span>
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
import { computed, ref, watch, nextTick } from "vue";
import RefuelIslandInfoBox from "./RefuelBox.vue";

const props = defineProps({
    selectedIsland: Object,
    player: Object,
    refuel: Object,
    travel: Object,
    infoBoxStyle: Object,
    isRefuelIsland: Boolean,
    isAdjacent: Boolean,
    loading: Boolean,
    travelError: String,
});

const emit = defineEmits(["buyFuel", "navigateToIsland", "travelToIsland"]);

const fuelInput = ref(null);
const fuelCount = ref(0);

const isCurrentIsland = computed(() => props.selectedIsland.id && props.player && props.selectedIsland.id === props.player.atIsland);

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

const getIconByType = (type) => {
    switch (type) {
        case "fuel":
            return "/images/icons/fuel.png";
        case "coin":
            return "/images/icons/coin.png"
    }
    return null;
}

watch(() => props.selectedIsland, () => {
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
