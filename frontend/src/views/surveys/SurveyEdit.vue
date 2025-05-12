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

      <div v-else>
        <div class="flex justify-between items-center mb-6">
          <h1 class="text-3xl font-bold text-gray-900">Edit Survey</h1>
          <router-link 
            :to="`/surveys/${surveyId}`" 
            class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition-colors"
          >
            <v-icon icon="mdi-arrow-left" class="mr-1"></v-icon> Back to Survey
          </router-link>
        </div>

        <!-- Form -->
        <div class="bg-white rounded-lg shadow-md">
          <!-- Basic Info Section -->
          <div class="p-6 border-b border-gray-200">
            <h2 class="text-xl font-semibold text-gray-800 mb-4">Basic Information</h2>
            
            <div class="mb-4">
              <label for="title" class="block text-sm font-medium text-gray-700 mb-1">
                Survey Title *
              </label>
              <input
                id="title"
                v-model="survey.title"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                :class="{ 'border-red-500': v$.survey.title.$error }"
                placeholder="Enter a descriptive title for your survey"
              />
              <div v-if="v$.survey.title.$error" class="text-red-500 text-xs mt-1">
                {{ v$.survey.title.$errors[0].$message }}
              </div>
            </div>
            
            <div class="mb-4">
              <label for="description" class="block text-sm font-medium text-gray-700 mb-1">
                Description
              </label>
              <textarea
                id="description"
                v-model="survey.description"
                rows="3"
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                placeholder="Describe the purpose of your survey"
              ></textarea>
            </div>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label for="status" class="block text-sm font-medium text-gray-700 mb-1">
                  Status
                </label>
                <select
                  id="status"
                  v-model="survey.status"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  :disabled="survey.status === 'closed'"
                >
                  <option value="draft">Draft</option>
                  <option value="active">Active</option>
                  <option value="closed">Closed</option>
                </select>
                <p v-if="survey.status === 'closed'" class="text-xs text-red-500 mt-1">
                  Closed surveys cannot be reopened to prevent data inconsistency.
                </p>
                <p v-else class="text-xs text-gray-500 mt-1">
                  Draft surveys are not visible to participants. Active surveys can be taken by participants.
                </p>
              </div>
              
              <div>
                <label for="expiresAt" class="block text-sm font-medium text-gray-700 mb-1">
                  Expiration Date (Optional)
                </label>
                <input
                  id="expiresAt"
                  v-model="survey.expiresAt"
                  type="date"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                />
                <p class="text-xs text-gray-500 mt-1">
                  Survey will automatically close on this date. Leave empty for no expiration.
                </p>
              </div>
            </div>
          </div>

          <!-- Questions Section -->
          <div class="p-6">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-xl font-semibold text-gray-800">Survey Questions</h2>
              <div v-if="survey.status !== 'active' && survey.status !== 'closed'">
                <v-btn
                  color="primary"
                  variant="outlined"
                  size="small"
                  @click="showQuestionTypeMenu = !showQuestionTypeMenu"
                  class="relative"
                >
                  <v-icon icon="mdi-plus" class="mr-1"></v-icon> Add Question
                </v-btn>
                
                <!-- Question type dropdown menu -->
                <div 
                  v-if="showQuestionTypeMenu" 
                  class="absolute right-0 mt-2 w-56 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-10"
                  style="top: 40px"
                >
                  <div class="py-1" role="menu" aria-orientation="vertical" aria-labelledby="options-menu">
                    <a 
                      href="#" 
                      @click.prevent="addQuestion('text')" 
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" 
                      role="menuitem"
                    >
                      <v-icon icon="mdi-form-textbox" class="mr-2"></v-icon> Text Input
                    </a>
                    <a 
                      href="#" 
                      @click.prevent="addQuestion('multiple_choice')" 
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" 
                      role="menuitem"
                    >
                      <v-icon icon="mdi-checkbox-marked-circle-outline" class="mr-2"></v-icon> Multiple Choice
                    </a>
                    <a 
                      href="#" 
                      @click.prevent="addQuestion('checkbox')" 
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" 
                      role="menuitem"
                    >
                      <v-icon icon="mdi-checkbox-marked-outline" class="mr-2"></v-icon> Checkbox (Multiple Answers)
                    </a>
                    <a 
                      href="#" 
                      @click.prevent="addQuestion('dropdown')" 
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" 
                      role="menuitem"
                    >
                      <v-icon icon="mdi-menu-down" class="mr-2"></v-icon> Dropdown
                    </a>
                    <a 
                      href="#" 
                      @click.prevent="addQuestion('rating')" 
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" 
                      role="menuitem"
                    >
                      <v-icon icon="mdi-star" class="mr-2"></v-icon> Rating Scale
                    </a>
                  </div>
                </div>
              </div>
              <div v-else>
                <div class="bg-amber-50 text-amber-800 px-3 py-1 rounded-md text-sm">
                  <v-icon icon="mdi-information-outline" class="mr-1"></v-icon>
                  Active or closed surveys cannot add/remove questions
                </div>
              </div>
            </div>
            
            <div v-if="survey.questions.length === 0" class="bg-gray-50 p-8 text-center rounded-lg">
              <v-icon icon="mdi-help-circle-outline" size="large" color="gray" class="mb-2"></v-icon>
              <p class="text-gray-500">No questions added yet. Click "Add Question" to start building your survey.</p>
            </div>
            
            <!-- List of questions -->
            <draggable 
              v-model="survey.questions" 
              item-key="id"
              handle=".drag-handle"
              ghost-class="bg-indigo-100"
              @start="drag = true"
              @end="drag = false"
              :animation="200"
              :disabled="survey.status === 'active' || survey.status === 'closed'"
            >
              <template #item="{ element, index }">
                <div class="border border-gray-200 rounded-lg p-4 mb-4 bg-white">
                  <div class="flex justify-between items-start mb-3">
                    <div class="flex items-center flex-1">
                      <div 
                        class="drag-handle cursor-move px-2 text-gray-400 hover:text-gray-600"
                        :class="{ 'cursor-not-allowed': survey.status === 'active' || survey.status === 'closed' }"
                      >
                        <v-icon icon="mdi-drag" size="small"></v-icon>
                      </div>
                      
                      <span class="bg-gray-200 text-gray-700 rounded-full h-6 w-6 flex items-center justify-center text-sm font-medium mr-2">
                        {{ index + 1 }}
                      </span>
                      
                      <div class="flex-1">
                        <input
                          v-model="element.text"
                          type="text"
                          class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                          placeholder="Enter question"
                          :disabled="survey.status === 'active' || survey.status === 'closed'"
                        />
                      </div>
                    </div>
                    
                    <div class="flex space-x-2">
                      <v-btn
                        v-if="survey.status !== 'active' && survey.status !== 'closed'"
                        icon
                        variant="text"
                        size="small"
                        color="red"
                        @click="removeQuestion(index)"
                      >
                        <v-icon icon="mdi-delete" size="small"></v-icon>
                      </v-btn>
                    </div>
                  </div>
                  
                  <!-- Question type specific options -->
                  <div v-if="element.type === 'text'" class="ml-10">
                    <div class="flex items-center mt-2">
                      <input
                        v-model="element.required"
                        type="checkbox"
                        class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                        :disabled="survey.status === 'active' || survey.status === 'closed'"
                      />
                      <label class="ml-2 block text-sm text-gray-700">Required</label>
                    </div>
                  </div>
                  
                  <div v-else-if="element.type === 'multiple_choice' || element.type === 'checkbox' || element.type === 'dropdown'" class="ml-10">
                    <div class="mt-3">
                      <label class="block text-sm font-medium text-gray-700 mb-1">Options</label>
                      <div 
                        v-for="(option, optIndex) in element.options" 
                        :key="optIndex"
                        class="flex items-center mb-2"
                      >
                        <div class="w-6 text-center mr-2">
                          <v-icon 
                            v-if="element.type === 'multiple_choice'" 
                            icon="mdi-radiobox-blank" 
                            size="small"
                          ></v-icon>
                          <v-icon 
                            v-else-if="element.type === 'checkbox'" 
                            icon="mdi-checkbox-blank-outline" 
                            size="small"
                          ></v-icon>
                          <span v-else>{{ optIndex + 1 }}.</span>
                        </div>
                        <input
                          v-model="element.options[optIndex]"
                          type="text"
                          class="flex-1 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                          placeholder="Option text"
                          :disabled="survey.status === 'active' || survey.status === 'closed'"
                        />
                        <button
                          v-if="survey.status !== 'active' && survey.status !== 'closed'"
                          @click="removeOption(element, optIndex)"
                          class="ml-2 text-red-600 hover:text-red-800"
                        >
                          <v-icon icon="mdi-close" size="small"></v-icon>
                        </button>
                      </div>
                      
                      <button
                        v-if="survey.status !== 'active' && survey.status !== 'closed'"
                        @click="addOption(element)"
                        class="mt-2 text-indigo-600 hover:text-indigo-800 text-sm flex items-center"
                      >
                        <v-icon icon="mdi-plus" size="small" class="mr-1"></v-icon> Add Option
                      </button>
                    </div>
                    
                    <div class="flex items-center mt-2">
                      <input
                        v-model="element.required"
                        type="checkbox"
                        class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                        :disabled="survey.status === 'active' || survey.status === 'closed'"
                      />
                      <label class="ml-2 block text-sm text-gray-700">Required</label>
                    </div>
                  </div>
                  
                  <div v-else-if="element.type === 'rating'" class="ml-10">
                    <div class="flex items-center mt-3 space-x-4">
                      <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Scale</label>
                        <select
                          v-model="element.scale"
                          class="px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                          :disabled="survey.status === 'active' || survey.status === 'closed'"
                        >
                          <option value="5">1-5</option>
                          <option value="10">1-10</option>
                        </select>
                      </div>
                      
                      <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">Low Label</label>
                        <input
                          v-model="element.lowLabel"
                          type="text"
                          class="px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                          placeholder="e.g., Poor"
                          :disabled="survey.status === 'active' || survey.status === 'closed'"
                        />
                      </div>
                      
                      <div>
                        <label class="block text-sm font-medium text-gray-700 mb-1">High Label</label>
                        <input
                          v-model="element.highLabel"
                          type="text"
                          class="px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                          placeholder="e.g., Excellent"
                          :disabled="survey.status === 'active' || survey.status === 'closed'"
                        />
                      </div>
                    </div>
                    
                    <div class="flex items-center mt-3">
                      <input
                        v-model="element.required"
                        type="checkbox"
                        class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                        :disabled="survey.status === 'active' || survey.status === 'closed'"
                      />
                      <label class="ml-2 block text-sm text-gray-700">Required</label>
                    </div>
                  </div>
                </div>
              </template>
            </draggable>
          </div>
          
          <!-- Action Buttons -->
          <div class="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-end space-x-3">
            <v-btn
              variant="outlined"
              color="default"
              @click="$router.push(`/surveys/${surveyId}`)"
            >
              Cancel
            </v-btn>
            
            <v-btn
              v-if="survey.status !== 'closed'"
              color="primary"
              @click="updateSurvey"
              :loading="saving"
              :disabled="saving"
            >
              Save Changes
            </v-btn>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { v4 as uuidv4 } from 'uuid';
