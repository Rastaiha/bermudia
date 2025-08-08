<template>
  <nav class="bg-white font-vazir shadow-lg">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        
        <!-- Brand -->
        <div class="flex items-center space-x-3">
          <router-link to="/" class="flex items-center space-x-3">
            <img 
              src="../../public/logo.png" 
              alt="Rasta Logo" 
              class="h-10 w-10 object-contain"
            />
            <span class="text-xl font-bold text-rasta-green">جمع علمی‑ترویجی رستا</span>
          </router-link>
        </div>

        <!-- Desktop Menu -->
        <div class="hidden md:flex items-center space-x-6">
          <router-link
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            class="px-3 py-2 rounded-md transition"
            :class="isActive(item.path) 
              ? 'bg-rasta-green text-white shadow-md' 
              : 'text-gray-700 hover:bg-rasta-green hover:text-white'"
          >
            {{ item.name }}
          </router-link>
        </div>

        <!-- Mobile menu button -->
        <div class="flex items-center md:hidden">
          <button @click="mobileOpen = !mobileOpen" class="focus:outline-none">
            <svg
              v-if="!mobileOpen"
              class="w-6 h-6 text-gray-700"
              fill="none" stroke="currentColor" stroke-width="2"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" 
                d="M4 6h16M4 12h16M4 18h16"/>
            </svg>
            <svg
              v-else
              class="w-6 h-6 text-gray-700"
              fill="none" stroke="currentColor" stroke-width="2"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round"
                d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Mobile Menu -->
    <div v-if="mobileOpen" class="md:hidden bg-white border-t border-gray-200">
      <div class="px-4 py-3 space-y-2">
        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="block px-3 py-2 rounded-md transition"
          :class="isActive(item.path) 
            ? 'bg-rasta-green text-white' 
            : 'text-gray-700 hover:bg-rasta-green hover:text-white'"
          @click="mobileOpen = false"
        >
          {{ item.name }}
        </router-link>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const mobileOpen = ref(false)

const navItems = [
  { name: 'خانه', path: '/' },
  { name: 'درباره ما', path: '/about' },
  { name: 'API', path: '/api' },
]

function isActive(path) {
  return route.path === path
}
</script>

<style scoped>
.bg-rasta-green { background-color: var(--color-rasta-green); }
.text-rasta-green { color: var(--color-rasta-green); }
</style>
