import { defineStore } from 'pinia';
import { ref } from 'vue';

export type NotificationType = 'success' | 'error' | 'warning' | 'info';

export interface Notification {
  id: string;
  type: NotificationType;
  message: string;
  duration?: number;
}

export const useUiStore = defineStore('ui', () => {
  // State
  const notifications = ref<Notification[]>([]);
  const isLoading = ref(false);
  const sidebarOpen = ref(false);

  // Actions
  const addNotification = (message: string, type: NotificationType = 'info', duration: number = 3000) => {
    const id = Math.random().toString(36).substring(2, 9);
    const notification: Notification = {
      id,
      type,
      message,
      duration,
    };

    notifications.value.push(notification);

    if (duration > 0) {
      setTimeout(() => {
        removeNotification(id);
      }, duration);
    }

    return id;
  };

  const removeNotification = (id: string) => {
    notifications.value = notifications.value.filter((n) => n.id !== id);
  };

  const clearNotifications = () => {
    notifications.value = [];
  };

  const showSuccess = (message: string, duration?: number) => {
    return addNotification(message, 'success', duration);
  };

  const showError = (message: string, duration?: number) => {
    return addNotification(message, 'error', duration || 5000);
  };

  const showWarning = (message: string, duration?: number) => {
    return addNotification(message, 'warning', duration);
  };

  const showInfo = (message: string, duration?: number) => {
    return addNotification(message, 'info', duration);
  };

  const setLoading = (loading: boolean) => {
    isLoading.value = loading;
  };

  const toggleSidebar = () => {
    sidebarOpen.value = !sidebarOpen.value;
  };

  const closeSidebar = () => {
    sidebarOpen.value = false;
  };

  const openSidebar = () => {
    sidebarOpen.value = true;
  };

  return {
    // State
    notifications,
    isLoading,
    sidebarOpen,

    // Actions
    addNotification,
    removeNotification,
    clearNotifications,
    showSuccess,
    showError,
    showWarning,
    showInfo,
    setLoading,
    toggleSidebar,
    closeSidebar,
    openSidebar,
  };
});
