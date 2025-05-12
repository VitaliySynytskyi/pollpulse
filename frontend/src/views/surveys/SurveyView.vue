<template>
  <div class="min-h-screen bg-gray-50 py-6">
    <div class="container mx-auto px-4">
      <!-- Loading state -->
      <div v-if="loading" class="bg-white rounded-lg shadow-md p-16 text-center">
        <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
        <p class="mt-4 text-gray-600">Loading survey...</p>
      </div>

      <!-- Error state -->
      <div v-else-if="error" class="bg-white rounded-lg shadow-md p-8">
        <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md mb-4">
          {{ error }}
        </div>
        <div class="flex justify-center mt-4">
          <router-link to="/surveys" class="text-indigo-600 hover:text-indigo-800">
            &larr; Back to Surveys
          </router-link>
        </div>
      </div>

      <!-- Survey content -->
      <div v-else>
        <div class="flex justify-between items-center mb-6">
          <h1 class="text-3xl font-bold text-gray-900">{{ survey.title }}</h1>
          <div class="space-x-2">
            <v-btn
              variant="outlined"
              color="default"
              @click="$router.go(-1)"
              class="mr-2"
            >
              <v-icon icon="mdi-arrow-left" class="mr-1"></v-icon>
              Back
            </v-btn>
            
            <v-btn
              v-if="survey.status === 'draft'"
              color="primary"
              @click="publishSurvey"
              :loading="publishLoading"
            >
              <v-icon icon="mdi-publish" class="mr-1"></v-icon>
              Publish
            </v-btn>
            
            <v-btn
              v-if="survey.status === 'active'"
              color="amber-darken-3"
              @click="closeSurvey"
              :loading="closeLoading"
            >
              <v-icon icon="mdi-close-circle" class="mr-1"></v-icon>
              Close
            </v-btn>
          </div>
        </div>

        <!-- Survey info card -->
        <div class="bg-white rounded-lg shadow-md p-6 mb-6">
          <div class="flex justify-between items-start">
            <div>
              <div class="flex items-center mb-4">
                <span 
                  class="px-2 py-1 text-xs font-semibold rounded-full mr-2" 
                  :class="{
                    'bg-green-100 text-green-800': survey.status === 'active',
                    'bg-gray-100 text-gray-800': survey.status === 'draft',
                    'bg-red-100 text-red-800': survey.status === 'closed'
                  }"
                >
                  {{ survey.status.charAt(0).toUpperCase() + survey.status.slice(1) }}
                </span>
                
                <span class="text-sm text-gray-500">
                  Created {{ formatDate(survey.createdAt) }}
                </span>
              </div>

              <p class="text-gray-700 mb-6">{{ survey.description }}</p>

              <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
                <div class="bg-gray-50 p-4 rounded-lg">
                  <div class="text-gray-500 text-sm mb-1">Responses</div>
                  <div class="text-2xl font-semibold">{{ survey.responseCount }}</div>
                </div>
                
                <div class="bg-gray-50 p-4 rounded-lg">
                  <div class="text-gray-500 text-sm mb-1">Questions</div>
                  <div class="text-2xl font-semibold">{{ survey.questions?.length || 0 }}</div>
                </div>
                
                <div class="bg-gray-50 p-4 rounded-lg">
                  <div class="text-gray-500 text-sm mb-1">Completion Rate</div>
                  <div class="text-2xl font-semibold">{{ survey.completionRate || 0 }}%</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Action buttons -->
        <div class="bg-white rounded-lg shadow-md p-6 mb-6">
          <h2 class="text-xl font-semibold text-gray-800 mb-4">Survey Actions</h2>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="border border-gray-200 rounded-lg p-4 flex flex-col items-center text-center">
              <v-icon icon="mdi-pencil" size="large" color="blue-darken-1" class="mb-2"></v-icon>
              <h3 class="font-medium mb-2">Edit Survey</h3>
              <p class="text-sm text-gray-600 mb-4">
                Modify survey questions, settings, and appearance
              </p>
              <router-link 
                :to="`/surveys/${survey.id}/edit`"
                class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors mt-auto"
              >
                Edit
              </router-link>
            </div>
            
            <div class="border border-gray-200 rounded-lg p-4 flex flex-col items-center text-center">
              <v-icon icon="mdi-eye" size="large" color="purple-darken-1" class="mb-2"></v-icon>
              <h3 class="font-medium mb-2">Preview Survey</h3>
              <p class="text-sm text-gray-600 mb-4">
                See how your survey looks to respondents
              </p>
              <router-link 
                :to="`/surveys/${survey.id}/preview`"
                class="px-4 py-2 bg-purple-600 text-white rounded-md hover:bg-purple-700 transition-colors mt-auto"
              >
                Preview
              </router-link>
            </div>
            
            <div class="border border-gray-200 rounded-lg p-4 flex flex-col items-center text-center">
              <v-icon icon="mdi-chart-bar" size="large" color="green-darken-1" class="mb-2"></v-icon>
              <h3 class="font-medium mb-2">View Results</h3>
              <p class="text-sm text-gray-600 mb-4">
                Analyze responses and view survey analytics
              </p>
              <router-link 
                :to="`/surveys/${survey.id}/results`"
                class="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors mt-auto"
              >
                Results
              </router-link>
            </div>
          </div>
        </div>

        <!-- Share survey -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <h2 class="text-xl font-semibold text-gray-800 mb-4">Share Survey</h2>
          
          <div class="mb-6">
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Survey Link
            </label>
            <div class="flex items-center">
              <input
                type="text"
                :value="surveyUrl"
                readonly
                class="flex-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              />
              <button
                type="button"
                @click="copyToClipboard"
                class="ml-2 px-4 py-2 border border-gray-300 bg-white text-gray-700 rounded-md hover:bg-gray-50 transition-colors"
              >
                <v-icon icon="mdi-content-copy" size="small"></v-icon>
                {{ copied ? 'Copied!' : 'Copy' }}
              </button>
            </div>
          </div>
          
          <div>
            <h3 class="font-medium text-gray-800 mb-3">Share on</h3>
            <div class="flex space-x-2">
              <button class="p-2 bg-blue-600 text-white rounded-full hover:bg-blue-700 transition-colors">
                <v-icon icon="mdi-facebook" size="small"></v-icon>
              </button>
              <button class="p-2 bg-sky-500 text-white rounded-full hover:bg-sky-600 transition-colors">
                <v-icon icon="mdi-twitter" size="small"></v-icon>
              </button>
              <button class="p-2 bg-blue-700 text-white rounded-full hover:bg-blue-800 transition-colors">
                <v-icon icon="mdi-linkedin" size="small"></v-icon>
              </button>
              <button class="p-2 bg-green-600 text-white rounded-full hover:bg-green-700 transition-colors">
                <v-icon icon="mdi-email" size="small"></v-icon>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useSurveyStore } from '@/stores/survey';
