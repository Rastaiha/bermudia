<template>
    <Transition name="popup-fade">
        <InfoBox
            :loading="loading"
            :info-box-style="infoBoxStyle"
            :title="selectedIsland.name"
            :button="buttonText"
            :error="errorText"
            :cost="cost"
            @action="actionOnClick"
        >
            <div
                v-if="isCurrentIsland && isRefuelIsland && refuel"
                class="w-full space-y-2"
            >
                <div class="flex justify-between items-center text-sm">
                    <span>قیمت هر واحد:</span>
                    <span class="font-semibold">{{
                        refuel.coinCostPerUnit
                    }}</span>
                </div>
                <div class="flex justify-between items-center text-sm">
                    <span>حداکثر واحد:</span>
                    <span class="font-semibold">{{
                        refuel.maxAvailableAmount
                    }}</span>
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
            </div>
            <div
                v-else-if="
                    isCurrentIsland && !isRefuelIsland && !player.anchored
                "
                class="w-full space-y-3"
            ></div>
            <div
                v-else-if="isCurrentIsland && !isRefuelIsland"
                class="w-full space-y-3"
            >
                <p class="text-center text-sm text-gray-800">
                    شما در این جزیره قرار دارید.
                </p>
            </div>
            <div v-else class="w-full space-y-3"></div>
        </InfoBox>
    </Transition>
</template>

