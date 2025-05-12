<template>
  <div class="min-h-screen bg-gray-50 py-6">
    <div class="container mx-auto px-4">
      <h1 class="text-3xl font-bold text-gray-900 mb-6">Dashboard</h1>

      <!-- Welcome Card -->
      <div class="bg-white rounded-lg shadow-md p-6 mb-6">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl font-semibold text-gray-800">Welcome back, {{ user?.name || 'User' }}!</h2>
            <p class="text-gray-600 mt-1">Here's an overview of your surveys and responses.</p>
          </div>
          <router-link 
            to="/surveys/create" 
            class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 transition-colors"
          >
            <v-icon icon="mdi-plus" class="mr-1"></v-icon> New Survey
          </router-link>
        </div>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div class="bg-white rounded-lg shadow-md p-6 flex items-center">
          <div class="w-12 h-12 rounded-full bg-blue-100 flex items-center justify-center mr-4">
            <v-icon icon="mdi-file-document-outline" size="large" color="blue-darken-2"></v-icon>
          </div>
          <div>
            <p class="text-gray-500 text-sm">Total Surveys</p>
            <p class="text-2xl font-semibold">{{ statsLoading ? '...' : stats.totalSurveys }}</p>
          </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6 flex items-center">
          <div class="w-12 h-12 rounded-full bg-green-100 flex items-center justify-center mr-4">
            <v-icon icon="mdi-check-circle-outline" size="large" color="green-darken-2"></v-icon>
          </div>
          <div>
            <p class="text-gray-500 text-sm">Active Surveys</p>
            <p class="text-2xl font-semibold">{{ statsLoading ? '...' : stats.activeSurveys }}</p>
          </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6 flex items-center">
          <div class="w-12 h-12 rounded-full bg-purple-100 flex items-center justify-center mr-4">
            <v-icon icon="mdi-poll" size="large" color="purple-darken-2"></v-icon>
          </div>
          <div>
            <p class="text-gray-500 text-sm">Total Responses</p>
            <p class="text-2xl font-semibold">{{ statsLoading ? '...' : stats.totalResponses }}</p>
          </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6 flex items-center">
          <div class="w-12 h-12 rounded-full bg-amber-100 flex items-center justify-center mr-4">
            <v-icon icon="mdi-account-group-outline" size="large" color="amber-darken-2"></v-icon>
          </div>
          <div>
            <p class="text-gray-500 text-sm">Average Response Rate</p>
            <p class="text-2xl font-semibold">{{ statsLoading ? '...' : stats.responseRate }}%</p>
          </div>
        </div>
      </div>

      <!-- Recent Surveys -->
      <div class="bg-white rounded-lg shadow-md p-6 mb-6">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold text-gray-800">Recent Surveys</h2>
          <router-link to="/surveys" class="text-indigo-600 hover:text-indigo-800">
            View all <v-icon icon="mdi-chevron-right" size="small"></v-icon>
          </router-link>
        </div>
        
        <div v-if="surveyLoading" class="text-center py-8">
          <v-progress-circular indeterminate color="primary"></v-progress-circular>
        </div>
        
        <div v-else-if="recentSurveys.length === 0" class="py-8 text-center">
          <p class="text-gray-500">You don't have any surveys yet.</p>
          <router-link 
            to="/surveys/create" 
            class="mt-4 inline-block px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 transition-colors"
          >
            Create your first survey
          </router-link>
        </div>
        
        <div v-else>
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Title
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Responses
                  </th>
                  <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Created
                  </th>
                  <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="survey in recentSurveys" :key="survey.id">
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm font-medium text-gray-900">{{ survey.title }}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span 
                      class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full" 
                      :class="{
                        'bg-green-100 text-green-800': survey.status === 'active',
                        'bg-gray-100 text-gray-800': survey.status === 'draft',
                        'bg-red-100 text-red-800': survey.status === 'closed'
                      }"
                    >
                      {{ survey.status.charAt(0).toUpperCase() + survey.status.slice(1) }}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {{ survey.responseCount }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {{ formatDate(survey.createdAt) }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <router-link :to="`/surveys/${survey.id}`" class="text-indigo-600 hover:text-indigo-900 mr-3">
                      <v-icon icon="mdi-eye" size="small"></v-icon>
                    </router-link>
                    <router-link :to="`/surveys/${survey.id}/edit`" class="text-blue-600 hover:text-blue-900 mr-3">
                      <v-icon icon="mdi-pencil" size="small"></v-icon>
                    </router-link>
                    <router-link :to="`/surveys/${survey.id}/results`" class="text-green-600 hover:text-green-900">
                      <v-icon icon="mdi-chart-bar" size="small"></v-icon>
                    </router-link>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Recent Responses -->
      <div class="bg-white rounded-lg shadow-md p-6">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold text-gray-800">Recent Responses</h2>
        </div>
        
        <div v-if="responseLoading" class="text-center py-8">
          <v-progress-circular indeterminate color="primary"></v-progress-circular>
        </div>
        
        <div v-else-if="recentResponses.length === 0" class="py-8 text-center">
          <p class="text-gray-500">No responses received yet.</p>
        </div>
        
        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div v-for="response in recentResponses" :key="response.id" class="border border-gray-200 rounded-lg p-4">
            <div class="flex justify-between items-start mb-2">
              <h3 class="text-md font-medium text-gray-900">{{ response.surveyTitle }}</h3>
              <span class="text-sm text-gray-500">{{ formatDate(response.createdAt) }}</span>
            </div>
            <p class="text-sm text-gray-600 mb-3">
              Respondent: {{ response.respondentEmail || 'Anonymous' }}
            </p>
            <router-link 
              :to="`/surveys/${response.surveyId}/results`" 
              class="text-sm text-indigo-600 hover:text-indigo-800"
            >
              View Survey Results <v-icon icon="mdi-arrow-right" size="small"></v-icon>
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useSurveyStore } from '@/stores/survey';
import { useResponseStore } from '@/stores/response';

const authStore = useAuthStore();
const surveyStore = useSurveyStore();
const responseStore = useResponseStore();

const user = computed(() => authStore.user);

// Stats
const statsLoading = ref(true);
const stats = ref({
  totalSurveys: 0,
  activeSurveys: 0,
  totalResponses: 0,
  responseRate: 0
});

// Surveys
const surveyLoading = ref(true);
const recentSurveys = ref([]);

// Responses
const responseLoading = ref(true);
const recentResponses = ref([]);

// Format date helper
const formatDate = (dateString) => {
  const options = { year: 'numeric', month: 'short', day: 'numeric' };
  return new Date(dateString).toLocaleDateString(undefined, options);
};

// Load dashboard data
onMounted(async () => {
  try {
    // Fetch stats
    statsLoading.value = true;
    const statsData = await surveyStore.fetchDashboardStats();
    stats.value = statsData;
    statsLoading.value = false;
    
    // Fetch recent surveys
    surveyLoading.value = true;
    const surveysData = await surveyStore.fetchSurveys({ limit: 5 });
    recentSurveys.value = surveysData.surveys;
    surveyLoading.value = false;
    
    // Fetch recent responses
    responseLoading.value = true;
    const responsesData = await responseStore.fetchRecentResponses({ limit: 4 });
    recentResponses.value = responsesData;
    responseLoading.value = false;
  } catch (error) {
    console.error('Error loading dashboard data:', error);
  }
});
</script> 