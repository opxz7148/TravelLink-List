<template>
  <div class="w-full max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <!-- Page Header -->
    <div class="relative overflow-hidden rounded-3xl border border-emerald-100 bg-gradient-to-br from-emerald-50 via-white to-sky-50 p-8 mb-12">
      <div class="pointer-events-none absolute -right-16 -top-16 h-44 w-44 rounded-full bg-emerald-200/50 blur-3xl"></div>
      <h1 class="tl-title text-4xl md:text-5xl font-bold text-slate-900 mb-2">Find Your Next Route</h1>
      <p class="text-slate-600 text-base md:text-lg">Discover hand-crafted itineraries shared by travelers around the world.</p>
    </div>

    <!-- Filters Section -->
    <div class="tl-surface flex flex-col gap-4 mx-auto my-12 p-6 md:p-8">
      <div class="relative flex items-center">
        <svg class="absolute left-3 w-5 h-5 text-slate-400 pointer-events-none" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <circle cx="11" cy="11" r="8" />
          <path d="m21 21-4.35-4.35" />
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search by destination..."
          @keyup.enter="search"
          class="w-full pl-10 pr-4 py-3 border border-slate-300 rounded-xl text-sm focus:outline-none focus:border-emerald-500 focus:ring-4 focus:ring-emerald-100"
        />
      </div>

      <div class="flex gap-3">
        <select v-model="sortBy" class="flex-1 px-3 py-2 border border-slate-300 rounded-xl text-sm bg-white focus:outline-none focus:border-emerald-500">
          <option value="recent">Recent</option>
          <option value="popular">Most Popular</option>
          <option value="rated">Highest Rated</option>
        </select>

        <button @click="search" class="flex-1 px-5 py-2 bg-gradient-to-r from-emerald-600 to-sky-600 text-white rounded-xl text-sm font-semibold cursor-pointer transition hover:from-emerald-700 hover:to-sky-700">Search</button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="tl-surface flex flex-col items-center justify-center py-20 gap-4 my-8">
      <div class="w-10 h-10 border-4 border-slate-200 border-t-emerald-600 rounded-full animate-spin"></div>
      <p class="text-slate-600">Loading plans...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="plans.length === 0" class="tl-surface flex flex-col items-center justify-center py-20 text-slate-600 my-8">
      <svg class="w-16 h-16 text-slate-300 mb-4" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <path d="M9 12h6m-6 4h6m2-16H7a2 2 0 00-2 2v16a2 2 0 002 2h10a2 2 0 002-2V2a2 2 0 00-2-2z" />
      </svg>
      <h2 class="tl-title text-xl font-semibold text-slate-900 mb-2">No plans found</h2>
      <p>Try adjusting your search filters</p>
    </div>

    <!-- Plans Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5 mb-12 mt-8">
      <PlanCard
        v-for="plan in sortedPlans"
        :key="plan.id"
        :plan="plan"
        @view="viewPlan(plan.id)"
      />
    </div>

    <!-- Pagination -->
    <div v-if="plans.length > 0" class="tl-surface flex flex-col sm:flex-row justify-center items-center gap-3 p-6 md:p-8 mt-12">
      <button
        @click="previousPage"
        :disabled="currentPage === 1"
        class="px-4 py-2 bg-slate-200 text-slate-900 rounded-lg text-sm font-medium cursor-pointer transition hover:bg-slate-300 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        ← Previous
      </button>

      <div class="text-sm text-slate-900">
        Page {{ currentPage }} of {{ totalPages }}
        <span class="text-xs text-slate-600 ml-2">({{ totalPlans }} total plans)</span>
      </div>

      <button
        @click="nextPage"
        :disabled="currentPage >= totalPages"
        class="px-4 py-2 bg-slate-200 text-slate-900 rounded-lg text-sm font-medium cursor-pointer transition hover:bg-slate-300 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Next →
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import PlanCard from '../components/PlanCard.vue';
import { planService } from '../services/plan_service';
import type { TravelPlan } from '../services/plan_service';

const router = useRouter();

const plans = ref<TravelPlan[]>([]);
const loading = ref(false);
const currentPage = ref(1);
const pageSize = ref(12);
const totalPlans = ref(0);
const searchQuery = ref('');
const sortBy = ref('recent');

const totalPages = computed(() => {
  return Math.max(1, Math.ceil((Number(totalPlans.value) || 0) / pageSize.value));
});

const sortedPlans = computed(() => {
  const sorted = [...plans.value];
  
  switch (sortBy.value) {
    case 'popular':
      // Sort by comment count (descending)
      return sorted.sort((a: TravelPlan, b: TravelPlan) => (b.comment_count || 0) - (a.comment_count || 0));
    case 'rated':
      // Sort by rating average (descending)
      return sorted.sort((a: TravelPlan, b: TravelPlan) => {
        const ratingA = (a.rating_count || 0) > 0 ? (a.rating_sum || 0) / (a.rating_count || 1) : 0;
        const ratingB = (b.rating_count || 0) > 0 ? (b.rating_sum || 0) / (b.rating_count || 1) : 0;
        return ratingB - ratingA;
      });
    case 'recent':
    default:
      // Sort by creation date (descending - newest first)
      return sorted.sort((a: TravelPlan, b: TravelPlan) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
  }
});

async function loadPlans(): Promise<void> {
  try {
    loading.value = true;

    const result = await planService.listPlans({
      page: currentPage.value,
      limit: pageSize.value,
    });

    plans.value = result.plans;
    totalPlans.value = Number(result.total) || 0;
  } catch (error) {
    console.error('Failed to load plans:', error);
    // TODO: Show error toast notification
  } finally {
    loading.value = false;
  }
}

async function search(): Promise<void> {
  if (!searchQuery.value.trim()) {
    currentPage.value = 1;
    loadPlans();
    return;
  }

  try {
    loading.value = true;
    currentPage.value = 1;

    const result = await planService.searchPlans({
      q: searchQuery.value,
      page: currentPage.value,
      limit: pageSize.value,
    });

    plans.value = result.plans;
    totalPlans.value = Number(result.total) || 0;
  } catch (error) {
    console.error('Search failed:', error);
    // TODO: Show error toast notification
  } finally {
    loading.value = false;
  }
}

function viewPlan(planId: string): void {
  router.push(`/plans/${planId}`);
}

function previousPage(): void {
  if (currentPage.value > 1) {
    currentPage.value--;
    loadPlans();
  }
}

function nextPage(): void {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
    loadPlans();
  }
}

onMounted(() => {
  loadPlans();
});
</script>
