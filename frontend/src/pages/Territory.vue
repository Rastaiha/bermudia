<script setup>
import TerritoryMap from '@/components/TerritoryMap.vue';

defineProps({
  id: {
    type: String,
    required: true,
  },
});
</script>

<template>
  <TerritoryMap :territory-id="id" />
</template>