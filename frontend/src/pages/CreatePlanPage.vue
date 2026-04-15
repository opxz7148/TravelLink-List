<template>
  <div class="w-full max-w-400 mx-auto px-4 sm:px-6 lg:px-8 py-8 md:py-12">
    <!-- Page Header -->
    <div class="text-center mb-8">
      <h1 class="text-4xl font-bold text-gray-900 mb-2">{{ isEditMode ? 'Edit Travel Plan' : 'Create a New Travel Plan' }}</h1>
      <p class="text-gray-600">{{ isEditMode ? 'Update your draft plan' : 'Design your perfect itinerary with our curated nodes' }}</p>
    </div>

    <!-- Loading existing plan -->
    <div v-if="isEditMode && !currentPlan && loadingNodes" class="min-h-[85vh] grid place-items-center">
      <div class="w-full max-w-md bg-white p-8 rounded-lg border border-gray-200 shadow-sm text-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
        <p class="text-gray-600">Loading your draft plan...</p>
      </div>
    </div>

    <!-- Step 1: Plan Basics -->
    <div v-else-if="!currentPlan" class="min-h-[85vh] grid place-items-center">
      <div class="w-full max-w-md bg-white p-10 rounded-lg border border-gray-200 shadow-sm">
        <h2 class="text-xl font-bold text-gray-900 mb-8">Step 1: Plan Basics</h2>

        <div class="flex flex-col gap-3 mb-7">
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

        <div class="flex flex-col gap-3 mb-8">
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
          @click="proceedToNodeSelection"
          :disabled="!planTitle.trim() || !planDestination.trim() || loadingNodes"
          class="w-full px-6 py-3.5 bg-blue-600 text-white rounded-md text-base font-medium cursor-pointer hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
        >
          {{ loadingNodes ? 'Loading nodes...' : 'Next: Select Nodes' }}
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
          <p class="text-sm text-gray-500 mt-1">Status: <span class="font-semibold text-orange-600">Draft</span></p>
        </div>

        <button @click="resetPlan" class="px-4 py-2 bg-gray-200 text-gray-900 rounded-md text-sm font-medium hover:bg-gray-300 transition-all">
          {{ isEditMode ? 'Cancel' : 'Start Over' }}
        </button>
      </div>

      <!-- Builder Grid: Itinerary Editor + Nodes Side by Side -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Left: Editor -->
        <div class="lg:col-span-1">
          <LinkedListEditor
            :selected-nodes="selectedNodes"
            :available-nodes="availableNodes"
            @reorder="reorderNodes"
            @remove="removeNodeAtIndex"
          />
        </div>

        <!-- Right: Node Selector (Attractions and Transitions Side by Side) -->
        <div class="lg:col-span-2">
          <div class="space-y-4">
            <!-- Create Node Button -->
            <button
              @click="showNodeCreationModal = true"
              class="w-full px-4 py-3 bg-blue-50 text-blue-600 border-2 border-dashed border-blue-300 rounded-lg font-medium hover:bg-blue-100 transition-all"
            >
              + Create Your Own Node
            </button>

            <NodeSelector
              ref="nodeSelectorRef"
              :selected-node-ids="selectedNodes"
              :is-loading="loadingNodes"
              :show-my-nodes-tab="true"
              @select="addNode"
              @deselect="removeNode"
            />
          </div>
        </div>
      </div>

      <!-- Node Creation Modal -->
      <NodeCreationModal
        :is-open="showNodeCreationModal"
        @close="showNodeCreationModal = false"
        @created="onNodeCreated"
      />

      <!-- Node Details Editor -->
      <div v-if="selectedNodes.length > 0" class="bg-white p-6 rounded-lg border border-gray-200 shadow-sm">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">📝 Customize Node Information</h3>
        <p class="text-sm text-gray-600 mb-5">Add plan-specific details for each node (optional)</p>

        <!-- Tabs/Accordion for each selected node -->
        <div class="space-y-3">
          <div v-for="(nodeId, index) in selectedNodes" :key="nodeId" class="border border-gray-200 rounded-lg overflow-hidden">
            <!-- Node Tab Header -->
            <button
              @click="startEditingNode(editingNodeId === nodeId ? null : nodeId)"
              class="w-full px-4 py-3 bg-gray-50 hover:bg-gray-100 flex items-center justify-between transition-colors"
            >
              <div class="flex items-center gap-3">
                <span
                  class="grid w-6 h-6 rounded-full text-white text-xs font-semibold place-items-center"
                  :class="getNodeType(nodeId) === 'attraction' ? 'bg-blue-500' : 'bg-green-500'"
                >
                  {{ index + 1 }}
                </span>
                <div class="text-left">
                  <p class="font-semibold text-gray-900">{{ getNodeName(nodeId) }}</p>
                  <p class="text-xs text-gray-600">{{ getNodeType(nodeId) }}</p>
                </div>
              </div>
              <svg
                :class="editingNodeId === nodeId ? 'rotate-180' : ''"
                class="w-5 h-5 text-gray-500 transition-transform"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3" />
              </svg>
            </button>

            <!-- Node Details Form (Expanded) -->
            <div v-if="editingNodeId === nodeId" class="px-4 py-4 bg-white border-t border-gray-200 space-y-4">
              <!-- Description -->
              <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-gray-700">Plan-Specific Notes (optional)</label>
                <textarea
                  :value="nodeDetails.get(nodeId)?.description || ''"
                  @input="updateNodeDescription(nodeId, ($event.target as HTMLTextAreaElement).value)"
                  placeholder="e.g., 'Try the house special pasta', 'Scenic metro route'"
                  maxlength="500"
                  rows="3"
                  class="px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100 font-inherit"
                />
                <p class="text-xs text-gray-500">
                  {{ (nodeDetails.get(nodeId)?.description || '').length }}/500 characters
                </p>
              </div>

              <!-- Estimated Price -->
              <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-gray-700">Estimated Cost (optional)</label>
                <div class="flex items-center gap-2">
                  <input
                    type="number"
                    :value="nodeDetails.get(nodeId)?.estimated_price_cents ? nodeDetails.get(nodeId)!.estimated_price_cents! / 100 : ''"
                    @input="updateNodePrice(nodeId, ($event.target as HTMLInputElement).value ? Math.round(parseFloat(($event.target as HTMLInputElement).value) * 100) : null)"
                    placeholder="0.00"
                    min="0"
                    step="0.01"
                    class="px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100 font-inherit flex-1"
                  />
                  <span class="text-sm font-semibold text-gray-600">USD</span>
                </div>
                <p class="text-xs text-gray-500">Cost of this leg in your travel plan</p>
              </div>

              <!-- Duration -->
              <div class="flex flex-col gap-2">
                <label class="text-sm font-semibold text-gray-700">Duration (optional)</label>
                <div class="flex items-center gap-2">
                  <input
                    type="number"
                    :value="nodeDetails.get(nodeId)?.duration_minutes || ''"
                    @input="updateNodeDuration(nodeId, ($event.target as HTMLInputElement).value ? parseInt(($event.target as HTMLInputElement).value) : null)"
                    placeholder="e.g., 90"
                    min="1"
                    class="px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100 font-inherit flex-1"
                  />
                  <span class="text-sm font-semibold text-gray-600">minutes</span>
                </div>
                <p class="text-xs text-gray-500">How long to spend at this location or travel here</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Save/Publish Actions Section -->
      <div v-if="selectedNodes.length > 0" class="p-6 bg-blue-50 rounded-lg border-2 border-blue-200">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">{{ isEditMode ? 'Update Your Plan?' : 'Ready to Save Your Plan?' }}</h3>
        
        <div class="flex flex-col gap-3">
          <!-- Save as Draft Button -->
          <button
            @click="saveDraft"
            :disabled="loading"
            class="w-full px-6 py-3 bg-blue-600 text-white rounded-md text-base font-medium hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
          >
            {{ loading ? (isEditMode ? 'Updating...' : 'Saving...') : (isEditMode ? 'Update Draft' : 'Save as Draft') }}
          </button>

          <!-- Publish Button (traveller and admin, only for new plans) -->
          <button
            v-if="!isEditMode && authStore.isTraveller"
            @click="publishPlan"
            :disabled="loading"
            class="w-full px-6 py-3 bg-green-600 text-white rounded-md text-base font-medium hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
          >
            {{ loading ? 'Publishing...' : 'Create & Publish Plan' }}
          </button>

          <!-- Info text for regular users -->
          <p v-if="!authStore.isTraveller" class="text-sm text-orange-700 bg-orange-50 px-3 py-2 rounded border border-orange-200 mt-2">
            💡 You can save and edit draft plans. To publish plans and make them visible to the community, please submit a promotion request to become a Traveller.
          </p>
          <p v-else class="text-sm text-gray-600 mt-2">
            💡 {{ isEditMode ? 'Save your changes to the draft.' : 'Save as Draft to keep editing later. Publish to make your plan visible to the community.' }}
          </p>
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
import { useRouter, useRoute } from 'vue-router';
import LinkedListEditor from '../components/LinkedListEditor.vue';
import NodeSelector from '../components/NodeSelector.vue';
import NodeCreationModal from '../components/NodeCreationModal.vue';
import { planService } from '../services/plan_service';
import type { TravelPlan, PlanDetail, NodeDetailForPlan } from '../services/plan_service';
import { nodeService, getNodeName as getNodeNameHelper } from '../services/node_service';
import type { Node } from '../services/node_service';
import { useAuthStore } from '../stores/auth_store';
import { useUiStore } from '../stores/ui_store';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();
const uiStore = useUiStore();

