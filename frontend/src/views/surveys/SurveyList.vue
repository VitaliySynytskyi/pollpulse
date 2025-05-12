<template>
  <div class="min-h-screen bg-gray-50 py-6">
    <div class="container mx-auto px-4">
      <div class="flex justify-between items-center mb-6">
        <h1 class="text-3xl font-bold text-gray-900">My Surveys</h1>
        <router-link 
          to="/surveys/create" 
          class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 transition-colors"
        >
          <v-icon icon="mdi-plus" class="mr-1"></v-icon> New Survey
        </router-link>
      </div>

      <!-- Filters -->
      <div class="bg-white rounded-lg shadow-md p-6 mb-6">
        <div class="flex flex-col md:flex-row gap-4">
          <div class="flex-1">
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <v-icon icon="mdi-magnify" class="text-gray-400"></v-icon>
              </div>
              <input
                type="text"
                v-model="searchQuery"
                placeholder="Search surveys..."
                class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              />
            </div>
          </div>
          
          <div class="md:w-1/4">
            <select
              v-model="statusFilter"
              class="block w-full pl-3 pr-10 py-2 text-base border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            >
              <option value="all">All Status</option>
              <option value="draft">Draft</option>
              <option value="active">Active</option>
              <option value="closed">Closed</option>
            </select>
          </div>
          
          <div class="md:w-1/4">
            <select
              v-model="sortBy"
              class="block w-full pl-3 pr-10 py-2 text-base border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            >
              <option value="created_desc">Newest First</option>
              <option value="created_asc">Oldest First</option>
              <option value="title_asc">Title (A-Z)</option>
              <option value="title_desc">Title (Z-A)</option>
              <option value="responses_desc">Most Responses</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="bg-white rounded-lg shadow-md p-16 text-center">
        <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
        <p class="mt-4 text-gray-600">Loading your surveys...</p>
      </div>

      <!-- Empty state -->
      <div v-else-if="filteredSurveys.length === 0" class="bg-white rounded-lg shadow-md p-16 text-center">
        <v-icon icon="mdi-poll-box" size="x-large" color="gray" class="mb-4"></v-icon>
        <h3 class="text-xl font-medium text-gray-900 mb-2">No surveys found</h3>
        <p class="text-gray-600 mb-6">
          {{ 
            searchQuery 
              ? `No surveys match your search "${searchQuery}".` 
              : 'You haven\'t created any surveys yet.' 
          }}
        </p>
        <router-link 
          to="/surveys/create" 
          class="px-6 py-3 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 transition-colors inline-block"
        >
          Create your first survey
        </router-link>
      </div>

      <!-- Survey list -->
      <div v-else>
        <div class="bg-white rounded-lg shadow-md overflow-hidden">
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
              <tr v-for="survey in paginatedSurveys" :key="survey.id">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <div class="text-sm font-medium text-gray-900">{{ survey.title }}</div>
                  </div>
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
                  <div class="flex justify-end space-x-2">
                    <router-link :to="`/surveys/${survey.id}`" class="text-indigo-600 hover:text-indigo-900">
                      <v-icon icon="mdi-eye" size="small"></v-icon>
                    </router-link>
                    
                    <router-link :to="`/surveys/${survey.id}/edit`" class="text-blue-600 hover:text-blue-900">
                      <v-icon icon="mdi-pencil" size="small"></v-icon>
                    </router-link>
                    
                    <router-link :to="`/surveys/${survey.id}/results`" class="text-green-600 hover:text-green-900">
                      <v-icon icon="mdi-chart-bar" size="small"></v-icon>
                    </router-link>
                    
                    <button 
                      @click="confirmDelete(survey)"
                      class="text-red-600 hover:text-red-900"
                    >
                      <v-icon icon="mdi-delete" size="small"></v-icon>
                    </button>
                    
                    <button
                      v-if="survey.status === 'draft'"
                      @click="publishSurvey(survey.id)"
                      class="text-purple-600 hover:text-purple-900"
                    >
                      <v-icon icon="mdi-publish" size="small"></v-icon>
                    </button>
                    
                    <button
                      v-if="survey.status === 'active'"
                      @click="closeSurvey(survey.id)"
                      class="text-orange-600 hover:text-orange-900"
                    >
                      <v-icon icon="mdi-close-circle" size="small"></v-icon>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <!-- Pagination -->
        <div class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6 mt-4 rounded-lg shadow-md">
          <div class="flex-1 flex justify-between sm:hidden">
            <button
              @click="currentPage--"
              :disabled="currentPage === 1"
              class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
              :class="{ 'opacity-50 cursor-not-allowed': currentPage === 1 }"
            >
              Previous
            </button>
            <button
              @click="currentPage++"
              :disabled="currentPage >= totalPages"
              class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
              :class="{ 'opacity-50 cursor-not-allowed': currentPage >= totalPages }"
            >
              Next
            </button>
          </div>
          <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
            <div>
              <p class="text-sm text-gray-700">
                Showing
                <span class="font-medium">{{ (currentPage - 1) * pageSize + 1 }}</span>
                to
                <span class="font-medium">{{ Math.min(currentPage * pageSize, filteredSurveys.length) }}</span>
                of
                <span class="font-medium">{{ filteredSurveys.length }}</span>
                results
              </p>
            </div>
            <div>
              <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
                <button
                  @click="currentPage--"
                  :disabled="currentPage === 1"
                  class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
                  :class="{ 'opacity-50 cursor-not-allowed': currentPage === 1 }"
                >
                  <span class="sr-only">Previous</span>
                  <v-icon icon="mdi-chevron-left" size="small"></v-icon>
                </button>
                
                <button
                  v-for="page in displayedPages"
                  :key="page"
                  @click="currentPage = page"
                  :class="[
                    page === currentPage ? 'z-10 bg-indigo-50 border-indigo-500 text-indigo-600' : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50',
                    'relative inline-flex items-center px-4 py-2 border text-sm font-medium'
                  ]"
                >
                  {{ page }}
                </button>
                
                <button
                  @click="currentPage++"
                  :disabled="currentPage >= totalPages"
                  class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
                  :class="{ 'opacity-50 cursor-not-allowed': currentPage >= totalPages }"
                >
                  <span class="sr-only">Next</span>
                  <v-icon icon="mdi-chevron-right" size="small"></v-icon>
                </button>
              </nav>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="deleteDialog" max-width="500px">
      <v-card>
        <v-card-title class="text-h5">Delete Survey</v-card-title>
        <v-card-text>
          Are you sure you want to delete the survey "{{ surveyToDelete?.title }}"? This action cannot be undone and all responses will be permanently deleted.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue-darken-1" variant="text" @click="deleteDialog = false">Cancel</v-btn>
          <v-btn color="red-darken-1" variant="text" @click="deleteSurvey" :loading="deleteLoading">Delete</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useSurveyStore } from '@/stores/survey';
