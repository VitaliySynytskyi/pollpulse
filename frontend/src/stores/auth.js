import { defineStore } from 'pinia';
import { authApi } from '@/api/auth';
import { useNotificationStore } from './notification';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: null,
    loading: false,
    error: null,
  }),
  
  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.user && state.user.roles && state.user.roles.includes('admin'),
    isSurveyCreator: (state) => {
      if (!state.user || !state.user.roles) return false;
      return state.user.roles.includes('admin') || state.user.roles.includes('survey_creator');
    },
    fullName: (state) => {
      if (!state.user) return '';
      return `${state.user.first_name} ${state.user.last_name}`;
    }
  },
  
  actions: {
    async login(username, password) {
      this.loading = true;
      this.error = null;
      const notificationStore = useNotificationStore();
      
      try {
        const response = await authApi.login(username, password);
        
        if (response.token && response.user) {
          this.setAuthData(response.token, response.user);
          notificationStore.addNotification({
            type: 'success',
            message: 'Logged in successfully!'
          });
          return true;
        } else {
          throw new Error('Invalid response from server');
        }
      } catch (error) {
        this.error = error.message || 'Failed to login';
        notificationStore.addNotification({
          type: 'error',
          message: this.error
        });
        return false;
      } finally {
        this.loading = false;
      }
    },
    
    async register(userData) {
      this.loading = true;
      this.error = null;
      const notificationStore = useNotificationStore();
      
      try {
        const response = await authApi.register(userData);
        
        if (response.token && response.user) {
          this.setAuthData(response.token, response.user);
          notificationStore.addNotification({
            type: 'success',
            message: 'Registered and logged in successfully!'
          });
          return true;
        } else {
          throw new Error('Invalid response from server');
        }
      } catch (error) {
        this.error = error.message || 'Failed to register';
        notificationStore.addNotification({
          type: 'error',
          message: this.error
        });
        return false;
      } finally {
        this.loading = false;
      }
    },
    
    async logout() {
      const notificationStore = useNotificationStore();
      this.clearAuthData();
      notificationStore.addNotification({
        type: 'info',
        message: 'You have been logged out.'
      });
    },
    
    checkAuth() {
      const token = localStorage.getItem('token');
      const userStr = localStorage.getItem('user');
      
      if (token && userStr) {
        try {
          const user = JSON.parse(userStr);
          this.setAuthData(token, user, false);
        } catch (e) {
          this.clearAuthData();
        }
      }
    },
    
    setAuthData(token, user, persist = true) {
      this.token = token;
      this.user = user;
      
      if (persist) {
        localStorage.setItem('token', token);
        localStorage.setItem('user', JSON.stringify(user));
      }
    },
    
    clearAuthData() {
      this.token = null;
      this.user = null;
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    },
    
    async updateProfile(userData) {
      this.loading = true;
      this.error = null;
      const notificationStore = useNotificationStore();
      
      try {
        const response = await authApi.updateProfile(userData);
        this.user = response;
        
        // Update the stored user data
        localStorage.setItem('user', JSON.stringify(response));
        
        notificationStore.addNotification({
          type: 'success',
          message: 'Profile updated successfully!'
        });
        return true;
      } catch (error) {
        this.error = error.message || 'Failed to update profile';
        notificationStore.addNotification({
          type: 'error',
          message: this.error
        });
        return false;
      } finally {
        this.loading = false;
      }
    },
  }
}); 