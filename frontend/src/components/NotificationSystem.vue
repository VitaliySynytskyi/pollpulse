<template>
  <div 
    class="fixed right-0 top-0 mr-4 mt-4 z-50 w-full max-w-md flex flex-col items-end space-y-2"
    style="pointer-events: none;"
  >
    <transition-group name="notification">
      <div 
        v-for="notification in notifications" 
        :key="notification.id" 
        class="flex items-center p-4 mb-2 rounded-lg shadow-lg transition-all transform duration-300"
        :class="notificationClasses(notification)"
        style="pointer-events: auto;"
      >
        <div class="mr-3">
          <v-icon 
            :icon="notificationIcon(notification)" 
            :color="notificationIconColor(notification)"
          ></v-icon>
        </div>
        <div class="flex-1">
          <p class="font-medium">{{ notification.message }}</p>
        </div>
        <button 
          @click="dismissNotification(notification.id)" 
          class="text-gray-400 hover:text-gray-800 ml-4 focus:outline-none"
        >
          <v-icon icon="mdi-close"></v-icon>
        </button>
      </div>
    </transition-group>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { storeToRefs } from 'pinia';
import { useNotificationStore } from '@/stores/notification';

const notificationStore = useNotificationStore();
const { notifications } = storeToRefs(notificationStore);

// Classes for notification based on type
const notificationClasses = (notification) => {
  const baseClasses = 'w-full';
  
  switch (notification.type) {
    case 'success':
      return `${baseClasses} bg-green-50 text-green-800 border-l-4 border-green-500`;
    case 'error':
      return `${baseClasses} bg-red-50 text-red-800 border-l-4 border-red-500`;
    case 'warning':
      return `${baseClasses} bg-amber-50 text-amber-800 border-l-4 border-amber-500`;
    case 'info':
    default:
      return `${baseClasses} bg-blue-50 text-blue-800 border-l-4 border-blue-500`;
  }
};

// Icon for notification based on type
const notificationIcon = (notification) => {
  switch (notification.type) {
    case 'success':
      return 'mdi-check-circle';
    case 'error':
      return 'mdi-alert-circle';
    case 'warning':
      return 'mdi-alert';
    case 'info':
    default:
      return 'mdi-information';
  }
};

// Icon color for notification based on type
const notificationIconColor = (notification) => {
  switch (notification.type) {
    case 'success':
      return 'green';
    case 'error':
      return 'red';
    case 'warning':
      return 'amber';
    case 'info':
    default:
      return 'blue';
  }
};

// Dismiss a notification
const dismissNotification = (id) => {
  notificationStore.dismissNotification(id);
};
</script>

<style scoped>
.notification-enter-active,
.notification-leave-active {
  transition: all 0.3s ease;
}

.notification-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.notification-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style> 