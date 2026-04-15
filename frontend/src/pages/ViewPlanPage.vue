<template>
  <div class="w-full max-w-[1400px] mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <!-- Plan Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 p-5 bg-white rounded-lg border border-gray-200 mb-12">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 mb-2">{{ plan?.title }}</h1>
        <p class="text-gray-600 mb-1">📍 {{ plan?.destination }}</p>
        <p class="text-sm text-gray-500">Status: {{ plan?.status }}</p>
      </div>
      <div class="flex gap-3 w-full sm:w-auto">
        <button v-if="isAuthor" @click="editPlan" class="flex-1 sm:flex-none px-4 py-2 bg-blue-600 text-white rounded-md text-sm font-medium hover:bg-blue-700 transition-all">
          Edit
        </button>
        <button v-if="isAdmin" @click="deletePlan" class="flex-1 sm:flex-none px-4 py-2 bg-red-600 text-white rounded-md text-sm font-medium hover:bg-red-700 transition-all">
          Delete
        </button>
      </div>
    </div>

    <!-- Nodes List -->
    <div v-if="plan?.nodes && plan.nodes.length > 0" class="space-y-4 mb-12">
      <h2 class="text-xl font-bold text-gray-900">Itinerary</h2>
      <div class="flex flex-col gap-3">
        <div v-for="(node, index) in plan.nodes" :key="node.id" class="border border-gray-200 p-4 rounded-lg bg-white hover:shadow-md transition-all">
          <div class="flex items-start gap-3">
            <!-- Position Badge -->
            <div class="flex-none w-8 h-8 rounded-full flex items-center justify-center text-white font-bold text-sm"
              :class="node.type === 'attraction' ? 'bg-blue-500' : 'bg-green-500'">
              {{ index + 1 }}
            </div>
            
            <!-- Node Details -->
            <div class="flex-1">
              <!-- Title and Type -->
              <div class="flex items-start justify-between gap-2">
                <div>
                  <h3 class="font-semibold text-gray-900 text-lg">{{ node.details?.name || node.details?.title || 'Unnamed Node' }}</h3>
                  <p class="text-xs text-gray-500 mt-1">
                    <span class="inline-block px-2 py-1 rounded-full"
                      :class="node.type === 'attraction' ? 'bg-blue-100 text-blue-700' : 'bg-green-100 text-green-700'">
                      {{ node.type === 'attraction' ? '📍 Attraction' : '🚗 Transition' }}
                    </span>
                  </p>
                </div>
              </div>
              
              <!-- Basic Description -->
              <p v-if="node.details?.description" class="text-sm text-gray-600 mt-2">{{ node.details.description }}</p>
              
              <!-- Node-Specific Information -->
              <div class="mt-3 space-y-1 text-sm">
                <!-- Attraction-specific info -->
                <div v-if="node.type === 'attraction'">
                  <p v-if="node.details?.location" class="text-gray-500">
                    <span class="font-medium">📍 Location:</span> {{ node.details.location }}
                  </p>
                  <p v-if="node.details?.category" class="text-gray-500 capitalize">
                    <span class="font-medium">🏷️ Category:</span> {{ node.details.category.replace(/_/g, ' ') }}
                  </p>
                </div>
                
                <!-- Plan-specific customization info -->
                <div v-if="node.duration_minutes || node.estimated_price_cents || node.description" class="border-t border-gray-200 pt-2 mt-2">
                  <p v-if="node.duration_minutes" class="text-gray-600">
                    <span class="font-medium">⏱️ Duration:</span> {{ node.duration_minutes }} min
                  </p>
                  <p v-if="node.estimated_price_cents" class="text-gray-600">
                    <span class="font-medium">💰 Price:</span> ${{ (node.estimated_price_cents / 100).toFixed(2) }}
                  </p>
                  <p v-if="node.description" class="text-gray-600 italic">
                    <span class="font-medium">📝 Plan Note:</span> {{ node.description }}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>


    <!-- Rating Section -->
    <div v-if="authStore.isAuthenticated" class="bg-white rounded-lg border border-gray-200 p-5 mb-6">
      <h2 class="text-lg font-bold text-gray-900 mb-4">Rate This Plan</h2>
      <div class="flex gap-2 mb-4">
        <button
          v-for="star in [1, 2, 3, 4, 5]"
          :key="star"
          @click="userRating = star"
          class="text-4xl bg-transparent border-0 cursor-pointer transition-all hover:scale-125"
          :style="{ opacity: star <= userRating ? 1 : 0.3 }"
        >
          ★
        </button>
      </div>
      <button
        @click="submitRating"
        :disabled="submittingRating || userRating === 0"
        class="px-4 py-2 bg-blue-600 text-white rounded-md text-sm font-medium hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
      >
        {{ submittingRating ? 'Submitting...' : 'Submit Rating' }}
      </button>
    </div>

    <div v-else class="bg-gray-100 rounded-lg border border-gray-300 p-4 mb-6">
      <p class="mb-2 text-gray-800">Want to rate this plan?</p>
      <router-link to="/login" class="text-blue-600 font-semibold hover:underline">Sign in</router-link>
    </div>

    <!-- Comments Section -->
    <div class="bg-white rounded-lg border border-gray-200 p-5">
      <h2 class="text-lg font-bold text-gray-900 mb-4">Comments</h2>

      <!-- Comment Form -->
      <div v-if="authStore.isAuthenticated" class="mb-6">
        <textarea
          v-model="newComment"
          placeholder="Share your thoughts about this plan..."
          class="w-full px-3 py-3 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100 font-inherit resize-y min-h-24"
        ></textarea>
        <button
          @click="submitComment"
          :disabled="submittingComment || !newComment.trim()"
          class="mt-3 px-4 py-2 bg-blue-600 text-white rounded-md text-sm font-medium hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
        >
          {{ submittingComment ? 'Posting...' : 'Post Comment' }}
        </button>
      </div>

      <div v-else class="bg-gray-100 rounded-lg border border-gray-300 p-4 mb-6">
        <p class="mb-2 text-gray-800">Want to comment?</p>
        <router-link to="/login" class="text-blue-600 font-semibold hover:underline">Sign in</router-link>
      </div>

      <!-- Comments List -->
      <div v-if="comments.length > 0" class="space-y-4">
        <div v-for="comment in comments" :key="comment.id" class="border border-gray-200 p-4 rounded-md">
          <div class="flex justify-between items-start mb-2">
            <span class="font-semibold text-gray-900">{{ comment.author?.username }}</span>
            <span class="text-xs text-gray-500">{{ formatDate(comment.created_at) }}</span>
          </div>
          <p class="text-gray-700 leading-relaxed">{{ comment.text }}</p>
        </div>
      </div>

      <div v-else class="text-center py-8 text-gray-600">
        <p>No comments yet. Be the first to comment!</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '../stores/auth_store';
