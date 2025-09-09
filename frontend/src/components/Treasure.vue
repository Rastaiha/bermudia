<template>
    <div class="w-full flex justify-center items-center">
        <transition name="popup-fade" mode="out-in">
            <div v-if="treasureData.unlocked">
                <InfoBox
                    title="صندوق گنج"
                    :info-box-style="infoBoxStyle"
                    :loading="false"
                    error-text="این صندوق قبلا باز شده‌است."
                >
                    <img src="/images/island/opened_treasure.png" />
                </InfoBox>
            </div>

            <div
                v-else-if="treasureFetchedInfo"
                key="form"
                class="w-full max-w-lg"
            >
                <InfoBox
                    title="صندوق گنج"
                    :info-box-style="infoBoxStyle"
                    :loading="false"
                    :button-enabled="treasureFetchedInfo.feasible"
                    button-text="بگشا"
                    :error-text="
                        treasureFetchedInfo.feasible
                            ? null
                            : treasureFetchedInfo.reason
                    "
                    :cost="treasureFetchedInfo.cost"
                    @action="unlockTreasure"
                >
                    <img src="/images/island/closed_treasure.png" />
                </InfoBox>
            </div>
        </transition>
    </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useModal } from 'vue-final-modal';
import { treasureCheck, treasureUnlock } from '@/services/api.js';
import emitter from '@/services/eventBus.js';
import InfoBox from './InfoBox.vue';
import TreasureRewardModal from './TreasureRewardModal.vue';

const props = defineProps({
    treasureData: { type: Object, required: true },
    modelValue: { type: Object, required: true },
});

const emit = defineEmits(['update:modelValue']);

const treasureFetchedInfo = ref(null);
const receivedRewards = ref([]);

const itemMap = {
    coins: { name: 'سکه', icon: '/images/icons/coin.png' },
    blueKeys: { name: 'کلید آبی', icon: '/images/icons/blueKeys.png' },
    redKeys: { name: 'کلید قرمز', icon: '/images/icons/redKeys.png' },
    goldenKeys: { name: 'کلید طلایی', icon: '/images/icons/goldenKeys.png' },
    books: { name: 'کتاب', icon: '/images/icons/book.png' },
};

const { open, close } = useModal({
    component: TreasureRewardModal,
    attrs: {
        rewards: receivedRewards,
        onClose: () => close(),
    },
});

function calculateRewards(oldPlayer, newPlayer) {
    const rewards = [];
    if (!oldPlayer || !newPlayer) return [];

    for (const key in itemMap) {
        const oldValue = Array.isArray(oldPlayer[key])
            ? oldPlayer[key].length
            : oldPlayer[key] || 0;
        const newValue = Array.isArray(newPlayer[key])
            ? newPlayer[key].length
            : newPlayer[key] || 0;
        if (newValue > oldValue) {
            rewards.push({ ...itemMap[key], quantity: newValue - oldValue });
        }
    }
    return rewards;
}

const unlockTreasure = async () => {
    try {
        await treasureUnlock(props.treasureData.id);
        emit('update:modelValue', { ...props.treasureData, unlocked: true });
    } catch (error) {
        console.error('Failed to unlock treasure:', error);
    }
};

const onTreasureUnlocked = ({ oldPlayerState, newPlayerState }) => {
    const rewards = calculateRewards(oldPlayerState, newPlayerState);
    if (rewards.length > 0) {
        receivedRewards.value = rewards;
        open();
    }
};

onMounted(() => {
    if (props.treasureData.unlocked) return;
    treasureCheck(props.treasureData.id).then(
        data => (treasureFetchedInfo.value = data)
    );
    emitter.on('treasure-unlocked', onTreasureUnlocked);
});

onUnmounted(() => {
    emitter.off('treasure-unlocked', onTreasureUnlocked);
});

const infoBoxStyle = computed(() => {
    return {
        margin: 'auto',
    };
});
</script>