const planTitle = ref('');
const planDestination = ref('');
const currentPlan = ref<PlanDetail | null>(null);
const selectedNodes = ref<string[]>([]);
const nodeDetails = ref<Map<string, NodeDetailForPlan>>(new Map()); // Track plan-specific details for each node
const availableNodes = ref<Node[]>([]);
const editingNodeId = ref<string | null>(null); // Track which node is being edited
const errorMessage = ref('');
const loading = ref(false);
const loadingNodes = ref(false);
const isEditMode = ref(false);
const editingPlanId = ref<string | null>(null);
const planNodeToNodeIdMap = ref<Map<string, string>>(new Map()); // Map PlanNode IDs to actual Node IDs (for edit mode)
const showNodeCreationModal = ref(false);
const nodeSelectorRef = ref<InstanceType<typeof NodeSelector> | null>(null);

/**
 * Proceed to node selection by loading available nodes
 */
async function proceedToNodeSelection(): Promise<void> {
  if (!planTitle.value.trim() || !planDestination.value.trim()) {
    errorMessage.value = 'Please fill in all required fields';
    return;
  }

  try {
    loadingNodes.value = true;
    errorMessage.value = '';

    // Load available nodes
    const { nodes } = await nodeService.listApprovedNodes({
      approved_only: true,
    });

    availableNodes.value = nodes;
    
    // Initialize plan as a draft object (not yet created on server for new plans)
    // For existing plans, use the loaded data
    if (!currentPlan.value || !isEditMode.value) {
      currentPlan.value = {
        id: '', // Will be assigned when created
        title: planTitle.value.trim(),
        destination: planDestination.value.trim(),
        author_id: authStore.user?.id || '',
        status: 'draft',
        rating_count: 0,
        rating_sum: 0,
        comment_count: 0,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        nodes: [],
      };
    }
  } catch (error: any) {
    console.error('Failed to load nodes:', error);
    errorMessage.value = 'Failed to load available nodes';
    currentPlan.value = null;
  } finally {
    loadingNodes.value = false;
  }
}

