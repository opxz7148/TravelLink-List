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
      <div class="flex items-center gap-2">
        <span v-if="showStatus" :class="[
          'inline-flex items-center px-2 py-1 rounded text-xs font-semibold',
          node.is_approved
            ? 'bg-green-100 text-green-700'
            : 'bg-orange-100 text-orange-700'
        ]">
          {{ node.is_approved ? '✓ Published' : '○ Draft' }}
        </span>
        <span v-if="showSequence" class="inline-flex items-center justify-center w-7 h-7 bg-blue-600 text-white rounded-full font-bold text-sm">{{ sequence }}</span>
      </div>
    </div>

    <div class="flex flex-col gap-2">
      <!-- Attraction Details -->
      <template v-if="node.type === 'attraction' && node.attraction">
        <h3 class="m-0 text-base font-semibold text-gray-800">
          {{ node.attraction.name }}
        </h3>
        <p class="m-0 text-xs text-gray-600">
          📍 {{ node.attraction.location }}
        </p>
        <p v-if="node.attraction.category" class="m-0 text-xs text-gray-500 capitalize">
          Category: {{ node.attraction.category.replace('_', ' ') }}
        </p>
        <p v-if="node.attraction.description" class="m-0 text-sm text-gray-700 leading-relaxed">
          {{ node.attraction.description }}
        </p>
        <p v-if="node.attraction.contact_info" class="m-0 text-xs text-gray-500">
          📞 {{ node.attraction.contact_info }}
        </p>
        <p v-if="node.attraction.hours_of_operation" class="m-0 text-xs text-gray-500">
          ⏰ {{ node.attraction.hours_of_operation }}
        </p>
      </template>

      <!-- Transition Details -->
      <template v-if="node.type === 'transition' && node.transition">
        <h3 class="m-0 text-base font-semibold text-gray-800">
          {{ node.transition.title }}
        </h3>
        <p class="m-0 text-xs text-gray-600 capitalize">
          Mode: {{ node.transition.mode.replace('_', ' ') }}
        </p>
        <p v-if="node.transition.description" class="m-0 text-sm text-gray-700 leading-relaxed">
          {{ node.transition.description }}
        </p>
        <p v-if="node.transition.hours_of_operation" class="m-0 text-xs text-gray-500">
          ⏰ {{ node.transition.hours_of_operation }}
        </p>
        <p v-if="node.transition.route_notes" class="m-0 text-xs text-gray-500 italic">
          📍 {{ node.transition.route_notes }}
        </p>
      </template>
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
  showStatus?: boolean;
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
  showStatus: false,
});

defineEmits<Emits>();
</script>
