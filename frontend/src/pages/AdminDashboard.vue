<template>
  <div class="min-h-screen bg-gray-50 px-4 py-12">
    <!-- Dashboard Header -->
    <div class="w-full max-w-400 mx-auto mb-12">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">Admin Dashboard</h1>
      <p class="text-gray-600">Manage content, users, and platform moderation</p>
    </div>

    <!-- Dashboard Grid -->
    <div class="w-full max-w-400 mx-auto grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Promotion Requests Section -->
      <div class="bg-white rounded-lg shadow-sm p-6">
        <h2 class="text-xl font-bold text-gray-900 mb-1 pb-3 border-b-2 border-gray-100">Promotion Requests</h2>
        <p class="text-gray-600 text-sm mb-4">Review user promotion requests to traveller status</p>
        <button 
          @click="loadPromotionRequests" 
          :disabled="isLoadingPromotions"
          class="px-4 py-2 bg-blue-600 text-white rounded-md text-sm font-medium hover:bg-blue-700 transition-all disabled:opacity-50"
        >
          {{ isLoadingPromotions ? 'Loading...' : 'Refresh Requests' }}
        </button>
        <div v-if="promotionRequests.length > 0" class="mt-4 space-y-3">
          <div v-for="request in promotionRequests" :key="request.id" class="border border-gray-200 p-4 rounded-md bg-gray-50 hover:bg-gray-100 transition-colors">
            <div class="flex justify-between items-start gap-4">
              <div class="flex-1">
                <h3 class="font-bold text-gray-900 mb-1">
                  <router-link 
                    :to="`/profile/${request.user?.username}`"
                    class="text-blue-600 hover:underline"
                  >
                    {{ request.user?.username || 'Unknown User' }}
                  </router-link>
                </h3>
                <p v-if="request.user?.email" class="text-sm text-gray-600 mb-1">
                  {{ request.user.email }}
                </p>
                <p v-if="request.user?.role" class="text-sm text-gray-600 mb-2">
                  Current Role: <span class="font-semibold capitalize">{{ request.user.role }}</span>
                </p>
                <p v-if="request.plan" class="text-sm text-gray-700 mb-2">
                  <strong>Referenced Plan:</strong>
                  <router-link 
                    :to="`/plans/${request.plan.id}`"
                    class="text-blue-600 hover:underline font-medium"
                  >
                    {{ request.plan.title }}
                  </router-link>
                </p>
                <p class="text-sm text-gray-600">
                  Submitted: {{ new Date(request.created_at).toLocaleDateString() }}
                </p>
              </div>
              <div class="text-right">
                <span 
                  class="inline-block px-3 py-1 rounded-full text-xs font-semibold"
                  :class="{
                    'bg-yellow-100 text-yellow-800': request.status === 'pending',
                    'bg-green-100 text-green-800': request.status === 'approved',
                    'bg-red-100 text-red-800': request.status === 'rejected'
                  }"
                >
                  {{ request.status.toUpperCase() }}
                </span>
              </div>
            </div>
            <div v-if="request.status === 'pending'" class="mt-4 flex gap-2">
              <button 
                @click="openPromotionReview(request)"
                class="flex-1 px-3 py-2 text-sm font-medium bg-green-600 text-white rounded-md hover:bg-green-700 transition-all"
              >
                Review & Approve
              </button>
              <button 
                @click="openPromotionReview(request)"
                class="flex-1 px-3 py-2 text-sm font-medium bg-red-600 text-white rounded-md hover:bg-red-700 transition-all"
              >
                Review & Reject
              </button>
            </div>
          </div>
        </div>
        <div v-else class="mt-4 text-center py-6 text-gray-600">
          <p>No pending promotion requests</p>
        </div>
      </div>
    </div>
  </div>

  <!-- Promotion Review Modal -->
  <div v-if="selectedRequestForReview" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
    <div class="bg-white rounded-lg shadow-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto">
      <!-- Modal Header -->
      <div class="border-b border-gray-200 p-6">
        <h2 class="text-2xl font-bold text-gray-900">Review Promotion Request</h2>
        <p class="text-gray-600 text-sm mt-1">Admin review and decision</p>
      </div>

      <!-- Modal Body -->
      <div class="p-6 space-y-6">
        <!-- User Information -->
        <div>
          <h3 class="text-lg font-semibold text-gray-900 mb-3">User Information</h3>
          <div class="bg-gray-50 p-4 rounded-md space-y-2">
            <div>
              <p class="text-sm text-gray-600">Username</p>
              <p class="font-semibold text-gray-900">{{ selectedRequestForReview.user?.username }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-600">Email</p>
              <p class="font-semibold text-gray-900">{{ selectedRequestForReview.user?.email }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-600">Current Role</p>
              <p class="font-semibold text-gray-900 capitalize">{{ selectedRequestForReview.user?.role }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-600">Request Date</p>
              <p class="font-semibold text-gray-900">{{ new Date(selectedRequestForReview.created_at).toLocaleString() }}</p>
            </div>
          </div>
        </div>

        <!-- Referenced Plan (if any) -->
        <div v-if="selectedRequestForReview.plan">
          <h3 class="text-lg font-semibold text-gray-900 mb-3">Referenced Travel Plan</h3>
          <div class="bg-gray-50 p-4 rounded-md space-y-2">
            <div>
              <p class="text-sm text-gray-600">Plan Title</p>
              <p class="font-semibold text-gray-900">{{ selectedRequestForReview.plan.title }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-600">Destination</p>
              <p class="font-semibold text-gray-900">{{ selectedRequestForReview.plan.destination }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-600">Status</p>
              <p class="font-semibold text-gray-900 capitalize">{{ selectedRequestForReview.plan.status }}</p>
            </div>
            <div v-if="formatRating(selectedRequestForReview.plan) !== '0.0'">
              <p class="text-sm text-gray-600">Rating</p>
              <p class="font-semibold text-gray-900">
                {{ formatRating(selectedRequestForReview.plan) }} / 5 
                ({{ selectedRequestForReview.plan.rating_count || 0 }} ratings)
              </p>
            </div>
          </div>
        </div>

        <!-- Admin Notes -->
        <div>
          <label for="review-notes" class="block text-sm font-semibold text-gray-900 mb-2">
            Admin Notes (Optional)
          </label>
          <textarea
            id="review-notes"
            v-model="reviewNotes"
            placeholder="Add notes about this promotion decision..."
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            rows="4"
            maxlength="500"
          ></textarea>
          <p class="text-xs text-gray-500 mt-1">{{ reviewNotes.length }}/500 characters</p>
        </div>

        <!-- Promotion Outcome -->
        <div class="bg-blue-50 border border-blue-200 p-4 rounded-md">
          <p class="text-sm font-semibold text-blue-900 mb-2">If Approved:</p>
          <ul class="text-sm text-blue-800 space-y-1 ml-4 list-disc">
            <li>User role will be upgraded to <strong>traveller</strong></li>
            <li>User will receive a notification of approval</li>
            <li>User can now create and publish travel plans</li>
          </ul>
        </div>
      </div>

      <!-- Modal Footer -->
      <div class="border-t border-gray-200 p-6 flex gap-3 justify-end">
        <button
          @click="closePromotionReview"
          class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 font-medium hover:bg-gray-50 transition-all"
        >
          Cancel
        </button>
        <button
          @click="rejectPromotion(selectedRequestForReview.id)"
          class="px-4 py-2 bg-red-600 text-white rounded-md font-medium hover:bg-red-700 transition-all"
        >
          Reject Request
        </button>
        <button
          @click="approvePromotion(selectedRequestForReview.id)"
          class="px-4 py-2 bg-green-600 text-white rounded-md font-medium hover:bg-green-700 transition-all"
        >
          Approve & Promote
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useUiStore } from '../stores/ui_store';
import { useAuthStore } from '../stores/auth_store';
import { promotionService, type PromotionRequest } from '../services/promotion_service';

const uiStore = useUiStore();
const authStore = useAuthStore();

const pendingNodes = ref<any[]>([]);
const promotionRequests = ref<PromotionRequest[]>([]);
const isLoadingPromotions = ref(false);
const selectedRequestForReview = ref<PromotionRequest | null>(null);
const reviewNotes = ref('');

const loadPendingNodes = async () => {
  try {
    // TODO: Implement API call to fetch pending nodes
    uiStore.showInfo('Feature coming soon');
  } catch (error) {
    uiStore.showError('Failed to load pending nodes');
  }
};

const loadPromotionRequests = async () => {
  try {
    isLoadingPromotions.value = true;
    const response = await promotionService.getPendingRequests(1);
    promotionRequests.value = response.requests;
    uiStore.showSuccess(`Loaded ${response.requests.length} pending requests`);
  } catch (error: any) {
    uiStore.showError('Failed to load promotion requests');
    console.error('Error loading promotion requests:', error);
  } finally {
    isLoadingPromotions.value = false;
  }
};

const approveNode = async (nodeId: string) => {
  // TODO: Implement approve logic
  uiStore.showSuccess('Node approved');
};

const rejectNode = async (nodeId: string) => {
  // TODO: Implement reject logic
  uiStore.showSuccess('Node rejected');
};

const openPromotionReview = (request: PromotionRequest) => {
  selectedRequestForReview.value = request;
  reviewNotes.value = '';
};

const closePromotionReview = () => {
  selectedRequestForReview.value = null;
  reviewNotes.value = '';
};

const formatRating = (plan: { rating_average?: number; rating_sum?: number; rating_count?: number }) => {
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

const approvePromotion = async (requestId: string) => {
  try {
    await promotionService.approvePromotion(requestId, reviewNotes.value);
    // Reload the list
    await loadPromotionRequests();
    // Refresh current user data to update role if they're logged in and were promoted
    if (authStore.isAuthenticated) {
      await authStore.getCurrentUser();
    }
    uiStore.showSuccess('Promotion request approved successfully');
    closePromotionReview();
  } catch (error: any) {
    uiStore.showError('Failed to approve promotion request');
    console.error('Error approving promotion:', error);
  }
};

const rejectPromotion = async (requestId: string) => {
  try {
    await promotionService.rejectPromotion(requestId, reviewNotes.value);
    // Reload the list
    await loadPromotionRequests();
    uiStore.showSuccess('Promotion request rejected');
    closePromotionReview();
  } catch (error: any) {
    uiStore.showError('Failed to reject promotion request');
    console.error('Error rejecting promotion:', error);
  }
};

onMounted(async () => {
  await loadPromotionRequests();
});
</script>
