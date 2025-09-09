<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue';

const isMobile = ref(false);
const isOpen = ref(false);

const menuItems = [
    { id: 'backpack', icon: '/images/icons/backpack.png', alt: 'Backpack' },
    { id: 'bookshelf', icon: '/images/icons/book.png', alt: 'Bookshelf' },
    { id: 'treasure', icon: '/images/icons/goldenKeys.png', alt: 'Treasure' },
    { id: 'settings', icon: '/images/icons/knowledge.png', alt: 'Settings' },
];

function handleItemClick(itemId) {
    console.log(`Clicked on: ${itemId}`);
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
                    @pointerdown="handleItemClick(item.id)"
                >
                    <img
                        :src="item.icon"
                        :alt="item.alt"
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
