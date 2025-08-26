<template>
  <div
    class="w-[99vw] min-h-screen p-8 box-border flex justify-center items-center bg-cover bg-fixed bg-center font-sans text-gray-200"
    :style="{ backgroundImage: `url(${backgroundImage})` }" @mousemove="updateMousePosition">
    <div v-if="chipBox"
      class="fixed top-6 left-1/2 -translate-x-1/2 bg-slate-900/80 text-gray-200 px-6 py-3 rounded-full shadow-lg z-[100] border border-sky-500/40 pointer-events-none whitespace-nowrap">
      {{ chipBox }}
    </div>

    <div v-if="tooltipText"
      class="fixed bg-slate-900/80 text-gray-200 px-4 py-2 rounded-md shadow-lg z-[100] border border-sky-500/40 pointer-events-none whitespace-nowrap text-sm"
      :style="tooltipStyle">
      {{ tooltipText }}
    </div>

    <div class="w-full max-w-4xl flex flex-col gap-10" v-if="isLoaded" :class="{ fullscreen: isFullscreen }">
      <router-link :to="{ name: 'Territory', params: { id: id } }" @mouseover="showTooltip('بازگشت به نقشه')"
        @mouseleave="hideTooltip">
        <div id="go-back">
          <img src="/images/ships/ship1.svg" alt="Go to the territory" />
        </div>
      </router-link>

      <template v-for="(comp, index) in components" :key="index">
        <div v-if="comp.iframe"
          class="relative bg-slate-900/80 border border-sky-500/40 rounded-3xl shadow-lg overflow-hidden transition-all duration-300 ease-linear">
          <button
            class="absolute top-4 left-4 z-10 bg-black/40 border border-white/20 rounded-full w-11 h-11 flex items-center justify-center cursor-pointer transition-all duration-200 hover:bg-black/70 hover:scale-110"
            @click="fullScreen($event.target)"
            @mouseover="showTooltip(isFullscreen ? 'خروج از تمام صفحه' : 'تمام صفحه')" @mouseleave="hideTooltip">
            <svg v-if="!isFullscreen" xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px"
              fill="#FFFFFF">
              <path d="M0 0h24v24H0V0z" fill="none" />
              <path d="M7 14H5v5h5v-2H7v-3zm-2-4h2V7h3V5H5v5zm12 7h-3v2h5v-5h-2v3zM14 5v2h3v3h2V5h-5z" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px"
              fill="#FFFFFF">
              <path d="M0 0h24v24H0V0z" fill="none" />
              <path d="M5 16h3v3h2v-5H5v2zm3-8H5v2h5V5H8v3zm6 11h2v-3h3v-2h-5v5zm2-11V5h-2v5h5V8h-3z" />
            </svg>
          </button>
          <iframe :src="comp.iframe.url" frameborder="0" allowfullscreen allow="fullscreen"
            class="w-full h-[550px] border-none block">
          </iframe>
        </div>

        <div v-else-if="comp.input"
          class="relative bg-slate-900/80 border border-sky-500/40 rounded-3xl shadow-lg overflow-hidden transition-all duration-300 ease-linear p-10 flex flex-col gap-6 items-center text-center"
          :class="`challenge-${inputsSubmissionState[comp.input.id].status}`">
          <p class="text-2xl font-medium leading-relaxed">{{ comp.input.description }}</p>
          <div class="flex w-full max-w-lg mt-4 flex-row-reverse">
            <input v-if="comp.input.type !== 'file'" :id="comp.input.id" :type="comp.input.type"
              v-model="inputsFormData[comp.input.id]"
              class="flex-grow border border-sky-500/40 bg-black/20 text-gray-200 px-4 py-3 rounded-l-xl text-base outline-none transition-shadow duration-200 focus:ring-2 focus:ring-sky-500/60"
              placeholder="پاسخ خود را اینجا وارد کنید..."
              :disabled="!inputsSubmissionState[comp.input.id].submittable" />
            <input v-else :id="comp.input.id" :type="comp.input.type" :accept="comp.input.accept?.join(',')"
              @change="handleFileChange($event, comp.input.id)"
              class="flex-grow border border-sky-500/40 bg-black/20 text-gray-200 px-4 py-3 rounded-l-xl text-base outline-none transition-shadow duration-200 focus:ring-2 focus:ring-sky-500/60"
              placeholder="پاسخ خود را اینجا وارد کنید..."
              :disabled="!inputsSubmissionState[comp.input.id].submittable" />
            <button v-if="inputsSubmissionState[comp.input.id].submittable" @click="submit(comp.input.id)"
              class="border-none bg-blue-500 text-white px-6 py-3 text-base font-semibold rounded-r-xl cursor-pointer transition-colors duration-200 hover:bg-blue-600">
              ارسال
            </button>
          </div>
        </div>
      </template>
    </div>
    <div v-else class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-gray-200 text-xl">
      {{ loadingMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick, reactive } from 'vue';
import { useTimeout } from './service/ChainedTimeout';
import { getIsland, submitAnswer } from "@/services/api";

const isLoaded = ref(false);
const isFullscreen = ref(false);
const backgroundImage = ref('');
const chipBox = ref(null);
const mousePosition = ref({ x: 0, y: 0 });
const loadingMessage = ref('Loading island data...');
const components = ref([]);
const tooltipText = ref('');
const inputsFormData = reactive({});
const inputsSubmissionState = reactive({});
const { startTimeout, clear } = useTimeout();

const props = defineProps({
  id: {
    type: String,
    required: true,
  },
  islandId: {
    type: String,
    required: true,
  },
});

const fetchIslandData = async (id) => {
  try {
    const rawData = await getIsland(id);
    isLoaded.value = true;
    backgroundImage.value = `/images/island/background.png`;
    components.value = rawData.components;
    rawData.components.forEach((c) => {
      if (c.input) {
        inputsFormData[c.input.id] = c.input.type === 'file' ? null : '';
        inputsSubmissionState[c.input.id] = c.input.submissionState;
      }
    });
    await nextTick();
  } catch (error) {
    console.error('Failed to load island data from API:', error);
    loadingMessage.value = `Error loading island: ${error.message}`;
  }
};

const updateMousePosition = (event) => { mousePosition.value = { x: event.clientX, y: event.clientY }; };

function fullScreen(clickedElement) {
  if (isFullscreen.value) {
    if (document.exitFullscreen) document.exitFullscreen();
    else if (document.webkitExitFullscreen) document.webkitExitFullscreen();
    else if (document.msExitFullscreen) document.msExitFullscreen();
  } else {
    if (isLoaded.value) {
      const el = clickedElement.closest('.relative');
      if (el) {
        if (el.requestFullscreen) el.requestFullscreen();
        else if (el.webkitRequestFullscreen) el.webkitRequestFullscreen();
        else if (el.msRequestFullscreen) el.msRequestFullscreen();
      }
    }
  }
}

const handleFileChange = (event, inputId) => {
  const file = event.target.files[0]
  inputsFormData[inputId] = file;
}

async function submit(inputId) {
  const inputValue = inputsFormData[inputId];
  if (!inputValue) {
    chipBox.value = `چیزی برای ارسال وارد نشده‌است!`;
    startTimeout(() => { chipBox.value = null; }, 5000);
    return;
  }
  const formData = new FormData();
  formData.append('data', inputValue);
  try {
    inputsSubmissionState[inputId] = await submitAnswer(inputId, formData);
    chipBox.value = `پاسخ ثبت شد. پس از بررسی نمره آن ثبت می‌شود.`;
  } catch (error) {
    console.error('Error submitting answer:', error);
    chipBox.value = 'خطا در ارسال پاسخ!';
  }
  startTimeout(() => { chipBox.value = null; }, 5000);
}

function onFullscreenChange() {
  isFullscreen.value = !!document.fullscreenElement;
}

const showTooltip = (text) => { tooltipText.value = text; };
const hideTooltip = () => { tooltipText.value = ''; };

const tooltipStyle = computed(() => ({
  position: 'fixed',
  top: `${mousePosition.value.y + 25}px`,
  left: `${mousePosition.value.x}px`,
  transform: 'translateX(-50%)',
}));

onMounted(() => {
  fetchIslandData(props.islandId);
  document.addEventListener('fullscreenchange', onFullscreenChange);
});

onUnmounted(() => {
  document.removeEventListener('fullscreenchange', onFullscreenChange);
  clear();
});
</script>

<style scoped>
/* Un-tailwindable styles */
@media (orientation: portrait) {
  .question {
    font-size: 1rem;
  }

  .challenge input {
    font-size: 0.8rem;
    width: inherit;
  }

  .challenge button {
    font-size: 0.8rem;
  }

  #go-back {
    width: 15vh;
    height: 15vh;
    bottom: 0;
    right: 0;
  }
}

