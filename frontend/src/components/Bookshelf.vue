<template>
    <VueFinalModal
        class="flex justify-center items-center"
        content-class="flex flex-col w-full md:w-1/2 mx-4 p-6 
                       bg-[#5C3A21] border-4 border-[#3E2A17] 
                       rounded-xl shadow-xl space-y-4"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div
            class="flex items-center justify-between border-b-2 border-[#3E2A17] pb-2 mb-4"
        >
            <h1 class="text-xl font-semibold text-amber-200">
                کتابخانه بنزوئیلا
            </h1>
            <button
                class="p-1 rounded-full hover:bg-[#3E2A17]"
                @click="handleClose"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-6 w-6 text-amber-200"
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
            class="w-full flex flex-col justify-between items-end space-y-2"
        >
            <div
                v-for="(booksInTerritory, territoryId) in booksByTerritory"
                :key="territoryId"
                class="w-full"
            >
                <div class="flex gap-1 items-end pb-2">
                    <div
                        v-for="(book, bookIndex) in booksInTerritory"
                        :key="bookIndex"
                        class="relative w-8 h-32 flex items-center justify-center bg-gradient-to-b from-yellow-600 via-yellow-700 to-yellow-800 border-l-2 border-yellow-900 rounded-sm shadow-md cursor-pointer hover:skew-x-[3deg] hover:skew-y-[3deg] transition-transform"
                        @click="navigateToIsland(territoryId, book.islandId)"
                    >
                        <span
                            class="transform -rotate-90 text-xs font-bold text-white whitespace-nowrap"
                        >
                            {{ book.islandId }}
                        </span>
                        <div
                            class="absolute left-0 top-0 h-full w-1 bg-yellow-900 rounded-l-sm"
                        ></div>
                    </div>
                </div>
                <div class="w-full h-1 bg-[#3E2A17] rounded-sm"></div>
            </div>
        </div>

        <div
            v-else
            class="w-full flex justify-center items-center h-48 text-amber-200 text-lg"
        >
            <p>هنوز هیچ کتابی دریافت نکرده اید</p>
        </div>
    </VueFinalModal>
</template>

<script setup>
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { VueFinalModal } from 'vue-final-modal';

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
