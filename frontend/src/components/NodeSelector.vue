<template>
  <div class="flex flex-col gap-4 p-5 bg-white rounded-lg border border-gray-200 h-full overflow-y-auto">
    <div class="flex items-center gap-2 border-b-2 border-gray-200 pb-3">
      <h3 class="m-0 text-base font-semibold text-gray-800">Available Nodes</h3>
      <span class="text-xs text-gray-400">({{ filteredNodes.length }})</span>
    </div>

    <div class="flex flex-col gap-2">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search nodes..."
        class="px-3 py-2 border border-gray-300 rounded text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100"
        @input="filterNodes"
      />

      <select v-model="typeFilter" class="px-3 py-2 border border-gray-300 rounded text-sm bg-white focus:outline-none focus:border-blue-500" @change="filterNodes">
        <option value="">All Types</option>
        <option value="attraction">🏛️ Attractions</option>
        <option value="transition">🛣️ Transitions</option>
      </select>
    </div>

    <div v-if="isLoading" class="text-center py-8 text-gray-400">
      <p>Loading nodes...</p>
    </div>

    <div v-else-if="filteredNodes.length === 0" class="text-center py-8 px-4 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300 text-gray-400">
      <p>No nodes found matching your criteria</p>
    </div>

    <div v-else class="grid grid-cols-1 gap-3">
      <NodeCard
        v-for="node in paginatedNodes"
        :key="node.id"
        :node="node"
        :selected="isNodeSelected(node.id)"
        show-actions
        @select="selectNode(node.id)"
        @deselect="deselectNode(node.id)"
      />
    </div>

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
import { ref, computed, onMounted } from 'vue';
import NodeCard from './NodeCard.vue';
import { nodeService, getNodeName, getNodeDescription } from '../services/node_service';
import type { Node } from '../services/node_service';

interface Props {
  selectedNodeIds?: string[];
  isLoading?: boolean;
}

interface Emits {
  (e: 'select', nodeId: string): void;
  (e: 'deselect', nodeId: string): void;
}

const props = withDefaults(defineProps<Props>(), {
  selectedNodeIds: () => [],
  isLoading: false,
});

const emit = defineEmits<Emits>();

const allNodes = ref<Node[]>([]);
const searchQuery = ref('');
const typeFilter = ref('');
const currentPage = ref(1);
const pageSize = ref(6);

onMounted(async () => {
  try {
    const result = await nodeService.listApprovedNodes({
      approved_only: true,
      limit: 100,
    });
    console.log('Fetched nodes:', result.nodes);
    allNodes.value = result.nodes;
  } catch (error) {
    console.error('Failed to load nodes:', error);
  }
});

const filteredNodes = computed(() => {
  let results = allNodes.value;

  // Filter by search query - search in node name and description via helper functions
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase();
    results = results.filter((node) => {
      const name = getNodeName(node);
      const desc = getNodeDescription(node);
      return name.toLowerCase().includes(query) || desc.toLowerCase().includes(query);
    });
  }

  // Filter by type
  if (typeFilter.value) {
    results = results.filter((node) => node.type === typeFilter.value);
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

function previousPage(): void {
  if (currentPage.value > 1) {
    currentPage.value--;
  }
}

function nextPage(): void {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
  }
}

function filterNodes(): void {
  currentPage.value = 1;
}
</script>
