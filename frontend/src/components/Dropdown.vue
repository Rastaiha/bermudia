<template>
    <div
        ref="dropdownRef"
        class="relative w-3/4 min-w-[80px] font-inherit ml-auto"
    >
        <div :class="dropdownHeaderClasses" @pointerdown="toggleDropdown">
            <span class="font-medium text-gray-200 text-xs drop-shadow-sm">{{
                title
            }}</span>
            <svg
                :class="chevronClasses"
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
            >
                <polyline points="6,9 12,15 18,9"></polyline>
            </svg>
        </div>

        <transition
            enter-active-class="transition-all duration-300 ease-out"
            leave-active-class="transition-all duration-300 ease-in"
            enter-from-class="opacity-0 -translate-y-2"
            leave-to-class="opacity-0 -translate-y-2"
        >
            <div
                v-if="isOpen"
                class="absolute top-full left-0 right-0 w-full bg-gray-800/95 border border-blue-500/70 border-t-0 rounded-b-md shadow-xl z-50 overflow-hidden backdrop-blur-md max-h-[300px] overflow-y-auto scrollbar-thin scrollbar-track-gray-700/50 scrollbar-thumb-gray-400/50 hover:scrollbar-thumb-gray-400/70"
            >
                <div v-if="items.length === 0" class="py-4 px-3 text-center">
                    <span class="text-gray-400 text-xs italic drop-shadow-sm"
                        >کیف خالی است</span
                    >
                </div>

                <div
                    v-for="(item, index) in items"
                    :key="`item-${index}`"
                    class="flex items-center py-2 px-3 transition-colors duration-200 border-b border-gray-600/30 last:border-b-0 dir-rtl animate-slide-down"
                    :style="{ 'animation-delay': `${index * 50}ms` }"
                >
                    <div
                        class="bg-blue-500/80 text-white px-2 py-0.5 rounded-xl text-xs font-semibold min-w-[24px] text-center drop-shadow-sm"
                    >
                        {{ item.quantity }}
                    </div>

                    <div class="flex-1 mx-3 text-right">
                        <span
                            class="text-xs text-gray-200 font-normal drop-shadow-sm"
                            >{{ item.name }}</span
                        >
                    </div>

                    <img
                        :src="item.icon"
                        :alt="item.name"
                        class="w-6 h-6 object-contain rounded bg-gray-700/50 p-0.5"
                        @error="handleImageError"
                    />
                </div>
            </div>
        </transition>
    </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';

defineProps({
    title: {
        type: String,
        default: 'کیف',
    },
    items: {
        type: Array,
        default: () => [],
    },
});

const isOpen = ref(false);
const dropdownRef = ref(null);

const dropdownHeaderClasses = computed(() => [
    'flex justify-between items-center px-3 py-1.5 bg-gray-700/80 border border-gray-400/30 rounded-md cursor-pointer transition-all duration-300 select-none backdrop-blur-sm gap-1.5 hover:border-blue-500/50 hover:bg-gray-700/90 touch-manipulation',
    isOpen.value ? 'border-blue-500/70 rounded-b-none bg-gray-700/95' : '',
]);

const chevronClasses = computed(() => [
    'transition-transform duration-300 text-gray-400 pointer-events-none',
    isOpen.value ? 'rotate-180' : '',
]);

function toggleDropdown(event) {
    event.preventDefault();
    isOpen.value = !isOpen.value;
}

function handleImageError(event) {
    event.target.src = '/images/icons/default-item.png';
}

function handleOutsideClick(event) {
    if (dropdownRef.value && !dropdownRef.value.contains(event.target)) {
        isOpen.value = false;
    }
}

onMounted(() => {
    document.addEventListener('pointerdown', handleOutsideClick);
});

onBeforeUnmount(() => {
    document.removeEventListener('pointerdown', handleOutsideClick);
});
</script>

<style scoped>
.animate-slide-in {
    animation: slide-in 0.3s ease forwards;
}

@keyframes slide-in {
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.touch-manipulation {
    touch-action: manipulation;
}
</style>
