<template>
  <div class="my-6 p-5 bg-gray-50 rounded-lg border border-gray-200">
    <div class="mb-4">
      <label for="comment-text" class="block mb-2 font-medium text-gray-800">Share your thoughts</label>
      <textarea
        id="comment-text"
        v-model="content"
        placeholder="What did you think about this travel plan?"
        rows="4"
        class="w-full px-3 py-3 border border-gray-300 rounded text-sm font-inherit focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100 resize-y"
        maxlength="1000"
      ></textarea>
      <div class="mt-1 text-right text-xs text-gray-500">{{ content.length }}/1000</div>
    </div>

    <div class="flex gap-3 justify-end">
      <button
        class="px-4 py-2 rounded text-sm font-medium cursor-pointer transition-all disabled:opacity-50 disabled:cursor-not-allowed bg-blue-600 text-white hover:bg-blue-700"
        :disabled="!content.trim() || isLoading"
        @click="submitComment"
      >
        {{ isLoading ? 'Posting...' : 'Post Comment' }}
      </button>
    </div>

    <div v-if="error" class="mt-3 p-3 bg-red-50 border-l-4 border-red-500 text-red-800 rounded text-sm">
      {{ error }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

interface Emits {
  (e: 'submit', content: string): Promise<void>;
}

const emit = defineEmits<Emits>();

const content = ref('');
const isLoading = ref(false);
const error = ref('');

async function submitComment(): Promise<void> {
  if (!content.value.trim()) return;

  isLoading.value = true;
  error.value = '';

  try {
    await emit('submit', content.value.trim());
    content.value = '';
  } catch (err: any) {
    error.value = err.message || 'Failed to post comment';
  } finally {
    isLoading.value = false;
  }
}
</script>
