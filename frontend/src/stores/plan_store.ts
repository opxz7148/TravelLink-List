import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { planService, type TravelPlan, type PlanDetail } from '../services/plan_service';

export const usePlanStore = defineStore('plan', () => {
  // State
  const plans = ref<TravelPlan[]>([]);
  const currentPlan = ref<PlanDetail | null>(null);
  const totalPlans = ref(0);
  const currentPage = ref(1);
  const pageSize = ref(10);
  const isLoading = ref(false);
  const error = ref<string>('');

  // Search/Filter state
  const searchQuery = ref('');
  const destinationFilter = ref('');
  const sortBy = ref<'recent' | 'popular' | 'rated'>('recent');

  // Computed
  const totalPages = computed(() => Math.ceil(totalPlans.value / pageSize.value));
  const hasNextPage = computed(() => currentPage.value < totalPages.value);
  const hasPrevPage = computed(() => currentPage.value > 1);

  // Actions
  const listPlans = async (page: number = 1) => {
    isLoading.value = true;
    error.value = '';
    try {
      const response = await planService.listPlans({
        page,
        limit: pageSize.value,
        sort: sortBy.value,
      });
      plans.value = response.plans;
      totalPlans.value = response.total;
      currentPage.value = page;
      return response;
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to load plans';
      throw err;
    } finally {
      isLoading.value = false;
    }
  };

  const searchPlans = async (query: string, destination?: string, page: number = 1) => {
    isLoading.value = true;
    error.value = '';
    try {
      const response = await planService.searchPlans({
        q: query,
        destination,
        page,
        limit: pageSize.value,
      });
      plans.value = response.plans;
      totalPlans.value = response.total;
      currentPage.value = page;
      searchQuery.value = query;
      destinationFilter.value = destination || '';
      return response;
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Search failed';
      throw err;
    } finally {
      isLoading.value = false;
    }
  };

  const getPlanDetail = async (planId: string) => {
    isLoading.value = true;
    error.value = '';
    try {
      const response = await planService.getPlanDetail(planId);
      currentPlan.value = response;
      return response;
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to load plan details';
      throw err;
    } finally {
      isLoading.value = false;
    }
  };

  const createPlan = async (title: string, destination: string) => {
    isLoading.value = true;
    error.value = '';
    try {
      const response = await planService.createDraftPlan(title, destination);
      currentPlan.value = response;
      return response;
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to create plan';
      throw err;
    } finally {
      isLoading.value = false;
    }
  };

  const updatePlan = async (planId: string, updates: any) => {
    isLoading.value = true;
    error.value = '';
    try {
      const response = await planService.updatePlan(planId, updates);
      currentPlan.value = response;
      // Update in plans list if present
      const index = plans.value.findIndex((p) => p.id === planId);
      if (index >= 0) {
        plans.value[index] = response;
      }
      return response;
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to update plan';
      throw err;
    } finally {
      isLoading.value = false;
    }
  };

  const publishPlan = async (planId: string) => {
    return updatePlan(planId, { status: 'published' });
  };

  const unpublishPlan = async (planId: string) => {
    return updatePlan(planId, { status: 'draft' });
  };

  const deletePlan = async (planId: string) => {
    isLoading.value = true;
    error.value = '';
    try {
      await planService.deletePlan(planId);
      plans.value = plans.value.filter((p) => p.id !== planId);
      totalPlans.value--;
      if (currentPlan.value?.id === planId) {
        currentPlan.value = null;
      }
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to delete plan';
      throw err;
    } finally {
      isLoading.value = false;
    }
  };

  const setSortBy = (sort: 'recent' | 'popular' | 'rated') => {
    sortBy.value = sort;
    currentPage.value = 1;
  };

  const clearCurrentPlan = () => {
    currentPlan.value = null;
  };

  const clearSearch = () => {
    searchQuery.value = '';
    destinationFilter.value = '';
  };

  return {
    // State
    plans,
    currentPlan,
    totalPlans,
    currentPage,
    pageSize,
    isLoading,
    error,
    searchQuery,
    destinationFilter,
    sortBy,

    // Computed
    totalPages,
    hasNextPage,
    hasPrevPage,

    // Actions
    listPlans,
    searchPlans,
    getPlanDetail,
    createPlan,
    updatePlan,
    publishPlan,
    unpublishPlan,
    deletePlan,
    setSortBy,
    clearCurrentPlan,
    clearSearch,
  };
});
