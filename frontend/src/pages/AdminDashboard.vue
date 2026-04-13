<template>
  <div class="min-h-screen bg-gray-50 px-4 py-8">
    <!-- Dashboard Header -->
    <div class="max-w-7xl mx-auto mb-8">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">Admin Dashboard</h1>
      <p class="text-gray-600">Manage content, users, and platform moderation</p>
    </div>

    <!-- Dashboard Grid -->
    <div class="max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Flagged Plans Section -->
      <div class="bg-white rounded-lg shadow-sm p-6">
        <h2 class="text-xl font-bold text-gray-900 mb-1 pb-3 border-b-2 border-gray-100">Flagged Plans</h2>
        <p class="text-gray-600 text-sm mb-4">Review and moderate flagged travel plans</p>
        <button @click="loadFlaggedPlans" class="px-4 py-2 bg-blue-600 text-white rounded-md text-sm font-medium hover:bg-blue-700 transition-all">
          Load Flagged Plans
        </button>
        <div v-if="flaggedPlans.length > 0" class="mt-4 space-y-3">
          <div v-for="plan in flaggedPlans" :key="plan.id" class="border border-gray-200 p-3 rounded-md bg-gray-50">
            <h3 class="font-bold text-gray-900 mb-1">{{ plan.title }}</h3>
            <p class="text-sm text-gray-600 mb-3">Author: {{ plan.author?.username }}</p>
            <div class="flex gap-2">
              <button @click="approvePlan(plan.id)" class="px-3 py-1.5 text-xs font-medium bg-green-600 text-white rounded-md hover:bg-green-700 transition-all">
                Approve
              </button>
              <button @click="rejectPlan(plan.id)" class="px-3 py-1.5 text-xs font-medium bg-red-600 text-white rounded-md hover:bg-red-700 transition-all">
                Reject
              </button>
            </div>
          </div>
        </div>
        <div v-else class="mt-4 text-center py-6 text-gray-600">
          <p>No flagged plans at the moment</p>
        </div>
      </div>

      <!-- Pending Nodes Section -->
      <div class="bg-white rounded-lg shadow-sm p-6">
        <h2 class="text-xl font-bold text-gray-900 mb-1 pb-3 border-b-2 border-gray-100">Pending Node Approvals</h2>
        <p class="text-gray-600 text-sm mb-4">Review and approve user-created nodes</p>
        <button @click="loadPendingNodes" class="px-4 py-2 bg-blue-600 text-white rounded-md text-sm font-medium hover:bg-blue-700 transition-all">
          Load Pending Nodes
        </button>
        <div v-if="pendingNodes.length > 0" class="mt-4 space-y-3">
          <div v-for="node in pendingNodes" :key="node.id" class="border border-gray-200 p-3 rounded-md bg-gray-50">
            <h3 class="font-bold text-gray-900 mb-1">{{ node.details?.name || 'Node' }}</h3>
            <p class="text-sm text-gray-600 mb-3">Type: {{ node.type }}</p>
            <div class="flex gap-2">
              <button @click="approveNode(node.id)" class="px-3 py-1.5 text-xs font-medium bg-green-600 text-white rounded-md hover:bg-green-700 transition-all">
                Approve
              </button>
              <button @click="rejectNode(node.id)" class="px-3 py-1.5 text-xs font-medium bg-red-600 text-white rounded-md hover:bg-red-700 transition-all">
                Reject
              </button>
            </div>
          </div>
        </div>
        <div v-else class="mt-4 text-center py-6 text-gray-600">
          <p>No pending nodes to review</p>
        </div>
      </div>

      <!-- Promotion Requests Section -->
      <div class="bg-white rounded-lg shadow-sm p-6">
        <h2 class="text-xl font-bold text-gray-900 mb-1 pb-3 border-b-2 border-gray-100">Promotion Requests</h2>
        <p class="text-gray-600 text-sm mb-4">Review user promotion requests</p>
        <button @click="loadPromotionRequests" class="px-4 py-2 bg-blue-600 text-white rounded-md text-sm font-medium hover:bg-blue-700 transition-all">
          Load Requests
        </button>
        <div v-if="promotionRequests.length > 0" class="mt-4 space-y-3">
          <div v-for="request in promotionRequests" :key="request.id" class="border border-gray-200 p-3 rounded-md bg-gray-50">
            <h3 class="font-bold text-gray-900 mb-1">{{ request.user?.username }}</h3>
            <p v-if="request.plan" class="text-sm text-gray-600 mb-1">Plan: {{ request.plan.title }}</p>
            <p class="text-sm font-semibold text-blue-600 mb-3">Status: {{ request.status }}</p>
            <div v-if="request.status === 'pending'" class="flex gap-2">
              <button @click="approvePromotion(request.id)" class="px-3 py-1.5 text-xs font-medium bg-green-600 text-white rounded-md hover:bg-green-700 transition-all">
                Approve
              </button>
              <button @click="rejectPromotion(request.id)" class="px-3 py-1.5 text-xs font-medium bg-red-600 text-white rounded-md hover:bg-red-700 transition-all">
                Reject
              </button>
            </div>
          </div>
        </div>
        <div v-else class="mt-4 text-center py-6 text-gray-600">
          <p>No pending promotion requests</p>
        </div>
      </div>

      <!-- Statistics Section -->
      <div class="bg-white rounded-lg shadow-sm p-6 lg:col-span-2">
        <h2 class="text-xl font-bold text-gray-900 mb-4 pb-3 border-b-2 border-gray-100">Platform Statistics</h2>
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
          <div class="bg-gradient-to-br from-indigo-500 to-purple-600 text-white p-4 rounded-lg">
            <p class="text-sm text-indigo-100 mb-1">Total Users</p>
            <p class="text-3xl font-bold">{{ stats.totalUsers }}</p>
          </div>
          <div class="bg-gradient-to-br from-blue-500 to-cyan-600 text-white p-4 rounded-lg">
            <p class="text-sm text-blue-100 mb-1">Total Plans</p>
            <p class="text-3xl font-bold">{{ stats.totalPlans }}</p>
          </div>
          <div class="bg-gradient-to-br from-orange-500 to-red-600 text-white p-4 rounded-lg">
            <p class="text-sm text-orange-100 mb-1">Pending Approvals</p>
            <p class="text-3xl font-bold">{{ stats.pendingApprovals }}</p>
          </div>
          <div class="bg-gradient-to-br from-rose-500 to-pink-600 text-white p-4 rounded-lg">
            <p class="text-sm text-rose-100 mb-1">Flagged Content</p>
            <p class="text-3xl font-bold">{{ stats.flaggedContent }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useUiStore } from '../stores/ui_store';
