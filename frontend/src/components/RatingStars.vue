<template>
  <div class="flex items-center gap-2">
    <div v-if="mode === 'display'" class="flex items-center gap-1">
      <span v-for="i in 5" :key="i" :class="['text-lg transition-opacity duration-200', i <= Math.round(rating) ? 'opacity-100' : 'opacity-30']">
        ⭐
      </span>
      <span class="ml-2 font-semibold text-gray-800">{{ rating.toFixed(1) }}/5</span>
    </div>

    <div v-else class="flex gap-1">
      <button
        v-for="star in 5"
        :key="star"
        class="bg-none border-none cursor-pointer text-2xl px-1 transition-transform duration-200 hover:scale-125"
        :class="{ 'opacity-100': star <= hoveredRating || star <= modelValue }"
        @click="$emit('update:modelValue', star)"
        @mouseenter="hoveredRating = star"
        @mouseleave="hoveredRating = 0"
      >
        {{ star <= hoveredRating || star <= modelValue ? '⭐' : '☆' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

interface Props {
  rating?: number;
  modelValue?: number;
  mode?: 'display' | 'input';
  readonly?: boolean;
}

interface Emits {
  (e: 'update:modelValue', value: number): void;
}

withDefaults(defineProps<Props>(), {
  rating: 0,
  modelValue: 0,
  mode: 'display',
  readonly: false,
});

defineEmits<Emits>();

const hoveredRating = ref(0);
</script>
