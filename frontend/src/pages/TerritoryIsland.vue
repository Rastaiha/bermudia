<script setup>
import Island from '@/components/Island.vue';

// Define the props that will be automatically passed by the router
// This is the definitive fix for the problem
defineProps({
  id: {
    type: String,
    required: true,
  },
  islandId: {
    type: String,
    required: true,
  }
});
</script>

<template>
  <div class="h-screen w-full">
    <Island :id="id" :island-id="islandId" />
  </div>
</template>