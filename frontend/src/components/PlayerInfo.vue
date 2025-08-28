<template>
    <div v-if="player" dir="rtl" class="fixed top-4 right-4 z-50 font-sans text-white p-3 rounded-lg w-64">

        <div class="flex justify-between items-center mb-3">
            <h2 class="text-lg font-bold text-gray-200 drop-shadow-lg">{{ username }}</h2>
            <button @click="logout"
                class="text-sm font-bold text-gray-200 transition-colors duration-300 transform cursor-pointer hover:text-red-400 drop-shadow-lg"
                title="خروج">
                خروج
            </button>
        </div>

        <div class="mb-3">
            <div class="flex justify-between px-1 mb-1 text-xs text-gray-300 drop-shadow-md">
                <span>دانش</span>
                <span>{{ knowledge.total }} / {{ Math.floor(knowledge.value) }}</span>
            </div>
            <div class="relative flex items-center w-full h-8 rounded-md overflow-hidden bg-black/30 shadow-inner">
                <img src="/images/icons/knowledge.png" alt="Knowledge Icon"
                    class="absolute right-1 top-1/2 -translate-y-1/2 w-6 h-6 z-10">
                <div class="absolute top-0 right-0 h-full rounded-md transition-[width] duration-500 ease-in-out bg-gradient-to-l from-[#ff7e5f] to-[#feb47b] knowledge-shadow"
                    :style="{ width: knowledgePercentage + '%' }"></div>
            </div>
        </div>

        <div>
            <div class="flex justify-between px-1 mb-1 text-xs text-gray-300 drop-shadow-md">
                <span>سوخت</span>
                <span>{{ player.fuelCap }} / {{ Math.floor(player.fuel) }}</span>
            </div>
            <div class="relative flex items-center w-full h-8 rounded-md overflow-hidden bg-black/30 shadow-inner">
                <img src="/images/icons/fuel.png" alt="Fuel Icon"
                    class="absolute right-1 top-1/2 -translate-y-1/2 w-6 h-6 z-10">
                <div class="absolute top-0 right-0 h-full rounded-md transition-[width] duration-500 ease-in-out bg-gradient-to-l from-gray-700 to-black fuel-shadow"
                    :style="{ width: fuelPercentage + '%' }"></div>
            </div>
        </div>

    </div>
</template>

<script setup>
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { logout as apiLogout } from '@/services/api';

// این کامپوننت داده‌ها را از والد خود (Territory.vue) می‌گیرد
const props = defineProps({
    player: {
        type: Object,
        required: true,
    },
    username: {
        type: String,
        default: '...',
    }
});

const router = useRouter();

const knowledge = computed(() => {
    if (!props.player || !props.player.knowledgeBars || props.player.knowledgeBars.length === 0) {
        return null;
    }
    return props.player.knowledgeBars.find(bar => bar.territoryId === props.player.atTerritory) || props.player.knowledgeBars[0];
});

const fuelPercentage = computed(() => {
    if (!props.player || props.player.fuelCap === 0) return 0;
    return Math.min((props.player.fuel / props.player.fuelCap) * 100, 100);
});

const knowledgePercentage = computed(() => {
    if (!knowledge.value || knowledge.value.total === 0) return 0;
    return Math.min((knowledge.value.value / knowledge.value.total) * 100, 100);
});

function logout() {
    apiLogout();
    router.push({ name: 'Login' });
}
</script>

<style scoped>
.drop-shadow-lg {
    filter: drop-shadow(0 4px 3px rgba(0, 0, 0, 0.5));
}

.drop-shadow-md {
    filter: drop-shadow(0 2px 2px rgba(0, 0, 0, 0.5));
}

.knowledge-shadow {
    box-shadow: 0 0 5px 0px #ff7e5f;
}

.fuel-shadow {
    box-shadow: 0 0 5px 0px #6B7280;
}
</style>