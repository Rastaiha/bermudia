<script setup>
import { useRoute } from 'vue-router';
import Island from '@/components/Island.vue';

const route = useRoute();
const islandId = route.query.islandId || 'no_island_id_found';
</script>

<template>
  <div class="h-screen w-full">
    <Island :island-id="islandId" />
  </div>
</template>