import { useNotificationStore } from '@/stores/notification';

const surveyStore = useSurveyStore();
const notificationStore = useNotificationStore();

// State
const loading = ref(true);
const surveys = ref([]);
const searchQuery = ref('');
const statusFilter = ref('all');
const sortBy = ref('created_desc');
const currentPage = ref(1);
const pageSize = ref(10);
const deleteDialog = ref(false);
const surveyToDelete = ref(null);
const deleteLoading = ref(false);

// Format date helper
const formatDate = (dateString) => {
  const options = { year: 'numeric', month: 'short', day: 'numeric' };
  return new Date(dateString).toLocaleDateString(undefined, options);
};

// Filter surveys based on search and status
const filteredSurveys = computed(() => {
  let result = [...surveys.value];
  
  // Apply search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter(survey => 
      survey.title.toLowerCase().includes(query) ||
      survey.description?.toLowerCase().includes(query)
    );
  }
  
  // Apply status filter
  if (statusFilter.value !== 'all') {
    result = result.filter(survey => survey.status === statusFilter.value);
  }
  
  // Apply sorting
  switch(sortBy.value) {
    case 'created_desc':
      result.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
      break;
    case 'created_asc':
      result.sort((a, b) => new Date(a.createdAt) - new Date(b.createdAt));
      break;
    case 'title_asc':
      result.sort((a, b) => a.title.localeCompare(b.title));
      break;
    case 'title_desc':
      result.sort((a, b) => b.title.localeCompare(a.title));
      break;
    case 'responses_desc':
      result.sort((a, b) => b.responseCount - a.responseCount);
      break;
  }
  
  return result;
});

