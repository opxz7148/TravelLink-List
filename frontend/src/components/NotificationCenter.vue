<script setup lang="ts">
import { computed } from 'vue';
import { useUiStore } from '../stores/ui_store';

const uiStore = useUiStore();

const notifications = computed(() => uiStore.notifications);

const getNotificationClasses = (type: string) => {
  const baseClasses = 'flex items-center justify-between p-4 rounded shadow-md animate-slideIn font-medium';
  const typeClasses: Record<string, string> = {
    success: 'bg-green-100 text-green-800 border-l-4 border-green-500',
    error: 'bg-red-100 text-red-800 border-l-4 border-red-600',
    warning: 'bg-yellow-100 text-yellow-800 border-l-4 border-yellow-500',
    info: 'bg-blue-100 text-blue-800 border-l-4 border-blue-500',
  };
  return `${baseClasses} ${typeClasses[type] || typeClasses.info}`;
};

const removeNotification = (id: string) => {
  uiStore.removeNotification(id);
};
</script>

<template>
  <div v-if="notifications.length > 0" class="fixed top-20 right-5 z-50 flex flex-col gap-3 max-w-md md:max-w-lg">
    <transition-group name="slide" tag="div" class="flex flex-col gap-3">
      <div
        v-for="notification in notifications"
        :key="notification.id"
        :class="getNotificationClasses(notification.type)"
      >
        <div class="flex-1 flex items-center gap-3">
          <span>{{ notification.message }}</span>
        </div>
        <button
          class="bg-transparent border-none text-2xl cursor-pointer opacity-70 hover:opacity-100 transition-opacity p-0 ml-4"
          @click="removeNotification(notification.id)"
        >
          &times;
        </button>
      </div>
    </transition-group>
  </div>
</template>

<style scoped>
@keyframes slideIn {
  from {
    transform: translateX(400px);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

.animate-slideIn {
  animation: slideIn 0.3s ease-out;
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}

.slide-enter-from {
  transform: translateX(400px);
  opacity: 0;
}

.slide-leave-to {
  transform: translateX(400px);
  opacity: 0;
}
</style>
