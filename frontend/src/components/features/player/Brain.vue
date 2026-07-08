<template>
    <VueFinalModal
        class="flex justify-center items-center"
        content-class="panel w-full md:w-1/2 max-h-[80vh]"
        overlay-transition="vfm-fade"
        content-transition="vfm-slide-up"
    >
        <div class="panel-header">
            <h1 class="panel-title">
                {{ glossary.brain }}
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
            v-if="knowledgeBars.length > 0"
            class="panel-body !gap-2 justify-between items-end"
        >
            <div
                v-for="(barData, index) in knowledgeBars"
                :key="index"
                class="w-full"
            >
                <PlayerInventoryBar
                    v-if="barData && migrateData"
                    :bar-data="adopt(barData)"
                />
            </div>
        </div>

        <div v-else class="panel-empty">
            <p>هنوز هیچ {{ glossary.book }}ی دریافت نکرده اید</p>
        </div>
    </VueFinalModal>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { VueFinalModal } from 'vue-final-modal';
import { glossary } from '@/services/glossary.js';
import PlayerInventoryBar from '@/components/common/PlayerInventoryBar.vue';
import { migrateCheck } from '@/services/api';
import { useToast } from 'vue-toastification';

defineProps({
    knowledgeBars: {
        type: Array,
        default: () => [],
    },
});

const migrateData = ref(null);
const toast = useToast();
const emit = defineEmits(['close']);

const adopt = barData => {
    if (!barData || !migrateData.value) return null;
    return {
        name:
            glossary.knowledge +
            ' ' +
            migrateData.value.territoryMigrationOptions.find(ele => {
                return ele.territoryId == barData.territoryId;
            }).territoryName,
        englishName: 'Knowledge for ' + barData.territoryId,
        total: barData.total,
        required:
            migrateData.value.knowledgeCriteriaTerritory == barData.territoryId
                ? migrateData.value.minAcceptableKnowledge
                : null,
        value: barData.value,
        icon: '/images/icons/knowledge.png',
        shadowColor: 'var(--gradient-meter-shadow)',
        gradientFrom: 'var(--gradient-meter-from)',
        gradientTo: 'var(--gradient-meter-to)',
    };
};

function handleClose() {
    emit('close');
}

onMounted(async () => {
    try {
        migrateData.value = await migrateCheck();
    } catch (err) {
        toast.error(err.message);
    }
});
</script>