<script setup>
import { computed, ref, watch, nextTick, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import {
    travelCheck,
    travelTo,
    refuelCheck,
    buyFuel,
    anchorCheck,
    dropAnchorAtIsland,
} from '@/services/api.js';
import InfoBox from './InfoBox.vue';

const props = defineProps({
    selectedIsland: Object,
    territoryId: String,
    player: Object,
    infoBoxStyle: Object,
    isAdjacent: Boolean,
    refuelIslands: Object,
});

const router = useRouter();
const travel = ref(null);
const anchor = ref(null);
const refuel = ref(null);
const refuelError = ref(null);
const travelError = ref(null);
const anchorError = ref(null);
const loading = ref(true);
const isRefuelIsland = ref(false);

const fuelCount = ref(0);
const fuelInput = ref(null);

const init = async () => {
    try {
        if (!props.selectedIsland) isRefuelIsland.value = false;
        else
            isRefuelIsland.value = props.refuelIslands.some(
                island => island.id === props.selectedIsland.id
            );
        const isCurrent = props.selectedIsland.id === props.player.atIsland;
        if (isCurrent) {
            if (isRefuelIsland.value) {
                await updateRefuel();
            } else if (!props.player.anchored) {
                await updateAnchor();
            } else {
                await nextTick();
            }
        } else {
            await updateTravel();
        }
    } finally {
        loading.value = false;
    }
};

const focusFuelInput = () => {
    nextTick(() => fuelInput.value?.focus());
};

const buyFuelFromIsland = () => {
    if (fuelCount.value > 0) buyFuel(fuelCount.value);
};

const isCurrentIsland = computed(
    () =>
        props.selectedIsland.id &&
        props.player &&
        props.selectedIsland.id === props.player.atIsland
);

const navigateToIsland = () => {
    router.push({
        name: 'Island',
        params: { id: props.territoryId, islandId: props.selectedIsland.id },
    });
};

const travelToIsland = async () => {
    const dest = props.selectedIsland.id;
    try {
        travelError.value = null;
        await travelTo(props.player.atIsland, dest);
    } catch (error) {
        travelError.value = error.message;
    }
};

const dropAnchor = async () => {
    try {
        anchorError.value = null;
        await dropAnchorAtIsland(props.player.atIsland);
    } catch (error) {
        anchorError.value = error.message;
    }
};

const updateTravel = async () => {
    if (!props.player || !props.selectedIsland) return;
    try {
        travel.value = await travelCheck(
            props.player.atIsland,
            props.selectedIsland.id
        );
        if (!travel.value.feasible) {
            travelError.value = travel.value.reason;
        } else {
            travelError.value = null;
        }
    } catch (err) {
        travelError.value = err.message;
    }
};

const updateAnchor = async () => {
    if (!props.player || !props.selectedIsland) return;
    try {
        anchor.value = await anchorCheck(props.player.atIsland);
        if (!anchor.value.feasible) {
            anchorError.value = anchor.value.reason;
        } else {
            anchorError.value = null;
        }
    } catch (err) {
        anchorError.value = err.message;
    }
};

const updateRefuel = async () => {
    try {
        refuel.value = await refuelCheck();
    } catch (err) {
        refuelError.value = err.message;
    }
};

const buttonText = computed(() => {
    if (isCurrentIsland.value && isRefuelIsland.value && !refuelError.value)
        return 'خرید سوخت';
    if (
        isCurrentIsland.value &&
        !isRefuelIsland.value &&
        !props.player.anchored &&
        anchor.value &&
        !anchorError.value
    )
        return 'لنگر بیندازید';
    if (isCurrentIsland.value && !isRefuelIsland.value) return 'ورود به جزیره';
    if (travel.value && !travelError.value) return 'سفر به جزیره';
    return null;
});

const errorText = computed(() => {
    if (isCurrentIsland.value && isRefuelIsland.value) {
        if (refuelError.value) return refuelError.value;
        else return null;
    }
    if (
        isCurrentIsland.value &&
        !isRefuelIsland.value &&
        !props.player.anchored
    ) {
        if (anchor.value && !anchorError.value) return null;
        else
            return anchorError.value
                ? anchorError.value
                : 'خطا در دریافت اطلاعات';
    }
    if (isCurrentIsland.value && !isRefuelIsland.value) return null;
    if (!(travel.value && !travelError.value))
        return travelError.value ? travelError.value : 'خطا در دریافت اطلاعات';
    return null;
});

const cost = computed(() => {
    if (isCurrentIsland.value && isRefuelIsland.value)
        return {
            items: [
                {
                    type: 'coin',
                    amount: refuel.value
                        ? refuel.value.coinCostPerUnit * fuelCount.value
                        : 0,
                },
            ],
        };
    if (
        isCurrentIsland.value &&
        !isRefuelIsland.value &&
        !props.player.anchored &&
        anchor.value &&
        !anchorError.value
    )
        return anchor.value.anchoringCost;
    if (isCurrentIsland.value && !isRefuelIsland.value) return null;
    if (travel.value && !travelError.value) return travel.value.travelCost;
    return null;
});

const actionOnClick = () => {
    if (isCurrentIsland.value && isRefuelIsland.value) {
        buyFuelFromIsland();
        return;
    }
    if (
        isCurrentIsland.value &&
        !isRefuelIsland.value &&
        !props.player.anchored &&
        anchor.value &&
        !anchorError.value
    ) {
        dropAnchor();
        return;
    }
    if (isCurrentIsland.value && !isRefuelIsland.value) {
        navigateToIsland();
        return;
    }
    if (travel.value && !travelError.value) {
        travelToIsland();
        return;
    }
    debugger;
    return;
};

watch(
    () => props.selectedIsland,
    () => {
        fuelCount.value = 0;
    }
);

watch(fuelCount, newValue => {
    if (!refuel.value) return;
    if (newValue > refuel.value.maxAvailableAmount) {
        fuelCount.value = refuel.value.maxAvailableAmount;
    } else if (newValue < 0) {
        fuelCount.value = 0;
    }
});

onMounted(() => {
    init();
});
</script>

<style scoped>
.popup-fade-enter-active,
.popup-fade-leave-active {
    transition:
        opacity 0.3s ease,
        transform 0.3s ease;
}

.popup-fade-enter-from,
.popup-fade-leave-to {
    opacity: 0;
    transform: translateY(10px);
}
</style>
