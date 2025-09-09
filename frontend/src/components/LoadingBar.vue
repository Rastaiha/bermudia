<template>
    <div class="loading-container w-full max-w-md mx-auto p-4 rounded-lg">
        <div class="relative w-full h-16 flex items-center justify-center">
            <img
                src="/images/ships/spaceship.png"
                alt="Loading spaceship"
                class="absolute z-10 transition-all duration-300 ease-linear object-contain"
                :style="{
                    left: progress + '%',
                    bottom: '50%',
                    transform: `translateX(-25%) translateY(50%)`,
                    height: '4rem',
                    width: '4rem',
                }"
            />

            <div
                class="relative h-6 w-full bg-gray-700 rounded-full overflow-hidden"
            >
                <div
                    class="absolute top-0 left-0 h-full bg-blue-500 rounded-full transition-all duration-300 ease-linear"
                    :style="{ width: progress + '%' }"
                ></div>
            </div>
        </div>

        <p class="text-center text-white mt-3 text-lg font-semibold">
            در حال بارگذاری نقشه... {{ Math.round(progress) }}%
        </p>
    </div>
</template>

<script setup>
defineProps({
    progress: {
        type: Number,
        required: true,
        default: 0,
        validator: value => value >= 0 && value <= 100,
    },
});
</script>

<style scoped>
.loading-container {
    background-color: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(5px);
}
</style>
