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
        default: 0.0001, // Stars per pixel (adjust this value to control density)
    },
});

const calculateStarCount = () => {
    if (!starsContainer.value) return 100; // fallback

    const rect = starsContainer.value.getBoundingClientRect();
    const screenArea = rect.width * rect.height;
    return Math.floor(screenArea * props.starDensity);
};

const paths = [
    `<path fill="#feffae" d="m90.89392,62.7246c8.16705,-0.20589 18.33738,-0.8922 31.02107,-1.6156c-12.68369,-0.7234 -22.85216,-1.40971 -31.01922,-1.6156c-13.77993,-5.44223 -22.67595,-14.33825 -28.11818,-28.11632c-0.20589,-8.16705 -0.89405,-18.33738 -1.6156,-31.02107c-0.7234,12.68369 -1.40971,22.85217 -1.6156,31.01922c-5.44037,13.77993 -14.33639,22.67595 -28.11818,28.11818c-8.16705,0.20589 -18.33552,0.89405 -31.01922,1.6156c12.68369,0.7234 22.85216,1.40971 31.02107,1.6156c13.77993,5.44037 22.67595,14.33639 28.11818,28.11818c0.20589,8.16705 0.8922,18.33738 1.6156,31.01922c0.7234,-12.68369 1.40971,-22.85402 1.6156,-31.02107c5.43852,-13.77993 14.33454,-22.67595 28.11447,-28.11632z"/>`,
    `<path fill="#feffae" d="M62.43,122.88h-1.98c0-16.15-6.04-30.27-18.11-42.34C30.27,68.47,16.16,62.43,0,62.43v-1.98 
            c16.16,0,30.27-6.04,42.34-18.14C54.41,30.21,60.45,16.1,60.45,0h1.98c0,16.15,6.04,30.27,18.11,42.34 
            c12.07,12.07,26.18,18.11,42.34,18.11v1.98c-16.15,0-30.27,6.04-42.34,18.11C68.47,92.61,62.43,106.72,62.43,122.88z"/>`,
    `<g>
  <path stroke="null" id="svg_48" fill="#feffae" d="m88.85967,89.15262c-4.20493,-6.38453 -8.76979,-14.03409 -13.57122,-23.13504c2.43807,-1.1765 4.978,-2.36115 7.56774,-3.55046c7.36174,-0.06523 16.8763,-0.75249 29.13681,-1.4712c-12.26051,-0.71871 -21.77394,-1.40714 -29.13568,-1.47237c-2.58974,-1.18931 -5.12967,-2.37396 -7.56774,-3.55046c4.80143,-9.10096 9.36516,-16.75051 13.57009,-23.13504c2.47202,-2.85854 5.12628,-5.94305 8.05105,-9.31645c-3.27792,3.00997 -6.27513,5.74154 -9.05276,8.28556c-6.20382,4.32741 -13.63687,9.02524 -22.48023,13.96653c-1.1432,-2.50908 -2.29432,-5.123 -3.44996,-7.78817c-0.06339,-7.57617 -0.73119,-17.36672 -1.43069,-29.98551c-0.69837,12.61763 -1.36731,22.40934 -1.43069,29.98551c-1.15565,2.66517 -2.30677,5.27909 -3.44996,7.78817c-8.84336,-4.94128 -16.27641,-9.63912 -22.48023,-13.96653c-2.77763,-2.54403 -5.77371,-5.27443 -9.05276,-8.2844c2.92477,3.37456 5.57903,6.45792 8.05105,9.31762c4.20493,6.38453 8.76979,14.03292 13.57009,23.13388c-2.43807,1.1765 -4.978,2.36115 -7.56774,3.55046c-7.36174,0.06523 -16.87517,0.75249 -29.13681,1.47237c12.26051,0.71871 21.77507,1.40714 29.13681,1.4712c2.58974,1.18931 5.12967,2.37396 7.56774,3.55046c-4.80143,9.10096 -9.36516,16.75051 -13.57009,23.13504c-2.47202,2.85854 -5.12628,5.94305 -8.05105,9.31762c3.27905,-3.00997 6.27513,-5.74154 9.05276,-8.28556c6.20382,-4.32741 13.63687,-9.02524 22.48023,-13.96653c1.1432,2.50908 2.29432,5.123 3.44996,7.78817c0.06339,7.57617 0.73119,17.36788 1.43069,29.98551c0.69837,-12.61763 1.36731,-22.40934 1.43069,-29.98551c1.15565,-2.66517 2.30677,-5.28026 3.44996,-7.78817c8.84336,4.94128 16.27641,9.63795 22.4791,13.96653c2.77763,2.54403 5.77484,5.2756 9.05389,8.28556c-2.92477,-3.37573 -5.57903,-6.46025 -8.05105,-9.31878z"/>
  <path id="svg_49" fill="#feffae" d="m44.975,91.666c3.257,3.257 3.257,5.817 0,9.074c3.257,-3.257 5.817,-3.257 9.074,0c-3.257,-3.257 -3.257,-5.817 0,-9.074c-3.258,3.258 -5.817,3.258 -9.074,0z"/>
  <path id="svg_50" fill="#feffae" d="m86.339,27.325c-3.257,-3.257 -3.257,-5.817 0,-9.074c-3.257,3.257 -5.817,3.257 -9.074,0c3.257,3.257 3.257,5.817 0,9.074c3.257,-3.257 5.817,-3.257 9.074,0z"/>
  <path id="svg_51" fill="#feffae" d="m32.398,51.228c-3.257,-3.257 -3.257,-5.817 0,-9.074c-3.257,3.257 -5.817,3.257 -9.074,0c3.257,3.257 3.257,5.817 0,9.074c3.257,-3.257 5.816,-3.257 9.074,0z"/>
 </g>`,
];

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
        filter: drop-shadow(0 0 5px rgba(255, 255, 201, 0.564)) blur(0.5px);
    `;
    const size = 8 + Math.random() * 10;
    star.style.width = `${size}px`;
    star.style.height = `${size}px`;

    star.innerHTML =
        `
        <svg viewBox="0 0 122.88 122.88" xmlns="http://www.w3.org/2000/svg">
        ` +
        paths[2] +
        `
        </svg>
    `;
    //Math.floor(Math.random() * paths.length)

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
