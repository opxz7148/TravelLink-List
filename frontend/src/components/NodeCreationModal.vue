<template>
  <div v-if="isOpen" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
    <div class="bg-white rounded-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto shadow-xl">
      <!-- Header -->
      <div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex justify-between items-center">
        <h2 class="text-xl font-bold text-gray-900">Create New Node</h2>
        <button @click="close" class="text-gray-500 hover:text-gray-700 text-2xl leading-none">&times;</button>
      </div>

      <!-- Content -->
      <div class="p-6 space-y-6">
        <!-- Node Type Selection -->
        <div>
          <h3 class="text-base font-semibold text-gray-900 mb-3">Node Type</h3>
          <div class="grid grid-cols-2 gap-3">
            <button
              @click="nodeType = 'attraction'"
              :class="[
                'p-4 border-2 rounded-lg text-center transition-all',
                nodeType === 'attraction'
                  ? 'border-blue-500 bg-blue-50'
                  : 'border-gray-200 bg-white hover:border-gray-300'
              ]"
            >
              <div class="text-2xl mb-2">🏛️</div>
              <p class="font-semibold text-gray-900">Attraction</p>
              <p class="text-xs text-gray-600">Tourist site, restaurant, hotel, museum, park, shopping, entertainment</p>
            </button>
            <button
              @click="nodeType = 'transition'"
              :class="[
                'p-4 border-2 rounded-lg text-center transition-all',
                nodeType === 'transition'
                  ? 'border-blue-500 bg-blue-50'
                  : 'border-gray-200 bg-white hover:border-gray-300'
              ]"
            >
              <div class="text-2xl mb-2">🛣️</div>
              <p class="font-semibold text-gray-900">Transition</p>
              <p class="text-xs text-gray-600">Walking, transit, driving route</p>
            </button>
          </div>
        </div>

        <!-- Attraction Form -->
        <div v-if="nodeType === 'attraction'" class="space-y-4">
          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-2">Name *</label>
            <input
              v-model="attractionForm.name"
              type="text"
              placeholder="e.g., Eiffel Tower"
              maxlength="100"
              class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100"
            />
          </div>

          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-2">Category *</label>
            <select
              v-model="attractionForm.category"
              class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm bg-white focus:outline-none focus:border-blue-500"
            >
              <option value="">Select a category</option>
              <option value="tourist_attraction">🏛️ Tourist Attraction</option>
              <option value="restaurant">🍽️ Restaurant</option>
              <option value="hotel">🏨 Hotel</option>
              <option value="museum">🎨 Museum</option>
              <option value="park">🌳 Park</option>
              <option value="shopping">🛍️ Shopping</option>
              <option value="entertainment">🎭 Entertainment</option>
              <option value="other">📍 Other</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-2">Location *</label>
            <input
              v-model="attractionForm.location"
              type="text"
              placeholder="e.g., Paris, France"
              maxlength="100"
              class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100"
            />
          </div>

          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-2">Description</label>
            <textarea
              v-model="attractionForm.description"
              placeholder="Describe this attraction..."
              maxlength="500"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100 font-inherit"
            />
            <p class="text-xs text-gray-500 mt-1">{{ (attractionForm.description || '').length }}/500</p>
          </div>

          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-2">Contact Info</label>
            <input
              v-model="attractionForm.contact_info"
              type="text"
              placeholder="Phone or website"
              maxlength="200"
              class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100"
            />
          </div>

          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-2">Hours of Operation</label>
            <input
              v-model="attractionForm.hours_of_operation"
              type="text"
              placeholder="e.g., 9AM - 5PM"
              maxlength="100"
              class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100"
            />
          </div>

          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-2">Estimated Visit Duration (minutes)</label>
            <input
              v-model="attractionForm.estimated_visit_duration_minutes"
              type="number"
              placeholder="e.g., 90"
              min="1"
              class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100"
            />
          </div>
        </div>

        <!-- Transition Form -->
        <div v-if="nodeType === 'transition'" class="space-y-4">
          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-2">Title *</label>
            <input
              v-model="transitionForm.title"
              type="text"
              placeholder="e.g., Walk to Louvre Museum"
              maxlength="100"
              class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100"
            />
          </div>

          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-2">Mode *</label>
            <select
              v-model="transitionForm.mode"
              class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm bg-white focus:outline-none focus:border-blue-500"
            >
              <option value="">Select travel mode</option>
              <option value="walking">🚶 Walking</option>
              <option value="transit">🚇 Public Transit</option>
              <option value="driving">🚗 Driving</option>
              <option value="cycling">🚴 Cycling</option>
              <option value="flight">✈️ Flight</option>
              <option value="train">🚂 Train</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-2">Description</label>
            <textarea
              v-model="transitionForm.description"
              placeholder="Describe the route or transition..."
              maxlength="500"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100 font-inherit"
            />
            <p class="text-xs text-gray-500 mt-1">{{ (transitionForm.description || '').length }}/500</p>
          </div>

          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-2">Hours</label>
            <input
              v-model="transitionForm.hours_of_operation"
              type="text"
              placeholder="e.g., Available 24/7"
              maxlength="100"
              class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100"
            />
          </div>

          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-2">Route Notes</label>
            <textarea
              v-model="transitionForm.route_notes"
              placeholder="Directions, landmarks, warnings..."
              maxlength="500"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-500 focus:ring-3 focus:ring-blue-100 font-inherit"
            />
            <p class="text-xs text-gray-500 mt-1">{{ (transitionForm.route_notes || '').length }}/500</p>
          </div>
        </div>

        <!-- Status Message -->
        <div v-if="nodeType === 'attraction'" class="p-3 bg-orange-50 text-orange-700 rounded border border-orange-200 text-sm">
          <p class="font-semibold mb-1">ℹ️ Node Status</p>
          <p>Your new attraction will be created as <strong>draft</strong>. It will be automatically published when you publish your plan (if you're a Traveller or Admin).</p>
        </div>
        <div v-if="nodeType === 'transition'" class="p-3 bg-orange-50 text-orange-700 rounded border border-orange-200 text-sm">
          <p class="font-semibold mb-1">ℹ️ Node Status</p>
          <p>Your new transition will be created as <strong>draft</strong>. It will be automatically published when you publish your plan (if you're a Traveller or Admin).</p>
        </div>

        <!-- Error Message -->
        <div v-if="errorMessage" class="p-3 bg-red-50 text-red-700 rounded border border-red-200 text-sm">
          {{ errorMessage }}
        </div>
      </div>

      <!-- Footer -->
      <div class="sticky bottom-0 bg-gray-50 border-t border-gray-200 px-6 py-4 flex gap-3 justify-end">
        <button
          @click="close"
          class="px-4 py-2 bg-gray-200 text-gray-900 rounded-md font-medium hover:bg-gray-300 transition-all"
        >
          Cancel
        </button>
        <button
          @click="createNode"
          :disabled="loading || !isFormValid"
          class="px-4 py-2 bg-blue-600 text-white rounded-md font-medium hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
        >
          {{ loading ? 'Creating...' : 'Create Node' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { nodeService } from '../services/node_service';

interface Props {
  isOpen: boolean;
}

interface Emits {
  (e: 'close'): void;
  (e: 'created'): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const nodeType = ref<'attraction' | 'transition'>('attraction');
const loading = ref(false);
const errorMessage = ref('');

// Attraction form
const attractionForm = ref({
  name: '',
  category: '',
  location: '',
  description: '',
  contact_info: '',
  hours_of_operation: '',
  estimated_visit_duration_minutes: '',
});

// Transition form
const transitionForm = ref({
  title: '',
  mode: '',
  description: '',
  hours_of_operation: '',
  route_notes: '',
});

const isFormValid = computed(() => {
  if (nodeType.value === 'attraction') {
    return (
      attractionForm.value.name.trim() &&
      attractionForm.value.category &&
      attractionForm.value.location.trim()
    );
  } else {
    return (
      transitionForm.value.title.trim() &&
      transitionForm.value.mode
    );
  }
});

async function createNode() {
  errorMessage.value = '';
  loading.value = true;

  try {
    if (nodeType.value === 'attraction') {
      const duration = attractionForm.value.estimated_visit_duration_minutes
        ? parseInt(attractionForm.value.estimated_visit_duration_minutes)
        : undefined;

      await nodeService.createAttractionNode({
        name: attractionForm.value.name,
        category: attractionForm.value.category || undefined,
        location: attractionForm.value.location,
        description: attractionForm.value.description || undefined,
        contact_info: attractionForm.value.contact_info || undefined,
        hours_of_operation: attractionForm.value.hours_of_operation || undefined,
      });
    } else {
      await nodeService.createTransitionNode({
        title: transitionForm.value.title,
        mode: transitionForm.value.mode,
        description: transitionForm.value.description || undefined,
        hours_of_operation: transitionForm.value.hours_of_operation || undefined,
        route_notes: transitionForm.value.route_notes || undefined,
      });
    }

    // Reset form
    attractionForm.value = {
      name: '',
      category: '',
      location: '',
      description: '',
      contact_info: '',
      hours_of_operation: '',
      estimated_visit_duration_minutes: '',
    };
    transitionForm.value = {
      title: '',
      mode: '',
      description: '',
      hours_of_operation: '',
      route_notes: '',
    };

    emit('created');
    close();
  } catch (error) {
    console.error('Failed to create node:', error);
    errorMessage.value = error instanceof Error ? error.message : 'Failed to create node. Please try again.';
  } finally {
    loading.value = false;
  }
}

function close() {
  errorMessage.value = '';
  emit('close');
}
</script>
