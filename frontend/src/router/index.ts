import { createRouter, createWebHistory } from 'vue-router';
import type { RouteRecordRaw } from 'vue-router';
import { useAuthStore } from '../stores/auth_store';

// Pages
import BrowsePage from '../pages/BrowsePage.vue';
import ViewPlanPage from '../pages/ViewPlanPage.vue';
import LoginPage from '../pages/LoginPage.vue';
import RegisterPage from '../pages/RegisterPage.vue';
import CreatePlanPage from '../pages/CreatePlanPage.vue';
import ProfilePage from '../pages/ProfilePage.vue';
import AdminDashboard from '../pages/AdminDashboard.vue';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/browse',
  },

  // Public routes
  {
    path: '/browse',
    name: 'Browse',
    component: BrowsePage,
  },

  {
    path: '/plans/:id',
    name: 'ViewPlan',
    component: ViewPlanPage,
  },

  {
    path: '/login',
    name: 'Login',
    component: LoginPage,
    meta: {
      requiresGuest: true,
    },
  },

  {
    path: '/register',
    name: 'Register',
    component: RegisterPage,
    meta: {
      requiresGuest: true,
    },
  },

  // Protected routes (authenticated users)
  {
    path: '/profile',
    name: 'Profile',
    component: ProfilePage,
    meta: {
      requiresAuth: true,
    },
  },

  // Protected routes (traveller only)
  {
    path: '/create-plan',
    name: 'CreatePlan',
    component: CreatePlanPage,
    meta: {
      requiresAuth: true,
      requiresRole: 'traveller',
    },
  },

  // Admin routes
  {
    path: '/admin',
    name: 'AdminDashboard',
    component: AdminDashboard,
    meta: {
      requiresAuth: true,
      requiresRole: 'admin',
    },
  },

  // Fallback 404
  {
    path: '/:pathMatch(.*)*',
    redirect: '/browse',
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

/**
 * Global navigation guards for authentication and authorization
 */
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  const meta = to.meta as any;

  // Restore session if not already done
  if (!authStore.token && !from.name) {
    authStore.restoreSession();
  }

  // If route requires authentication
  if (meta?.requiresAuth && !authStore.isAuthenticated) {
    next('/login');
    return;
  }

  // If route requires guest (not authenticated)
  if (meta?.requiresGuest && authStore.isAuthenticated) {
    next('/browse');
    return;
  }

  // If route requires specific role
  if (meta?.requiresRole) {
    if (meta.requiresRole === 'traveller' && !authStore.isTraveller) {
      next('/browse');
      return;
    }
    if (meta.requiresRole === 'admin' && !authStore.isAdmin) {
      next('/browse');
      return;
    }
  }

  next();
});

export default router;
