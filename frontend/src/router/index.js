import { createRouter, createWebHistory } from 'vue-router';
import { getToken } from '@/services/api/index.js';

const routes = [
    {
        path: '/',
        redirect: { name: 'Login' },
    },
    {
        path: '/territory/:id',
        name: 'Territory',
        component: () => import('../pages/Territory.vue'),
        props: true,
        meta: { requiresAuth: true },
    },
    {
        path: '/territory/:id/:islandId',
        name: 'Island',
        component: () => import('../pages/TerritoryIsland.vue'),
        props: true,
        meta: { requiresAuth: true },
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('../pages/Login.vue'),
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

router.beforeEach((to, from, next) => {
    const isLoggedIn = !!getToken();

    if (to.meta.requiresAuth && !isLoggedIn) {
        next({ name: 'Login' });
    } else {
        next();
    }
});

export default router;
