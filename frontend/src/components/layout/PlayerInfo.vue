<template>
    <div
        v-if="player"
        dir="rtl"
        class="fixed top-20 right-4 z-50 font-sans text-white w-64 pointer-events-none space-y-2"
    >
        <PlayerInventoryBar
            v-if="knowledgeBar"
            :bar-data="knowledgeBar"
        ></PlayerInventoryBar>
        <PlayerInventoryBar
            v-if="fuelBar"
            :bar-data="fuelBar"
        ></PlayerInventoryBar>
        <PlayerInventoryBar
            v-if="coinBar"
            :bar-data="coinBar"
        ></PlayerInventoryBar>
    </div>
</template>

<script setup>
import { computed } from 'vue';
import PlayerInventoryBar from '@/components/common/PlayerInventoryBar.vue';
import { COST_ITEMS_INFO } from '@/services/cost';

const props = defineProps({
    player: {
        type: Object,
        required: true,
    },
});

const knowledgeBar = computed(() => {
    if (
        !props.player ||
        !props.player.knowledgeBars ||
        props.player.knowledgeBars.length === 0
    ) {
        return null;
    }
    const fetchedKnowledgeBar =
        props.player.knowledgeBars.find(
            bar => bar.territoryId === props.player.atTerritory
        ) || props.player.knowledgeBars[0];
    return {
        name: 'دانش',
        englishName: 'Knowledge',
        total: fetchedKnowledgeBar.total,
        value: fetchedKnowledgeBar.value,
        icon: '/images/icons/knowledge.png',
        shadowColor: '#ff7e5f',
        gradientFrom: '#b65f69',
        gradientTo: '#feb47b',
    };
});

const fuelBar = computed(() => {
    if (!props.player) return null;
    return {
        name: COST_ITEMS_INFO['fuel'].name,
        englishName: 'Fuel',
        total: props.player.fuelCap,
        value: props.player.fuel,
        icon: COST_ITEMS_INFO['fuel'].icon,
        shadowColor: '#6B7280',
        gradientFrom: '#364153',
        gradientTo: '#000',
    };
});

const coinBar = computed(() => {
    if (!props.player) return null;
    return {
        name: COST_ITEMS_INFO['coin'].name,
        englishName: 'Coin',
        total: -1,
        value: props.player.coin,
        width: 0.75,
        icon: COST_ITEMS_INFO['coin'].icon,
        shadowColor: '#6B7280',
        gradientFrom: '#364153',
        gradientTo: '#000',
    };
});
</script>
