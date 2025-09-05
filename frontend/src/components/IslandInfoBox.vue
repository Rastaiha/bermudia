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
                <img src="../../public/images/territories/territory1.jfif" />
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
                <PlayerInventoryBar
                    v-if="knowledgeBar"
                    :bar-data="knowledgeBar"
                ></PlayerInventoryBar>
                <div class="w-full flex justify-around">
                    <div
                        v-for="territory in migrate.territoryMigrationOptions"
                        :key="territory.territoryId"
                        style="margin-bottom: 0"
                    >
                        <IslandInfoBox
                            v-if="territory"
                            :info-box-style="migrateBoxStyle"
                            :selected-island="selectedIsland"
                            :player="player"
                            :refuel-islands="refuelIslands"
                            :terminal-islands="terminalIslands"
                            :territory-id="territory.territoryId"
                            :migrate-info="territory"
                        />
                    </div>
                </div>
            </div>
            <div
                v-else-if="checkNonAnchoredIsland() && !player.anchored"
                class="w-full space-y-3"
            ></div>
            <div v-else-if="checkAnchoredIsland()" class="w-full space-y-3">
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
    migrateCheck,
    migrateTo,
} from '@/services/api.js';
import InfoBox from './InfoBox.vue';
import PlayerInventoryBar from './PlayerInventoryBar.vue';

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

const migrateBoxStyle = computed(() => {
    return {
        height: `100%`,
    };
});

const isCurrentIsland = computed(
    () =>
        props.selectedIsland.id &&
        props.player &&
        props.selectedIsland.id === props.player.atIsland
);

const calculateWidth = () => {
    if (checkSubMigration()) return 39;
    else if (checkMigrate()) return 160;
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
                    return 'سفر به قلمرو';
                case 'resident':
                    return 'قلمروی فعلی!';
                case 'untouched':
                    return 'مهاجرت علمی';
            }
            debugger;
        }
        return null;
    } else if (checkRefuelIsland() && !refuelError.value) return 'خرید سوخت';
    else if (checkMigrate() && !migrateError.value) return null;
    else if (checkNonAnchoredIsland() && anchor.value && !anchorError.value)
        return 'لنگر بیندازید';
    else if (checkAnchoredIsland()) return 'ورود به جزیره';
    else if (checkTravel() && !travelError.value) return 'سفر به جزیره';
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

const knowledgeBar = computed(() => {
    const fetchedKnowledgeBar = props.player.knowledgeBars.find(
        bar => bar.territoryId === migrate.value.knowledgeCriteriaTerritory
    );
    return {
        name: 'دانش',
        englishName: 'Knowledge',
        total: fetchedKnowledgeBar.total,
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
