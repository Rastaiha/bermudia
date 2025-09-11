<template>
    <router-link
        :to="{ name: 'Territory', params: { id: territoryId } }"
        @mouseover="$emit('showTooltip', 'بازگشت به نقشه')"
        @mouseleave="$emit('hideTooltip')"
    >
        <div id="go-back">
            <img src="/images/icons/back.png" alt="Go to the territory" />
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
    width: 10vw;
    height: 10vw;
    top: 1.5rem;
    left: 1.5rem;
    filter: drop-shadow(0 5px 20px rgba(0, 0, 0, 0.5));
}

#go-back img {
    object-fit: contain;
    width: 100%;
    height: 100%;
    transition: transform 0.3s ease-in-out;
}

#go-back:hover img {
    transform: scale(1.15);
}

@media (min-width: 1024px) {
    #go-back {
        width: 5vw;
        height: 5vw;
    }
}
</style>
