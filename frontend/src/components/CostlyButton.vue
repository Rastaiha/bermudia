<template>
    <button
        :disabled="loading || !enabled"
        class="btn-hover w-full p-2 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed text-xs"
        :style="{ backgroundColor: backgroundColor, color: textColor }"
        @pointerdown.stop="onClick"
    >
        <div
            v-if="!loading"
            class="w-full h-full flex items-center"
            :class="cost && label ? 'justify-between' : 'justify-evenly'"
        >
            <div v-if="label">{{ label }}</div>
            <div v-if="cost">
                <div
                    v-for="(costItem, index) in cost.items"
                    :key="index"
                    class="flex justify-between items-center flex-col"
                >
                    <div class="flex items-center gap-x-1">
                        <span class="font-bold" :style="{ color: textColor }">{{
                            costItem.amount
                        }}</span>
                        <img
                            :src="getIconByType(costItem.type)"
                            :alt="costItem.type + ' Icon'"
                            class="w-5 h-5"
                        />
                    </div>
                </div>
            </div>
        </div>
    </button>
</template>

<script setup>
defineProps({
    onClick: Function,
    cost: Object,
    label: String,
    enabled: Boolean,
    loading: Boolean,
    backgroundColor: {
        type: String,
        default: 'green',
    },
    textColor: {
        type: String,
        default: 'white',
    },
});

const getIconByType = type => {
    switch (type) {
        case 'blueKey':
            return '/images/icons/blueKeys.png';
        case 'goldenKey':
            return '/images/icons/goldenKeys.png';
        case 'redKey':
            return '/images/icons/redKeys.png';
    }
    return '/images/icons/' + type + '.png';
};
</script>
<style>
button {
    transition:
        transform 0.2s ease,
        filter 0.2s ease;
}

button:hover:not(:disabled) {
    transform: scale(1.05);
    filter: brightness(1.1);
}
</style>
