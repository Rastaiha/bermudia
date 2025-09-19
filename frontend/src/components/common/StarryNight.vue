<template>
    <div
        ref="starsContainer"
        class="fixed top-0 left-0 w-full h-full -z-10"
        :style="{
            background: 'linear-gradient(180deg, #692A47 0%, #123952 100%)',
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
        default: 0.0001, // Stars per pixel (adjust this value to control density)
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
    star.className = 'star';
    star.style.cssText = `
        position: absolute;
        background-color: #fff;
        width: 1px;
        height: 1px;
        border-radius: 50%;
        top: ${Math.random() * 100}%;
        left: ${Math.random() * 100}%;
    `;

    // Add the ::before pseudo-element effect using a real element
    const twinkle = document.createElement('div');
    twinkle.style.cssText = `
        position: absolute;
        top: -5px;
        left: -5px;
        width: 10px;
        height: 10px;
        border-radius: 50%;
        background-color: rgba(255, 255, 255, 0.8);
        opacity: 0;
        animation: starTwinkle 2s infinite;
        animation-delay: ${Math.random() * 2}s;
    `;

    star.appendChild(twinkle);
    return star;
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