import { useVuelidate } from '@vuelidate/core';
import { required, minLength } from '@vuelidate/validators';
import { useSurveyStore } from '@/stores/survey';
import { useNotificationStore } from '@/stores/notification';
import draggable from 'vuedraggable';

const router = useRouter();
const route = useRoute();
const surveyId = route.params.id;
const surveyStore = useSurveyStore();
const notificationStore = useNotificationStore();

// State
const loading = ref(true);
const saving = ref(false);
const error = ref('');
const drag = ref(false);
const showQuestionTypeMenu = ref(false);

// Survey data - initialized with default values and updated by the fetchSurvey
const survey = reactive({
  id: '',
  title: '',
  description: '',
  status: 'draft',
  expiresAt: '',
  questions: [],
  createdAt: '',
  updatedAt: ''
});

// Fetch the survey data
onMounted(async () => {
  try {
    loading.value = true;
    error.value = '';
    
    const surveyData = await surveyStore.fetchSurveyById(surveyId);
    
    if (!surveyData) {
      throw new Error('Survey not found');
    }
    
    // Update our reactive survey object with fetched data
    Object.assign(survey, surveyData);
    
    // Make sure each question has an id for draggable functionality
    survey.questions = survey.questions.map(q => ({
      ...q,
      id: q.id || uuidv4() // Use existing id or generate a new one
    }));
    
    // Format expiration date as YYYY-MM-DD for input
    if (survey.expiresAt) {
      const date = new Date(survey.expiresAt);
      survey.expiresAt = date.toISOString().split('T')[0];
    }
  } catch (err) {
    error.value = err.message || 'Failed to load survey';
    console.error('Error loading survey:', err);
  } finally {
    loading.value = false;
  }
});

