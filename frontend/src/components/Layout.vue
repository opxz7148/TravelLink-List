<template>
  <div class="relative flex min-h-screen flex-col overflow-hidden">
    <div class="pointer-events-none absolute inset-0 -z-10">
      <div class="absolute -left-20 top-24 h-64 w-64 rounded-full bg-emerald-200/40 blur-3xl"></div>
      <div class="absolute right-0 top-0 h-72 w-72 rounded-full bg-sky-200/40 blur-3xl"></div>
    </div>
    <!-- Navigation Header -->
    <Navigation />

    <!-- Notification Center -->
    <NotificationCenter />

    <!-- Loading Overlay -->
    <div v-if="isLoading" class="fixed inset-0 bg-black/30 flex flex-col justify-center items-center z-999">
      <div class="w-8 h-8 border-4 border-gray-100 border-t-blue-600 rounded-full animate-spin mb-4"></div>
      <p class="text-white font-medium">Loading...</p>
    </div>

    <!-- Main Content -->
    <main class="flex-1 w-full max-w-[1700px] mx-auto px-4 sm:px-6 lg:px-10 py-6 md:py-10">
      <router-view />
    </main>

    <!-- Footer -->
    <footer class="border-t border-slate-200/80 bg-white/80 px-6 py-8 backdrop-blur">
      <div class="w-full max-w-[1600px] mx-auto flex flex-col md:flex-row justify-between items-center gap-4 text-center md:text-left text-slate-600">
        <p class="text-sm">&copy; 2026 TravelLink. All rights reserved.</p>
        <div class="flex justify-center gap-6">
          <a href="#" class="no-underline transition-colors hover:text-emerald-700">Privacy</a>
          <a href="#" class="no-underline transition-colors hover:text-emerald-700">Terms</a>
          <a href="#" class="no-underline transition-colors hover:text-emerald-700">Contact</a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import Navigation from './Navigation.vue';
import NotificationCenter from './NotificationCenter.vue';
import { useUiStore } from '../stores/ui_store';

const uiStore = useUiStore();

const isLoading = computed(() => uiStore.isLoading);
</script>
