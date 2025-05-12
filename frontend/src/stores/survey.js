import { defineStore } from 'pinia';
import { surveyApi } from '@/api/survey';
import { useNotificationStore } from './notification';

export const useSurveyStore = defineStore('survey', {
  state: () => ({
    surveys: [],
    activeSurvey: null,
    loading: false,
    error: null,
    pagination: {
      total: 0,
      page: 1,
      perPage: 10
    }
  }),
  
  getters: {
    getSurveyById: (state) => (id) => {
      return state.surveys.find(survey => survey.id === id) || null;
    },
    getSurveysByStatus: (state) => (status) => {
      return state.surveys.filter(survey => survey.status === status);
    },
    totalPages: (state) => {
      return Math.ceil(state.pagination.total / state.pagination.perPage);
    },
    draftSurveys: (state) => {
      return state.surveys.filter(survey => survey.status === 'draft');
    },
    publishedSurveys: (state) => {
      return state.surveys.filter(survey => survey.status === 'published');
    },
    closedSurveys: (state) => {
      return state.surveys.filter(survey => survey.status === 'closed');
    }
  },
  
  actions: {
    async fetchSurveys(page = 1, perPage = 10) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await surveyApi.getSurveys(page, perPage);
        this.surveys = response.data;
        this.pagination = {
          total: response.meta.total,
          page: response.meta.page,
          perPage: response.meta.per_page
        };
      } catch (error) {
        this.error = error.message || 'Failed to fetch surveys';
        useNotificationStore().error(this.error);
      } finally {
        this.loading = false;
      }
    },
    
    async fetchSurveyById(id) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await surveyApi.getSurveyById(id);
        this.activeSurvey = response;
        
        // Also update the survey in the list if it exists
        const index = this.surveys.findIndex(s => s.id === id);
        if (index !== -1) {
          this.surveys[index] = response;
        }
        
        return response;
      } catch (error) {
        this.error = error.message || `Failed to fetch survey with ID ${id}`;
        useNotificationStore().error(this.error);
        return null;
      } finally {
        this.loading = false;
      }
    },
    
    async createSurvey(surveyData) {
      this.loading = true;
      this.error = null;
      const notificationStore = useNotificationStore();
      
      try {
        const response = await surveyApi.createSurvey(surveyData);
        
        // Add the new survey to the list
        this.surveys.unshift(response);
        this.activeSurvey = response;
        
        notificationStore.success('Survey created successfully!');
        return response;
      } catch (error) {
        this.error = error.message || 'Failed to create survey';
        notificationStore.error(this.error);
        return null;
      } finally {
        this.loading = false;
      }
    },
    
    async updateSurvey(id, surveyData) {
      this.loading = true;
      this.error = null;
      const notificationStore = useNotificationStore();
      
      try {
        const response = await surveyApi.updateSurvey(id, surveyData);
        
        // Update the survey in the list
        const index = this.surveys.findIndex(s => s.id === id);
        if (index !== -1) {
          this.surveys[index] = response;
        }
        
        this.activeSurvey = response;
        
        notificationStore.success('Survey updated successfully!');
        return response;
      } catch (error) {
        this.error = error.message || `Failed to update survey with ID ${id}`;
        notificationStore.error(this.error);
        return null;
      } finally {
        this.loading = false;
      }
    },
    
    async deleteSurvey(id) {
      this.loading = true;
      this.error = null;
      const notificationStore = useNotificationStore();
      
      try {
        await surveyApi.deleteSurvey(id);
        
        // Remove the survey from the list
        this.surveys = this.surveys.filter(s => s.id !== id);
        
        // Clear active survey if it's the one being deleted
        if (this.activeSurvey && this.activeSurvey.id === id) {
          this.activeSurvey = null;
        }
        
        notificationStore.success('Survey deleted successfully!');
        return true;
      } catch (error) {
        this.error = error.message || `Failed to delete survey with ID ${id}`;
        notificationStore.error(this.error);
        return false;
      } finally {
        this.loading = false;
      }
    },
    
    async updateSurveyStatus(id, status) {
      this.loading = true;
      this.error = null;
      const notificationStore = useNotificationStore();
      
      try {
        const response = await surveyApi.updateSurveyStatus(id, status);
        
        // Update the survey in the list
        const index = this.surveys.findIndex(s => s.id === id);
        if (index !== -1) {
          this.surveys[index].status = status;
        }
        
        // Update active survey if it's the one being modified
        if (this.activeSurvey && this.activeSurvey.id === id) {
          this.activeSurvey.status = status;
        }
        
        const statusText = status.charAt(0).toUpperCase() + status.slice(1);
        notificationStore.success(`Survey ${statusText} successfully!`);
        return response;
      } catch (error) {
        this.error = error.message || `Failed to update survey status`;
        notificationStore.error(this.error);
        return null;
      } finally {
        this.loading = false;
      }
    },
    
    resetActiveSurvey() {
      this.activeSurvey = null;
    }
  }
}); 