/**
 * Handle when a new node is created in the modal
 * Reload the available nodes to show the newly created node in "My Nodes" tab
 */
async function onNodeCreated(): Promise<void> {
  try {
    // Reload available nodes and user's draft nodes
    const [approvedNodes, userNodes] = await Promise.all([
      nodeService.listApprovedNodes({ approved_only: true }),
      nodeService.getUserNodes(),
    ]);
    
    // Combine both lists for display
    availableNodes.value = [...approvedNodes.nodes, ...userNodes];
    
    // Also refresh the NodeSelector's user nodes list
    if (nodeSelectorRef.value) {
      await nodeSelectorRef.value.refreshUserNodes();
    }
    
    uiStore.showSuccess('Node created successfully! Find it in the "My Nodes" tab to add it to your plan.');
  } catch (error: any) {
    console.error('Failed to reload nodes:', error);
    errorMessage.value = 'Node created but failed to reload the node list. Please refresh the page.';
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
  
  // Initialize empty details for this node (user can fill in optional fields)
  if (!nodeDetails.value.has(nodeId)) {
    nodeDetails.value.set(nodeId, {
      node_id: nodeId,
      // description, estimated_price_cents, duration_minutes are optional
    });
  }
  
  errorMessage.value = '';
}

function removeNode(nodeId: string): void {
  const index = selectedNodes.value.indexOf(nodeId);
  if (index !== -1) {
    selectedNodes.value.splice(index, 1);
    nodeDetails.value.delete(nodeId);
    // Clear editing state if the removed node was being edited
    if (editingNodeId.value === nodeId) {
      editingNodeId.value = null;
    }
  }
}

function removeNodeAtIndex(index: number): void {
  const nodeId = selectedNodes.value[index];
  selectedNodes.value.splice(index, 1);
  if (nodeId) {
    nodeDetails.value.delete(nodeId);
    // Clear editing state if the removed node was being edited
    if (editingNodeId.value === nodeId) {
      editingNodeId.value = null;
    }
  }
}

function reorderNodes(newOrder: string[]): void {
  selectedNodes.value = newOrder;
}

/**
 * Start editing details for a specific node (toggle expand/collapse)
 */
function startEditingNode(nodeId: string | null): void {
  editingNodeId.value = nodeId;
}

/**
 * Update description for a node
 */
function updateNodeDescription(nodeId: string, description: string | null): void {
  const detail = nodeDetails.value.get(nodeId) || { node_id: nodeId };
  detail.description = description && description.trim() ? description.trim() : undefined;
  nodeDetails.value.set(nodeId, detail);
}

/**
 * Update estimated price for a node (in cents)
 */
function updateNodePrice(nodeId: string, priceCents: number | null): void {
  const detail = nodeDetails.value.get(nodeId) || { node_id: nodeId };
  detail.estimated_price_cents = priceCents && priceCents > 0 ? priceCents : undefined;
  nodeDetails.value.set(nodeId, detail);
}

/**
 * Update duration for a node (in minutes)
 */
function updateNodeDuration(nodeId: string, durationMinutes: number | null): void {
  const detail = nodeDetails.value.get(nodeId) || { node_id: nodeId };
  detail.duration_minutes = durationMinutes && durationMinutes > 0 ? durationMinutes : undefined;
  nodeDetails.value.set(nodeId, detail);
}

/**
 * Get node name by ID using the node service helper
 * In edit mode, takes name from backend response (details.name or details.title)
 */
function getNodeName(nodeId: string): string {
  // First check if this node is in the existing plan (edit mode)
  if (isEditMode.value && currentPlan.value?.nodes) {
    const existingNode = currentPlan.value.nodes.find((n) => n.id === nodeId);
    if (existingNode?.details) {
      if (existingNode.type === 'attraction' && (existingNode.details as any).name) {
        return (existingNode.details as any).name;
      } else if (existingNode.type === 'transition' && (existingNode.details as any).title) {
        return (existingNode.details as any).title;
      }
    }
  }
  // Otherwise, look up from availableNodes
  const node = availableNodes.value.find((n) => n.id === nodeId);
  return node ? getNodeNameHelper(node) : 'Unknown Node';
}

/**
 * Get node type by ID
 */
function getNodeType(nodeId: string): string {
  const node = availableNodes.value.find((n) => n.id === nodeId);
  return node?.type || 'unknown';
}

function resetPlan(): void {
  if (isEditMode.value) {
    // In edit mode, go back to my plans
    planNodeToNodeIdMap.value.clear();
    router.push('/my-plans');
  } else {
    // In create mode, reset the form
    currentPlan.value = null;
    selectedNodes.value = [];
    nodeDetails.value.clear();
    planNodeToNodeIdMap.value.clear();
    editingNodeId.value = null;
    planTitle.value = '';
    planDestination.value = '';
    errorMessage.value = '';
  }
}

/**
 * Save plan and nodes, keeping it as draft (or update existing draft)
 */
async function saveDraft(): Promise<void> {
  if (!currentPlan.value || selectedNodes.value.length === 0) {
    errorMessage.value = 'Please add at least one node to your plan before saving';
    return;
  }

  try {
    loading.value = true;
    errorMessage.value = '';

    if (isEditMode.value && editingPlanId.value) {
      // EDIT MODE: Use atomic editPlan endpoint to update plan + replace nodes
      
      // Build NodeDetailForPlan array from selected nodes
      // Map PlanNode IDs back to actual Node IDs using the stored mapping
      const nodesToUpdate: NodeDetailForPlan[] = selectedNodes.value.map((planNodeId) => {
        const detail = nodeDetails.value.get(planNodeId);
        // Use the mapping to get the real Node ID, fallback to planNodeId if not found
        const realNodeId = planNodeToNodeIdMap.value.get(planNodeId) || planNodeId;
        return {
          node_id: realNodeId,
          description: detail?.description,
          estimated_price_cents: detail?.estimated_price_cents,
          duration_minutes: detail?.duration_minutes,
        };
      });

      // Call atomic editPlan endpoint
      const updatedPlan = await planService.editPlan(
        editingPlanId.value,
        planTitle.value.trim(),
        planDestination.value.trim(),
        currentPlan.value?.description || '',
        nodesToUpdate
      );

      currentPlan.value = updatedPlan;
      uiStore.showSuccess('Plan updated successfully!');
      
      // Redirect to my plans
      setTimeout(() => {
        router.push('/my-plans');
      }, 1500);
    } else {
      // CREATE MODE: Create new draft plan with nodes
      
      // Build NodeDetailForPlan array from selected nodes
      const nodesToAdd: NodeDetailForPlan[] = selectedNodes.value.map((nodeId) => {
        const detail = nodeDetails.value.get(nodeId);
        return {
          node_id: nodeId,
          description: detail?.description,
          estimated_price_cents: detail?.estimated_price_cents,
          duration_minutes: detail?.duration_minutes,
        };
      });

      // Create plan with nodes in one request
      const createdPlan = await planService.createPlanWithNodes(
        planTitle.value.trim(),
        planDestination.value.trim(),
        nodesToAdd,
        'draft'
      );

      currentPlan.value = createdPlan;
      editingPlanId.value = createdPlan.id;

      // Show success
      uiStore.showSuccess('Plan saved as draft! You can edit it later or publish it whenever you\'re ready.');
      
      // Redirect to my plans after a short delay
      setTimeout(() => {
        router.push('/my-plans');
      }, 1500);
    }
  } catch (error: any) {
    console.error('Failed to save plan:', error);
    errorMessage.value = error.message || 'Failed to save plan';
    uiStore.showError(errorMessage.value);
  } finally {
    loading.value = false;
  }
}

/**
 * Create and publish plan (traveller only)
 */
async function publishPlan(): Promise<void> {
  if (!currentPlan.value || selectedNodes.value.length === 0) {
    errorMessage.value = 'Please add at least one node to your plan before publishing';
    return;
  }

  try {
    loading.value = true;
    errorMessage.value = '';

    // Build NodeDetailForPlan array from selected nodes
    const nodesToAdd: NodeDetailForPlan[] = selectedNodes.value.map((nodeId) => {
      const detail = nodeDetails.value.get(nodeId);
      return {
        node_id: nodeId,
        description: detail?.description,
        estimated_price_cents: detail?.estimated_price_cents,
        duration_minutes: detail?.duration_minutes,
      };
    });

    // Create plan with nodes and publish in one request
    const createdPlan = await planService.createPlanWithNodes(
      planTitle.value.trim(),
      planDestination.value.trim(),
      nodesToAdd,
      'published'
    );

    // Redirect to plan view
    router.push(`/plans/${createdPlan.id}`);
  } catch (error: any) {
    console.error('Failed to publish plan:', error);
    errorMessage.value = error.message || 'Failed to publish plan';
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  // Ensure user is authenticated
  if (!authStore.isAuthenticated) {
    router.push('/login');
    return;
  }

  // Prevent admin users from accessing create plan
  if (authStore.userRole === 'admin') {
    router.push('/my-plans');
    return;
  }

  // Check if we're in edit mode (loading existing draft plan)
  const editPlanId = route.query.edit as string | undefined;
  if (editPlanId) {
    try {
      isEditMode.value = true;
      editingPlanId.value = editPlanId;
      loadingNodes.value = true;

      // Load the existing plan and available nodes in parallel
      const [existingPlan, { nodes }] = await Promise.all([
        planService.getPlanDetail(editPlanId),
        nodeService.listApprovedNodes({ approved_only: true }),
      ]);

      // Verify the plan is owned by the current user
      if (existingPlan.author_id !== authStore.user?.id) {
        errorMessage.value = 'You can only edit your own plans.';
        loadingNodes.value = false;
        return;
      }

      // Pre-fill form with existing plan data
      planTitle.value = existingPlan.title;
      planDestination.value = existingPlan.destination;
      
      // In edit mode, we have enriched nodes already - use them for display
      if (existingPlan.nodes && existingPlan.nodes.length > 0) {
        // Convert enriched nodes to Node-like objects for LinkedListEditor
        const enrichedNodes = existingPlan.nodes.map((pn: any) => {
          const nodeId = pn.id; // This is PlanNode ID
          
          // Find matching approved node to get the real Node ID
          const detailsKey = pn.type === 'attraction' ? pn.details?.name : pn.details?.title;
          const realNode = nodes.find((n: any) => {
            const name = n.type === 'attraction' ? n.attraction?.name : n.transition?.title;
            return name === detailsKey && n.type === pn.type;
          });
          
          // Store the mapping: PlanNode ID -> Real Node ID
          if (realNode) {
            planNodeToNodeIdMap.value.set(nodeId, realNode.id);
          }
          
          return {
            id: nodeId, // Keep PlanNode ID for selectedNodes
            type: pn.type,
            created_by: '',
            is_approved: true,
            created_at: '',
            attraction: pn.type === 'attraction' && pn.details ? {
              name: pn.details.name || 'Unknown',
              description: pn.details.description || '',
              location: pn.details.location || '',
              category: pn.details.category || '',
              node_id: nodeId,
              created_by: '',
              created_at: '',
            } : null,
            transition: pn.type === 'transition' && pn.details ? {
              title: pn.details.title || 'Unknown',
              description: pn.details.description || '',
              node_id: nodeId,
              created_by: '',
              created_at: '',
            } : null,
          } as any;
        });
        availableNodes.value = enrichedNodes;
      }
      
      // Also include newly approved nodes for adding new nodes
      availableNodes.value = [...availableNodes.value, ...nodes];

      // Pre-populate selected nodes and their details
      selectedNodes.value = existingPlan.nodes?.map((n) => n.id) || [];
      existingPlan.nodes?.forEach((node) => {
        const detail: NodeDetailForPlan = {
          node_id: node.id,
        };
        if (node.description) detail.description = node.description;
        if (node.estimated_price_cents) detail.estimated_price_cents = node.estimated_price_cents;
        if (node.duration_minutes) detail.duration_minutes = node.duration_minutes;
        nodeDetails.value.set(node.id, detail);
      });

      // Initialize the plan object for UI
      currentPlan.value = existingPlan;
    } catch (error: any) {
      console.error('Failed to load draft plan:', error);
      errorMessage.value = error.message || 'Failed to load your draft plan';
      isEditMode.value = false;
      editingPlanId.value = null;
      planNodeToNodeIdMap.value.clear();
    } finally {
      loadingNodes.value = false;
    }
  }
});
</script>
