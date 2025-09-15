<template>
    <div
        class="relative bg-slate-900/90 border border-slate-700 rounded-2xl shadow-xl p-8 transition-all duration-500"
    >
        <div class="text-center">
            <p
                class="text-xl md:text-2xl font-light text-gray-200 leading-relaxed"
            >
                {{ challenge.description }}
            </p>
        </div>

        <div
            class="mt-8 flex items-center justify-center"
            style="min-height: 100px"
        >
            <transition name="popup-fade" mode="out-in">
                <div
                    v-if="
                        challenge.submissionState.status === 'empty' ||
                        challenge.submissionState.status === 'wrong'
                    "
                    key="form"
                    class="w-full max-w-lg"
                >
                    <p
                        v-if="challenge.submissionState.status === 'wrong'"
                        class="text-center text-red-400 mb-4"
                    >
                        پاسخ اشتباه بود، دوباره تلاش کنید.
                    </p>

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
                        class="mt-2 text-sm text-gray-300 max-w-md"
                    >
                        {{ challenge.submissionState.feedback }}
                    </p>
                </div>
            </transition>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue';

const props = defineProps({
    challenge: {
        type: Object,
        required: true,
    },
});

const emit = defineEmits(['submit']);

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
</style>
