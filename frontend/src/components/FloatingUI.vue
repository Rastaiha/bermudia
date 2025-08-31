<template>
  <div>
    <div v-if="chipBoxText"
      class="fixed top-6 left-1/2 -translate-x-1/2 bg-slate-900/80 text-gray-200 px-6 py-3 rounded-full shadow-lg z-[100] border border-sky-500/40 pointer-events-none whitespace-nowrap">
      {{ chipBoxText }}
    </div>

    <div v-if="tooltipText"
      class="fixed bg-slate-900/80 text-gray-200 px-4 py-2 rounded-md shadow-lg z-[100] border border-sky-500/40 pointer-events-none whitespace-nowrap text-sm"
      :style="tooltipStyle">
      {{ tooltipText }}
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  chipBoxText: {
    type: String,
    default: '',
  },
  tooltipText: {
    type: String,
    default: '',
  },
  mousePosition: {
    type: Object,
    default: () => ({ x: 0, y: 0 }),
  },
});

const tooltipStyle = computed(() => ({
  position: 'fixed',
  top: `${props.mousePosition.y + 25}px`,
  left: `${props.mousePosition.x}px`,
  transform: 'translateX(-50%)',
}));
</script>
