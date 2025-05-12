import { defineStore } from 'pinia';
import { v4 as uuidv4 } from 'uuid';

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    notifications: []
  }),
  
  actions: {
    addNotification({ type, message, timeout = 5000 }) {
      const id = uuidv4();
      
      // Add notification to the array
      this.notifications.push({
        id,
        type,
        message,
        show: true,
        timeout
      });
      
      // Auto-remove notification after timeout
      if (timeout > 0) {
        setTimeout(() => {
          this.removeNotification(id);
        }, timeout);
      }
      
      return id;
    },
    
    removeNotification(id) {
      const index = this.notifications.findIndex(n => n.id === id);
      if (index !== -1) {
        // First set show to false to trigger animation
        this.notifications[index].show = false;
        
        // Remove after animation completes
        setTimeout(() => {
          this.notifications = this.notifications.filter(n => n.id !== id);
        }, 300);
      }
    },
    
    clearNotifications() {
      this.notifications = [];
    },
    
    // Convenience methods for different notification types
    success(message, timeout = 5000) {
      return this.addNotification({ type: 'success', message, timeout });
    },
    
    error(message, timeout = 0) { // Errors don't auto-dismiss by default
      return this.addNotification({ type: 'error', message, timeout });
    },
    
    warning(message, timeout = 8000) { // Warnings stay longer
      return this.addNotification({ type: 'warning', message, timeout });
    },
    
    info(message, timeout = 5000) {
      return this.addNotification({ type: 'info', message, timeout });
    },
  }
}); 