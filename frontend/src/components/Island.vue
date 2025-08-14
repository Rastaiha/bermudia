<template>
  <div class="island-container" :style="{backgroundImage: `url(${backgroundImage})`}" @mousemove="updateMousePosition">
    <div v-if="chipBox" class="info-box">
      {{ chipBox }}
    </div>
    
    <div v-if="tooltipText" class="tooltip" :style="tooltipStyle">
      {{ tooltipText }}
    </div>

    <div class="svg-wrapper" v-if="isLoaded" :class="{ fullscreen: isFullscreen }">
      <router-link :to="`/territory/${id}`" @mouseover="showTooltip('بازگشت به نقشه')" @mouseleave="hideTooltip"> 
        <img id="go-back" src="/images/ships/ship1.svg" alt="Go to the territory" />
      </router-link>

      <template v-for="(comp, index) in components" :key="index">
        <div v-if="comp.iframe" class="iframe-holder">
          <button class="fullscreen-button" @click="fullScreen($event.target)" @mouseover="showTooltip(isFullscreen ? 'خروج از تمام صفحه' : 'تمام صفحه')" @mouseleave="hideTooltip">
            <svg v-if="!isFullscreen" xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#FFFFFF"><path d="M0 0h24v24H0V0z" fill="none"/><path d="M7 14H5v5h5v-2H7v-3zm-2-4h2V7h3V5H5v5zm12 7h-3v2h5v-5h-2v3zM14 5v2h3v3h2V5h-5z"/></svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 0 24 24" width="24px" fill="#FFFFFF"><path d="M0 0h24v24H0V0z" fill="none"/><path d="M5 16h3v3h2v-5H5v2zm3-8H5v2h5V5H8v3zm6 11h2v-3h3v-2h-5v5zm2-11V5h-2v5h5V8h-3z"/></svg>
          </button>
          <iframe
            :src="comp.iframe.url"
            frameborder="0"
            allowfullscreen
            allow="fullscreen"
          >
          </iframe>
        </div>

        <div v-else-if="comp.input" class="challenge">
          <p class="question">{{ comp.input.description }}</p>
          <div class="input-group">
            <input
              :id="comp.input.id"
              :type="comp.input.type"
              :accept="comp.input.accept?.join(',')"
              v-model="comp.input.answer"
              class="challenge-input"
              placeholder="پاسخ خود را اینجا وارد کنید..."
            />
            <button @click="submit(comp.input)" class="submit-button">ارسال</button>
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
import { ref, onMounted, onUnmounted, computed, nextTick } from 'vue';
import {useTimeout} from './service/ChainedTimeout';

const mapWidth = ref(window.innerWidth);
const mapHeight = ref(window.innerHeight);

// --- Define reactive state ---
const BASE_URL = 'http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1';
const svgRef = ref(null);
const isLoaded = ref(false);
const isFullscreen = ref(false);
const backgroundImage = ref('');
const chipBox = ref(null);
const iframe = ref(null);
const challenge = ref(null);
const mousePosition = ref({ x: 0, y: 0 });
const loadingMessage = ref('Loading island data...');
const components = ref([]);
const { startTimeout, clear } = useTimeout()
// THE CHANGE IS HERE: Added state for the new tooltip
const tooltipText = ref('');

// --- Define component props ---
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

// --- Fetch and process data from the REAL API ---
const fetchIslandData = async (id) => {
  const apiUrl = `${BASE_URL}/islands/${id}`;
  
  try {
    loadingMessage.value = 'Fetching island data from server...';
    console.log(`Initiating fetch request to REAL API endpoint: ${apiUrl}`);
    const response = await fetch(apiUrl);
    const data = await response.json();

    if (!data.ok || !data.result) {
      throw new Error(data.error || 'Invalid API response format');
    }

    isLoaded.value = true;

    const rawData = data.result;
    backgroundImage.value = `/images/island/background.png`;
    components.value = rawData.components;

    
    await nextTick();
  } catch (error) {
    console.error('Failed to load island data from API:', error);
    loadingMessage.value = `Error loading island: ${error.message}`;
  }
};


// --- Helper Functions ---
const updateMousePosition = (event) => { mousePosition.value = { x: event.clientX, y: event.clientY }; };

