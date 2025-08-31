<template>
  <div
    class="w-full min-h-screen p-8 box-border flex justify-center items-center bg-cover bg-fixed bg-center font-sans text-gray-200"
    :style="{ backgroundImage: `url(${backgroundImage})` }" @mousemove="updateMousePosition">

    <FloatingUI :tooltipText="tooltipText" :mousePosition="mousePosition" />
    <BackButton :territoryId="id" @showTooltip="showTooltip" @hideTooltip="hideTooltip" v-if="isLoaded" />

    <div v-if="isLoaded" class="w-full max-w-4xl flex flex-col gap-10">
      <template v-for="(componentData) in components" :key="componentData.id">
        <Iframe v-if="componentData.iframe" :url="componentData.iframe.url" @showTooltip="showTooltip"
          @hideTooltip="hideTooltip" />
        <challengeBox v-else-if="componentData.input" :challenge="componentData.input"
          @submit="handleChallengeSubmit" />
      </template>
    </div>

    <div v-else class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-gray-200 text-xl">
      {{ loadingMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useToast } from 'vue-toastification';
import { getIsland, submitAnswer } from "@/services/api";

import Iframe from '@/components/Iframe.vue';
import challengeBox from '@/components/challengeBox.vue';
import BackButton from '@/components/BackButton.vue';
import FloatingUI from '@/components/FloatingUI.vue';

const props = defineProps({
  id: { type: String, required: true },
  islandId: { type: String, required: true },
});

const isLoaded = ref(false);
const backgroundImage = ref('');
const tooltipText = ref('');
const mousePosition = ref({ x: 0, y: 0 });
const loadingMessage = ref('درحال بارگذاری اطلاعات جزیره...');
const components = ref([]);
const toast = useToast();

const fetchIslandData = async (id) => {
  try {
    const rawData = await getIsland(id);
    backgroundImage.value = `/images/island/background.png`;
    components.value = rawData.components;
    isLoaded.value = true;
  } catch (error) {
    console.error('Failed to load island data from API:', error);
    loadingMessage.value = `خطا در بارگذاری جزیره: ${error.message}`;
    toast.error(error.message || 'خطا در دریافت اطلاعات جزیره');
  }
};

const handleChallengeSubmit = async ({ inputId, data }) => {
  if (data === '' || data === null || data === undefined) {
    toast.error(`چیزی برای ارسال وارد نشده‌است!`);
    return;
  }

  const formData = new FormData();
  formData.append('data', data);

  try {
    await submitAnswer(inputId, formData);
    toast.success(`پاسخ شما با موفقیت ثبت شد.`);
    await fetchIslandData(props.islandId);
  } catch (error) {
    console.error('Error submitting answer:', error);
    toast.error(error.message || 'خطا در ارسال پاسخ!');
  }
};

const updateMousePosition = (event) => {
  mousePosition.value = { x: event.clientX, y: event.clientY };
};

const showTooltip = (text) => { tooltipText.value = text; };
const hideTooltip = () => { tooltipText.value = ''; };

onMounted(() => {
  fetchIslandData(props.islandId);
});
</script>

<style scoped>
.w-full {
  width: 100%;
}
</style>
