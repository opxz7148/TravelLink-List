<template>
  <div class="flex flex-col gap-4 p-5 bg-white rounded-lg border border-gray-200 h-full overflow-y-auto">
    <div class="flex items-center gap-2 border-b-2 border-gray-200 pb-3">
      <h3 class="m-0 text-base font-semibold text-gray-800">Available Nodes</h3>
      <span class="text-xs text-gray-400">({{ filteredNodes.length }})</span>
    </div>

    <!-- Tab System -->
    <div v-if="showMyNodesTab" class="flex gap-2 border-b border-gray-200">
      <button
        @click="activeTab = 'attraction'"
        :class="[
          'px-4 py-2 font-medium text-sm transition-colors',
          activeTab === 'attraction'
            ? 'text-blue-600 border-b-2 border-blue-600'
            : 'text-gray-600 hover:text-gray-900'
        ]"
      >
        🏛️ Attractions
      </button>
      <button
        @click="activeTab = 'transition'"
        :class="[
          'px-4 py-2 font-medium text-sm transition-colors',
          activeTab === 'transition'
            ? 'text-blue-600 border-b-2 border-blue-600'
            : 'text-gray-600 hover:text-gray-900'
        ]"
      >
        🛣️ Transitions
      </button>
      <button
        @click="activeTab = 'my-nodes'"
        :class="[
          'px-4 py-2 font-medium text-sm transition-colors',
          activeTab === 'my-nodes'
            ? 'text-blue-600 border-b-2 border-blue-600'
            : 'text-gray-600 hover:text-gray-900'
        ]"
      >
        My Nodes <span v-if="userNodes.length > 0" class="ml-1 text-xs bg-orange-100 text-orange-700 px-2 py-0.5 rounded-full">{{ userNodes.length }}</span>
      </button>
    </div>

    <!-- Search -->
    <div class="flex flex-col gap-2">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search nodes..."
        class="px-3 py-2 border border-gray-300 rounded text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100"
        @input="filterNodes"
      />
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="text-center py-8 text-gray-400">
      <p>Loading nodes...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredNodes.length === 0" class="text-center py-8 px-4 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300 text-gray-400">
      <p v-if="activeTab === 'attraction'">No attractions found matching your criteria</p>
      <p v-else-if="activeTab === 'transition'">No transitions found matching your criteria</p>
      <p v-else-if="activeTab === 'my-nodes'">
        You haven't created any nodes yet. Create one using the button above to make it available here.
      </p>
    </div>

    <!-- Nodes Grid -->
    <div v-else class="grid grid-cols-1 gap-3">
      <NodeCard
        v-for="node in paginatedNodes"
        :key="node.id"
        :node="node"
        :selected="isNodeSelected(node.id)"
        :show-status="activeTab === 'my-nodes'"
        show-actions
        @select="selectNode(node.id)"
        @deselect="deselectNode(node.id)"
      />
    </div>

    <!-- Pagination -->
    <div v-if="filteredNodes.length > 0" class="flex justify-center items-center gap-3 pt-3 border-t border-gray-200">
      <button
        class="px-3 py-1 rounded text-sm font-medium bg-gray-600 text-white hover:bg-gray-700 transition-colors cursor-pointer border-none disabled:opacity-50 disabled:cursor-not-allowed"
        :disabled="currentPage === 1"
        @click="previousPage"
      >
        ← Previous
      </button>

      <span class="text-xs text-gray-400 whitespace-nowrap">
        Page {{ currentPage }} of {{ totalPages }}
      </span>

      <button
        class="px-3 py-1 rounded text-sm font-medium bg-gray-600 text-white hover:bg-gray-700 transition-colors cursor-pointer border-none disabled:opacity-50 disabled:cursor-not-allowed"
        :disabled="currentPage >= totalPages"
        @click="nextPage"
      >
        Next →
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import NodeCard from './NodeCard.vue';
import { nodeService, getNodeName, getNodeDescription } from '../services/node_service';
import type { Node } from '../services/node_service';
import { useAuthStore } from '../stores/auth_store';

interface Props {
  selectedNodeIds?: string[];
  isLoading?: boolean;
  showMyNodesTab?: boolean;
}

interface Emits {
  (e: 'select', nodeId: string): void;
  (e: 'deselect', nodeId: string): void;
}

const props = withDefaults(defineProps<Props>(), {
  selectedNodeIds: () => [],
  isLoading: false,
  showMyNodesTab: false,
});

const emit = defineEmits<Emits>();
const authStore = useAuthStore();

const allSystemNodes = ref<Node[]>([]);
const userNodes = ref<Node[]>([]);
const activeTab = ref<'attraction' | 'transition' | 'my-nodes'>('attraction');
const searchQuery = ref('');
const currentPage = ref(1);
const pageSize = ref(6);

// Refetch user nodes from API
async function refreshUserNodes(): Promise<void> {
  if (props.showMyNodesTab && authStore.isAuthenticated) {
    try {
      userNodes.value = await nodeService.getUserNodes();
    } catch (error) {
      console.warn('Failed to load user nodes:', error);
      userNodes.value = [];
    }
  }
}

onMounted(async () => {
  try {
    // Load system nodes (approved only)
    const systemResult = await nodeService.listApprovedNodes({
      approved_only: true,
      limit: 100,
    });
    allSystemNodes.value = systemResult.nodes;

    // Load user's draft nodes if authenticated
    await refreshUserNodes();
  } catch (error) {
    console.error('Failed to load nodes:', error);
  }
});

// Refetch user nodes whenever the "my-nodes" tab is clicked
watch(activeTab, async (newTab) => {
  if (newTab === 'my-nodes') {
    await refreshUserNodes();
  }
});

// Get nodes based on active tab
const tabNodes = computed(() => {
  if (activeTab.value === 'my-nodes') {
    return userNodes.value;
  }
  // Filter system nodes by type (attraction or transition)
  return allSystemNodes.value.filter((node) => node.type === activeTab.value);
});

const filteredNodes = computed(() => {
  let results = tabNodes.value;

  // Filter by search query - search in node name and description via helper functions
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase();
    results = results.filter((node) => {
      const name = getNodeName(node);
      const desc = getNodeDescription(node);
      return name.toLowerCase().includes(query) || desc.toLowerCase().includes(query);
    });
  }

  return results;
});

const paginatedNodes = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return filteredNodes.value.slice(start, end);
});

const totalPages = computed(() => {
  return Math.ceil(filteredNodes.value.length / pageSize.value);
});

function isNodeSelected(nodeId: string): boolean {
  return props.selectedNodeIds.includes(nodeId);
}

function selectNode(nodeId: string): void {
  if (!isNodeSelected(nodeId)) {
    emit('select', nodeId);
  }
}

function deselectNode(nodeId: string): void {
  if (isNodeSelected(nodeId)) {
    emit('deselect', nodeId);
  }
}

function filterNodes(): void {
  currentPage.value = 1;
}

function nextPage(): void {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
  }
}

function previousPage(): void {
  if (currentPage.value > 1) {
    currentPage.value--;
  }
}

// Expose refresh method so parent component can manually trigger refresh
defineExpose({
  refreshUserNodes,
});
</script>

