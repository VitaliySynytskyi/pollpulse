<template>
  <div class="min-h-screen bg-gray-100 py-6">
    <div class="container mx-auto px-4 max-w-3xl">
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
          <router-link to="/" class="text-indigo-600 hover:text-indigo-800">
            &larr; Return to Home
          </router-link>
        </div>
      </div>

      <!-- Survey already completed -->
      <div v-else-if="completed" class="bg-white rounded-lg shadow-md p-8 text-center">
        <v-icon icon="mdi-check-circle" size="x-large" color="success" class="mb-4"></v-icon>
        <h2 class="text-2xl font-bold text-gray-900 mb-4">Thank You!</h2>
        <p class="text-gray-600 mb-6">
          Your responses have been recorded. We appreciate your participation in this survey.
        </p>
        <router-link 
          to="/" 
          class="px-6 py-3 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 transition-colors inline-block"
        >
          Return to Home
        </router-link>
      </div>

      <!-- Survey closed or expired -->
      <div v-else-if="survey.status === 'closed' || isExpired" class="bg-white rounded-lg shadow-md p-8 text-center">
        <v-icon icon="mdi-lock-outline" size="x-large" color="warning" class="mb-4"></v-icon>
        <h2 class="text-2xl font-bold text-gray-900 mb-4">Survey Unavailable</h2>
        <p class="text-gray-600 mb-6">
          This survey is no longer accepting responses. It has been closed or has expired.
        </p>
        <router-link 
          to="/" 
          class="px-6 py-3 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 transition-colors inline-block"
        >
          Return to Home
        </router-link>
      </div>

      <!-- Survey content -->
      <div v-else>
        <!-- Survey Container -->
        <div class="bg-white rounded-lg shadow-md overflow-hidden">
          <!-- Progress bar -->
          <div class="bg-gray-50 h-2">
            <div 
              class="bg-indigo-600 h-2 transition-all duration-300 ease-in-out"
              :style="`width: ${progress}%`"
            ></div>
          </div>
          
          <!-- Survey Header -->
          <div class="bg-indigo-600 py-6 px-6 text-white">
            <h1 class="text-2xl font-bold mb-2">{{ survey.title }}</h1>
            <p v-if="survey.description" class="text-indigo-100">{{ survey.description }}</p>
          </div>
          
          <!-- Survey Form -->
          <form @submit.prevent="submitSurvey" class="p-6">
            <div class="flex justify-between items-center mb-4 text-sm text-gray-500">
              <span>{{ currentStep }} of {{ totalSteps }}</span>
              <span>
                <span class="text-red-500">*</span> Required
              </span>
            </div>

            <div v-if="currentStep === 1" class="mb-6">
              <!-- Optional: Collect respondent info -->
              <h2 class="text-lg font-semibold text-gray-800 mb-4">Your Information (Optional)</h2>
              
              <div class="mb-4">
                <label for="respondentName" class="block text-sm font-medium text-gray-700 mb-1">
                  Name
                </label>
                <input
                  id="respondentName"
                  v-model="respondent.name"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  placeholder="Your name (optional)"
                />
              </div>
              
              <div class="mb-4">
                <label for="respondentEmail" class="block text-sm font-medium text-gray-700 mb-1">
                  Email
                </label>
                <input
                  id="respondentEmail"
                  v-model="respondent.email"
                  type="email"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  placeholder="Your email (optional)"
                />
                <p v-if="v$.respondent.email.$error" class="mt-1 text-sm text-red-600">
                  {{ v$.respondent.email.$errors[0].$message }}
                </p>
              </div>
            </div>
            
            <!-- Questions by group -->
            <div v-else>
              <div v-for="question in currentQuestions" :key="question.id" class="mb-8">
                <div class="flex items-start mb-3">
                  <span class="bg-indigo-100 text-indigo-800 rounded-full h-6 w-6 flex items-center justify-center text-sm font-medium mr-2 mt-1">
                    {{ survey.questions.indexOf(question) + 1 }}
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
                    v-model="answers[survey.questions.indexOf(question)]"
                    rows="3"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                    :class="{ 'border-red-500': hasError(question) }"
                    :placeholder="question.required ? 'Required' : 'Optional'"
                  ></textarea>
                  <p v-if="hasError(question)" class="mt-1 text-sm text-red-600">
                    This field is required
                  </p>
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
                      :id="`question-${survey.questions.indexOf(question)}-option-${optIndex}`"
                      :name="`question-${survey.questions.indexOf(question)}`"
                      :value="option"
                      v-model="answers[survey.questions.indexOf(question)]"
                      :class="{ 'border-red-500': hasError(question) }"
                      class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300"
                    />
                    <label :for="`question-${survey.questions.indexOf(question)}-option-${optIndex}`" class="ml-3 block text-sm text-gray-700">
                      {{ option }}
                    </label>
                  </div>
                  <p v-if="hasError(question)" class="mt-1 text-sm text-red-600">
                    Please select an option
                  </p>
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
                      :id="`question-${survey.questions.indexOf(question)}-option-${optIndex}`"
                      :value="option"
                      v-model="answers[survey.questions.indexOf(question)]"
                      :class="{ 'border-red-500': hasError(question) }"
                      class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                    />
                    <label :for="`question-${survey.questions.indexOf(question)}-option-${optIndex}`" class="ml-3 block text-sm text-gray-700">
                      {{ option }}
                    </label>
                  </div>
                  <p v-if="hasError(question)" class="mt-1 text-sm text-red-600">
                    Please select at least one option
                  </p>
                </div>
                
                <!-- Dropdown -->
                <div v-else-if="question.type === 'dropdown'" class="ml-8">
                  <select
                    v-model="answers[survey.questions.indexOf(question)]"
                    :class="{ 'border-red-500': hasError(question) }"
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
                  <p v-if="hasError(question)" class="mt-1 text-sm text-red-600">
                    Please select an option
                  </p>
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
                            :class="[
                              answers[survey.questions.indexOf(question)] === n ? 'bg-indigo-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200',
                              { 'ring-2 ring-red-500': hasError(question) }
                            ]"
                            @click="answers[survey.questions.indexOf(question)] = n"
                          >
                            {{ n }}
                          </button>
                        </div>
                      </template>
                    </div>
                    <div class="text-sm text-gray-500 w-24 text-right">{{ question.highLabel || 'Excellent' }}</div>
                  </div>
                  <p v-if="hasError(question)" class="mt-1 text-sm text-red-600 text-center">
                    Please select a rating
                  </p>
                </div>
              </div>
            </div>
            
            <!-- Navigation/Submit Buttons -->
            <div class="mt-8 flex justify-between">
              <v-btn
                v-if="currentStep > 1"
                variant="outlined"
                color="default"
                @click="prevStep"
              >
                <v-icon icon="mdi-arrow-left" class="mr-1"></v-icon> Back
              </v-btn>
              <div v-else></div>
              
              <v-btn
                v-if="currentStep < totalSteps"
                color="primary"
                @click="nextStep"
              >
                Next <v-icon icon="mdi-arrow-right" class="ml-1"></v-icon>
              </v-btn>
              
              <v-btn
                v-else
                type="submit"
                color="success"
                :loading="submitting"
                :disabled="submitting"
              >
                Submit Survey
              </v-btn>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useVuelidate } from '@vuelidate/core';
