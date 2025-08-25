import { createRouter, createWebHistory } from 'vue-router';
import { getToken } from '../services/api'; // For checking authentication

const routes = [
  {
    path: '/',
    redirect: { name: 'Login' } // Redirect root to login
  },
  {
    path: '/territory/:id',
    name: 'Territory',
    component: () => import('../pages/Territory.vue'),
    props: true,
    meta: { requiresAuth: true }, // This route requires login
  },
  {
    path: '/territory/:id/:islandId',
    name: 'Island',
    component: () => import('../pages/TerritoryIsland.vue'),
    props: true,
    meta: { requiresAuth: true }, // This route requires login
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../pages/Login.vue'),
  },
  {
    path: '/user_page',
    name: 'UserPage',
    component: () => import('../pages/UserPage.vue'),
    meta: { requiresAuth: true }, // This route requires login
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Navigation Guard to protect routes
router.beforeEach((to, from, next) => {
  const isLoggedIn = !!getToken(); // Check if user token exists

  if (to.meta.requiresAuth && !isLoggedIn) {
    // If the route requires auth and user is not logged in, redirect to login page
    next({ name: 'Login' });
  } else {
    // Otherwise, allow navigation
    next();
  }
});

export default router;
