<template>
  <div class="min-h-screen bg-gray-100 py-6">
    <div class="container mx-auto px-4 max-w-3xl">
      <!-- Loading state -->
      <div v-if="loading" class="bg-white rounded-lg shadow-md p-16 text-center">
        <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
        <p class="mt-4 text-gray-600">Loading survey preview...</p>
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

      <div v-else>
        <!-- Preview notice -->
        <div class="bg-amber-50 border border-amber-200 text-amber-800 px-4 py-3 rounded-md mb-6 flex items-center">
          <v-icon icon="mdi-eye" class="mr-2" size="small"></v-icon>
          <span>
            <strong>Preview Mode:</strong> This is how your survey will appear to respondents. Responses in preview mode are not recorded.
          </span>
        </div>
        
        <!-- Navigation buttons -->
        <div class="flex justify-between mb-6">
          <router-link 
            :to="`/surveys/${surveyId}`" 
            class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition-colors"
          >
            <v-icon icon="mdi-arrow-left" class="mr-1"></v-icon> Back to Survey
          </router-link>
          
          <router-link 
            v-if="survey.status === 'draft'"
            :to="`/surveys/${surveyId}/edit`" 
            class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
          >
            <v-icon icon="mdi-pencil" class="mr-1"></v-icon> Edit Survey
          </router-link>
        </div>

        <!-- Survey Preview Container -->
        <div class="bg-white rounded-lg shadow-md overflow-hidden">
          <!-- Survey Header -->
          <div class="bg-indigo-600 py-6 px-6 text-white">
            <h1 class="text-2xl font-bold mb-2">{{ survey.title }}</h1>
            <p v-if="survey.description" class="text-indigo-100">{{ survey.description }}</p>
          </div>
          
          <!-- Survey Form -->
          <form @submit.prevent="submitPreview" class="p-6">
            <!-- Questions -->
            <div 
              v-for="(question, index) in survey.questions" 
              :key="index"
              class="mb-8 py-4 border-b border-gray-200 last:border-0"
            >
              <div class="flex items-start mb-3">
                <span class="bg-indigo-100 text-indigo-800 rounded-full h-6 w-6 flex items-center justify-center text-sm font-medium mr-2 mt-1">
                  {{ index + 1 }}
                </span>
                <div class="flex-1">
                  <h3 class="text-lg font-medium text-gray-900 mb-1">
                    {{ question.text }}
                    <span v-if="question.required" class="text-red-500">*</span>
                  </h3>
                </div>
              </div>
              
              <!-- Question by type -->
              <!-- Text input -->
              <div v-if="question.type === 'text'" class="ml-8">
                <textarea
                  v-model="answers[index]"
                  rows="3"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  :placeholder="question.required ? 'Required' : 'Optional'"
                ></textarea>
              </div>
              
              <!-- Multiple choice -->
              <div v-else-if="question.type === 'multiple_choice'" class="ml-8 space-y-2">
                <div
                  v-for="(option, optIndex) in question.options"
                  :key="optIndex"
                  class="flex items-center"
                >
                  <input
                    type="radio"
                    :id="`question-${index}-option-${optIndex}`"
                    :name="`question-${index}`"
                    :value="option"
                    v-model="answers[index]"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300"
                  />
                  <label :for="`question-${index}-option-${optIndex}`" class="ml-3 block text-sm text-gray-700">
                    {{ option }}
                  </label>
                </div>
              </div>
              
              <!-- Checkbox -->
              <div v-else-if="question.type === 'checkbox'" class="ml-8 space-y-2">
                <div
                  v-for="(option, optIndex) in question.options"
                  :key="optIndex"
                  class="flex items-center"
                >
                  <input
                    type="checkbox"
                    :id="`question-${index}-option-${optIndex}`"
                    :value="option"
                    v-model="answers[index]"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label :for="`question-${index}-option-${optIndex}`" class="ml-3 block text-sm text-gray-700">
                    {{ option }}
                  </label>
                </div>
              </div>
              
              <!-- Dropdown -->
              <div v-else-if="question.type === 'dropdown'" class="ml-8">
                <select
                  v-model="answers[index]"
                  class="block w-full pl-3 pr-10 py-2 text-base border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                >
                  <option value="" disabled selected>Please select an option</option>
                  <option
                    v-for="(option, optIndex) in question.options"
                    :key="optIndex"
                    :value="option"
                  >
                    {{ option }}
                  </option>
                </select>
              </div>
              
              <!-- Rating scale -->
              <div v-else-if="question.type === 'rating'" class="ml-8">
                <div class="flex items-center">
                  <div class="text-sm text-gray-500 w-24">{{ question.lowLabel || 'Poor' }}</div>
                  <div class="flex-1 flex justify-center space-x-3">
                    <template v-for="n in parseInt(question.scale || 5)" :key="n">
                      <div class="flex flex-col items-center space-y-1">
                        <button
                          type="button"
                          class="w-10 h-10 rounded-full flex items-center justify-center focus:outline-none transition-colors"
                          :class="answers[index] === n ? 'bg-indigo-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
                          @click="answers[index] = n"
                        >
                          {{ n }}
                        </button>
                      </div>
                    </template>
                  </div>
                  <div class="text-sm text-gray-500 w-24 text-right">{{ question.highLabel || 'Excellent' }}</div>
                </div>
              </div>
            </div>
            
            <!-- Submit Button -->
            <div class="mt-8 flex justify-center">
              <v-btn
                type="submit"
                color="primary"
                size="large"
                :loading="submitting"
                :disabled="submitting"
              >
                Submit Survey (Preview)
              </v-btn>
            </div>
          </form>
        </div>
      </div>
    </div>
    
    <!-- Preview Submission Dialog -->
    <v-dialog v-model="previewDialog" max-width="600px">
      <v-card>
        <v-card-title class="text-h5 bg-indigo-50 py-4">
          <v-icon icon="mdi-information-outline" class="mr-2"></v-icon>
          Preview Mode
        </v-card-title>
        <v-card-text class="py-6">
          <p class="mb-4">
            This is a preview of the survey. In preview mode, your responses are not recorded and no data is submitted to the server.
          </p>
          <p class="mb-4">
            If this were a real submission, the form would be validated and responses would be recorded in the database.
          </p>
          <div class="bg-gray-50 p-4 rounded-md">
            <p class="text-sm font-medium text-gray-700 mb-2">Form validation summary:</p>
            <ul class="text-sm text-gray-600 space-y-1 ml-4 list-disc">
              <li v-if="hasRequiredFieldsMissing" class="text-red-600">
                Some required fields are missing
              </li>
              <li v-else class="text-green-600">
                All required fields are filled
              </li>
            </ul>
          </div>
        </v-card-text>
        <v-card-actions>
          <v-btn color="primary" variant="text" block @click="closePreviewDialog">
            Close Preview
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useSurveyStore } from '@/stores/survey';

