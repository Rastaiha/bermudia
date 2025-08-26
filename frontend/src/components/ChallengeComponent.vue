<template>
    <div class="relative bg-slate-900/80 border border-sky-500/40 rounded-3xl shadow-lg overflow-hidden transition-all duration-300 ease-linear p-10 flex flex-col gap-6 items-center text-center"
        :class="`challenge-${challenge.submissionState.status}`">
        <p class="text-2xl font-medium leading-relaxed">{{ challenge.description }}</p>
        <div class="flex w-full max-w-lg mt-4 flex-row-reverse">
            <input v-if="challenge.type !== 'file'" :type="challenge.type" v-model="inputValue"
                class="flex-grow border border-sky-500/40 bg-black/20 text-gray-200 px-4 py-3 rounded-l-xl text-base outline-none transition-shadow duration-200 focus:ring-2 focus:ring-sky-500/60"
                placeholder="پاسخ خود را اینجا وارد کنید..." :disabled="!challenge.submissionState.submittable" />
            <input v-else :type="challenge.type" :accept="challenge.accept?.join(',')" @change="handleFileChange"
                class="flex-grow border border-sky-500/40 bg-black/20 text-gray-200 px-4 py-3 rounded-l-xl text-base outline-none transition-shadow duration-200 focus:ring-2 focus:ring-sky-500/60"
                :disabled="!challenge.submissionState.submittable" />
            <button v-if="challenge.submissionState.submittable" @click="submit"
                class="border-none bg-blue-500 text-white px-6 py-3 text-base font-semibold rounded-r-xl cursor-pointer transition-colors duration-200 hover:bg-blue-600">
                ارسال
            </button>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue';

const props = defineProps({
    challenge: {
        type: Object,
        required: true,
    },
});

const emit = defineEmits(['submit']);

const inputValue = ref(props.challenge.type === 'file' ? null : '');

const handleFileChange = (event) => {
    const file = event.target.files[0];
    if (file) {
        inputValue.value = file;
    }
};

const submit = () => {
    emit('submit', {
        inputId: props.challenge.id,
        value: inputValue.value,
    });
};
</script>

<style scoped>
.challenge-correct {
    border-color: #22c55e;
}

.challenge-incorrect {
    border-color: #ef4444;
}
</style>