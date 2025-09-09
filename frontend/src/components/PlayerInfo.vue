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

        <div class="flex items-center gap-2 mt-2">
            <button
                class="transition-transform duration-200 hover:scale-110 pointer-events-auto"
                title="کتابخانه"
                @pointerdown="openBookshelf"
            >
                <img
                    src="/images/icons/book.png"
                    class="w-12 h-12 drop-shadow-lg"
                    alt="کتاب‌ها"
                />
            </button>
            <button
                class="transition-transform duration-200 hover:scale-110 pointer-events-auto"
                title="بازار"
                @pointerdown="openMarket"
            >
                <img
                    src="/images/icons/market.png"
                    class="w-12 h-12 drop-shadow-lg"
                    alt="داد و ستد"
                />
            </button>
            <button
                class="transition-transform duration-200 hover:scale-110 pointer-events-auto"
                title="کوله پشتی"
                @pointerdown="openBackpack"
            >
                <img
                    src="/images/icons/backpack.png"
                    class="w-12 h-12 drop-shadow-lg"
                    alt="کوله پشتی"
                />
            </button>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { logout as apiLogout } from '@/services/api';
import PlayerInventoryBar from './PlayerInventoryBar.vue';
import { useModal } from 'vue-final-modal';
import ConfirmModal from './ConfirmModal.vue';
import Bookshelf from './Bookshelf.vue';
import Backpack from './Backpack.vue';
import Market from './Market.vue';

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

const { open: openBookshelf, close: closeBookshelf } = useModal({
    component: Bookshelf,
    attrs: {
        books: props.player.books,
        onClose() {
            closeBookshelf();
        },
    },
});

const { open: openMarket, close: closeMarket } = useModal({
    component: Market,
    attrs: {
        player: props.player,
        username: props.username,
        onClose() {
            closeMarket();
        },
    },
});

const { open: openBackpack, close: closeBackpack } = useModal({
    component: Backpack,
    attrs: {
        inventoryItems: computed(() => {
            if (!props.player) return [];
            const items = ['goldenKeys', 'redKeys', 'blueKeys'];
            return items.map(key => ({
                icon: `/images/icons/${key}.png`,
                name: getKeyDisplayName(key),
                quantity: props.player[key],
            }));
        }),
        onClose() {
            closeBackpack();
        },
    },
});

function getKeyDisplayName(key) {
    const displayNames = {
        goldenKeys: 'کلید طلایی',
        redKeys: 'کلید قرمز',
        blueKeys: 'کلید آبی',
    };
    return displayNames[key] || key;
}

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