import { email } from '@vuelidate/validators';
import { useSurveyStore } from '@/stores/survey';
import { useResponseStore } from '@/stores/response';
import { useNotificationStore } from '@/stores/notification';

const route = useRoute();
const router = useRouter();
const surveyStore = useSurveyStore();
const responseStore = useResponseStore();
const notificationStore = useNotificationStore();

const surveyId = route.params.id;
const loading = ref(true);
const error = ref('');
const submitting = ref(false);
const completed = ref(false);
const validationErrors = ref({});
const currentStep = ref(1);
const questionsPerStep = 5; // Number of questions to show per step

// Survey data
const survey = ref({
  id: '',
  title: '',
  description: '',
  status: '',
  expiresAt: null,
  questions: []
});

// Respondent information
const respondent = reactive({
  name: '',
  email: ''
});

// Validation rules
const rules = {
  respondent: {
    email: { email }
  }
};

const v$ = useVuelidate(rules, { respondent });

// Answers array - index corresponds to question index
const answers = reactive([]);

// Check if survey is expired
const isExpired = computed(() => {
  if (!survey.value.expiresAt) return false;
  const expiryDate = new Date(survey.value.expiresAt);
  return expiryDate < new Date();
});

// Calculate total steps (respondent info + questions divided into groups)
const totalSteps = computed(() => {
  if (!survey.value.questions || survey.value.questions.length === 0) return 1;
  return 1 + Math.ceil(survey.value.questions.length / questionsPerStep);
});

