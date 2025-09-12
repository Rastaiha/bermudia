<script setup>
import Login from '@/components/layout/LoginTemplate.vue';
import StarryNight from '@/components/common/StarryNight.vue';
</script>

<template>
    <div class="relative min-h-screen overflow-hidden">
        <StarryNight :star-count="200" :shooting-star-interval="6000" />
        <Login />
    </div>
</template>