// THE ONLY CHANGE IS THIS FUNCTION: It's now robust against clicks on the SVG icon
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
    // Check if content is ready before making fullscreen
    if(isLoaded.value) {
        // Use .closest() to reliably find the iframe-holder, no matter what was clicked
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

async function submit(input) {
  const inputValue = input.answer;
  const formData = new FormData();
  
  if (!inputValue) {
    chipBox.value = `چیزی برای ارسال وارد نشده‌است!`;
    startTimeout(() => {
      chipBox.value = null;
    }, 5000);
    return;
  }
  
  if (typeof inputValue === 'object' && inputValue instanceof FileList) {
    for (let i = 0; i < inputValue.length; i++) {
      formData.append(input.id, inputValue[i]);
    }
  } else {
    formData.append(input.id, inputValue);
  }

  try {
    const apiUrl = `${BASE_URL}/answer/${input.id}`;
    const response = await fetch(apiUrl, {
      method: 'POST',
      body: formData,
    });
    const data = await response.json();
    console.log('Submit response:', data);
    if (data.ok) {
      chipBox.value = `پاسخ ثبت شد. پس از بررسی نمره آن ثبت می‌شود.`;
    } else {
      chipBox.value = data.error; 
    }
    startTimeout(() => {
      chipBox.value = null;
    }, 5000);
  } catch (error) {
    console.error('Error submitting answer:', error);
  }
}

function onFullscreenChange() {
  if (!document.fullscreenElement) {
    isFullscreen.value = false;
  } else {
    isFullscreen.value = true;
  }
}

// THE CHANGE IS HERE: Added functions and computed property for the new tooltip
const showTooltip = (text) => {
  tooltipText.value = text;
};
const hideTooltip = () => {
  tooltipText.value = '';
};
const tooltipStyle = computed(() => ({
  position: 'fixed',
  top: `${mousePosition.value.y + 25}px`,
  left: `${mousePosition.value.x}px`,
  transform: 'translateX(-50%)',
}));

// --- Lifecycle Hooks ---
onMounted(() => {
  fetchIslandData(props.islandId);
  document.addEventListener('fullscreenchange', onFullscreenChange);
});
onUnmounted(() => {
  document.removeEventListener('fullscreenchange', onFullscreenChange);
});
</script>

<style scoped>
/* --- THE ONLY CHANGE IS THIS ENTIRE STYLE BLOCK --- */

/* --- 1. Define Color Palette & Variables --- */
.island-container {
  --color-primary: #3498db;
  --color-primary-dark: #2980b9;
  --color-text-light: #ecf0f1;
  --color-text-dark: #2c3e50;
  --color-bg-dark: #2c3e50;
  --color-bg-card: rgba(22, 32, 42, 0.85); /* Slightly darker card */
  --color-border: rgba(52, 152, 219, 0.4);
  --shadow-medium: 0 8px 30px rgba(0, 0, 0, 0.4);
  --border-radius-lg: 24px; /* Softer, rounder corners */
  --border-radius-md: 12px;
  --font-family: sans-serif;
}

/* --- 2. Main Layout --- */
.island-container {
  width: 100vw;
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

/* --- 3. Content Cards (Iframe & Challenge) --- */
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

/* --- 4. Buttons --- */
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

/* --- 5. Fullscreen Fix --- */
.svg-wrapper.fullscreen .iframe-holder,
.svg-wrapper.fullscreen {
  width: 100%;
  height: 100%;
  max-width: 100%;
  margin: 0;
  padding: 0;
  border-radius: 0;
  background-color: var(--color-bg-dark); 
}
.svg-wrapper.fullscreen iframe {
  height: 100%;
}
.svg-wrapper.fullscreen .fullscreen-button {
  top: 1.5rem;
  left: 1.5rem;
}

/* --- 6. Go Back Ship Fix --- */
#go-back {
  position: fixed;
  z-index: 20;
  width: 150px;
  height: 150px;
  bottom: 1.5rem;
  right: 1.5rem;
  filter: drop-shadow(0 5px 20px rgba(0,0,0,0.5));
  transition: transform 0.3s cubic-bezier(0.25, 1, 0.5, 1);
  object-fit: contain;
}
#go-back:hover {
  transform: scale(1.15) rotate(-5deg);
}

/* --- 7. Floating Info & Loading & NEW TOOLTIP --- */
.info-box, .tooltip {
  position: fixed;
  background-color: var(--color-bg-dark);
  color: var(--color-text-light);
  padding: 0.75rem 1.5rem;
  border-radius: 2rem;
  box-shadow: var(--shadow-medium);
  z-index: 100;
  font-size: 1rem;
  border: 1px solid var(--color-border);
  pointer-events: none; /* Important for tooltips */
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

.loading-message {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: var(--color-text-light);
  font-size: 1.25rem;
}
</style>