// Get current questions to display based on current step
const currentQuestions = computed(() => {
  if (currentStep.value === 1) return []; // First step is respondent info
  
  const questionStartIndex = (currentStep.value - 2) * questionsPerStep;
  const questionEndIndex = questionStartIndex + questionsPerStep;
  
  return survey.value.questions.slice(questionStartIndex, questionEndIndex);
});

// Calculate progress percentage
const progress = computed(() => {
  return ((currentStep.value - 1) / (totalSteps.value - 1)) * 100;
});

// Check if a question has a validation error
const hasError = (question) => {
  const index = survey.value.questions.indexOf(question);
  return validationErrors.value[index] === true;
};

// Load survey data
onMounted(async () => {
  try {
    loading.value = true;
    error.value = '';
    
    const surveyData = await surveyStore.fetchSurveyById(surveyId);
    
    if (!surveyData) {
      throw new Error('Survey not found');
    }
    
    // Check if survey is active
    if (surveyData.status !== 'active' && !isExpired.value) {
      error.value = 'This survey is not currently active.';
      loading.value = false;
      return;
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
    console.error('Error loading survey:', err);
  } finally {
    loading.value = false;
  }
});

// Move to next step
const nextStep = async () => {
  // If on respondent info step, validate email if provided
  if (currentStep.value === 1 && respondent.email) {
    const isEmailValid = await v$.value.respondent.email.$validate();
    if (!isEmailValid) return;
  }
  
  // Validate current questions before proceeding
  if (currentStep.value > 1) {
    // Clear previous errors
    validationErrors.value = {};
    
    // Check required questions
    let hasValidationErrors = false;
    
    for (const question of currentQuestions.value) {
      const index = survey.value.questions.indexOf(question);
      
      if (question.required) {
        let isValid = true;
        
        if (question.type === 'checkbox') {
          isValid = answers[index] && answers[index].length > 0;
        } else {
          isValid = !!answers[index];
        }
        
        if (!isValid) {
          validationErrors.value[index] = true;
          hasValidationErrors = true;
        }
      }
    }
    
    if (hasValidationErrors) {
      notificationStore.showNotification({
        message: 'Please fill in all required fields',
        type: 'error'
      });
      return;
    }
  }
  
  currentStep.value++;
};

// Move to previous step
const prevStep = () => {
  if (currentStep.value > 1) {
    currentStep.value--;
  }
};

// Submit the survey
const submitSurvey = async () => {
  try {
    // Validate email if provided
    if (respondent.email) {
      const isEmailValid = await v$.value.respondent.email.$validate();
      if (!isEmailValid) return;
    }
    
    // Validate required questions
    validationErrors.value = {};
    let hasValidationErrors = false;
    
    for (const question of survey.value.questions) {
      const index = survey.value.questions.indexOf(question);
      
      if (question.required) {
        let isValid = true;
        
        if (question.type === 'checkbox') {
          isValid = answers[index] && answers[index].length > 0;
        } else {
          isValid = !!answers[index];
        }
        
        if (!isValid) {
          validationErrors.value[index] = true;
          hasValidationErrors = true;
          
          // If the invalid question is not on current page, go to that page
          const questionStep = Math.floor(index / questionsPerStep) + 2;
          if (questionStep !== currentStep.value) {
            currentStep.value = questionStep;
          }
        }
      }
    }
    
    if (hasValidationErrors) {
      notificationStore.showNotification({
        message: 'Please fill in all required fields',
        type: 'error'
      });
      return;
    }
    
    // Prepare response data
    const responseData = {
      surveyId: survey.value.id,
      respondentName: respondent.name || null,
      respondentEmail: respondent.email || null,
      answers: survey.value.questions.map((question, index) => ({
        questionId: question.id || index.toString(), // Fall back to index if no ID
        questionText: question.text,
        questionType: question.type,
        answer: answers[index]
      }))
    };
    
    submitting.value = true;
    
    // Submit response
    const result = await responseStore.submitResponse(survey.value.id, responseData);
    
    if (result) {
      // Show success and reset form
      completed.value = true;
    } else {
      throw new Error('Failed to submit survey response');
    }
  } catch (err) {
    error.value = err.message || 'An error occurred while submitting your response';
    notificationStore.showNotification({
      message: error.value,
      type: 'error'
    });
    console.error('Error submitting survey:', err);
  } finally {
    submitting.value = false;
  }
};
</script> 