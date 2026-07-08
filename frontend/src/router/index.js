import { createRouter, createWebHistory } from 'vue-router';
import { getToken } from '@/services/api/index.js';
import { adminRoutes, adminGuard } from '@/admin/router.js';

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
    ...adminRoutes,
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

router.beforeEach((to, from, next) => {
    // Admin routes have their own auth flow, independent of the player token.
    if (to.path.startsWith('/admin')) {
        if (adminGuard(to, from, next)) {
            next();
        }
        return;
    }

    const isLoggedIn = !!getToken();

    if (to.meta.requiresAuth && !isLoggedIn) {
        next({ name: 'Login' });
    } else {
        next();
    }
});

export default router;