import { useUiStore } from '../stores/ui_store';
import { planService, type TravelPlan, type PlanDetail } from '../services/plan_service';
import { commentService, type Comment } from '../services/comment_service';
import { ratingService } from '../services/rating_service';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();
const uiStore = useUiStore();

const plan = ref<PlanDetail | null>(null);
const comments = ref<Comment[]>([]);
const newComment = ref('');
const userRating = ref(0);
const submittingComment = ref(false);
const submittingRating = ref(false);

const isAuthor = ref(false);
const isAdmin = ref(false);

const formatDate = (dateStr?: string) => {
  if (!dateStr) return 'Unknown';
  return new Date(dateStr).toLocaleDateString();
};

async function loadPlan(): Promise<void> {
  try {
    const planId = route.params.id as string;
    if (!planId) {
      uiStore.showError('No plan ID provided');
      router.push('/browse');
      return;
    }
    plan.value = await planService.getPlanDetail(planId);
    if (!plan.value) {
      uiStore.showError('Plan not found');
      router.push('/browse');
      return;
    }
    isAuthor.value = plan.value?.author_id === authStore.user?.id;
    isAdmin.value = authStore.isAdmin;
    await loadComments();
  } catch (error) {
    console.error('Failed to load plan:', error);
    uiStore.showError('Failed to load plan');
    router.push('/browse');
  }
}

async function loadComments(): Promise<void> {
  try {
    const planId = route.params.id as string;
    if (!planId) {
      console.error('No plan ID for loading comments');
      return;
    }
    const result = await commentService.getComments(planId);
    comments.value = result?.comments || [];
  } catch (error) {
    console.error('Failed to load comments:', error);
    comments.value = [];
  }
}

async function submitComment(): Promise<void> {
  if (!newComment.value.trim()) return;

  try {
    submittingComment.value = true;
    const planId = route.params.id as string;
    if (!planId) {
      uiStore.showError('Invalid plan ID');
      return;
    }
    await commentService.createComment(planId, newComment.value.trim());
    newComment.value = '';
    uiStore.showSuccess('Comment posted');
    await loadComments();
  } catch (error) {
    console.error('Comment submission failed:', error);
    uiStore.showError('Failed to post comment');
  } finally {
    submittingComment.value = false;
  }
}

async function submitRating(): Promise<void> {
  if (userRating.value === 0) return;

  try {
    submittingRating.value = true;
    const planId = route.params.id as string;
    if (!planId) {
      uiStore.showError('Invalid plan ID');
      return;
    }
    await ratingService.submitRating(planId, userRating.value);
    userRating.value = 0;
    uiStore.showSuccess('Rating submitted');
  } catch (error) {
    console.error('Rating submission failed:', error);
    uiStore.showError('Failed to submit rating');
  } finally {
    submittingRating.value = false;
  }
}

function editPlan(): void {
  router.push(`/create-plan?edit=${plan.value?.id}`);
}

async function deletePlan(): Promise<void> {
  if (!plan.value?.id) {
    uiStore.showError('Invalid plan');
    return;
  }
  if (!confirm('Are you sure you want to delete this plan?')) return;

  try {
    await planService.deletePlan(plan.value.id);
    uiStore.showSuccess('Plan deleted');
    router.push('/browse');
  } catch (error) {
    console.error('Failed to delete plan:', error);
    uiStore.showError('Failed to delete plan');
  }
}

onMounted(() => {
  loadPlan();
});
</script>
