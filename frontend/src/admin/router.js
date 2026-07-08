import { getAdminToken } from './services/api.js';

// Admin routes are lazy-loaded so the player app bundle is unaffected.
export const adminRoutes = [
    {
        path: '/admin/login',
        name: 'AdminLogin',
        component: () => import('./pages/AdminLogin.vue'),
    },
    {
        path: '/admin',
        component: () => import('./layout/AdminLayout.vue'),
        meta: { requiresAdmin: true },
        children: [
            {
                path: '',
                name: 'AdminSettings',
                component: () => import('./pages/AdminSettings.vue'),
            },
            {
                path: 'users',
                name: 'AdminUsers',
                component: () => import('./pages/AdminUsers.vue'),
            },
            {
                path: 'map',
                name: 'AdminMap',
                component: () => import('./pages/AdminMap.vue'),
            },
            {
                path: 'islands',
                name: 'AdminIslands',
                component: () => import('./pages/AdminIslands.vue'),
            },
        ],
    },
];

// Navigation guard for admin routes. Registered from the main router.
export const adminGuard = (to, from, next) => {
    if (to.meta.requiresAdmin && !getAdminToken()) {
        next({ name: 'AdminLogin' });
        return false;
    }
    // If already logged in, skip the login page.
    if (to.name === 'AdminLogin' && getAdminToken()) {
        next({ name: 'AdminSettings' });
        return false;
    }
    return true;
};
