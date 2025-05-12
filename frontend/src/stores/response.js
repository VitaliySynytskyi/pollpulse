import { defineStore } from 'pinia';
import { responseApi } from '@/api/response';
import { useNotificationStore } from './notification';

export const useResponseStore = defineStore('response', {
  state: () => ({
    responses: [],
    surveyResults: null,
    currentResponse: null,
    loading: false,
    error: null,
    pagination: {
      total: 0,
      page: 1,
      perPage: 20
    }
  }),
  
  getters: {
    completionRate: (state) => {
      if (!state.surveyResults) return 0;
      return state.surveyResults.completion_rate;
    },
    totalResponses: (state) => {
      if (!state.surveyResults) return 0;
      return state.surveyResults.response_count;
    },
    questionResults: (state) => {
      if (!state.surveyResults) return [];
      return state.surveyResults.questions;
    }
  },
  
  actions: {
    async fetchResponses(surveyId, page = 1, perPage = 20) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await responseApi.getResponses(surveyId, page, perPage);
        this.responses = response.data;
        this.pagination = {
          total: response.meta.total,
          page: response.meta.page,
          perPage: response.meta.per_page
        };
      } catch (error) {
        this.error = error.message || 'Failed to fetch responses';
        useNotificationStore().error(this.error);
      } finally {
        this.loading = false;
      }
    },
    
    async fetchResponseById(id) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await responseApi.getResponseById(id);
        this.currentResponse = response;
        return response;
      } catch (error) {
        this.error = error.message || `Failed to fetch response with ID ${id}`;
        useNotificationStore().error(this.error);
        return null;
      } finally {
        this.loading = false;
      }
    },
    
    async submitResponse(surveyId, responseData) {
      this.loading = true;
      this.error = null;
      const notificationStore = useNotificationStore();
      
      try {
        const response = await responseApi.submitResponse(surveyId, responseData);
        notificationStore.success('Survey submitted successfully!');
        this.currentResponse = response;
        return response;
      } catch (error) {
        this.error = error.message || 'Failed to submit response';
        notificationStore.error(this.error);
        return null;
      } finally {
        this.loading = false;
      }
    },
    
    async fetchSurveyResults(surveyId) {
      this.loading = true;
      this.error = null;
      
      try {
        const results = await responseApi.getSurveyResults(surveyId);
        this.surveyResults = results;
        return results;
      } catch (error) {
        this.error = error.message || `Failed to fetch results for survey ${surveyId}`;
        useNotificationStore().error(this.error);
        return null;
      } finally {
        this.loading = false;
      }
    },
    
    async exportResults(surveyId, format) {
      this.loading = true;
      this.error = null;
      const notificationStore = useNotificationStore();
      
      try {
        const exportData = await responseApi.exportResults(surveyId, format);
        notificationStore.success(`Survey results exported successfully!`);
        return exportData;
      } catch (error) {
        this.error = error.message || 'Failed to export results';
        notificationStore.error(this.error);
        return null;
      } finally {
        this.loading = false;
      }
    },
    
    async getAnalytics(surveyId, period = 'all', startDate = null, endDate = null) {
      this.loading = true;
      this.error = null;
      
      try {
        const analytics = await responseApi.getAnalytics(surveyId, period, startDate, endDate);
        return analytics;
      } catch (error) {
        this.error = error.message || 'Failed to fetch analytics';
        useNotificationStore().error(this.error);
        return null;
      } finally {
        this.loading = false;
      }
    },
    
    resetCurrentResponse() {
      this.currentResponse = null;
    }
  }
}); 