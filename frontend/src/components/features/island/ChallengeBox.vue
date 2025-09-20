<template>
    <div
        class="relative bg-slate-900/90 border border-slate-700 rounded-2xl shadow-xl p-8 transition-all duration-500"
    >
        <div class="absolute top-4 right-4">
            <button
                :class="[helpButtonClass, helpButtonAnimationClass]"
                :disabled="isHelpButtonDisabled"
                class="p-2 rounded-full"
                @click="handleHelpClick"
            >
                <UserIcon
                    v-if="helpButtonIcon === 'QuestionMarkCircleIcon'"
                    class="w-12 h-12"
                />
                <ArrowTopRightOnSquareIcon v-else class="w-12 h-12" />
            </button>
        </div>
        <div class="text-center pt-12">
            <p
                class="text-xl md:text-2xl font-light text-gray-200 leading-relaxed"
            >
                {{ challenge.description }}
            </p>
        </div>

        <div
            class="mt-8 flex items-center justify-center flex-col"
            style="min-height: 100px"
        >
            <transition name="popup-fade" mode="out-in">
                <div
                    v-if="challenge.submissionState.status === 'empty'"
                    key="form"
                ></div>

                <div
                    v-else-if="challenge.submissionState.status === 'wrong'"
                    key="wrong"
                    class="flex flex-col items-center justify-center text-center"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-12 w-12 text-red-400"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                    </svg>
                    <p class="mt-4 text-xl font-semibold text-red-400">
                        پاسخ شما اشتباه است
                    </p>
                    <p
                        v-if="challenge.submissionState.feedback"
                        class="mt-2 text-sm text-gray-300 max-w-md text-justify"
                    >
                        {{ challenge.submissionState.feedback }}
                    </p>
                </div>

                <div
                    v-else-if="challenge.submissionState.status === 'pending'"
                    key="pending"
                    class="flex flex-col items-center justify-center text-center"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-12 w-12 text-sky-400"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                    </svg>
                    <p class="mt-4 text-xl font-semibold text-gray-200">
                        پاسخ شما ثبت شد
                    </p>
                    <p class="mt-2 text-sm text-gray-400">
                        نتیجه تا چند دقیقه دیگر اعلام می‌شود. می‌توانید به بازی
                        ادامه دهید.
                    </p>
                </div>

                <div
                    v-else-if="challenge.submissionState.status === 'correct'"
                    key="correct"
                    class="flex flex-col items-center justify-center text-green-400"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-12 w-12"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                    </svg>
                    <p class="mt-4 text-xl font-semibold">پاسخ شما صحیح است!</p>
                </div>

                <div
                    v-else-if="
                        challenge.submissionState.status === 'half-correct'
                    "
                    key="half-correct"
                    class="flex flex-col items-center justify-center text-center"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-12 w-12 text-amber-400"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                        stroke-width="2"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                    </svg>
                    <p class="mt-4 text-xl font-semibold text-amber-400">
                        پاسخ شما نیمه‌درست است
                    </p>
                    <p
                        v-if="challenge.submissionState.feedback"
                        class="mt-2 text-sm text-gray-300 max-w-md text-justify"
                    >
                        {{ challenge.submissionState.feedback }}
                    </p>
                </div>
            </transition>
            <div
                v-if="challenge.submissionState.submittable"
                class="w-full flex justify-center"
            >
                <div
                    v-if="challenge.type !== 'file'"
                    class="flex items-center gap-3"
                >
                    <input
                        v-model="inputValue"
                        :type="challenge.type"
                        placeholder="پاسخ..."
                        class="w-full p-3 text-lg text-center text-gray-100 bg-slate-800/70 rounded-lg border-2 border-slate-600 focus:border-cyan-500 focus:ring-0 outline-none transition-colors"
                        :disabled="!challenge.submissionState.submittable"
                        @keyup.enter="submit"
                    />
                    <button
                        :disabled="!challenge.submissionState.submittable"
                        class="btn-hover px-6 py-3 text-lg font-semibold text-white bg-[#07458bb5] rounded-lg shrink-0"
                        @click="submit"
                    >
                        ارسال
                    </button>
                </div>

                <div v-else class="flex items-center justify-center gap-3">
                    <label
                        :for="fileInputId"
                        class="btn-hover flex-grow text-center px-5 py-3 text-lg font-medium text-gray-200 bg-slate-700/80 rounded-lg border-2 border-transparent hover:border-cyan-500 cursor-pointer"
                    >
                        <span v-if="!selectedFileName">انتخاب فایل</span>
                        <span v-else class="text-cyan-400">{{
                            selectedFileName
                        }}</span>
                    </label>
                    <input
                        :id="fileInputId"
                        type="file"
                        class="hidden"
                        :accept="challenge.accept?.join(',')"
                        :disabled="!challenge.submissionState.submittable"
                        @change="handleFileChange"
                    />
                    <button
                        :disabled="
                            !challenge.submissionState.submittable ||
                            !inputValue
                        "
                        class="btn-hover px-6 py-3 text-lg font-semibold text-white bg-green-600 rounded-lg disabled:bg-gray-600 disabled:opacity-50 disabled:transform-none disabled:filter-none shrink-0"
                        @click="submit"
                    >
                        ارسال
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useModal } from 'vue-final-modal';
import ConfirmModal from '@/components/common/ConfirmModal.vue';
import { ArrowTopRightOnSquareIcon, UserIcon } from '@heroicons/vue/24/outline';

