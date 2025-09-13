<template>
    <div class="fixed bottom-4 left-4 lg:bottom-6 lg:left-6 z-20">
        <div
            class="rounded-xl bg-gray-900 bg-opacity-75 p-2 lg:p-3 shadow-2xl border border-gray-700"
        >
            <div v-if="treasureData.unlocked" class="flex justify-center">
                <img
                    src="/images/island/opened_treasure.png"
                    alt="گنج باز شده"
                    class="h-20 w-20 lg:h-28 lg:w-28 object-contain cursor-not-allowed"
                />
            </div>

            <div v-else>
                <div v-if="!treasureFetchedInfo" class="flex justify-center">
                    <img
                        src="/images/island/closed_treasure.png"
                        alt="در حال بارگذاری..."
                        class="h-20 w-20 lg:h-28 lg:w-28 object-contain opacity-60"
                    />
                </div>

                <div v-else>
                    <div class="flex justify-center">
                        <img
                            src="/images/island/closed_treasure.png"
                            alt="گنج بسته"
                            class="h-20 w-20 lg:h-28 lg:w-28 object-contain transition-transform duration-300 cursor-pointer hover:scale-115"
                            @click="handleTreasureClick"
                        />
                    </div>

                    <div class="mt-2 border-t-2 border-gray-600 pt-2">
                        <div
                            class="flex items-center justify-center gap-x-2 lg:gap-x-4 px-1 lg:px-2"
                        >
                            <div
                                v-for="req in treasureFetchedInfo.cost.items"
                                :key="req.type"
                                class="flex flex-col items-center gap-y-1"
                            >
                                <img
                                    :src="COST_ITEMS_INFO[req.type].icon"
                                    :alt="req.type + ' Icon'"
                                    class="h-7 w-7 lg:h-9 lg:w-9 object-contain"
                                />
                                <span
                                    class="text-sm lg:text-md font-bold text-white text-shadow"
                                    >x{{ req.amount }}</span
                                >
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useModal } from 'vue-final-modal';
import { useToast } from 'vue-toastification';
import { treasureCheck, treasureUnlock } from '@/services/api/index.js';
import TreasureRewardModal from '@/components/features/island/TreasureRewardModal.vue';
import COST_ITEMS_INFO from '@/services/cost.js';

const props = defineProps({
    treasureData: { type: Object, required: true },
    modelValue: { type: Object, required: true },
});

const emit = defineEmits(['update:modelValue']);
const toast = useToast();

const treasureFetchedInfo = ref(null);
const receivedRewards = ref([]);

const { open, close } = useModal({
    component: TreasureRewardModal,
    attrs: {
        rewards: receivedRewards,
        onClose: () => close(),
    },
});

const handleTreasureClick = async () => {
    try {
        const treasure = await treasureUnlock(props.treasureData.id);

        if (treasure.unlocked) {
            receivedRewards.value = treasure.reward;
            emit('update:modelValue', {
                ...props.treasureData,
                unlocked: true,
            });
            open();
        } else {
            toast.error('گنج باز نشد.');
        }
    } catch (error) {
        toast.error(error.message || 'شرایط لازم برای باز کردن گنج را ندارید.');
        console.error('Failed to unlock treasure:', error);
    }
};

onMounted(() => {
    if (props.treasureData.unlocked) return;
    treasureCheck(props.treasureData.id)
        .then(data => (treasureFetchedInfo.value = data))
        .catch(error => {
            console.error('Failed to fetch treasure info:', error);
        });
});
</script>

<style scoped>
.text-shadow {
    text-shadow: 1px 1px 3px rgba(0, 0, 0, 0.8);
}
</style>
