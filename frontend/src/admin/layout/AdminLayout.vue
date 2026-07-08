<script setup>
import { RouterView, useRouter } from 'vue-router';
import { adminLogout } from '../services/api.js';

const router = useRouter();

const nav = [
    { name: 'AdminSettings', label: 'General Settings', icon: '⚙️' },
    { name: 'AdminUsers', label: 'Users', icon: '👥' },
    { name: 'AdminMap', label: 'Map Editor', icon: '🗺️' },
    { name: 'AdminIslands', label: 'Island Content', icon: '🏝️' },
];

const isActive = name => router.currentRoute.value.name === name;

const go = name => router.push({ name });

const onLogout = () => adminLogout();
</script>

<template>
    <div class="admin-root" dir="ltr">
        <aside class="admin-sidebar">
            <div class="admin-brand">
                <span class="admin-brand-mark">🌀</span>
                <span class="admin-brand-text">Bermudia Admin</span>
            </div>
            <nav class="admin-nav">
                <button
                    v-for="item in nav"
                    :key="item.name"
                    class="admin-nav-item"
                    :class="{ 'is-active': isActive(item.name) }"
                    @click="go(item.name)"
                >
                    <span class="admin-nav-icon">{{ item.icon }}</span>
                    <span>{{ item.label }}</span>
                </button>
            </nav>
            <button class="admin-logout" @click="onLogout">
                <span class="admin-nav-icon">⎋</span>
                <span>Log out</span>
            </button>
        </aside>
        <main class="admin-main">
            <RouterView />
        </main>
    </div>
</template>

<style scoped>
.admin-root {
    display: flex;
    min-height: 100vh;
    background: #0f172a;
    color: #e2e8f0;
    font-family:
        ui-sans-serif,
        system-ui,
        -apple-system,
        'Segoe UI',
        Roboto,
        sans-serif;
}

.admin-sidebar {
    width: 240px;
    flex-shrink: 0;
    background: #111827;
    border-right: 1px solid #1f2937;
    padding: 20px 14px;
    display: flex;
    flex-direction: column;
    gap: 8px;
    position: sticky;
    top: 0;
    height: 100vh;
}

.admin-brand {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 10px 18px;
    font-weight: 700;
    font-size: 16px;
}
.admin-brand-mark {
    font-size: 22px;
}

.admin-nav {
    display: flex;
    flex-direction: column;
    gap: 4px;
    flex: 1;
}

.admin-nav-item,
.admin-logout {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 12px;
    border-radius: 10px;
    background: transparent;
    border: none;
    color: #cbd5e1;
    font-size: 14px;
    text-align: left;
    cursor: pointer;
    transition:
        background 0.15s,
        color 0.15s;
    width: 100%;
}
.admin-nav-item:hover,
.admin-logout:hover {
    background: #1f2937;
    color: #f8fafc;
}
.admin-nav-item.is-active {
    background: #2563eb;
    color: #fff;
}
.admin-nav-icon {
    width: 22px;
    text-align: center;
    font-size: 16px;
}

.admin-logout {
    color: #f87171;
    margin-top: auto;
}
.admin-logout:hover {
    background: #7f1d1d33;
    color: #fca5a5;
}

.admin-main {
    flex: 1;
    padding: 32px 40px;
    overflow-x: hidden;
    max-width: 100%;
}

@media (max-width: 720px) {
    .admin-root {
        flex-direction: column;
    }
    .admin-sidebar {
        width: 100%;
        height: auto;
        position: static;
        flex-direction: row;
        flex-wrap: wrap;
        align-items: center;
    }
    .admin-nav {
        flex-direction: row;
        flex-wrap: wrap;
    }
    .admin-logout {
        margin-top: 0;
    }
    .admin-main {
        padding: 20px;
    }
}
</style>
