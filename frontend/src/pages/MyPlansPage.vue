<template>
  <div class="w-full max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <!-- Header -->
    <div class="mb-12">
      <h1 class="text-4xl font-bold text-gray-900 mb-2">My Travel Plans</h1>
      <p class="text-gray-600">View and manage all your travel plans (draft and published)</p>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-20 my-8">
      <div class="text-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
        <p class="text-gray-600">Loading your plans...</p>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="plans.length === 0" class="bg-white rounded-lg border border-gray-200 p-8 text-center my-8">
      <div class="text-gray-500 mb-4">
        <svg class="mx-auto h-12 w-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">No plans yet</h3>
      <p v-if="authStore.userRole !== 'admin'" class="text-gray-600 mb-4">Get started by creating your first travel plan</p>
      <p v-else class="text-gray-600 mb-4">No plans to moderate</p>
      <router-link v-if="authStore.userRole !== 'admin'" to="/create-plan" class="text-blue-600 font-semibold hover:underline">
        Create a Plan →
      </router-link>
    </div>

    <!-- Plans List -->
    <div v-else class="space-y-6 mt-8">
      <!-- Filter Tabs -->
      <div class="flex gap-2 border-b border-gray-200">
        <button
          @click="filterStatus = 'all'"
          :class="[
            'px-4 py-2 font-medium transition-all',
            filterStatus === 'all'
              ? 'text-blue-600 border-b-2 border-blue-600'
              : 'text-gray-600 hover:text-gray-900',
          ]"
        >
          All Plans ({{ plans.length }})
        </button>
        <button
          @click="filterStatus = 'draft'"
          :class="[
            'px-4 py-2 font-medium transition-all',
            filterStatus === 'draft'
              ? 'text-blue-600 border-b-2 border-blue-600'
              : 'text-gray-600 hover:text-gray-900',
          ]"
        >
          Drafts ({{ draftCount }})
        </button>
        <button
          @click="filterStatus = 'published'"
          :class="[
            'px-4 py-2 font-medium transition-all',
            filterStatus === 'published'
              ? 'text-blue-600 border-b-2 border-blue-600'
              : 'text-gray-600 hover:text-gray-900',
          ]"
        >
          Published ({{ publishedCount }})
        </button>
      </div>

      <!-- Plans Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="plan in filteredPlans"
          :key="plan.id"
          class="bg-white rounded-lg border border-gray-200 p-5 hover:shadow-lg transition-all hover:border-blue-300"
        >
          <!-- Status Badge -->
          <div class="flex justify-between items-start mb-3">
            <h3 class="text-lg font-semibold text-gray-900 flex-1">{{ plan.title }}</h3>
            <span
              :class="[
                'inline-block px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap ml-2',
                plan.status === 'draft'
                  ? 'bg-yellow-100 text-yellow-800'
                  : plan.status === 'published'
                    ? 'bg-green-100 text-green-800'
                    : 'bg-red-100 text-red-800',
              ]"
            >
              {{ plan.status.charAt(0).toUpperCase() + plan.status.slice(1) }}
            </span>
          </div>

          <!-- Destination -->
          <p class="text-gray-600 mb-2">
            <span class="font-medium">📍 Destination:</span> {{ plan.destination }}
          </p>

          <!-- Description (truncated) -->
          <p v-if="plan.description" class="text-sm text-gray-600 mb-3 line-clamp-2">
            {{ plan.description }}
          </p>

          <!-- Plan Stats -->
          <div class="grid grid-cols-3 gap-2 py-3 border-y border-gray-200 mb-3">
            <div class="text-center">
              <p class="text-lg font-semibold text-gray-900">{{ plan.node_count }}</p>
              <p class="text-xs text-gray-500">Nodes</p>
            </div>
            <div class="text-center">
              <p class="text-lg font-semibold text-gray-900">{{ plan.rating_count }}</p>
              <p class="text-xs text-gray-500">Ratings</p>
            </div>
            <div class="text-center">
              <p class="text-lg font-semibold text-gray-900">{{ plan.comment_count }}</p>
              <p class="text-xs text-gray-500">Comments</p>
            </div>
          </div>

          <!-- Rating -->
          <div v-if="formatRating(plan) !== '0.0'" class="mb-3">
            <p class="text-sm text-gray-600">
              ⭐ <span class="font-semibold">{{ formatRating(plan) }}</span> Average Rating
            </p>
          </div>

          <!-- Created Date -->
          <p class="text-xs text-gray-500 mb-4">
            Created {{ formatDate(plan.created_at) }}
          </p>

          <!-- Action Buttons -->
          <div class="flex gap-2 flex-wrap">
            <router-link
              :to="`/plans/${plan.id}`"
              class="flex-1 min-w-fit px-3 py-2 bg-blue-600 text-white rounded-md text-sm font-medium hover:bg-blue-700 transition-all text-center"
            >
              View
            </router-link>
            <button
              v-if="plan.status === 'draft'"
              @click="editPlan(plan.id)"
              class="flex-1 min-w-fit px-3 py-2 bg-gray-200 text-gray-900 rounded-md text-sm font-medium hover:bg-gray-300 transition-all"
            >
              Edit
            </button>
            <!-- Simple users: Submit to Promote -->
            <button
              v-if="plan.status === 'draft' && authStore.userRole === 'simple'"
              @click="submitPromotePlan(plan.id)"
              :disabled="promotingId === plan.id"
              class="flex-1 min-w-fit px-3 py-2 bg-purple-600 text-white rounded-md text-sm font-medium hover:bg-purple-700 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ promotingId === plan.id ? 'Submitting...' : 'Submit to Promote' }}
            </button>
            <!-- Admin/Traveller: Publish -->
            <button
              v-else-if="plan.status === 'draft' && (authStore.userRole === 'traveller' || authStore.userRole === 'admin')"
              @click="publishPlan(plan.id)"
              class="flex-1 min-w-fit px-3 py-2 bg-green-600 text-white rounded-md text-sm font-medium hover:bg-green-700 transition-all"
            >
              Publish
            </button>
            <button
              @click="deletePlan(plan.id)"
              :disabled="deletingId === plan.id"
              class="px-3 py-2 bg-red-600 text-white rounded-md text-sm font-medium hover:bg-red-700 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              :title="`Delete ${plan.title}`"
            >
              {{ deletingId === plan.id ? '...' : '✕' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalCount > limit" class="flex justify-center gap-2 mt-12 py-6">
        <button
          @click="previousPage"
          :disabled="currentPage === 1"
          class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          ← Previous
        </button>
        <span class="px-4 py-2 text-gray-600">
          Page {{ currentPage }} of {{ totalPages }}
        </span>
        <button
          @click="nextPage"
          :disabled="currentPage >= totalPages"
          class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Next →
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useUiStore } from '../stores/ui_store';
import { useAuthStore } from '../stores/auth_store';
import { planService, type TravelPlan } from '../services/plan_service';
import { promotionService } from '../services/promotion_service';

const router = useRouter();
const uiStore = useUiStore();
const authStore = useAuthStore();

const loading = ref(true);
const plans = ref<TravelPlan[]>([]);
const filterStatus = ref<'all' | 'draft' | 'published'>('all');
const currentPage = ref(1);
const limit = ref(12);
const totalCount = ref(0);
const deletingId = ref<string | null>(null);
const promotingId = ref<string | null>(null);

const draftCount = computed(() => plans.value.filter((p) => p.status === 'draft').length);
const publishedCount = computed(() => plans.value.filter((p) => p.status === 'published').length);

const filteredPlans = computed(() => {
  if (filterStatus.value === 'all') {
    return plans.value;
  }
  return plans.value.filter((p) => p.status === filterStatus.value);
});

const totalPages = computed(() => Math.max(1, Math.ceil((Number(totalCount.value) || 0) / limit.value)));

const formatRating = (plan: TravelPlan) => {
  const average = Number(plan.rating_average);
  if (Number.isFinite(average) && average >= 0) {
    return average.toFixed(1);
  }

  const ratingSum = Number(plan.rating_sum);
  const ratingCount = Number(plan.rating_count);

  if (Number.isFinite(ratingSum) && Number.isFinite(ratingCount) && ratingCount > 0) {
    return (ratingSum / ratingCount).toFixed(1);
  }

  return '0.0';
};

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
};

