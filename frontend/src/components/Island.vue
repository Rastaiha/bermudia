<template>
  <div class="island-container" :style="{ backgroundImage: `url(${backgroundImage})` }" @mousemove="updateMousePosition">
    <div v-if="chipBox" class="info-box">
      {{ chipBox }}
    </div>

    <div v-if="tooltipText" class="tooltip" :style="tooltipStyle">
      {{ tooltipText }}
    </div>

    <div class="svg-wrapper" v-if="isLoaded" :class="{ fullscreen: isFullscreen }">
      <router-link :to="{ name: 'Territory', params: { id: id } }" @mouseover="showTooltip('بازگشت به نقشه')"
        @mouseleave="hideTooltip">
        <div id="go-back">
          <img src="/images/ships/ship1.svg" alt="Go to the territory" />
        </div>
      </router-link>

      <template v-for="(comp, index) in components" :key="index">
        <div v-if="comp.iframe" class="iframe-holder">
          <button class="fullscreen-button" @click="fullScreen($event.target)"
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
          <iframe :src="comp.iframe.url" frameborder="0" allowfullscreen allow="fullscreen">
          </iframe>
        </div>

        <div v-else-if="comp.input" class="challenge"
          :class="`challenge-${inputsSubmissionState[comp.input.id].status}`">
          <p class="question">{{ comp.input.description }}</p>
          <div class="input-group">
            <input v-if="comp.input.type !== 'file'" :id="comp.input.id" :type="comp.input.type"
              v-model="inputsFormData[comp.input.id]" class="challenge-input"
              placeholder="پاسخ خود را اینجا وارد کنید..."
              :disabled="!inputsSubmissionState[comp.input.id].submittable" />
            <input v-else :id="comp.input.id" :type="comp.input.type" :accept="comp.input.accept?.join(',')"
              @change="handleFileChange($event, comp.input.id)" class="challenge-input"
              placeholder="پاسخ خود را اینجا وارد کنید..."
              :disabled="!inputsSubmissionState[comp.input.id].submittable" />
            <button v-if="inputsSubmissionState[comp.input.id].submittable" @click="submit(comp.input.id)"
              class="submit-button">
              ارسال
            </button>
          </div>
        </div>
      </template>
    </div>
    <div v-else class="loading-message">
      {{ loadingMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick, reactive } from 'vue';
import { useTimeout } from './service/ChainedTimeout';
import { getIsland, submitAnswer } from "@/services/api";

// --- State ---
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

// --- Props ---
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

// --- Methods ---
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
    if (document.exitFullscreen) {
      document.exitFullscreen();
    } else if (document.webkitExitFullscreen) { // Safari
      document.webkitExitFullscreen();
    } else if (document.msExitFullscreen) { // IE11
      document.msExitFullscreen();
    }
  } else {
    if (isLoaded.value) {
      const el = clickedElement.closest('.iframe-holder');
      if (el && el.requestFullscreen) {
        el.requestFullscreen();
      } else if (el && el.webkitRequestFullscreen) { // Safari
        el.webkitRequestFullscreen();
      } else if (el && el.msRequestFullscreen) { // IE11
        el.msRequestFullscreen();
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

// --- Lifecycle ---
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
/* Variables */
.island-container {
  --color-primary: #3498db;
  --color-primary-dark: #2980b9;
  --color-text-light: #ecf0f1;
  --color-bg-card: rgba(22, 32, 42, 0.85);
  --color-border: rgba(52, 152, 219, 0.4);
  --shadow-medium: 0 8px 30px rgba(0, 0, 0, 0.4);
  --border-radius-lg: 24px;
  --border-radius-md: 12px;
  --font-family: sans-serif;
}

/* Layout */
.island-container {
  width: 99vw;
  min-height: 100vh;
  padding: 2rem;
  box-sizing: border-box;
  display: flex;
  justify-content: center;
  align-items: center;
  background-size: cover;
  background-attachment: fixed;
  background-position: center;
  font-family: var(--font-family);
  color: var(--color-text-light);
}

.svg-wrapper {
  width: 100%;
  max-width: 900px;
  display: flex;
  flex-direction: column;
  gap: 2.5rem;
}

.loading-message {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: var(--color-text-light);
  font-size: 1.25rem;
}

/* Components */
.iframe-holder,
.challenge {
  position: relative;
  background-color: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius-lg);
  box-shadow: var(--shadow-medium);
  overflow: hidden;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.iframe-holder {
  padding: 0;
}

.iframe-holder iframe {
  width: 100%;
  height: 550px;
  border: none;
  display: block;
}

.challenge {
  padding: 2.5rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  align-items: center;
  text-align: center;
}

.question {
  font-size: 1.5rem;
  font-weight: 500;
  line-height: 1.6;
}

.input-group {
  display: flex;
  width: 100%;
  max-width: 500px;
  margin-top: 1rem;
  flex-direction: row-reverse;
}

.challenge-input {
  flex-grow: 1;
  border: 1px solid var(--color-border);
  background-color: rgba(0, 0, 0, 0.2);
  color: var(--color-text-light);
  padding: 0.8rem 1rem;
  border-radius: var(--border-radius-md) 0 0 var(--border-radius-md);
  font-size: 1rem;
  outline: none;
  transition: box-shadow 0.2s ease;
}

.challenge-input:focus {
  box-shadow: 0 0 0 3px var(--color-border);
}

.submit-button {
  border: none;
  background-color: var(--color-primary);
  color: white;
  padding: 0.8rem 1.5rem;
  font-size: 1rem;
  font-weight: 600;
  border-radius: 0 var(--border-radius-md) var(--border-radius-md) 0;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.submit-button:hover {
  background-color: var(--color-primary-dark);
}

.fullscreen-button {
  position: absolute;
  top: 1rem;
  left: 1rem;
  z-index: 10;
  background-color: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color 0.2s ease, transform 0.2s ease;
}

.fullscreen-button:hover {
  background-color: rgba(0, 0, 0, 0.7);
  transform: scale(1.1);
}

.info-box,
.tooltip {
  position: fixed;
  background-color: var(--color-bg-card);
  color: var(--color-text-light);
  padding: 0.75rem 1.5rem;
  border-radius: 2rem;
  box-shadow: var(--shadow-medium);
  z-index: 100;
  font-size: 1rem;
  border: 1px solid var(--color-border);
  pointer-events: none;
  white-space: nowrap;
}

.info-box {
  top: 1.5rem;
  left: 50%;
  transform: translateX(-50%);
}

.tooltip {
  padding: 0.5rem 1rem;
  font-size: 0.9rem;
  border-radius: 6px;
}

@media (orientation: portrait) {
  .question {
    font-size: 1rem;
  }

  .challenge-input {
    font-size: 0.8rem;
    width: inherit;
  }

  .submit-button {
    font-size: 0.8rem;
  }

  #go-back {
    width: 15vh;
    height: 15vh;
    bottom: 0;
    right: 0;
  }
}

/* Animations & Effects */
.svg-wrapper.fullscreen .iframe-holder,
.svg-wrapper.fullscreen {
  width: 100%;
  height: 100%;
  max-width: 100%;
  margin: 0;
  padding: 0;
  border-radius: 0;
}

.svg-wrapper.fullscreen iframe {
  height: 100%;
}

.svg-wrapper.fullscreen .fullscreen-button {
  top: 1.5rem;
  left: 1.5rem;
}

#go-back {
  position: fixed;
  z-index: 20;
  width: 15vw;
  height: 15vw;
  bottom: 1.5rem;
  right: 1.5rem;
  filter: drop-shadow(0 5px 20px rgba(0, 0, 0, 0.5));
  animation: 10s linear boat-animation infinite;
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
    transform: translate(10px, 5px) rotate(-10deg)
  }

  70% {
    transform: translate(-10px, 5px) rotate(3deg)
  }

  100% {
    transform: translate(0, 0) rotate(10deg);
  }
}
</style>