import { type TravelPlan } from '../services/plan_service';

const uiStore = useUiStore();

const flaggedPlans = ref<TravelPlan[]>([]);
const pendingNodes = ref<any[]>([]);
const promotionRequests = ref<any[]>([]);

const stats = ref({
  totalUsers: 0,
  totalPlans: 0,
  pendingApprovals: 0,
  flaggedContent: 0,
});

const loadFlaggedPlans = async () => {
  try {
    // TODO: Implement API call to fetch flagged plans
    uiStore.showInfo('Feature coming soon');
  } catch (error) {
    uiStore.showError('Failed to load flagged plans');
  }
};

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
    // TODO: Implement API call to fetch promotion requests
    uiStore.showInfo('Feature coming soon');
  } catch (error) {
    uiStore.showError('Failed to load promotion requests');
  }
};

const approvePlan = async (planId: string) => {
  // TODO: Implement approve logic
  uiStore.showSuccess('Plan approved');
};

const rejectPlan = async (planId: string) => {
  // TODO: Implement reject logic
  uiStore.showSuccess('Plan rejected');
};

const approveNode = async (nodeId: string) => {
  // TODO: Implement approve logic
  uiStore.showSuccess('Node approved');
};

const rejectNode = async (nodeId: string) => {
  // TODO: Implement reject logic
  uiStore.showSuccess('Node rejected');
};

const approvePromotion = async (requestId: string) => {
  // TODO: Implement approve logic
  uiStore.showSuccess('Promotion approved');
};

const rejectPromotion = async (requestId: string) => {
  // TODO: Implement reject logic
  uiStore.showSuccess('Promotion rejected');
};
</script>
