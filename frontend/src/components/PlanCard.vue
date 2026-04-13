<template>
  <div class="group relative overflow-hidden rounded-2xl border border-slate-200 bg-white p-5 shadow-sm transition-all hover:-translate-y-0.5 hover:shadow-lg hover:border-emerald-300">
    <div class="pointer-events-none absolute -right-10 -top-10 h-24 w-24 rounded-full bg-emerald-100/60 blur-2xl"></div>
    <div class="flex justify-between items-start gap-2">
      <h3 class="tl-title text-xl font-semibold m-0 flex-1 text-slate-900">{{ plan.title }}</h3>
      <span
        v-if="plan.status"
        :class="[
          'px-3 py-1 rounded-full text-xs font-semibold whitespace-nowrap',
          plan.status === 'published' ? 'bg-green-100 text-green-700' : '',
          plan.status === 'draft' ? 'bg-amber-100 text-amber-700' : '',
          plan.status === 'archived' ? 'bg-gray-100 text-gray-600' : '',
        ]"
      >
        {{ formatStatus(plan.status) }}
      </span>
    </div>

    <div class="mt-3 flex gap-4 text-sm text-slate-600 flex-col sm:flex-row">
      <div class="flex items-center gap-1.5">
        <svg class="w-4 h-4 text-emerald-600" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path d="M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10 10-4.5 10-10S17.5 2 12 2z" />
        </svg>
        <span>{{ plan.destination }}</span>
      </div>
      <div class="flex items-center gap-1.5">
        <svg class="w-4 h-4 text-slate-400" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2" />
          <circle cx="12" cy="7" r="4" />
        </svg>
        <span>{{ plan.author?.username || 'Unknown' }}</span>
      </div>
    </div>

    <div class="mt-3 flex gap-4 py-2.5 border-y border-slate-100 text-sm">
      <div class="flex items-center gap-1.5 text-slate-600">
        <svg class="w-4 h-4 text-amber-500" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
        </svg>
        <span>{{ planRating }} ({{ plan.rating_count }})</span>
      </div>
      <div class="flex items-center gap-1.5 text-slate-600">
        <svg class="w-4 h-4 text-sky-500" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" />
        </svg>
        <span>{{ plan.comment_count }} comments</span>
      </div>
    </div>

    <p v-if="plan.destination" class="mt-3 text-sm text-slate-600 m-0 leading-relaxed">
      A travel plan to {{ plan.destination }}
    </p>

    <div class="mt-4 flex gap-2 flex-col sm:flex-row">
      <button class="px-4 py-2 rounded-lg text-sm font-semibold cursor-pointer transition-all bg-emerald-600 text-white hover:bg-emerald-700" @click="viewPlan">View Plan</button>
      <button v-if="canEdit" class="px-4 py-2 rounded-lg text-sm font-semibold cursor-pointer transition-all bg-slate-700 text-white hover:bg-slate-800" @click="editPlan">Edit</button>
      <button v-if="canDelete" class="px-4 py-2 rounded-lg text-sm font-semibold cursor-pointer transition-all bg-rose-600 text-white hover:bg-rose-700" @click="deletePlan">Delete</button>
    </div>

    <div class="mt-3 text-xs text-slate-400">
      Created {{ formatDate(plan.created_at) }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { TravelPlan } from '../services/plan_service';

interface Props {
  plan: TravelPlan;
  canEdit?: boolean;
  canDelete?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  canEdit: false,
  canDelete: false,
});

const emit = defineEmits<{
  view: [];
  edit: [];
  delete: [];
}>();

const planRating = computed(() => {
  if (props.plan.rating_count === 0) return '0';
  return (props.plan.rating_sum / props.plan.rating_count).toFixed(1);
});

function formatStatus(status: string): string {
  const statusMap: Record<string, string> = {
    draft: 'Draft',
    published: 'Published',
    archived: 'Archived',
  };
  return statusMap[status] || status;
}

function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

function viewPlan(): void {
  emit('view');
}

function editPlan(): void {
  emit('edit');
}

function deletePlan(): void {
  if (confirm('Are you sure you want to delete this plan?')) {
    emit('delete');
  }
}
</script>