import { useNotificationStore } from '@/stores/notification';

const route = useRoute();
const router = useRouter();
const surveyStore = useSurveyStore();
const notificationStore = useNotificationStore();

// State
const loading = ref(true);
const error = ref('');
const survey = ref({});
const publishLoading = ref(false);
const closeLoading = ref(false);
const copied = ref(false);

// Survey URL for sharing
const surveyUrl = computed(() => {
  const baseUrl = window.location.origin;
  return `${baseUrl}/surveys/${survey.value.id}/take`;
});

// Format date helper
const formatDate = (dateString) => {
  const options = { year: 'numeric', month: 'short', day: 'numeric' };
  return new Date(dateString).toLocaleDateString(undefined, options);
};

// Load survey details
onMounted(async () => {
  try {
    loading.value = true;
    const surveyId = route.params.id;
    const data = await surveyStore.fetchSurvey(surveyId);
    survey.value = data;
  } catch (err) {
    error.value = err.message || 'Failed to load survey. Please try again.';
    console.error('Error loading survey:', err);
  } finally {
    loading.value = false;
  }
});

// Publish survey
const publishSurvey = async () => {
  try {
    publishLoading.value = true;
    await surveyStore.updateSurveyStatus(survey.value.id, 'active');
    
    survey.value.status = 'active';
    
    notificationStore.showNotification({
      message: 'Survey published successfully',
      type: 'success'
    });
  } catch (error) {
    notificationStore.showNotification({
      message: error.message || 'Failed to publish survey',
      type: 'error'
    });
    console.error('Error publishing survey:', error);
  } finally {
    publishLoading.value = false;
  }
};

// Close survey
const closeSurvey = async () => {
  try {
    closeLoading.value = true;
    await surveyStore.updateSurveyStatus(survey.value.id, 'closed');
    
    survey.value.status = 'closed';
    
    notificationStore.showNotification({
      message: 'Survey closed successfully',
      type: 'success'
    });
  } catch (error) {
    notificationStore.showNotification({
      message: error.message || 'Failed to close survey',
      type: 'error'
    });
    console.error('Error closing survey:', error);
  } finally {
    closeLoading.value = false;
  }
};

// Copy survey URL to clipboard
const copyToClipboard = () => {
  navigator.clipboard.writeText(surveyUrl.value).then(() => {
    copied.value = true;
    
    // Reset "Copied!" text after 2 seconds
    setTimeout(() => {
      copied.value = false;
    }, 2000);
    
    notificationStore.showNotification({
      message: 'Survey URL copied to clipboard',
      type: 'success'
    });
  }).catch(err => {
    console.error('Failed to copy:', err);
    notificationStore.showNotification({
      message: 'Failed to copy URL',
      type: 'error'
    });
  });
};
</script> 