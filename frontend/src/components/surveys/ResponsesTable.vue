<template>
  <div>
    <!-- Filter and search controls -->
    <div class="flex flex-col md:flex-row justify-between mb-4 space-y-2 md:space-y-0">
      <div class="relative md:w-1/3">
        <input
          v-model="searchInput"
          type="text"
          placeholder="Search responses..."
          class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          @input="handleSearch"
        />
        <v-icon
          icon="mdi-magnify"
          class="absolute left-3 top-2.5 text-gray-400"
        ></v-icon>
      </div>
      
      <div class="flex space-x-2">
        <select
          v-model="sortingOption"
          class="px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          @change="handleSort"
        >
          <option value="date_desc">Newest First</option>
          <option value="date_asc">Oldest First</option>
          <option value="completion_time_asc">Fastest Completion</option>
          <option value="completion_time_desc">Slowest Completion</option>
        </select>
      </div>
    </div>
    
    <!-- Responses table -->
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Respondent
            </th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Submitted
            </th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Completion Time
            </th>
            <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="(response, index) in paginatedResponses" :key="index" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium text-gray-900">{{ response.respondentName || 'Anonymous' }}</div>
              <div v-if="response.respondentEmail" class="text-sm text-gray-500">{{ response.respondentEmail }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm text-gray-900">{{ formatDate(response.submittedAt, true) }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm text-gray-900">{{ formatTime(response.timeToComplete) }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <button 
                @click="viewDetails(response)"
                class="text-indigo-600 hover:text-indigo-900"
              >
                View Details
              </button>
            </td>
          </tr>
          
          <!-- No results message -->
          <tr v-if="paginatedResponses.length === 0">
            <td colspan="4" class="px-6 py-10 text-center text-gray-500">
              <v-icon icon="mdi-file-search-outline" size="large" class="mb-2"></v-icon>
              <p>No matching responses found</p>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <!-- Pagination controls -->
    <div v-if="filteredResponses.length > itemsPerPage" class="flex justify-between items-center mt-4 px-2">
      <div class="text-sm text-gray-600">
        Showing {{ paginationStart + 1 }}-{{ Math.min(paginationStart + itemsPerPage, filteredResponses.length) }} of {{ filteredResponses.length }} responses
      </div>
      <div class="flex space-x-1">
        <button 
          @click="prevPage" 
          :disabled="currentPage === 1" 
          class="px-3 py-1 rounded-md text-sm"
          :class="currentPage === 1 ? 'text-gray-400 cursor-not-allowed' : 'text-gray-700 hover:bg-gray-100'"
        >
          <v-icon icon="mdi-chevron-left" size="small"></v-icon>
        </button>
        <button 
          v-for="page in totalPages" 
          :key="page"
          @click="goToPage(page)"
          class="px-3 py-1 rounded-md text-sm"
          :class="currentPage === page ? 'bg-indigo-600 text-white' : 'text-gray-700 hover:bg-gray-100'"
        >
          {{ page }}
        </button>
        <button 
          @click="nextPage" 
          :disabled="currentPage === totalPages" 
          class="px-3 py-1 rounded-md text-sm"
          :class="currentPage === totalPages ? 'text-gray-400 cursor-not-allowed' : 'text-gray-700 hover:bg-gray-100'"
        >
          <v-icon icon="mdi-chevron-right" size="small"></v-icon>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';

const props = defineProps({
  responses: {
    type: Array,
    required: true,
    default: () => []
  }
});

const emit = defineEmits(['view-response']);

// Pagination
const currentPage = ref(1);
const itemsPerPage = 10;

// Search and sort
const searchInput = ref('');
const sortingOption = ref('date_desc');

// Filtered responses based on search input
const filteredResponses = computed(() => {
  if (!props.responses) return [];
  
  let filtered = [...props.responses];
  
  // Apply search filter
  if (searchInput.value) {
    const term = searchInput.value.toLowerCase();
    filtered = filtered.filter(response => {
      return (
        (response.respondentName && response.respondentName.toLowerCase().includes(term)) ||
        (response.respondentEmail && response.respondentEmail.toLowerCase().includes(term))
      );
    });
  }
  
  // Apply sorting
  filtered.sort((a, b) => {
    if (sortingOption.value === 'date_asc') {
      return new Date(a.submittedAt).getTime() - new Date(b.submittedAt).getTime();
    } else if (sortingOption.value === 'date_desc') {
      return new Date(b.submittedAt).getTime() - new Date(a.submittedAt).getTime();
    } else if (sortingOption.value === 'completion_time_asc') {
      return (a.timeToComplete || 0) - (b.timeToComplete || 0);
    } else if (sortingOption.value === 'completion_time_desc') {
      return (b.timeToComplete || 0) - (a.timeToComplete || 0);
    }
    return 0;
  });
  
  return filtered;
});

// Paginated responses
const paginationStart = computed(() => (currentPage.value - 1) * itemsPerPage);
const paginatedResponses = computed(() => {
  return filteredResponses.value.slice(paginationStart.value, paginationStart.value + itemsPerPage);
});

// Total pages
const totalPages = computed(() => {
  return Math.ceil(filteredResponses.value.length / itemsPerPage);
});

// Format date with time
const formatDate = (dateStr, includeTime = false) => {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  if (includeTime) {
    return date.toLocaleString(undefined, { 
      month: 'short', 
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
  return date.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
};

// Format time (seconds to minutes:seconds)
const formatTime = (seconds) => {
  if (!seconds) return '0:00';
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = Math.floor(seconds % 60);
  return `${minutes}:${remainingSeconds < 10 ? '0' + remainingSeconds : remainingSeconds}`;
};

// Pagination methods
const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
  }
};

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--;
  }
};

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page;
  }
};

// Search and sort handlers
const handleSearch = () => {
  currentPage.value = 1; // Reset to first page on search
};

const handleSort = () => {
  // No need to reset page on sort
};

// View response details
const viewDetails = (response) => {
  emit('view-response', response);
};

// Reset pagination when responses change
watch(() => props.responses, () => {
  currentPage.value = 1;
});
</script> 