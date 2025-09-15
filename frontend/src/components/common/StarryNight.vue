<template>
    <div
        ref="starsContainer"
        class="fixed top-0 left-0 w-full h-full -z-10"
        :style="{
            background:
                'linear-gradient(180deg, rgb(1, 4, 31) 50%, rgb(8, 43, 51) 100%)',
        }"
    ></div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue';

const starsContainer = ref(null);
let resizeObserver = null;

const props = defineProps({
    starDensity: {
        type: Number,
        default: 0.0005, // Stars per pixel (adjust this value to control density)
    },
});

const calculateStarCount = () => {
    if (!starsContainer.value) return 100; // fallback

    const rect = starsContainer.value.getBoundingClientRect();
    const screenArea = rect.width * rect.height;
    return Math.floor(screenArea * props.starDensity);
};

const createStar = () => {
    const star = document.createElement('div');
    const starContainer = document.createElement('div');
    starContainer.appendChild(star);
    star.className = 'star';
    starContainer.className = 'starContainer';
    starContainer.style.cssText = `
        position: absolute;
        top: ${Math.random() * 100}%;
        left: ${Math.random() * 100}%;
    `;
    const size = 6 + Math.random() * 8; // between 6px and 14px
    star.style.width = `${size}px`;
    star.style.height = `${size}px`;

    star.innerHTML = `
        <svg viewBox="0 0 122.88 122.88" xmlns="http://www.w3.org/2000/svg">
            <path fill="white" d="M62.43,122.88h-1.98c0-16.15-6.04-30.27-18.11-42.34C30.27,68.47,16.16,62.43,0,62.43v-1.98 
            c16.16,0,30.27-6.04,42.34-18.14C54.41,30.21,60.45,16.1,60.45,0h1.98c0,16.15,6.04,30.27,18.11,42.34 
            c12.07,12.07,26.18,18.11,42.34,18.11v1.98c-16.15,0-30.27,6.04-42.34,18.11C68.47,92.61,62.43,106.72,62.43,122.88z"/>
        </svg>
    `;

    star.style.animation = `starTwinkle 3s infinite`;
    star.style.animationDelay = `${Math.random() * 2}s`;
    return starContainer;
};

const clearStars = () => {
    if (!starsContainer.value) return;

    const existingStars = starsContainer.value.querySelectorAll('.star');
    existingStars.forEach(star => star.remove());
};

const generateStars = () => {
    if (!starsContainer.value) return;

    clearStars();
    const starCount = calculateStarCount();

    for (let i = 0; i < starCount; i++) {
        const star = createStar();
        starsContainer.value.appendChild(star);
    }
};

const addStyles = () => {
    // Check if styles already exist
    if (document.getElementById('starry-night-styles')) return;

    const style = document.createElement('style');
    style.id = 'starry-night-styles';
    style.textContent = `
        @keyframes starTwinkle {
            0% {
                transform: scale(0.3);
                opacity: 0;
            }
            50% {
                transform: scale(0.5);
                opacity: 1;
            }
            100% {
                transform: scale(0.5);
                opacity: 0;
            }
        }
    `;

    document.head.appendChild(style);
};

const cleanupStyles = () => {
    const existingStyles = document.getElementById('starry-night-styles');
    if (existingStyles) {
        document.head.removeChild(existingStyles);
    }
};

const handleResize = () => {
    // Debounce the resize to avoid too many recalculations
    clearTimeout(handleResize.timeoutId);
    handleResize.timeoutId = setTimeout(() => {
        generateStars();
    }, 100);
};

onMounted(async () => {
    addStyles();

    await nextTick();

    generateStars();

    // Set up ResizeObserver to watch for container size changes
    if (window.ResizeObserver && starsContainer.value) {
        resizeObserver = new ResizeObserver(handleResize);
        resizeObserver.observe(starsContainer.value);
    } else {
        // Fallback to window resize event for older browsers
        window.addEventListener('resize', handleResize);
    }
});

onUnmounted(() => {
    if (resizeObserver) {
        resizeObserver.disconnect();
    } else {
        window.removeEventListener('resize', handleResize);
    }

    if (handleResize.timeoutId) {
        clearTimeout(handleResize.timeoutId);
    }
    cleanupStyles();
});
</script>
<style>
.star {
    filter: drop-shadow(0 0 5px rgba(255, 255, 201, 0.564)) blur(1px);
}
.starContainer {
    padding: 10px;
    border-radius: 50%;
    background: radial-gradient(
        #8c8c8530,
        #8c8c8501,
        #8c8c8501,
        #8c8c8501,
        transparent
    );
}
</style>
