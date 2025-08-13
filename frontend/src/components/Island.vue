<template>
  <div class="island-container" :style="{backgroundImage: `url(${backgroundImage})`}">
    <div v-if="chipBox" class="info-box">
      {{ chipBox }}
    </div>

    <div class="svg-wrapper" v-if="isLoaded" :class="{ fullscreen: isFullscreen }">
      <a href="/territory"> 
        <img id="go-back" src="/images/ships/ship1.svg" alt="Go to the territory" />
      </a>

      <template v-for="(comp, index) in components" :key="index">
        <div v-if="comp.iframe" class="iframe-holder">
          <button class="fullscreen-button background-moving-animation" @click="fullScreen($event.target)">
            حالت تمام صفحه
          </button>
          <iframe
            :src="comp.iframe.url"
            frameborder="0"
            allowfullscreen
            allow="fullscreen"
          >
          </iframe>
        </div>

        <div v-else-if="comp.input" class="challenge background-moving-animation">
          <div class="question">{{ comp.input.description }}</div>
          <input
            :id="comp.input.id"
            :type="comp.input.type"
            :accept="comp.input.accept?.join(',')"
            v-model="comp.input.answer"
            class="background-moving-animation"
          />
          <button @click="submit(comp.input)" class="background-moving-animation">تایید پاسخ!</button>

        </div>
      </template>
    </div>
    <div v-else class="loading-message">
      <span> آیدی جزیره: {{props.islandId}} </span> <br/>
      <span> {{ loadingMessage }} </span>
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

// --- Define component props ---
const props = defineProps({
  islandId: {
    type: [String],
    required: true,
  },
});

// --- Fetch and process data from the REAL API ---
const fetchTerritoryData = async (id) => {
  const apiUrl = `${BASE_URL}/islands/${id}`;
  
  try {
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
    console.error('Failed to load territory data from API:', error);
    loadingMessage.value = `Error loading island: ${error.message}`;
  }
};


// --- Helper Functions ---
const updateMousePosition = (event) => { mousePosition.value = { x: event.clientX, y: event.clientY }; };

function fullScreen(button) {
  if (isFullscreen.value) {
    isFullscreen.value = false;
    if (document.exitFullscreen) {
      document.exitFullscreen();
    } else if (document.webkitExitFullscreen) { // Safari
      document.webkitExitFullscreen();
    } else if (document.msExitFullscreen) { // IE11
      document.msExitFullscreen();
    }
  } else {
    isFullscreen.value = true;
    const el = button.parentElement;
    if (el.requestFullscreen) {
      el.requestFullscreen();
    } else if (el.webkitRequestFullscreen) { // Safari
      el.webkitRequestFullscreen();
    } else if (el.msRequestFullscreen) { // IE11
      el.msRequestFullscreen();
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
      chipBox.value = e.error;
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

// --- Lifecycle Hooks ---
onMounted(() => {
  fetchTerritoryData(props.islandId);
  document.addEventListener('fullscreenchange', onFullscreenChange);
});
onUnmounted(() => {
  if (panzoomInstance) panzoomInstance.dispose();
  document.removeEventListener('fullscreenchange', onFullscreenChange);
});
</script>

<style scoped>
/* Styles remain the same */
.island-container {
  --color-edge: #ffffff;
  --color-info-box-bg: #123456df;
  --color-info-box-text: white;
  --color-loading-text: #ddd;
  --color-bg-fallback: #0c2036;
  --iframe-width: 65%;

  width: 98vw;
  min-height: 100vh;
  position: relative;
  background-color: var(--color-bg-fallback);
  overflow: hidden;
  background-size: 100vw 100vh;
  background-attachment: fixed;
  background-repeat: no-repeat;
}

.svg-wrapper iframe {
  margin: 3rem auto;
  width: var(--iframe-width);
  height: 500px;
  border-radius: 20px;
  filter: drop-shadow(2px 4px 6px black);
  border: 5px ridge #C2DFFF;
}

.info-box {
  background-color: var(--color-info-box-bg);
  color: var(--color-info-box-text);
  padding: 8px 12px;
  border-radius: 6px;
  font-family: sans-serif;
  font-size: 14px;
  pointer-events: none;
  z-index: 100;
  max-width: 98vw;
  position: fixed;
  top: 10px;
  left: 10px;
}
.loading-message {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 1.5rem;
  color: var(--color-loading-text);
  font-family: sans-serif;
  background-color: var(--color-info-box-bg);
  padding: 1rem 2rem;
  border-radius: 0.5rem;
  text-align: center;
}

.background-moving-animation {
  background: linear-gradient(45deg, #7BCCB5, #C2DFFF, #7BCCB5, #C2DFFF, #7BCCB5, #C2DFFF, #7BCCB5, #C2DFFF, #7BCCB5);
  background-size: 400% 100%;
  background-repeat: no-repeat;
  animation: 15s cubic-bezier(1, 0, 0.73, 1) background-moving-d6e042fc 5s infinite;
}

button.fullscreen-button, .challenge, .challenge input, .challenge button {
  filter: drop-shadow(2px 4px 6px black);
  border-radius: 20px;
  border: 5px ridge #C2DFFF;
  color: #123456;
  font-family: sans-serif;
  opacity: 0.98;  
  max-width: -webkit-fill-available;
}

button.fullscreen-button {
  padding: 8px;
  position: absolute;
  left: calc(50% - var(--iframe-width) / 2);
  margin: 3px;
  z-index: 1;
}

.challenge {
  width: var(--iframe-width);
  min-height: 200px;
  margin: 3rem auto;
  padding: 20px;
}

.challenge input, .challenge button {
  min-width: 50%;
  display: block;
  margin: 50px auto;
  filter: drop-shadow(2px 4px 6px white) invert(1);
  height: 50px;
  text-align: center;
  padding: 5px;
}

.svg-wrapper.fullscreen iframe {
  width: 100%;
  height: 100%;
  margin: 0;
}

.svg-wrapper.fullscreen button.fullscreen-button {
  left: 0;
  top: 0;
}

#go-back {
    position: fixed;
    z-index: 2;
    width: 20vw;
    bottom: 5vh;
    filter: drop-shadow(6px 7px 13px black);
}

button:hover, a>img:hover {
    transform: translate(0, 5px);
    filter: hue-rotate(45deg) !important;
}

@media (orientation: portrait) {
  .svg-wrapper {
    --iframe-width: 90%;
  }

  .island-container {
    background-size: 180vw 100vh;
    background-position: -40vw 0;
  }

  button.fullscreen-button {
    margin: auto;
    left: 50%;
    transform: translate(-50%, -35%);
  }

  .svg-wrapper.fullscreen button.fullscreen-button {
    left: 50%;
  }

  #go-back {
    bottom: -3vh;
  }

  .info-box {
    left: 50%;
    transform: translateX(-50%);
    white-space: normal;
    overflow-wrap: break-word; 
    width: 98vw;
  }
}

@keyframes background-moving {
  0% {
    background-position: 0 0;
  }
  10% {
    background-position: 100% 0;
  }
  80% {
    background-position: 0% 0;
  }
}

</style>