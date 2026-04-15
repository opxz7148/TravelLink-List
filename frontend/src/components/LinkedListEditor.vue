<template>
  <div class="flex flex-col gap-4 p-5 bg-gray-50 rounded-lg">
    <div class="flex justify-between items-center border-b-2 border-gray-300 pb-3">
      <h3 class="m-0 text-base font-semibold text-gray-800">Your Itinerary ({{ selectedNodes.length }} stops)</h3>
      <span v-if="selectedNodes.length > 0" class="text-xs text-gray-400 italic">
        Drag to reorder
      </span>
    </div>

    <div v-if="selectedNodes.length === 0" class="text-center py-8 px-4 text-gray-400 bg-white rounded-lg border-2 border-dashed border-gray-300">
      <p>👈 Select nodes from the right panel to build your itinerary</p>
    </div>

    <div v-else class="flex flex-col" @dragover.prevent @drop.prevent>
      <div
        v-for="(nodeId, index) in selectedNodes"
        :key="nodeId"
        class="relative p-4 bg-white border border-gray-300 cursor-move select-none transition-all hover:bg-gray-50 hover:shadow-md"
        :class="{ 'opacity-50': index === draggingIndex }"
        draggable="true"
        @dragstart="startDrag(index)"
        @dragenter="dragEnter(index)"
        @dragend="endDrag"
      >
        <div class="flex items-center gap-3 mb-2">
          <span class="inline-flex items-center justify-center w-8 h-8 bg-blue-600 text-white rounded-full font-bold text-sm flex-shrink-0">{{ index + 1 }}</span>
          <span class="font-semibold text-gray-800 flex-1">{{ getNodeName(nodeId) }}</span>
        </div>

        <p v-if="getNodeDescription(nodeId)" class="m-0 ml-11 pt-0 text-xs text-gray-600 leading-relaxed">
          {{ getNodeDescription(nodeId) }}
        </p>

        <button
          class="absolute top-3 right-3 bg-transparent border-0 text-red-600 cursor-pointer text-lg px-1 rounded hover:bg-red-50 transition-colors"
          title="Remove from itinerary"
          @click="removeNode(index)"
        >
          ✕
        </button>

        <div v-if="index < selectedNodes.length - 1" class="text-center text-gray-400 text-3xl py-0.5 select-none">
          ↓
        </div>
      </div>
    </div>

    <div v-if="selectedNodes.length > 0" class="grid grid-cols-3 gap-3 p-3 bg-white rounded border border-gray-200">
      <div class="flex flex-col gap-1 text-center">
        <span class="text-xs text-gray-400 font-medium">Total Stops:</span>
        <span class="text-lg font-bold text-blue-600">{{ selectedNodes.length }}</span>
      </div>
      <div class="flex flex-col gap-1 text-center">
        <span class="text-xs text-gray-400 font-medium">Attractions:</span>
        <span class="text-lg font-bold text-blue-600">{{ attractionCount }}</span>
      </div>
      <div class="flex flex-col gap-1 text-center">
        <span class="text-xs text-gray-400 font-medium">Transitions:</span>
        <span class="text-lg font-bold text-blue-600">{{ transitionCount }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { getNodeName as getNodeNameHelper, getNodeDescription as getNodeDescriptionHelper } from '../services/node_service';
import type { Node } from '../services/node_service';

interface Props {
  selectedNodes: string[];
  availableNodes: Node[];
}

interface Emits {
  (e: 'reorder', nodes: string[]): void;
  (e: 'remove', index: number): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const draggingIndex = ref<number | null>(null);
const dragOverIndex = ref<number | null>(null);

const attractionCount = computed(() => {
  return props.selectedNodes.filter((nodeId) => {
    const node = props.availableNodes.find((n) => n.id === nodeId);
    return node?.type === 'attraction';
  }).length;
});

const transitionCount = computed(() => {
  return props.selectedNodes.filter((nodeId) => {
    const node = props.availableNodes.find((n) => n.id === nodeId);
    return node?.type === 'transition';
  }).length;
});

function getNode(nodeId: string): Node | undefined {
  return props.availableNodes.find((n) => n.id === nodeId);
}

/**
 * Get display name for a node - uses embedded detail structure
 */
function getNodeName(nodeId: string): string {
  const node = getNode(nodeId);
  console.log('Getting name for nodeId:', nodeId, 'Found node:', node);
  if (!node) return 'Unknown';
  return getNodeNameHelper(node);
}

/**
 * Get display description for a node - uses embedded detail structure
 */
function getNodeDescription(nodeId: string): string {
  const node = getNode(nodeId);
  if (!node) return '';
  return getNodeDescriptionHelper(node);
}

function startDrag(index: number): void {
  draggingIndex.value = index;
}

function dragEnter(index: number): void {
  dragOverIndex.value = index;
}

function endDrag(): void {
  if (draggingIndex.value !== null && dragOverIndex.value !== null && draggingIndex.value !== dragOverIndex.value) {
    // Reorder nodes
    const newOrder = [...props.selectedNodes];
    const draggedNode = newOrder[draggingIndex.value];
    newOrder.splice(draggingIndex.value, 1);
    newOrder.splice(dragOverIndex.value, 0, draggedNode as string);
    emit('reorder', newOrder);
  }

  draggingIndex.value = null;
  dragOverIndex.value = null;
}

function removeNode(index: number): void {
  emit('remove', index);
}
</script>