// Validation rules
const rules = {
  survey: {
    title: { required, minLength: minLength(3) }
  }
};

const v$ = useVuelidate(rules, { survey });

// Computed property to determine if questions can be edited
const canEditQuestions = computed(() => {
  return survey.status !== 'active' && survey.status !== 'closed';
});

// Methods
const addQuestion = (type) => {
  if (!canEditQuestions.value) return;
  
  const newQuestion = {
    id: uuidv4(),
    type,
    text: '',
    required: false
  };
  
  // Add type-specific properties
  if (type === 'multiple_choice' || type === 'checkbox' || type === 'dropdown') {
    newQuestion.options = ['', ''];
  } else if (type === 'rating') {
    newQuestion.scale = '5';
    newQuestion.lowLabel = 'Poor';
    newQuestion.highLabel = 'Excellent';
  }
  
  survey.questions.push(newQuestion);
  showQuestionTypeMenu.value = false;
};

const removeQuestion = (index) => {
  if (!canEditQuestions.value) return;
  survey.questions.splice(index, 1);
};

const addOption = (question) => {
  if (!canEditQuestions.value) return;
  question.options.push('');
};

const removeOption = (question, index) => {
  if (!canEditQuestions.value) return;
  
  if (question.options.length > 1) {
    question.options.splice(index, 1);
  } else {
    notificationStore.showNotification({
      message: 'A question must have at least one option',
      type: 'warning'
    });
  }
};

