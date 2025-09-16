<template>
    <Transition name="popup-fade">
        <InfoBox
            :loading="loading"
            :info-box-style="infoBoxStyle"
            :title="
                checkSubMigration()
                    ? migrateInfo.territoryName
                    : selectedIsland.name
            "
            :button-text="buttonText"
            :error-text="errorText"
            :button-enabled="isButtonEnabled"
            :cost="cost"
            :box-width="calculateWidth()"
            @action="actionOnClick"
        >
            <div v-if="checkSubMigration()" class="w-full">
                <img
                    v-if="migrateInfo.status == 'untouched'"
                    :src="`/images/territories/${migrateInfo.territoryId}_locked.png`"
                />
                <img
                    v-else
                    :src="`/images/territories/${migrateInfo.territoryId}_unlocked.png`"
                />
            </div>
            <div
                v-else-if="checkRefuelIsland() && refuel"
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
            <div v-else-if="checkMigrate()" class="w-full space-y-3">
                <div class="flex justify-center">
                    <div class="w-11/12 md:w-full">
                        <PlayerInventoryBar
                            v-if="knowledgeBar"
                            :bar-data="knowledgeBar"
                        ></PlayerInventoryBar>
                    </div>
                </div>
                <div class="grid grid-cols-2 gap-4">
                    <button
                        v-for="option in migrate.territoryMigrationOptions"
                        :key="option.territoryId"
                        class="relative rounded-lg overflow-hidden group text-center border-0 p-0 bg-transparent"
                        :class="
                            option.feasible
                                ? 'cursor-pointer'
                                : 'cursor-not-allowed'
                        "
                        :disabled="!option.feasible"
                        @pointerdown="handleMigrationClick($event, option)"
                    >
                        <img
                            v-if="option.status == 'untouched'"
                            :src="`/images/territories/${option.territoryId}_locked.png`"
                            class="w-full h-24 object-cover"
                        />
                        <img
                            v-else
                            :src="`/images/territories/${option.territoryId}_unlocked.png`"
                            class="w-full h-24 object-cover"
                        />
                        <div
                            class="absolute inset-0 transition-colors pointer-events-none"
                            :class="
                                option.feasible
                                    ? 'bg-black/40 md:bg-transparent md:group-hover:bg-black/40'
                                    : 'bg-black/60'
                            "
                        ></div>
                        <div
                            class="absolute inset-0 flex flex-col items-center justify-center text-white p-2 transition-opacity pointer-events-none"
                            :class="
                                option.feasible
                                    ? 'md:opacity-0 group-hover:opacity-100'
                                    : ''
                            "
                        >
                            <h4 class="font-bold text-base mb-1">
                                {{ option.territoryName }}
                            </h4>
                            <div v-if="option.feasible">
                                <p v-if="option.mustPayCost" class="text-sm">
                                    مهاجرت پولی
                                </p>
                                <p
                                    v-else-if="option.status === 'untouched'"
                                    class="text-sm"
                                >
                                    مهاجرت علمی
                                </p>
                                <p
                                    v-else-if="option.status === 'visited'"
                                    class="text-sm"
                                >
                                    سفر به منظومه
                                </p>
                                <p
                                    v-else-if="option.status === 'resident'"
                                    class="text-sm"
                                >
                                    منظومه‌ی فعلی!
                                </p>
                                <div
                                    v-if="option.mustPayCost"
                                    class="flex items-center justify-center text-xs mt-1"
                                >
                                    <span>{{
                                        option.migrationCost.items[0].amount
                                    }}</span>
                                    <img
                                        src="/images/icons/coin.png"
                                        class="w-4 h-4 mr-1"
                                    />
                                </div>
                            </div>
                            <p v-else class="text-red-400 text-xs px-1">
                                {{ option.reason }}
                            </p>
                        </div>
                    </button>
                </div>
            </div>
            <div
                v-else-if="checkNonAnchoredIsland() && !player.anchored"
                class="w-full space-y-3"
            ></div>
            <div v-else-if="checkAnchoredIsland()" class="w-full space-y-3">
                <p class="text-center text-sm text-gray-800">
                    شما در این سیاره قرار دارید.
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
    migrateCheck,
    migrateTo,
} from '@/services/api/index.js';
import InfoBox from '@/components/common/InfoBox.vue';
import PlayerInventoryBar from '@/components/common/PlayerInventoryBar.vue';
import { glossary } from '@/services/glossary.js';

defineOptions({
    name: 'IslandInfoBox',
});

const props = defineProps({
    selectedIsland: Object,
    territoryId: String,
    player: Object,
    infoBoxStyle: Object,
    refuelIslands: Object,
    terminalIslands: Object,
    migrateInfo: Object,
});

const router = useRouter();
const travel = ref(null);
const anchor = ref(null);
const refuel = ref(null);
const migrate = ref(null);
const refuelError = ref(null);
const travelError = ref(null);
const anchorError = ref(null);
const migrateError = ref(null);
const loading = ref(true);
const isRefuelIsland = ref(false);
const isTerminalIsland = ref(false);

const fuelCount = ref(0);
const fuelInput = ref(null);

