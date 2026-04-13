<template>
  <div class="max-w-7xl mx-auto px-4 py-6 min-h-screen">
    <!-- Page Header -->
    <div class="text-center mb-8">
      <h1 class="text-4xl font-bold text-gray-900 mb-2">Create a New Travel Plan</h1>
      <p class="text-gray-600">Design your perfect itinerary with our curated nodes</p>
    </div>

    <!-- Step 1: Plan Basics -->
    <div v-if="!currentPlan" class="max-w-md mx-auto">
      <div class="bg-white p-8 rounded-lg border border-gray-200 shadow-sm">
        <h2 class="text-xl font-bold text-gray-900 mb-6">Step 1: Plan Basics</h2>

        <div class="flex flex-col gap-2 mb-5">
          <label for="title" class="text-sm font-semibold text-gray-700">Travel Plan Title *</label>
          <input
            id="title"
            v-model="planTitle"
            type="text"
            placeholder="e.g., Summer European Adventure"
            maxlength="100"
            class="px-3 py-3 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100 transition-all font-inherit"
          />
        </div>

        <div class="flex flex-col gap-2 mb-5">
          <label for="destination" class="text-sm font-semibold text-gray-700">Destination *</label>
          <input
            id="destination"
            v-model="planDestination"
            type="text"
            placeholder="e.g., Paris, France"
            maxlength="100"
            class="px-3 py-3 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100 transition-all font-inherit"
          />
        </div>

        <button
          @click="createDraft"
          :disabled="!planTitle.trim() || !planDestination.trim() || loading"
          class="w-full px-6 py-3.5 bg-blue-600 text-white rounded-md text-base font-medium cursor-pointer hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
        >
          {{ loading ? 'Creating...' : 'Create Draft Plan' }}
        </button>
      </div>

      <div v-if="errorMessage" class="mt-6 p-4 bg-red-50 text-red-800 rounded-lg border-l-4 border-red-500">
        {{ errorMessage }}
      </div>
    </div>

    <!-- Step 2: Plan Builder -->
    <div v-else class="space-y-6">
      <!-- Builder Header -->
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center p-5 bg-white rounded-lg border border-gray-200 gap-4">
        <div>
          <h2 class="text-2xl font-bold text-gray-900 mb-2">{{ currentPlan.title }}</h2>
          <p class="text-sm text-gray-600">📍 {{ currentPlan.destination }}</p>
        </div>

        <button @click="resetPlan" class="px-4 py-2 bg-gray-200 text-gray-900 rounded-md text-sm font-medium hover:bg-gray-300 transition-all">
          Start Over
        </button>
      </div>

      <!-- Builder Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 items-start">
        <!-- Editor Column -->
        <div class="space-y-4">
          <LinkedListEditor
            :selected-nodes="selectedNodes"
            :available-nodes="availableNodes"
            @reorder="reorderNodes"
            @remove="removeNodeAtIndex"
          />

          <!-- Publish Section -->
          <div v-if="selectedNodes.length > 0" class="p-5 bg-gray-100 rounded-lg text-center">
            <button
              @click="publishPlan"
              :disabled="loading"
              class="w-full px-6 py-3 bg-green-600 text-white rounded-md text-base font-medium hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
            >
              {{ loading ? 'Publishing...' : 'Publish Plan' }}
            </button>
            <p class="text-xs text-gray-600 mt-3">Your plan will be visible to other users once published</p>
          </div>
        </div>

        <!-- Selector Column -->
        <div>
          <NodeSelector
            :selected-node-ids="selectedNodes"
            :is-loading="loadingNodes"
            @select="addNode"
            @deselect="removeNode"
          />
        </div>
      </div>

      <!-- Error Alert -->
      <div v-if="errorMessage" class="p-4 bg-red-50 text-red-800 rounded-lg border-l-4 border-red-500">
        {{ errorMessage }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import LinkedListEditor from '../components/LinkedListEditor.vue';
import NodeSelector from '../components/NodeSelector.vue';
import { planService } from '../services/plan_service';
import type { TravelPlan } from '../services/plan_service';
import { nodeService } from '../services/node_service';
import type { Node } from '../services/node_service';
import { useAuthStore } from '../stores/auth_store';

const router = useRouter();
const authStore = useAuthStore();

const planTitle = ref('');
const planDestination = ref('');
const currentPlan = ref<TravelPlan | null>(null);
const selectedNodes = ref<string[]>([]);
const availableNodes = ref<Node[]>([]);
const errorMessage = ref('');
const loading = ref(false);
const loadingNodes = ref(false);

async function createDraft(): Promise<void> {
  if (!planTitle.value.trim() || !planDestination.value.trim()) {
    errorMessage.value = 'Please fill in all required fields';
    return;
  }

  try {
    loading.value = true;
    errorMessage.value = '';

    currentPlan.value = await planService.createDraftPlan(
      planTitle.value.trim(),
      planDestination.value.trim()
    );

    // Load available nodes
    await loadNodes();
  } catch (error: any) {
    console.error('Failed to create draft plan:', error);
    errorMessage.value = error.message || 'Failed to create draft plan';
  } finally {
    loading.value = false;
  }
}

async function loadNodes(): Promise<void> {
  try {
    loadingNodes.value = true;

    const { nodes } = await nodeService.listApprovedNodes({
      approved_only: true,
    });

    availableNodes.value = nodes;
  } catch (error) {
    console.error('Failed to load nodes:', error);
    errorMessage.value = 'Failed to load available nodes';
  } finally {
    loadingNodes.value = false;
  }
}

function addNode(nodeId: string): void {
  if (selectedNodes.value.includes(nodeId)) return;

  // Validation: no 2 consecutive attractions
  if (selectedNodes.value.length > 0) {
    const lastNodeId = selectedNodes.value[selectedNodes.value.length - 1];
    const lastNode = availableNodes.value.find((n) => n.id === lastNodeId);
    const newNode = availableNodes.value.find((n) => n.id === nodeId);

    if (
      lastNode &&
      newNode &&
      lastNode.type === 'attraction' &&
      newNode.type === 'attraction'
    ) {
      errorMessage.value = 'Cannot add two consecutive attractions. Please add a transition node first.';
      return;
    }
  }

  selectedNodes.value.push(nodeId);
  errorMessage.value = '';
}

function removeNode(nodeId: string): void {
  const index = selectedNodes.value.indexOf(nodeId);
  if (index !== -1) {
    selectedNodes.value.splice(index, 1);
  }
}

function removeNodeAtIndex(index: number): void {
  selectedNodes.value.splice(index, 1);
}

function reorderNodes(newOrder: string[]): void {
  selectedNodes.value = newOrder;
}

function resetPlan(): void {
  currentPlan.value = null;
  selectedNodes.value = [];
  planTitle.value = '';
  planDestination.value = '';
  errorMessage.value = '';
}

async function publishPlan(): Promise<void> {
  if (!currentPlan.value || selectedNodes.value.length === 0) {
    errorMessage.value = 'Please add at least one node to your plan before publishing';
    return;
  }

  try {
    loading.value = true;
    errorMessage.value = '';

    // Add nodes to plan
    await planService.addNodesToPlan(currentPlan.value.id, selectedNodes.value);

    // Publish plan
    await planService.publishPlan(currentPlan.value.id);

    // Redirect to plan view
    router.push(`/plans/${currentPlan.value.id}`);
  } catch (error: any) {
    console.error('Failed to publish plan:', error);
    errorMessage.value = error.message || 'Failed to publish plan';
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  // Ensure user is authenticated and is traveller or admin
  if (!authStore.isAuthenticated) {
    router.push('/login');
    return;
  }

  if (!authStore.isTraveller) {
    errorMessage.value = 'You must be a traveller to create plans. Submit a promotion request to upgrade your account.';
  }

  // Load nodes
  loadNodes();
});
</script>
