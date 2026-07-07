/// <reference types="vitest/config" />
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';
import tailwindcss from '@tailwindcss/vite';
import svgLoader from 'vite-svg-loader';

export default defineConfig({
    resolve: {
        alias: {
            '@': path.resolve(__dirname, './src'),
        },
    },
    plugins: [vue(), tailwindcss(), svgLoader()],
    test: {
        environment: 'jsdom',
        globals: true,
        restoreMocks: true,
        include: ['tests/**/*.{test,spec}.{js,mjs,jsx}'],
    },
});
