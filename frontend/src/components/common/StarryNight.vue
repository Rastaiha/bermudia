<template>
    <div
        ref="starsContainer"
        class="fixed top-0 left-0 w-full h-full -z-10"
        :style="{
            background: 'linear-gradient(90deg, #01041F 0%, #082B33 100%)',
        }"
    ></div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';

const starsContainer = ref(null);
let animationIntervals = [];

const props = defineProps({
    starCount: {
        type: Number,
        default: 100,
    },
    shootingStarInterval: {
        type: Number,
        default: 8000,
    },
});

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

const createShootingStar = () => {
    const shootingStar = document.createElement('div');
    shootingStar.className = 'shooting-star';
    shootingStar.style.cssText = `
    position: absolute;
    width: 2px;
    height: 10px;
    background-color: #fff;
    opacity: 1;
    z-index: 1;
    top: ${Math.random() * 100}%;
    left: ${Math.random() * 100}%;
    animation: shootingStarAnimation 0.5s linear forwards;
  `;

    document.body.appendChild(shootingStar);

    setTimeout(() => {
        if (document.body.contains(shootingStar)) {
            document.body.removeChild(shootingStar);
        }
    }, 3000);
};

const generateStars = () => {
    if (!starsContainer.value) return;

    for (let i = 0; i < props.starCount; i++) {
        const star = createStar();
        starsContainer.value.appendChild(star);
    }
};

const randomizeShootingStarInterval = () => {
    const interval = Math.random() * props.shootingStarInterval;
    const timeoutId = setTimeout(() => {
        createShootingStar();
        randomizeShootingStarInterval();
    }, interval);

    animationIntervals.push(timeoutId);
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
    
    @keyframes shootingStarAnimation {
      0% {
        opacity: 1;
        transform: translateX(-200px) translateY(-200px) rotate(70deg);
      }
      100% {
        opacity: 0;
        transform: translateX(200px) translateY(-300px) rotate(70deg);
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

onMounted(() => {
    addStyles();
    generateStars();

    // Initial shooting star
    const initialDelay = setTimeout(() => {
        createShootingStar();
    }, Math.random() * props.shootingStarInterval);

    animationIntervals.push(initialDelay);

    // Start the randomized shooting star intervals
    randomizeShootingStarInterval();
});

onUnmounted(() => {
    // Clear all timeouts
    animationIntervals.forEach(clearTimeout);
    animationIntervals = [];

    // Clean up shooting stars that might still be in the DOM
    const shootingStars = document.querySelectorAll('.shooting-star');
    shootingStars.forEach(star => {
        if (document.body.contains(star)) {
            document.body.removeChild(star);
        }
    });

    // Clean up styles when component is destroyed
    cleanupStyles();
});
</script>
