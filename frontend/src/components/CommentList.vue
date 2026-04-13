<template>
  <div class="my-6">
    <div v-if="comments.length === 0" class="text-center py-8 text-gray-400 bg-gray-50 rounded-lg">
      <p>No comments yet. Be the first to share your thoughts!</p>
    </div>

    <div v-else class="flex flex-col gap-3">
      <div v-for="comment in comments" :key="comment.id" class="p-4 bg-white border border-gray-200 rounded-lg transition-all hover:shadow-md">
        <div class="flex justify-between items-start mb-3 gap-2">
          <div class="flex items-center gap-3">
            <span class="font-semibold text-gray-800">{{ comment.author?.username || 'Anonymous' }}</span>
            <span class="text-xs text-gray-400">{{ formatDate(comment.created_at) }}</span>
          </div>

          <div v-if="canEditComment(comment)" class="flex gap-2">
            <button class="bg-none border-none text-blue-600 cursor-pointer text-xs underline hover:text-blue-700 p-0" @click="editComment(comment)">Edit</button>
            <button class="bg-none border-none text-red-600 cursor-pointer text-xs underline hover:text-red-700 p-0" @click="deleteComment(comment.id)">Delete</button>
          </div>
        </div>

        <div v-if="editingId === comment.id" class="mt-3 p-3 bg-gray-50 rounded">
          <textarea v-model="editText" class="w-full px-2 py-2 border border-gray-300 rounded text-sm font-inherit focus:outline-none focus:border-blue-500 resize-y" rows="3"></textarea>
          <div class="flex gap-2 mt-2">
            <button class="px-3 py-1 rounded text-sm font-medium bg-blue-600 text-white hover:bg-blue-700 transition-colors cursor-pointer border-none" @click="saveEdit(comment.id)">Save</button>
            <button class="px-3 py-1 rounded text-sm font-medium bg-gray-600 text-white hover:bg-gray-700 transition-colors cursor-pointer border-none" @click="cancelEdit">Cancel</button>
          </div>
        </div>

        <p v-else class="text-gray-800 leading-relaxed m-0 break-words">{{ comment.text }}</p>
      </div>
    </div>

    <div v-if="isLoading" class="text-center py-4 text-gray-400 text-sm">
      <p>Loading more comments...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Comment } from '../services/comment_service';
import { useAuthStore } from '../stores/auth_store';

interface Props {
  comments: Comment[];
  isLoading?: boolean;
  currentUserId?: string;
}

interface Emits {
  (e: 'edit', comment: Comment, text: string): Promise<void>;
  (e: 'delete', commentId: string): Promise<void>;
}

const props = withDefaults(defineProps<Props>(), {
  isLoading: false,
  currentUserId: '',
});

const emit = defineEmits<Emits>();

const authStore = useAuthStore();
const editingId = ref('');
const editText = ref('');

const canEditComment = computed(() => {
  return (comment: Comment) => {
    return authStore.isAuthenticated && (authStore.user?.id === comment.author_id || authStore.isAdmin);
  };
});

function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

function editComment(comment: Comment): void {
  editingId.value = comment.id;
  editText.value = comment.text;
}

function cancelEdit(): void {
  editingId.value = '';
  editText.value = '';
}

async function saveEdit(commentId: string): Promise<void> {
  const comment = props.comments.find((c) => c.id === commentId);
  if (!comment) return;

  try {
    await emit('edit', comment, editText.value);
    editingId.value = '';
    editText.value = '';
  } catch (error) {
    console.error('Failed to edit comment:', error);
  }
}

async function deleteComment(commentId: string): Promise<void> {
  if (!confirm('Delete this comment?')) return;

  try {
    await emit('delete', commentId);
  } catch (error) {
    console.error('Failed to delete comment:', error);
  }
}
</script>
