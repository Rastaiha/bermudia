<template>
    <div
        v-if="player"
        dir="rtl"
        class="fixed top-4 right-4 z-50 font-sans text-white p-3 rounded-lg w-64 pointer-events-none"
    >
        <div class="flex justify-between items-center mb-3">
            <h2 class="text-lg font-bold text-gray-200 drop-shadow-lg">
                {{ username }}
            </h2>
            <div class="flex items-center gap-3">
                <button
                    class="text-sm font-bold text-gray-200 transition-colors duration-300 transform cursor-pointer hover:text-red-400 drop-shadow-lg pointer-events-auto"
                    @pointerdown="openLogoutModal"
                >
                    خروج
                </button>
            </div>
        </div>

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
import { useRouter } from 'vue-router';
import { logout as apiLogout } from '@/services/api';
import PlayerInventoryBar from './PlayerInventoryBar.vue';
import { useModal } from 'vue-final-modal';
import ConfirmModal from './ConfirmModal.vue';

const props = defineProps({
    player: {
        type: Object,
        required: true,
    },
    username: {
        type: String,
        default: '...',
    },
});

const router = useRouter();

const { open: openLogoutModal, close: closeLogoutModal } = useModal({
    component: ConfirmModal,
    attrs: {
        title: 'خروج از حساب',
        onConfirm() {
            apiLogout();
            router.push({ name: 'Login' });
            closeLogoutModal();
        },
        onCancel() {
            closeLogoutModal();
        },
    },
    slots: {
        content: '<p>آیا برای خروج از حساب کاربری خود اطمینان دارید؟</p>',
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
        name: 'سوخت',
        englishName: 'Fuel',
        total: props.player.fuelCap,
        value: props.player.fuel,
        icon: '/images/icons/fuel.png',
        shadowColor: '#6B7280',
        gradientFrom: '#364153',
        gradientTo: '#000',
    };
});

const coinBar = computed(() => {
    if (!props.player) return null;
    return {
        name: 'سکه',
        englishName: 'Coin',
        total: -1,
        value: props.player.coins,
        width: 0.75,
        icon: '/images/icons/coin.png',
        shadowColor: '#6B7280',
        gradientFrom: '#364153',
        gradientTo: '#000',
    };
});
</script>