// Calculate total pages for pagination
const totalPages = computed(() => 
  Math.ceil(filteredSurveys.value.length / pageSize.value)
);

// Get current page of surveys
const paginatedSurveys = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return filteredSurveys.value.slice(start, end);
});

// Calculate which page numbers to display
const displayedPages = computed(() => {
  const total = totalPages.value;
  const current = currentPage.value;
  const pages = [];
  
  if (total <= 7) {
    // If there are 7 or fewer pages, show all
    for (let i = 1; i <= total; i++) {
      pages.push(i);
    }
  } else {
    // Always include first page
    pages.push(1);
    
    // Add ellipsis if needed
    if (current > 3) {
      pages.push('...');
    }
    
    // Add pages around current page
    const start = Math.max(2, current - 1);
    const end = Math.min(total - 1, current + 1);
    
    for (let i = start; i <= end; i++) {
      pages.push(i);
    }
    
    // Add ellipsis if needed
    if (current < total - 2) {
      pages.push('...');
    }
    
    // Always include last page
    pages.push(total);
  }
  
  return pages;
});

// Reset to first page when filters change
watch([searchQuery, statusFilter, sortBy], () => {
  currentPage.value = 1;
});

// Fetch surveys on component mount
onMounted(async () => {
  try {
    loading.value = true;
    const response = await surveyStore.fetchSurveys();
    surveys.value = response.surveys;
  } catch (error) {
    notificationStore.showNotification({
      message: 'Failed to load surveys. Please try again.',
      type: 'error'
    });
    console.error('Error fetching surveys:', error);
  } finally {
    loading.value = false;
  }
});

// Open delete confirmation dialog
const confirmDelete = (survey) => {
  surveyToDelete.value = survey;
  deleteDialog.value = true;
};

// Delete survey
const deleteSurvey = async () => {
  try {
    deleteLoading.value = true;
    await surveyStore.deleteSurvey(surveyToDelete.value.id);
    
    // Remove from local state
    surveys.value = surveys.value.filter(s => s.id !== surveyToDelete.value.id);
    
    notificationStore.showNotification({
      message: 'Survey deleted successfully',
      type: 'success'
    });
  } catch (error) {
    notificationStore.showNotification({
      message: 'Failed to delete survey. Please try again.',
      type: 'error'
    });
    console.error('Error deleting survey:', error);
  } finally {
    deleteLoading.value = false;
    deleteDialog.value = false;
    surveyToDelete.value = null;
  }
};

// Publish a draft survey
const publishSurvey = async (id) => {
  try {
    await surveyStore.updateSurveyStatus(id, 'active');
    
    // Update local state
    surveys.value = surveys.value.map(s => {
      if (s.id === id) {
        return { ...s, status: 'active' };
      }
      return s;
    });
    
    notificationStore.showNotification({
      message: 'Survey published successfully',
      type: 'success'
    });
  } catch (error) {
    notificationStore.showNotification({
      message: 'Failed to publish survey. Please try again.',
      type: 'error'
    });
    console.error('Error publishing survey:', error);
  }
};

// Close an active survey
const closeSurvey = async (id) => {
  try {
    await surveyStore.updateSurveyStatus(id, 'closed');
    
    // Update local state
    surveys.value = surveys.value.map(s => {
      if (s.id === id) {
        return { ...s, status: 'closed' };
      }
      return s;
    });
    
    notificationStore.showNotification({
      message: 'Survey closed successfully',
      type: 'success'
    });
  } catch (error) {
    notificationStore.showNotification({
      message: 'Failed to close survey. Please try again.',
      type: 'error'
    });
    console.error('Error closing survey:', error);
  }
};
</script> 