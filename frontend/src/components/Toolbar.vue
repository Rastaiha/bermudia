<template>
    <div
        class="fixed top-1/2 left-0 z-50 -translate-y-1/2"
        :class="{
            'transition-transform duration-300 ease-in-out': isMobile,
            '-translate-x-full': !shouldBeOpen,
            'translate-x-0': shouldBeOpen,
        }"
    >
        <div class="relative">
            <div
                class="flex flex-col gap-2 rounded-r-xl bg-slate-800/80 p-2 shadow-lg backdrop-blur-sm"
            >
                <button
                    v-for="item in menuItems"
                    :key="item.id"
                    class="h-12 w-12 rounded-lg p-2 transition-colors duration-200 hover:bg-slate-700/90"
                    :title="item.alt"
                    @pointerdown="handleItemClick(item)"
                >
                    <img
                        :src="item.icon"
                        :alt="item.alt"
                        class="h-full w-full object-contain"
                    />
                </button>

                <MuteButton />

                <button
                    class="h-12 w-12 rounded-lg p-2 transition-colors duration-200 hover:bg-slate-700/90"
                    title="خروج"
                    @pointerdown="openLogoutModal"
                >
                    <img
                        src="/images/icons/logout.png"
                        alt="خروج"
                        class="h-full w-full object-contain"
                    />
                </button>
            </div>

            <button
                v-if="isMobile"
                class="absolute top-1/2 left-full h-16 w-8 -translate-y-1/2 rounded-r-xl bg-slate-800/80 p-1 shadow-lg"
                @pointerdown="isOpen = !isOpen"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke-width="2"
                    stroke="currentColor"
                    class="h-6 w-6 text-white transition-transform duration-300"
                    :class="{ 'rotate-180': isOpen }"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="m8.25 4.5 7.5 7.5-7.5 7.5"
                    />
                </svg>
            </button>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useModal } from 'vue-final-modal';
import { logout as apiLogout } from '@/services/api';
import Backpack from './Backpack.vue';
import Bookshelf from './Bookshelf.vue';
import Market from './Market.vue';
import ConfirmModal from './ConfirmModal.vue';
import MuteButton from './MuteButton.vue';

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
const isMobile = ref(false);
const isOpen = ref(false);

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
        books: computed(() => (props.player ? props.player.books : [])),
        onClose() {
            closeBookshelf();
        },
    },
});

const { open: openMarket, close: closeMarket } = useModal({
    component: Market,
    attrs: {
        player: computed(() => props.player),
        username: computed(() => props.username),
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
                quantity: props.player[key] || 0,
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

const menuItems = [
    {
        id: 'backpack',
        icon: '/images/icons/backpack.png',
        alt: 'کوله پشتی',
        action: openBackpack,
    },
    {
        id: 'bookshelf',
        icon: '/images/icons/book.png',
        alt: 'کتابخانه',
        action: openBookshelf,
    },
    {
        id: 'market',
        icon: '/images/icons/market.png',
        alt: 'بازار',
        action: openMarket,
    },
];

function handleItemClick(item) {
    if (item.action) {
        item.action();
    }
}

const checkScreenSize = () => {
    const newIsMobile = window.innerWidth < 1024;
    if (newIsMobile !== isMobile.value) {
        isMobile.value = newIsMobile;
        if (!isMobile.value) {
            isOpen.value = true;
        }
    }
};

onMounted(() => {
    checkScreenSize();
    window.addEventListener('resize', checkScreenSize);
});

onUnmounted(() => {
    window.removeEventListener('resize', checkScreenSize);
});

const shouldBeOpen = computed(() => {
    return isOpen.value || !isMobile.value;
});
</script>
