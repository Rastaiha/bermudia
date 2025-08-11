<template>
  <div class="island-container" :style="{backgroundImage: `url(${backgroundImage})`}">
    <div v-if="chipBox" class="info-box" :style="infoBoxStyle">
      {{ chipBox.name }}
    </div>

    <div class="svg-wrapper" v-if="isLoaded" :class="{ fullscreen: isFullscreen }">
      <template v-for="(comp, index) in components" :key="index">
        <div v-if="comp.iframe" class="iframe-holder">
          <button class="fullscreen-button" @click="fullScreen($event.target)">
            حالت تمام صفحه
          </button>
          <iframe
            :src="comp.iframe.url"
            frameborder="0"
          >
          </iframe>
        </div>

        <input
          v-else-if="comp.input"
          :id="comp.input.id"
          :type="comp.input.type"
          :placeholder="comp.input.description"
          :accept="comp.input.accept?.join(',')"
        />
      </template>
    </div>
    <div v-else class="loading-message">
      {{ loadingMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick } from 'vue';

const mapWidth = ref(window.innerWidth);
const mapHeight = ref(window.innerHeight);

// --- Define reactive state ---
const svgRef = ref(null);
const isLoaded = true;
const isFullscreen = ref(false);
const backgroundImage = ref('');
const chipBox = ref(null);
const iframe = ref(null);
const challenge = ref(null);
const mousePosition = ref({ x: 0, y: 0 });
const loadingMessage = ref('Loading island data...');
const components = ref([]);

// --- Define component props ---
const props = defineProps({
  islandId: {
    type: [String],
    required: true,
  },
});

// --- Fetch and process data from the REAL API ---
const fetchTerritoryData = async (id) => {
  const BASE_URL = 'http://97590f57-b983-48f8-bb0a-c098bed1e658.hsvc.ir:30254/api/v1';
  const apiUrl = `${BASE_URL}/islands/${id}`;
  
  try {
    console.log(`Initiating fetch request to REAL API endpoint: ${apiUrl}`);
    const response = await fetch(apiUrl);
    const data = await response.json();

    if (!data.ok || !data.result) {
      throw new Error(data.error || 'Invalid API response format');
    }

    const rawData = data.result; // Extract data from the 'result' key
    //TODO get the image path from the api
    backgroundImage.value = `/images/island/background.png`;
    components.value = rawData.components;

    
    await nextTick();
  } catch (error) {
    console.error('Failed to load territory data from API:', error);
    loadingMessage.value = `Error loading map: ${error.message}`;
  }
};


// --- Helper Functions ---
const updateMousePosition = (event) => { mousePosition.value = { x: event.clientX, y: event.clientY }; };
const showInfoBox = (node) => { hoveredNode.value = node; };
const hideInfoBox = () => { hoveredNode.value = null; };
const infoBoxStyle = computed(() => ({
  position: 'fixed',
  top: `${mousePosition.value.y + 20}px`,
  left: `${mousePosition.value.x}px`,
  transform: 'translateX(-50%)',
}));

function resolveComponent(comp) {
  if (comp.iframe) return 'iframe'
  if (comp.input) return 'input'
  return 'div'
}

function getProps(comp) {
  if (comp.iframe) {
    return {
      src: comp.iframe.url,
      width: '100%',
      height: '300',
      frameborder: '0'
    }
  }
  if (comp.input) {
    return {
      id: comp.input.id,
      type: comp.input.type,
      placeholder: comp.input.description,
      accept: comp.input.accept?.join(',')
    }
  }
  return {}
}

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

// --- Lifecycle Hooks ---
onMounted(() => {
  fetchTerritoryData(props.islandId);
});
onUnmounted(() => {
  if (panzoomInstance) panzoomInstance.dispose();
});
</script>

<style scoped>
/* Styles remain the same */
.island-container {
  --color-edge: #ffffff;
  --color-info-box-bg: rgba(0, 0, 0, 0.8);
  --color-info-box-text: white;
  --color-loading-text: #ddd;
  --color-bg-fallback: #0c2036;

  width: 100vw;
  position: relative;
  background-color: var(--color-bg-fallback);
  overflow: hidden;
  background-size: 100vw 100vh;
  background-repeat: no-repeat;
}

.svg-wrapper iframe {
  margin: 3rem auto;
  width: 65%;
  height: 500px;
  border-radius: 20px;
  filter: drop-shadow(2px 4px 6px black);
  border: 5px ridge #C2DFFF;
}

.info-box {
  position: absolute;
  background-color: var(--color-info-box-bg);
  color: var(--color-info-box-text);
  padding: 8px 12px;
  border-radius: 6px;
  font-family: sans-serif;
  font-size: 14px;
  pointer-events: none;
  z-index: 100;
  white-space: nowrap;
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
}

button.fullscreen-button {
  background: linear-gradient(45deg, #7BCCB5, #C2DFFF, #7BCCB5);
  padding: 8px;
  border-radius: 20px;
  filter: drop-shadow(2px 4px 6px black);
  border: 5px ridge #C2DFFF;
  position: absolute;
  left: calc(50% - 65% / 2);
  margin: 3px;
  z-index: 1;
  color: #123456;
  font-family: sans-serif;
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

</style>