const init = async () => {
    try {
        if (!props.selectedIsland) {
            isRefuelIsland.value = false;
        } else {
            isRefuelIsland.value = props.refuelIslands.some(
                island => island.id === props.selectedIsland.id
            );
        }
        isTerminalIsland.value = props.terminalIslands.some(
            island => island.id === props.selectedIsland.id
        );

        const isCurrent = props.selectedIsland.id === props.player.atIsland;
        if (isCurrent) {
            if (isRefuelIsland.value) {
                await updateRefuel();
            } else if (isTerminalIsland.value) {
                await updateMigration();
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

const calculateWidth = () => {
    if (checkSubMigration()) return 39;
    else if (checkMigrate()) return 85;
    else return 60;
};

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

const updateMigration = async () => {
    try {
        migrate.value = await migrateCheck();
    } catch (err) {
        migrateError.value = err.message;
    }
};

const checkSubMigration = () => {
    return props.migrateInfo && 1;
};
const checkRefuelIsland = () => {
    return isCurrentIsland.value && isRefuelIsland.value;
};
const checkNonAnchoredIsland = () => {
    return (
        isCurrentIsland.value &&
        !isRefuelIsland.value &&
        !isTerminalIsland.value &&
        !props.player.anchored
    );
};
const checkAnchoredIsland = () => {
    return (
        isCurrentIsland.value &&
        !isRefuelIsland.value &&
        !isTerminalIsland.value
    );
};
const checkTravel = () => {
    return travel.value && 1;
};
const checkMigrate = () => {
    return isTerminalIsland.value && migrate.value;
};
const buttonText = computed(() => {
    if (checkSubMigration()) {
        if (props.migrateInfo.feasible) {
            if (props.migrateInfo.mustPayCost) {
                return 'مهاجرت پولی';
            }
            switch (props.migrateInfo.status) {
                case 'visited':
                    return 'سفر به منظومه';
                case 'resident':
                    return 'منظومه‌ی فعلی!';
                case 'untouched':
                    return 'مهاجرت علمی';
            }
            debugger;
        }
        return null;
    } else if (checkRefuelIsland() && !refuelError.value)
        return glossary.buyFuel;
    else if (checkMigrate() && !migrateError.value) return null;
    else if (checkNonAnchoredIsland() && anchor.value && !anchorError.value)
        return 'فرود آمدن';
    else if (checkAnchoredIsland()) return 'ورود به سیاره';
    else if (checkTravel() && !travelError.value) return 'سفر به سیاره';
    return null;
});

const errorText = computed(() => {
    if (checkSubMigration()) {
        if (props.migrateInfo.feasible) return null;
        return props.migrateInfo.reason;
    } else if (checkRefuelIsland()) {
        if (refuelError.value) return refuelError.value;
        else return null;
    } else if (checkMigrate()) {
        if (migrateError.value) return migrateError.value;
        else return null;
    } else if (checkNonAnchoredIsland()) {
        if (anchor.value && !anchorError.value) return null;
        else
            return anchorError.value
                ? anchorError.value
                : 'خطا در دریافت اطلاعات';
    } else if (checkAnchoredIsland()) return null;
    else if (!(checkTravel() && !travelError.value))
        return travelError.value ? travelError.value : 'خطا در دریافت اطلاعات';
    return null;
});

const isButtonEnabled = computed(() => {
    if (checkMigrate()) {
        if (props.selectedIsland.id == props.player.atIsland) return true;
        return false;
    }
    return true;
});

const cost = computed(() => {
    if (checkSubMigration()) {
        if (props.migrateInfo.feasible) {
            if (props.migrateInfo.mustPayCost) {
                return props.migrateInfo.migrationCost;
            }
        }
        return null;
    } else if (checkRefuelIsland()) {
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
    } else if (checkMigrate()) {
        return null;
    } else if (checkNonAnchoredIsland() && anchor.value && !anchorError.value)
        return anchor.value.anchoringCost;
    else if (checkAnchoredIsland()) return null;
    else if (checkTravel() && !travelError.value)
        return travel.value.travelCost;
    return null;
});

const actionOnClick = () => {
    if (checkSubMigration()) {
        migrateTo(props.migrateInfo.territoryId);
        return;
    } else if (checkRefuelIsland()) {
        buyFuelFromIsland();
        return;
    } else if (checkMigrate()) {
        return;
    } else if (checkNonAnchoredIsland() && anchor.value && !anchorError.value) {
        dropAnchor();
        return;
    } else if (checkAnchoredIsland()) {
        navigateToIsland();
        return;
    } else if (checkTravel() && !travelError.value) {
        travelToIsland();
        return;
    }
    debugger;
    return;
};

const handleMigrationClick = (event, option) => {
    if (!option.feasible) return;
    event.preventDefault();
    event.stopPropagation();
    migrateTo(option.territoryId);
};

const knowledgeBar = computed(() => {
    const fetchedKnowledgeBar = props.player.knowledgeBars.find(
        bar => bar.territoryId === migrate.value.knowledgeCriteriaTerritory
    );
    return {
        name: glossary.requiredKnowledge,
        englishName: 'Knowledge',
        total: fetchedKnowledgeBar.total,
        required: migrate.value.minAcceptableKnowledge,
        value: migrate.value.knowledgeValue,
        icon: '/images/icons/knowledge.png',
        shadowColor: '#ff7e5f',
        gradientFrom: '#b65f69',
        gradientTo: '#feb47b',
    };
});

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