const loadUserPlans = async () => {
  try {
    loading.value = true;
    const offset = (currentPage.value - 1) * limit.value;
    const result = await planService.getUserPlans(offset, limit.value);
    plans.value = result.plans;
    totalCount.value = Number(result.total) || 0;
  } catch (error) {
    console.error('Failed to load user plans:', error);
    uiStore.showError('Failed to load your plans');
  } finally {
    loading.value = false;
  }
};

const checkAdminAccess = () => {
  if (authStore.userRole === 'admin') {
    router.push('/admin');
  }
};

const editPlan = (planId: string) => {
  router.push(`/create-plan?edit=${planId}`);
};

const publishPlan = async (planId: string) => {
  try {
    await planService.publishPlan(planId);
    uiStore.showSuccess('Plan published successfully!');
    await loadUserPlans();
  } catch (error) {
    console.error('Failed to publish plan:', error);
    uiStore.showError('Failed to publish plan');
  }
};

const submitPromotePlan = async (planId: string) => {
  const plan = plans.value.find((p) => p.id === planId);
  if (!plan) return;

  try {
    promotingId.value = planId;
    await promotionService.submitPromotionRequest('', planId);
    uiStore.showSuccess('Plan submitted for promotion! Admins will review and decide.');
    await loadUserPlans();
  } catch (error) {
    console.error('Failed to submit plan for promotion:', error);
    uiStore.showError('Failed to submit plan for promotion');
  } finally {
    promotingId.value = null;
  }
};

const deletePlan = async (planId: string) => {
  const plan = plans.value.find((p) => p.id === planId);
  if (!plan) return;

  // Confirmation dialog
  const confirmed = window.confirm(`Are you sure you want to delete "${plan.title}"? This action cannot be undone.`);
  if (!confirmed) return;

  try {
    deletingId.value = planId;
    await planService.deletePlan(planId);
    uiStore.showSuccess('Plan deleted successfully!');
    await loadUserPlans();
  } catch (error) {
    console.error('Failed to delete plan:', error);
    uiStore.showError('Failed to delete plan');
  } finally {
    deletingId.value = null;
  }
};

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
    loadUserPlans();
  }
};

const previousPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--;
    loadUserPlans();
  }
};

onMounted(() => {
  checkAdminAccess();
  loadUserPlans();
});
</script>
