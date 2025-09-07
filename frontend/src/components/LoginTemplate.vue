<template>
    <div
        class="relative min-h-screen flex items-center justify-center font-vazirmatn"
    >
        <NightSky />
        <div class="animate-fade-in-scale w-full max-w-md mx-4 relative z-10">
            <div
                class="card-hover animate-fade-in-scale w-full max-w-md mx-4 relative z-10 bg-[rgba(20,20,30,0.85)] backdrop-blur-md border border-white/10 shadow-2xl rounded-2xl p-8 text-gray-100"
            >
                <div class="p-6 sm:p-8">
                    <div class="flex justify-center mb-8">
                        <div class="bg-cyan-100 p-4 rounded-full">
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-12 w-12 text-cyan-600"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                                />
                            </svg>
                        </div>
                    </div>

                    <h1
                        class="text-2xl font-bold text-center text-cyan-800 mb-1"
                    >
                        ورود به سیستم
                    </h1>

                    <form class="space-y-6" @submit.prevent="handleLogin">
                        <div>
                            <label
                                for="username"
                                class="block text-sm font-medium text-cyan-700 mb-1"
                                >نام کاربری</label
                            >
                            <div class="relative">
                                <div
                                    class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none"
                                >
                                    <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        class="h-5 w-5 text-cyan-400"
                                        fill="none"
                                        viewBox="0 0 24 24"
                                        stroke="currentColor"
                                    >
                                        <path
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            stroke-width="2"
                                            d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                                        />
                                    </svg>
                                </div>
                                <input
                                    id="username"
                                    v-model="username"
                                    autocomplete="on"
                                    type="text"
                                    class="input-focus text-left [direction:ltr] bg-blue-50 border border-cyan-200 text-cyan-800 text-sm sm:text-base rounded-lg block w-full pl-10 p-2 sm:p-3 sm:pl-10"
                                    placeholder="Username"
                                    required
                                />
                            </div>
                        </div>

                        <div>
                            <label
                                for="password"
                                class="block text-sm font-medium text-cyan-700 mb-1"
                                >رمز عبور</label
                            >
                            <div class="relative">
                                <div
                                    class="absolute inset-y-0 left-0 flex items-center pl-3 cursor-pointer"
                                    @click="showPassword = !showPassword"
                                >
                                    <svg
                                        v-if="!showPassword"
                                        xmlns="http://www.w3.org/2000/svg"
                                        class="h-5 w-5 text-cyan-400"
                                        fill="none"
                                        viewBox="0 0 24 24"
                                        stroke="currentColor"
                                    >
                                        <path
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            stroke-width="2"
                                            d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                                        />
                                        <path
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            stroke-width="2"
                                            d="M2.458 12C3.732 7.943 7.523 5 12 5c4.477 0 8.268 2.943 9.542 7-1.274 
                                4.057-5.065 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                                        />
                                    </svg>
                                    <svg
                                        v-else
                                        xmlns="http://www.w3.org/2000/svg"
                                        class="h-5 w-5 text-cyan-400"
                                        fill="none"
                                        viewBox="0 0 24 24"
                                        stroke="currentColor"
                                    >
                                        <path
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            stroke-width="2"
                                            d="M13.875 18.825A10.05 10.05 0 0112 19c-4.477 0-8.268-2.943-9.542-7
                                a9.956 9.956 0 012.783-4.419M6.228 6.228A9.956 9.956 0 0112 5c4.477 0
                                8.268 2.943 9.542 7a9.96 9.96 0 01-4.132 5.411M15 12a3 3 0 
                                11-6 0 3 3 0 016 0z"
                                        />
                                        <path
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            stroke-width="2"
                                            d="M3 3l18 18"
                                        />
                                    </svg>
                                </div>
                                <input
                                    id="password"
                                    v-model="password"
                                    :type="showPassword ? 'text' : 'password'"
                                    autocomplete="on"
                                    class="input-focus text-left [direction:ltr] bg-blue-50 border border-cyan-200 text-cyan-800 text-sm sm:text-base rounded-lg block w-full pl-10 p-2 sm:p-3 sm:pl-10"
                                    placeholder="********"
                                    required
                                />
                            </div>
                        </div>

                        <div class="flex items-center justify-between">
                            <div class="flex items-start">
                                <div class="flex items-center h-5">
                                    <input
                                        id="remember"
                                        v-model="remember"
                                        type="checkbox"
                                        class="w-4 h-4 bg-blue-50 rounded border border-cyan-300 focus:ring-3 focus:ring-cyan-200"
                                    />
                                </div>
                                <div class="mr-3 text-sm">
                                    <label
                                        for="remember"
                                        class="font-medium text-cyan-700"
                                        >مرا به خاطر بسپار</label
                                    >
                                </div>
                            </div>
                            <div></div>
                        </div>

                        <button
                            type="submit"
                            :disabled="isLoading"
                            class="btn-hover w-full text-white bg-cyan-600 hover:bg-cyan-700 font-medium rounded-lg text-sm px-5 py-2.5 text-center disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            <span v-if="isLoading">در حال ورود...</span>
                            <span v-else>ورود</span>
                        </button>
                    </form>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from 'vue-toastification';
import { login, getPlayer } from '../services/api.js';
import NightSky from './NightSky.vue';

const username = ref('');
const password = ref('');
const remember = ref(false);
const showPassword = ref(false);
const isLoading = ref(false);

const router = useRouter();
const toast = useToast();

async function handleLogin() {
    if (!username.value || !password.value) {
        toast.warning('لطفا نام کاربری و رمز عبور را وارد کنید.');
        return;
    }

    isLoading.value = true;
    try {
        const result = await login(username.value, password.value);

        if (remember.value) {
            localStorage.setItem('authToken', result.token);
        } else {
            sessionStorage.setItem('authToken', result.token);
        }

        toast.success('ورود با موفقیت انجام شد! در حال انتقال...');

        const playerData = await getPlayer();
        if (playerData && playerData.atTerritory) {
            router.push({
                name: 'Territory',
                params: { id: playerData.atTerritory },
            });
        } else {
            throw new Error('اطلاعات تریتوری از سرور دریافت نشد.');
        }
    } catch (err) {
        if (err.message && err.message.includes('invalid credentials')) {
            toast.error('نام کاربری یا رمز عبور اشتباه است.');
        } else {
            toast.error(
                err.message ||
                    'خطایی در ارتباط با سرور رخ داد. لطفا دوباره تلاش کنید.'
            );
        }
    } finally {
        isLoading.value = false;
    }
}
</script>

<style scoped>
@keyframes fadeInScale {
    0% {
        opacity: 0;
        transform: scale(0.95);
    }

    100% {
        opacity: 1;
        transform: scale(1);
    }
}

.animate-fade-in-scale {
    animation: fadeInScale 0.6s ease-out forwards;
}

.card-hover {
    transition:
        all 0.3s ease,
        transform 0.3s ease;
}

.card-hover:hover {
    transform: scale(1.02);
    box-shadow: 0 10px 25px -5px rgba(0, 105, 146, 0.2);
}

.input-focus {
    transition: all 0.3s ease;
}

.input-focus:focus {
    border-color: #38bdf8;
    box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.2);
}

.btn-hover {
    transition:
        all 0.3s ease,
        transform 0.2s ease;
}

.btn-hover:hover:not(:disabled) {
    background-color: #0ea5e9;
    transform: scale(1.03);
}
</style>
