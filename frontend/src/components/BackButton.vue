<template>
    <router-link :to="{ name: 'Territory', params: { id: territoryId } }"
        @mouseover="$emit('showTooltip', 'بازگشت به نقشه')" @mouseleave="$emit('hideTooltip')">
        <div id="go-back">
            <img src="/images/ships/ship1.svg" alt="Go to the territory" />
        </div>
    </router-link>
</template>

<script setup>
defineProps({
    territoryId: {
        type: String,
        required: true,
    },
});

defineEmits(['showTooltip', 'hideTooltip']);
</script>

<style scoped>
#go-back {
    position: fixed;
    z-index: 20;
    width: 15vw;
    height: 15vw;
    bottom: 1.5rem;
    right: 1.5rem;
    filter: drop-shadow(0 5px 20px rgba(0, 0, 0, 0.5));
    animation: boat-animation 10s linear infinite;
}

#go-back img {
    object-fit: contain;
    width: 100%;
    height: 100%;
    transition: transform 2s cubic-bezier(0.25, 1, 0.5, 1);
}

#go-back:hover img {
    transform: scale(1.15) rotate(-5deg);
}

@keyframes boat-animation {
    0% {
        transform: translate(0, 0) rotate(10deg);
    }

    35% {
        transform: translate(10px, 5px) rotate(-10deg);
    }

    70% {
        transform: translate(-10px, 5px) rotate(3deg);
    }

    100% {
        transform: translate(0, 0) rotate(10deg);
    }
}
</style>