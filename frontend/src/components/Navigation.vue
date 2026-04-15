<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth_store';
import { useUiStore } from '../stores/ui_store';

const router = useRouter();
const authStore = useAuthStore();
const uiStore = useUiStore();

const isAuthenticated = computed(() => authStore.isAuthenticated);
const username = computed(() => authStore.user?.username || 'User');
const isTraveller = computed(() => authStore.isTraveller);
const isAdmin = computed(() => authStore.isAdmin);

const handleLogout = () => {
  authStore.logout();
  uiStore.showSuccess('Logged out successfully');
  router.push('/browse');
};

const navigateTo = (path: string) => {
  router.push(path);
};
</script>

<template>
  <nav class="sticky top-0 z-50 border-b border-slate-200/80 bg-white/85 shadow-sm backdrop-blur-lg">
    <div class="max-w-7xl mx-auto px-4 lg:px-8 py-3">
      <!-- Logo -->
      <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div class="flex-shrink-0">
          <router-link to="/browse" class="tl-title no-underline text-slate-900 font-bold text-xl hover:text-emerald-700 transition-colors flex items-center gap-2">
            <span class="inline-flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-emerald-500 to-sky-500 text-white shadow">✈</span>
            <span>TravelLink</span>
          </router-link>
        </div>

        <!-- Navigation Links -->
        <div class="flex flex-wrap items-center gap-2 md:gap-3 md:justify-end">
        <!-- Public Links -->
          <router-link to="/browse" class="no-underline text-slate-700 font-medium transition-colors px-3 py-2 rounded-lg hover:text-emerald-700 hover:bg-emerald-50 text-sm md:text-base">Browse Plans</router-link>

        <!-- Authenticated Links -->
          <template v-if="isAuthenticated">
          <!-- Create Plan (for all authenticated users except admin) -->
            <router-link v-if="!isAdmin" to="/create-plan" class="no-underline text-slate-700 font-medium transition-colors px-3 py-2 rounded-lg hover:text-emerald-700 hover:bg-emerald-50 text-sm md:text-base">Create Plan</router-link>

          <!-- Admin Links -->
            <template v-if="isAdmin">
              <router-link to="/admin" class="no-underline text-amber-700 font-semibold transition-colors px-3 py-2 rounded-lg hover:text-amber-800 hover:bg-amber-50 text-sm md:text-base">Admin Dashboard</router-link>
            </template>

          <!-- User Menu -->
            <div class="flex items-center gap-2 border-l border-slate-200 pl-3 md:pl-4">
              <span class="hidden sm:inline rounded-lg bg-slate-100 px-2.5 py-1 text-xs font-semibold text-slate-700">{{ username }}</span>
              <button v-if="!isAdmin" @click="navigateTo('/my-plans')" class="text-slate-700 font-medium transition-colors px-3 py-2 rounded-lg hover:text-blue-700 hover:bg-blue-50 text-sm md:text-base border-none bg-transparent cursor-pointer">My Plans</button>
              <button @click="navigateTo('/profile')" class="text-slate-700 font-medium transition-colors px-3 py-2 rounded-lg hover:text-emerald-700 hover:bg-emerald-50 text-sm md:text-base border-none bg-transparent cursor-pointer">Profile</button>
              <button @click="handleLogout" class="bg-rose-600 text-white px-3.5 py-2 rounded-lg hover:bg-rose-700 transition-colors text-sm font-semibold border-none cursor-pointer">Logout</button>
            </div>
          </template>

        <!-- Unauthenticated Links -->
          <template v-else>
            <router-link to="/login" class="no-underline text-slate-700 font-medium transition-colors px-3 py-2 rounded-lg hover:text-emerald-700 hover:bg-emerald-50 text-sm md:text-base">Login</router-link>
            <router-link to="/register" class="bg-gradient-to-r from-emerald-600 to-sky-600 text-white hover:from-emerald-700 hover:to-sky-700 rounded-lg px-4 py-2 transition-colors font-semibold no-underline text-sm md:text-base shadow">Register</router-link>
          </template>
        </div>
      </div>
    </div>
  </nav>
</template>
