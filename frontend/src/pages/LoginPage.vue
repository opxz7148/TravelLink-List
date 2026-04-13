<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-emerald-100 via-sky-100 to-cyan-100 px-4 py-8">
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-10 max-w-5xl w-full items-center">
      <!-- Login Card -->
      <div class="rounded-3xl border border-white/80 bg-white/90 shadow-xl p-8 md:p-10 backdrop-blur">
        <h1 class="tl-title text-4xl font-bold text-slate-900 mb-2">Welcome Back</h1>
        <p class="text-slate-600 mb-6">Sign in to continue your travel planning.</p>

        <form @submit.prevent="handleLogin" class="space-y-4 mb-6">
          <div class="flex flex-col gap-1.5">
            <label for="email" class="text-sm font-medium text-slate-700">Email</label>
            <input
              id="email"
              v-model="email"
              type="email"
              placeholder="you@example.com"
              required
              class="px-3 py-2.5 border border-slate-300 rounded-xl text-sm focus:outline-none focus:border-emerald-500 focus:ring-4 focus:ring-emerald-100 transition-all"
            />
          </div>

          <div class="flex flex-col gap-1.5">
            <label for="password" class="text-sm font-medium text-slate-700">Password</label>
            <input
              id="password"
              v-model="password"
              type="password"
              placeholder="••••••••"
              required
              class="px-3 py-2.5 border border-slate-300 rounded-xl text-sm focus:outline-none focus:border-emerald-500 focus:ring-4 focus:ring-emerald-100 transition-all"
            />
          </div>

          <div v-if="authStore.error" class="p-3 bg-red-50 text-red-800 rounded-md text-sm">
            {{ authStore.error }}
          </div>

          <button type="submit" class="w-full px-4 py-2.5 bg-gradient-to-r from-emerald-600 to-sky-600 text-white rounded-xl text-sm font-semibold cursor-pointer transition-all hover:from-emerald-700 hover:to-sky-700 hover:shadow-lg disabled:opacity-70 disabled:cursor-not-allowed" :disabled="authStore.isLoading">
            <span v-if="authStore.isLoading">Signing in...</span>
            <span v-else>Sign In</span>
          </button>
        </form>

        <div class="relative my-6 text-center">
          <div class="absolute inset-x-0 top-1/2 h-px bg-slate-300"></div>
          <span class="relative bg-white px-2 text-slate-400 text-sm">or</span>
        </div>

        <router-link to="/register" class="block w-full px-4 py-2.5 bg-slate-100 text-slate-900 rounded-xl text-sm font-medium text-center cursor-pointer transition-all hover:bg-slate-200">
          Create new account
        </router-link>

        <div class="text-center mt-4">
          <router-link to="/browse" class="text-emerald-700 text-sm font-semibold no-underline hover:underline">Continue as guest</router-link>
        </div>
      </div>

      <!-- Illustration -->
      <div class="hidden lg:flex flex-col items-center justify-center">
        <svg viewBox="0 0 200 200" fill="none" stroke="currentColor" class="w-40 h-40 mb-6 opacity-90 text-emerald-800">
          <circle cx="100" cy="100" r="80" stroke-width="2" />
          <path d="M 100 40 L 140 100 L 120 100 L 120 160 L 80 160 L 80 100 L 60 100 Z" />
        </svg>
        <p class="text-slate-700 text-lg text-center leading-relaxed max-w-md">Discover routes, build journeys, and turn ideas into memorable adventures.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth_store';
import { useUiStore } from '../stores/ui_store';

const router = useRouter();
const authStore = useAuthStore();
const uiStore = useUiStore();

const email = ref('');
const password = ref('');

async function handleLogin(): Promise<void> {
  if (!email.value || !password.value) {
    uiStore.showError('Please enter both email and password');
    return;
  }

  try {
    await authStore.login(email.value, password.value);
    uiStore.showSuccess('Logged in successfully!');
    router.push('/browse');
  } catch (error: any) {
    // Error is already set in authStore
    uiStore.showError(authStore.error || 'Login failed');
  }
}
</script>