const updateSurvey = async () => {
  try {
    // Validate form
    const isValid = await v$.value.$validate();
    if (!isValid) {
      return;
    }
    
    // Additional validations
    if (survey.questions.length === 0) {
      notificationStore.showNotification({
        message: 'Please add at least one question to your survey',
        type: 'error'
      });
      return;
    }
    
    // Validate each question has required fields
    let isQuestionsValid = true;
    survey.questions.forEach((question, index) => {
      if (!question.text.trim()) {
        notificationStore.showNotification({
          message: `Question ${index + 1} is missing text`,
          type: 'error'
        });
        isQuestionsValid = false;
      }
      
      if ((question.type === 'multiple_choice' || question.type === 'checkbox' || question.type === 'dropdown') &&
          question.options.some(opt => !opt.trim())) {
        notificationStore.showNotification({
          message: `Question ${index + 1} has empty options`,
          type: 'error'
        });
        isQuestionsValid = false;
      }
    });
    
    if (!isQuestionsValid) {
      return;
    }
    
    saving.value = true;
    error.value = '';
    
    // Prepare questions data - strip client-side uuid which is only needed for frontend
    const questionsForSubmission = survey.questions.map(q => {
      const { id, ...questionData } = q;
      return questionData;
    });
    
    // Update survey with API
    const result = await surveyStore.updateSurvey(surveyId, {
      ...survey,
      questions: questionsForSubmission
    });
    
    if (result) {
      notificationStore.showNotification({
        message: 'Survey updated successfully',
        type: 'success'
      });
      
      router.push(`/surveys/${surveyId}`);
    } else {
      throw new Error('Failed to update survey');
    }
  } catch (err) {
    error.value = err.message || 'An error occurred while updating the survey';
    console.error('Error updating survey:', err);
  } finally {
    saving.value = false;
  }
};

// Close question type menu when clicking outside
onMounted(() => {
  document.addEventListener('click', (e) => {
    if (showQuestionTypeMenu.value && !e.target.closest('.relative')) {
      showQuestionTypeMenu.value = false;
    }
  });
});
</script> 