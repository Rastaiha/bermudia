<template>
    <VueFinalModal
        class="flex justify-center items-center"
        content-class="panel w-full md:w-1/2 max-h-[80vh]"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div class="panel-header">
            <h1 class="panel-title">
                {{ glossary.benzuelaLibrary }}
            </h1>
            <button class="panel-close" @click="handleClose">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-6 w-6"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M6 18L18 6M6 6l12 12"
                    />
                </svg>
            </button>
        </div>

        <div
            v-if="books.length > 0"
            class="panel-body !gap-2 justify-between items-end"
        >
            <div
                v-for="(booksInTerritory, territoryId) in booksByTerritory"
                :key="territoryId"
                class="w-full"
            >
                <h3 class="text-accent-strong font-semibold mb-2 text-right">
                    {{ booksInTerritory[0].territoryName }}
                </h3>
                <div class="flex gap-1 items-end pb-2 overflow-auto">
                    <div
                        v-for="(book, bookIndex) in booksInTerritory"
                        :key="bookIndex"
                        class="relative min-w-8 w-min h-[7rem] flex items-center justify-center bg-gradient-to-b from-book-from via-book-via to-book-to border-l-2 border-book-edge rounded-sm shadow-md cursor-pointer hover:skew-x-[3deg] hover:skew-y-[3deg] transition-transform"
                        @click="navigateToIsland(territoryId, book.islandId)"
                    >
                        <div
                            class="transform -rotate-90 text-[0.7rem] font-bold p-[0.1px 1rem] text-white text-center"
                        >
                            {{ book.name }}
                        </div>
                        <div
                            class="absolute left-0 top-0 h-full w-1 bg-book-edge rounded-l-sm"
                        ></div>
                    </div>
                </div>
                <div class="w-full h-1 bg-border rounded-sm"></div>
            </div>
        </div>

        <div v-else class="panel-empty">
            <p>هنوز هیچ {{ glossary.book }}ی دریافت نکرده اید</p>
        </div>
    </VueFinalModal>
</template>

<script setup>
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { VueFinalModal } from 'vue-final-modal';
import { glossary } from '@/services/glossary.js';

const props = defineProps({
    books: {
        type: Array,
        default: () => [],
    },
});

const emit = defineEmits(['close']);
const router = useRouter();

const navigateToIsland = (territoryId, islandId) => {
    router.push({
        name: 'Island',
        params: { id: territoryId, islandId: islandId },
    });
};

const booksByTerritory = computed(() => {
    const grouped = {};

    props.books.forEach(book => {
        const territory = book.territoryId;
        if (!grouped[territory]) {
            grouped[territory] = [];
        }
        grouped[territory].push(book);
    });

    return grouped;
});

function handleClose() {
    emit('close');
}
</script>
