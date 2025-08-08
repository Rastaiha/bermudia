<template>
  <nav class="bg-white shadow p-4 flex space-x-6">
    <router-link 
      to="/" 
      class="text-gray-700 hover:text-blue-600" 
      :class="{ 'font-bold underline': isActive('/') }"
    >
      Home
    </router-link>
    <router-link 
      to="/about" 
      class="text-gray-700 hover:text-blue-600"
      :class="{ 'font-bold underline': isActive('/about') }"
    >
      About
    </router-link>
    <router-link 
      to="/api" 
      class="text-gray-700 hover:text-blue-600"
      :class="{ 'font-bold underline': isActive('/path') }"
    >
      API
    </router-link>
  </nav>
</template>

<script setup>
import { useRoute } from 'vue-router'

const route = useRoute()

function isActive(path) {
  return route.path === path
}
</script>
