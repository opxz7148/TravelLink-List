<template>
  <div
    :class="[
      'flex flex-col gap-3 p-4 border-2 rounded-lg bg-white transition-all cursor-pointer',
      node.type === 'attraction' ? 'border-l-4 border-l-amber-500 border-gray-200' : 'border-l-4 border-l-blue-500 border-gray-200',
      selected ? 'border-green-500 bg-green-50' : 'hover:shadow-md hover:border-blue-600',
    ]"
  >
    <div class="flex justify-between items-center">
      <span class="inline-flex items-center gap-1.5 px-2 py-1 bg-gray-100 rounded text-xs font-semibold text-gray-800">
        {{ node.type === 'attraction' ? '🏛️' : '🛣️' }}
        {{ node.type === 'attraction' ? 'Attraction' : 'Transition' }}
      </span>
      <span v-if="showSequence" class="inline-flex items-center justify-center w-7 h-7 bg-blue-600 text-white rounded-full font-bold text-sm">{{ sequence }}</span>
    </div>

    <div class="flex flex-col gap-2">
      <h3 v-if="node.type === 'attraction'" class="m-0 text-base font-semibold text-gray-800">
        {{ (node as any).name }}
      </h3>
      <p v-if="node.type === 'attraction'" class="m-0 text-xs text-gray-600">
        📍 {{ (node as any).location }}
      </p>
      <p v-if="node.type === 'attraction'" class="m-0 text-sm text-gray-700 leading-relaxed">
        {{ (node as any).description }}
      </p>

      <div v-if="node.type === 'transition'" class="flex flex-col gap-1 text-sm text-gray-600">
        <span>⏱️ {{ (node as any).duration_minutes }} min</span>
        <span class="text-gray-500 italic"> {{ (node as any).description }}</span>
      </div>
    </div>

    <div v-if="showActions" class="flex justify-end gap-2">
      <button
        v-if="!selected"
        class="px-3 py-1 rounded text-sm font-medium bg-blue-600 text-white hover:bg-blue-700 transition-colors cursor-pointer border-none"
        @click="$emit('select')"
      >
        Add
      </button>
      <button
        v-else
        class="px-3 py-1 rounded text-sm font-medium bg-gray-600 text-white hover:bg-gray-700 transition-colors cursor-pointer border-none"
        @click="$emit('deselect')"
      >
        Remove
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Node } from '../services/node_service';

interface Props {
  node: Node;
  selected?: boolean;
  sequence?: number;
  showSequence?: boolean;
  showActions?: boolean;
}

interface Emits {
  (e: 'select'): void;
  (e: 'deselect'): void;
}

withDefaults(defineProps<Props>(), {
  selected: false,
  sequence: 0,
  showSequence: false,
  showActions: true,
});

defineEmits<Emits>();
</script>
