<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from 'vue-toastification';
import { adminLogin } from '../services/api.js';

const router = useRouter();
const toast = useToast();

const username = ref('');
const password = ref('');
const loading = ref(false);

const submit = async () => {
    if (loading.value) return;
    loading.value = true;
    try {
        await adminLogin(username.value.trim(), password.value);
        toast.success('Welcome back!');
        router.push({ name: 'AdminSettings' });
    } catch (e) {
        toast.error(e.message || 'Login failed');
    } finally {
        loading.value = false;
    }
};
</script>

<template>
    <div class="login-wrap" dir="ltr">
        <form class="login-card" @submit.prevent="submit">
            <div class="login-brand">
                <span class="login-mark">🌀</span>
                <h1>Bermudia Admin</h1>
                <p>Sign in to manage the game</p>
            </div>

            <label class="login-field">
                <span>Username</span>
                <input
                    v-model="username"
                    type="text"
                    autocomplete="username"
                    required
                />
            </label>

            <label class="login-field">
                <span>Password</span>
                <input
                    v-model="password"
                    type="password"
                    autocomplete="current-password"
                    required
                />
            </label>

            <button class="login-btn" type="submit" :disabled="loading">
                {{ loading ? 'Signing in…' : 'Sign in' }}
            </button>
        </form>
    </div>
</template>

<style scoped>
.login-wrap {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background:
        radial-gradient(1200px 600px at 50% -10%, #1e3a8a55, transparent),
        #0f172a;
    padding: 20px;
    font-family:
        ui-sans-serif,
        system-ui,
        -apple-system,
        'Segoe UI',
        Roboto,
        sans-serif;
}

.login-card {
    width: 100%;
    max-width: 380px;
    background: #111827;
    border: 1px solid #1f2937;
    border-radius: 18px;
    padding: 32px 28px;
    display: flex;
    flex-direction: column;
    gap: 18px;
    box-shadow: 0 20px 60px #00000055;
}

.login-brand {
    text-align: center;
    color: #e2e8f0;
    margin-bottom: 6px;
}
.login-mark {
    font-size: 34px;
}
.login-brand h1 {
    font-size: 20px;
    font-weight: 700;
    margin: 8px 0 4px;
}
.login-brand p {
    font-size: 13px;
    color: #94a3b8;
}

.login-field {
    display: flex;
    flex-direction: column;
    gap: 6px;
    color: #cbd5e1;
    font-size: 13px;
}
.login-field input {
    background: #0b1220;
    border: 1px solid #334155;
    border-radius: 10px;
    padding: 11px 12px;
    color: #f8fafc;
    font-size: 14px;
    outline: none;
    transition: border-color 0.15s;
}
.login-field input:focus {
    border-color: #2563eb;
}

.login-btn {
    margin-top: 6px;
    background: #2563eb;
    color: #fff;
    border: none;
    border-radius: 10px;
    padding: 12px;
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    transition:
        background 0.15s,
        opacity 0.15s;
}
.login-btn:hover:not(:disabled) {
    background: #1d4ed8;
}
.login-btn:disabled {
    opacity: 0.6;
    cursor: default;
}
</style>
