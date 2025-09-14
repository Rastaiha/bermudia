<template>
    <div
        ref="treasureContainer"
        class="fixed bottom-4 left-4 z-20 lg:bottom-6 lg:left-6"
        :class="{
            'transition-transform duration-300 ease-in-out': isMobile,
            'translate-y-[calc(100%+1rem)]': !shouldBeOpen,
            'translate-y-0': shouldBeOpen,
        }"
    >
        <div class="relative">
            <div
                class="rounded-xl bg-gray-900 bg-opacity-75 p-2 lg:p-3 shadow-2xl border border-gray-700 w-fit"
            >
                <div v-if="treasureData.unlocked" class="flex justify-center">
                    <img
                        src="/images/icons/opened_treasure.png"
                        alt="گنج باز شده"
                        class="h-20 w-20 lg:h-28 lg:w-28 object-contain cursor-not-allowed"
                    />
                </div>

                <div v-else>
                    <div
                        v-if="!treasureFetchedInfo"
                        class="flex justify-center"
                    >
                        <img
                            src="/images/icons/closed_treasure.png"
                            alt="در حال بارگذاری..."
                            class="h-20 w-20 lg:h-28 lg:w-28 object-contain opacity-60"
                        />
                    </div>

                    <div v-else>
                        <div class="flex justify-center">
                            <img
                                src="/images/icons/closed_treasure.png"
                                alt="گنج بسته"
                                class="h-20 w-20 lg:h-28 lg:w-28 object-contain transition-transform duration-300 cursor-pointer hover:scale-115"
                                @pointerdown="handleTreasureClick"
                            />
                        </div>

                        <div class="mt-2 border-t-2 border-gray-600 pt-2">
                            <div
                                class="flex items-center justify-center gap-x-2 lg:gap-x-4 px-1 lg:px-2"
                            >
                                <div
                                    v-for="req in treasureFetchedInfo.cost
                                        .items"
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
            <button
                v-if="isMobile"
                ref="toggleButton"
                class="absolute top-0 left-1/2 h-10 w-20 -translate-x-1/2 -translate-y-full rounded-t-xl bg-gray-900/75 p-1 shadow-lg"
                @click.stop="isOpen = !isOpen"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke-width="2"
                    stroke="currentColor"
                    class="h-6 w-6 text-white transition-transform duration-300 mx-auto"
                    :class="{ 'rotate-180': isOpen }"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="m4.5 15.75 7.5-7.5 7.5 7.5"
                    />
                </svg>
            </button>
        </div>
    </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref, computed } from 'vue';
import { useModal } from 'vue-final-modal';
import { useToast } from 'vue-toastification';
import { treasureCheck, treasureUnlock } from '@/services/api/index.js';
import TreasureRewardModal from '@/components/features/island/TreasureRewardModal.vue';
import { COST_ITEMS_INFO } from '@/services/cost.js';

const props = defineProps({
    treasureData: { type: Object, required: true },
    modelValue: { type: Object, required: true },
});

const emit = defineEmits(['update:modelValue']);
const toast = useToast();

const isMobile = ref(false);
const isOpen = ref(false);
const treasureContainer = ref(null);
const toggleButton = ref(null);

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

const checkScreenSize = () => {
    const newIsMobile = window.innerWidth < 1024;
    if (newIsMobile !== isMobile.value) {
        isMobile.value = newIsMobile;
        if (!isMobile.value) {
            isOpen.value = true;
        }
    }
};

const handleClickOutside = event => {
    if (
        treasureContainer.value &&
        !treasureContainer.value.contains(event.target) &&
        isOpen.value &&
        isMobile.value
    ) {
        isOpen.value = false;
    }
};

onMounted(() => {
    if (props.treasureData.unlocked) {
        isOpen.value = true;
    }
    treasureCheck(props.treasureData.id)
        .then(data => {
            treasureFetchedInfo.value = data;
            if (!props.treasureData.unlocked && toggleButton.value) {
                toggleButton.value.classList.add('shake-animation');
                setTimeout(() => {
                    toggleButton.value.classList.remove('shake-animation');
                }, 1000);
            }
        })
        .catch(error => {
            console.error('Failed to fetch treasure info:', error);
        });

    checkScreenSize();
    window.addEventListener('resize', checkScreenSize);
    document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
    window.removeEventListener('resize', checkScreenSize);
    document.removeEventListener('click', handleClickOutside);
});

const shouldBeOpen = computed(() => {
    return isOpen.value || !isMobile.value;
});
</script>

<style scoped>
.text-shadow {
    text-shadow: 1px 1px 3px rgba(0, 0, 0, 0.8);
}

@keyframes shake {
    0%,
    100% {
        transform: translateX(0);
    }
    10%,
    30%,
    50%,
    70%,
    90% {
        transform: translateX(-5px);
    }
    20%,
    40%,
    60%,
    80% {
        transform: translateX(5px);
    }
}

.shake-animation {
    animation: shake 1s ease-in-out;
}
</style>
