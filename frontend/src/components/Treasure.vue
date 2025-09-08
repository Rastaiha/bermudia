<template>
    <div
        class="relative bg-slate-900/90 border border-slate-700 rounded-2xl shadow-xl p-8 transition-all duration-500"
    >
        <div class="text-center">
            <p
                class="text-xl md:text-2xl font-light text-gray-200 leading-relaxed"
            >
                گنج
            </p>
        </div>

        <div
            class="mt-8 flex items-center justify-center"
            style="min-height: 100px"
        >
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
                        :button-enabled="true"
                        :button-text="
                            treasureFetchedInfo.feasible ? 'بگشا' : null
                        "
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
    </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { treasureCheck, treasureUnlock } from '@/services/api.js';
import InfoBox from './InfoBox.vue';

const props = defineProps({
    treasureData: {
        type: Object,
        required: true,
    },
    modelValue: { type: Object, required: true },
});

const emit = defineEmits(['update:modelValue', 'treasureOpened']);

const treasureFetchedInfo = ref(null);

const unlockTreasure = async () => {
    try {
        const resp = await treasureUnlock(props.treasureData.id);

        emit('update:modelValue', resp);

        if (resp.rewards && resp.rewards.length > 0) {
            emit('treasureOpened', resp.rewards);
        }
    } catch (error) {
        console.error('Failed to unlock treasure:', error);
    }
};

const infoBoxStyle = computed(() => {
    return {
        margin: 'auto',
    };
});

onMounted(async () => {
    if (props.treasureData.unlocked) return;
    treasureFetchedInfo.value = await treasureCheck(props.treasureData.id);
});
</script>
