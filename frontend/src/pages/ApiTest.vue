<script setup>
import SimpleLayout from "../layouts/SimpleLayout.vue";
import { ref, onMounted } from "vue";
import { fetchPhotos } from "../services/api";

const photos = ref([]);
const page = ref(1);
const isLoading = ref(false);
const hasError = ref(false);

const loadPhotos = async () => {
  isLoading.value = true;
  hasError.value = false;
  try {
    const newPhotos = await fetchPhotos(page.value, 12);
    photos.value.push(...newPhotos);
    page.value += 1;
  } catch (err) {
    console.error("Failed to fetch photos:", err);
    hasError.value = true;
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  loadPhotos();
});
</script>

<template>
  <SimpleLayout>
    <section class="max-w-7xl mx-auto px-6 py-8">
      <h1
        class="text-4xl font-extrabold text-gray-900 mb-8 text-center"
      >
        Random Photo Gallery
      </h1>
      <div
        class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6"
      >
        <div
          v-for="photo in photos"
          :key="photo.id"
          class="rounded-lg overflow-hidden shadow-lg hover:shadow-2xl transition-shadow duration-300 bg-white"
        >
          <img
            :src="`https://picsum.photos/id/${photo.id}/400/250`"
            :alt="photo.author"
            class="w-full h-56 object-cover"
            loading="lazy"
          />
          <div
            class="p-4 text-center text-gray-700 font-medium text-sm select-none"
          >
            {{ photo.author }}
          </div>
        </div>
      </div>

      <div class="mt-8 flex justify-center">
        <button
          @click="loadPhotos"
          :disabled="isLoading"
          class="px-6 py-3 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-blue-300 transition"
        >
          <span v-if="isLoading">Loading...</span>
          <span v-else>Load More</span>
        </button>
      </div>

      <p
        v-if="hasError"
        class="mt-4 text-center text-red-600 font-semibold"
      >
        Failed to load more photos. Please try again.
      </p>
    </section>
  </SimpleLayout>
</template>