const props = defineProps({
    challenge: {
        type: Object,
        required: true,
    },
});

const emit = defineEmits(['submit', 'help-requested']);

const { open: openHelpModal, close: closeHelpModal } = useModal({
    component: ConfirmModal,
    attrs: {
        title: 'درخواست کمک',
        onConfirm() {
            emit('help-requested', props.challenge);
            closeHelpModal();
        },
        onCancel() {
            closeHelpModal();
        },
    },
    slots: {
        content: '<p>با درخواست کمک دانش دریافتی از حل این سوال نصف میشه</p>',
    },
});

const helpButtonIcon = computed(() => {
    return props.challenge.submissionState.hasRequestedHelp
        ? 'ArrowTopRightOnSquareIcon'
        : 'QuestionMarkCircleIcon';
});

const helpButtonClass = computed(() => {
    if (isHelpButtonDisabled.value) {
        return 'text-gray-500';
    }
    if (props.challenge.submissionState.hasRequestedHelp) {
        return 'text-green-400';
    }
    return 'text-blue-400';
});

const helpButtonAnimationClass = computed(() => {
    if (isHelpButtonDisabled.value) {
        return '';
    }
    if (props.challenge.submissionState.hasRequestedHelp) {
        return 'animate-pulse-green';
    }
    return 'animate-pulse-blue';
});

const isHelpButtonDisabled = computed(() => {
    return true;
    // return !props.challenge.submissionState.showHelp;
});

const handleHelpClick = async () => {
    if (props.challenge.submissionState.hasRequestedHelp) {
        emit('help-requested', props.challenge);
    } else {
        openHelpModal();
    }
};

const inputValue = ref(
    props.challenge.type === 'file'
        ? null
        : props.challenge.type === 'number'
          ? null
          : ''
);
const fileInputId = computed(() => `file-upload-${props.challenge.id}`);

const selectedFileName = computed(() => {
    if (inputValue.value instanceof File) {
        return inputValue.value.name;
    }
    return '';
});

const handleFileChange = event => {
    const file = event.target.files[0];
    if (file) {
        inputValue.value = file;
    }
};

const submit = () => {
    if (props.challenge.submissionState.submittable) {
        emit('submit', {
            inputId: props.challenge.id,
            data: inputValue.value,
        });
    }
};
</script>

<style scoped>
.popup-fade-enter-active,
.popup-fade-leave-active {
    transition:
        opacity 0.3s ease,
        transform 0.3s ease;
}

.popup-fade-enter-from,
.popup-fade-leave-to {
    opacity: 0;
    transform: translateY(10px);
}

.btn-hover {
    transition:
        transform 0.2s ease,
        filter 0.2s ease;
}

.btn-hover:hover:not(:disabled) {
    transform: scale(1.05);
    filter: brightness(1.1);
}

@keyframes pulse-blue {
    0% {
        box-shadow: 0 0 0 0 rgba(96, 165, 250, 0.7);
    }
    100% {
        box-shadow: 0 0 0 1rem rgba(96, 165, 250, 0);
    }
}

@keyframes pulse-green {
    0% {
        box-shadow: 0 0 0 0 rgba(74, 222, 128, 0.7);
    }
    100% {
        box-shadow: 0 0 0 1rem rgba(74, 222, 128, 0);
    }
}

.animate-pulse-blue {
    animation: pulse-blue 2s infinite;
}

.animate-pulse-green {
    animation: pulse-green 2s infinite;
}
</style>