const route = useRoute();
const router = useRouter();
const surveyStore = useSurveyStore();

const surveyId = route.params.id;
const loading = ref(true);
const error = ref('');
const submitting = ref(false);
const previewDialog = ref(false);

// Survey data
const survey = ref({
  title: '',
  description: '',
  questions: []
});

// Answers array - index corresponds to question index
const answers = reactive([]);

// Load survey data
onMounted(async () => {
  try {
    loading.value = true;
    error.value = '';
    
    const surveyData = await surveyStore.fetchSurveyById(surveyId);
    
    if (!surveyData) {
      throw new Error('Survey not found');
    }
    
    survey.value = surveyData;
    
    // Initialize answers array with appropriate default values
    surveyData.questions.forEach((question, index) => {
      if (question.type === 'checkbox') {
        answers[index] = [];
      } else {
        answers[index] = '';
      }
    });
    
  } catch (err) {
    error.value = err.message || 'Failed to load survey';
    console.error('Error loading survey preview:', err);
  } finally {
    loading.value = false;
  }
});

// Computed property to check if required fields are missing
const hasRequiredFieldsMissing = computed(() => {
  if (!survey.value.questions) return false;
  
  return survey.value.questions.some((question, index) => {
    if (!question.required) return false;
    
    if (question.type === 'checkbox') {
      return !answers[index] || answers[index].length === 0;
    } else {
      return !answers[index];
    }
  });
});

// Submit preview form
const submitPreview = () => {
  submitting.value = true;
  
  // Simulate network delay for a more realistic preview
  setTimeout(() => {
    submitting.value = false;
    previewDialog.value = true;
  }, 500);
};

// Close preview dialog
const closePreviewDialog = () => {
  previewDialog.value = false;
};
</script> 