/* Fullscreen Overrides */
.fullscreen {
  width: 100% !important;
  height: 100% !important;
  max-width: 100% !important;
  margin: 0 !important;
  padding: 0 !important;
  border-radius: 0 !important;
  background-color: #2c3e50;
}

.fullscreen iframe {
  height: 100% !important;
}

.fullscreen .fullscreen-button {
  top: 1.5rem;
  left: 1.5rem;
}

/* Go Back Ship */
#go-back {
  position: fixed;
  z-index: 20;
  width: 15vw;
  height: 15vw;
  bottom: 1.5rem;
  right: 1.5rem;
  filter: drop-shadow(0 5px 20px rgba(0, 0, 0, 0.5));
  animation: boat-animation 10s linear infinite;
}

#go-back img {
  object-fit: contain;
  width: 100%;
  height: 100%;
  transition: transform 2s cubic-bezier(0.25, 1, 0.5, 1);
}

#go-back:hover img {
  transform: scale(1.15) rotate(-5deg);
}

@keyframes boat-animation {
  0% {
    transform: translate(0, 0) rotate(10deg);
  }

  35% {
    transform: translate(10px, 5px) rotate(-10deg);
  }

  70% {
    transform: translate(-10px, 5px) rotate(3deg);
  }

  100% {
    transform: translate(0, 0) rotate(10deg);
  }
}
</style>