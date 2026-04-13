<template>
  <div class="min-h-screen bg-gray-50 px-4 py-8">
    <div class="max-w-3xl mx-auto">
      <!-- Profile Header -->
      <div class="bg-white rounded-lg shadow-sm p-8 mb-6 text-center">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">@{{ authStore.user?.username || 'User' }}</h1>
        <p class="text-gray-600 mb-1">{{ authStore.user?.email }}</p>
        <p class="text-blue-600 font-semibold">Role: {{ authStore.userRole }}</p>
      </div>

      <!-- Profile Content -->
      <div class="space-y-6">
        <!-- User Info Section -->
        <div class="bg-white rounded-lg shadow-sm p-6">
          <h2 class="text-xl font-bold text-gray-900 mb-4 pb-3 border-b-2 border-gray-100">Profile Information</h2>
          <div class="space-y-3">
            <div class="flex justify-between items-center pb-3 border-b border-gray-100">
              <label class="font-semibold text-gray-700">Username:</label>
              <span class="text-gray-600">{{ authStore.user?.username }}</span>
            </div>
            <div class="flex justify-between items-center pb-3 border-b border-gray-100">
              <label class="font-semibold text-gray-700">Email:</label>
              <span class="text-gray-600">{{ authStore.user?.email }}</span>
            </div>
            <div class="flex justify-between items-center">
              <label class="font-semibold text-gray-700">Member Since:</label>
              <span class="text-gray-600">{{ formatDate(authStore.user?.created_at) }}</span>
            </div>
          </div>
          <button @click="navigateTo('/profile/edit')" class="mt-6 px-4 py-2 bg-blue-600 text-white rounded-md text-sm font-medium hover:bg-blue-700 transition-all">
            Edit Profile
          </button>
        </div>

        <!-- My Plans Section -->
        <div class="bg-white rounded-lg shadow-sm p-6">
          <h2 class="text-xl font-bold text-gray-900 mb-4 pb-3 border-b-2 border-gray-100">My Travel Plans</h2>
          <div v-if="userPlans.length > 0" class="space-y-3">
            <div v-for="plan in userPlans" :key="plan.id" class="border border-gray-200 p-4 rounded-md bg-gray-50 hover:shadow-md hover:border-blue-600 transition-all">
              <h3 class="font-bold text-gray-900 mb-2">{{ plan.title }}</h3>
              <p class="text-gray-600 mb-1">{{ plan.destination }}</p>
              <p class="font-semibold text-blue-600 mb-3">Status: {{ plan.status }}</p>
              <router-link :to="`/plans/${plan.id}`" class="inline-block px-3 py-1.5 text-xs font-medium bg-blue-600 text-white rounded-md no-underline hover:bg-blue-700 transition-all">
                View
              </router-link>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-600">
            <p class="mb-4">You haven't created any travel plans yet.</p>
            <router-link to="/create-plan" class="inline-block px-4 py-2 bg-blue-600 text-white rounded-md text-sm font-medium no-underline hover:bg-blue-700 transition-all">
              Create Your First Plan
            </router-link>
          </div>
        </div>

        <!-- Actions Section -->
        <div class="bg-white rounded-lg shadow-sm p-6">
          <h2 class="text-xl font-bold text-gray-900 mb-4 pb-3 border-b-2 border-gray-100">Actions</h2>
          <div class="space-y-3">
            <button @click="handleChangePassword" class="w-full px-4 py-2 bg-gray-100 text-gray-900 rounded-md text-sm font-medium border border-gray-300 hover:bg-gray-200 transition-all">
              Change Password
            </button>
            <button @click="handleLogout" class="w-full px-4 py-2 bg-red-600 text-white rounded-md text-sm font-medium hover:bg-red-700 transition-all">
              Logout
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth_store';
import { useUiStore } from '../stores/ui_store';
import { planService, type TravelPlan } from '../services/plan_service';

const router = useRouter();
const authStore = useAuthStore();
const uiStore = useUiStore();

const userPlans = ref<TravelPlan[]>([]);

const formatDate = (dateStr?: string) => {
  if (!dateStr) return 'Unknown';
  return new Date(dateStr).toLocaleDateString();
};

const navigateTo = (path: string) => {
  router.push(path);
};

const handleLogout = () => {
  authStore.logout();
  uiStore.showSuccess('Logged out');
  router.push('/browse');
};

const handleChangePassword = () => {
  // TODO: Open password change modal
  uiStore.showInfo('Password change feature coming soon');
};

onMounted(async () => {
  try {
    const plans = await planService.getUserPlans();
    userPlans.value = plans;
  } catch (error) {
    uiStore.showError('Failed to load your plans');
  }
});
</script>
