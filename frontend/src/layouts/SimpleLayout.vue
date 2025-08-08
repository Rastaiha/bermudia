<template>
  <div class="min-h-screen flex flex-col">
    <SimpleNavbar />
    <main class="flex-grow container mx-auto p-6">
      <slot />
    </main>
    <SimpleFooter />
  </div>
</template>

<script setup>
import SimpleNavbar from '../components/SimpleNavbar.vue'
import SimpleFooter from '../components/SimpleFooter.vue'
</script>
