<template>
  <div>
    <!-- Response details dialog -->
    <v-dialog v-model="showDialog" max-width="800px">
      <v-card>
        <v-card-title class="bg-indigo-50 py-4 px-6">
          <div class="flex justify-between items-center w-full">
            <div class="flex items-center">
              <v-icon icon="mdi-clipboard-text" class="mr-2"></v-icon>
              <span class="text-xl">Response Details</span>
            </div>
            <v-btn icon size="small" variant="text" @click="closeDialog">
              <v-icon icon="mdi-close"></v-icon>
            </v-btn>
          </div>
        </v-card-title>
        
        <v-card-text class="py-6 px-6">
          <!-- Respondent info -->
          <div class="mb-6 pb-4 border-b border-gray-200">
            <h3 class="text-lg font-medium text-gray-800 mb-3">Respondent Information</h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <p class="text-sm text-gray-500">Name</p>
                <p class="font-medium">{{ response.respondentName || 'Anonymous' }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500">Email</p>
                <p class="font-medium">{{ response.respondentEmail || 'Not provided' }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500">Submitted On</p>
                <p class="font-medium">{{ formatDateTime(response.submittedAt) }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-500">Time to Complete</p>
                <p class="font-medium">{{ formatTime(response.timeToComplete) }}</p>
              </div>
            </div>
          </div>
          
          <!-- Response answers -->
          <h3 class="text-lg font-medium text-gray-800 mb-4">Responses</h3>
          
          <div v-for="(answer, index) in response.answers" :key="index" class="mb-6 pb-4 border-b border-gray-200 last:border-0">
            <div class="flex items-start">
              <span class="bg-indigo-100 text-indigo-800 rounded-full h-6 w-6 flex items-center justify-center text-sm font-medium mr-2 mt-1">
                {{ index + 1 }}
              </span>
              <div>
                <h4 class="text-md font-medium text-gray-900">{{ answer.questionText }}</h4>
                
                <!-- Different display based on question type -->
                <div class="mt-2 ml-1">
                  <!-- Text question -->
                  <div v-if="answer.questionType === 'text'" class="bg-gray-50 p-3 rounded-md">
                    <p v-if="answer.answer" class="text-gray-800">{{ answer.answer }}</p>
                    <p v-else class="text-gray-500 italic">No answer provided</p>
                  </div>
                  
                  <!-- Multiple choice / Dropdown -->
                  <div v-else-if="answer.questionType === 'multiple_choice' || answer.questionType === 'dropdown'">
                    <p v-if="answer.answer" class="text-gray-800">{{ answer.answer }}</p>
                    <p v-else class="text-gray-500 italic">No option selected</p>
                  </div>
                  
                  <!-- Checkbox (multiple answers) -->
                  <div v-else-if="answer.questionType === 'checkbox'">
                    <div v-if="answer.answer && answer.answer.length > 0">
                      <div v-for="(option, optIndex) in answer.answer" :key="optIndex" class="flex items-center mb-1">
                        <v-icon icon="mdi-check" size="small" class="mr-1 text-green-600"></v-icon>
                        <span>{{ option }}</span>
                      </div>
                    </div>
                    <p v-else class="text-gray-500 italic">No options selected</p>
                  </div>
                  
                  <!-- Rating -->
                  <div v-else-if="answer.questionType === 'rating'">
                    <div class="flex items-center">
                      <div class="flex">
                        <template v-for="n in 10" :key="n">
                          <v-icon
                            :icon="n <= parseInt(answer.answer) ? 'mdi-star' : 'mdi-star-outline'"
                            :color="n <= parseInt(answer.answer) ? 'amber' : 'gray'"
                            size="small"
                          ></v-icon>
                        </template>
                      </div>
                      <span class="ml-2 font-medium">{{ answer.answer || 'No rating' }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </v-card-text>
        
        <v-card-actions class="px-6 py-4 bg-gray-50">
          <v-spacer></v-spacer>
          <v-btn color="primary" variant="text" @click="closeDialog">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, defineProps, defineEmits } from 'vue';

const props = defineProps({
  response: {
    type: Object,
    required: true,
    default: () => ({
      id: '',
      surveyId: '',
      respondentName: '',
      respondentEmail: '',
      submittedAt: '',
      timeToComplete: 0,
      answers: []
    })
  },
  modelValue: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:modelValue']);

// Dialog visibility controlled by v-model
const showDialog = ref(false);

// Watch for changes from parent
const updateDialog = (value) => {
  showDialog.value = value;
  if (!value) {
    emit('update:modelValue', false);
  }
};

// Watch for v-model changes
const closeDialog = () => {
  showDialog.value = false;
  emit('update:modelValue', false);
};

// Format time (seconds to minutes:seconds)
const formatTime = (seconds) => {
  if (!seconds) return '0:00';
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = Math.floor(seconds % 60);
  return `${minutes}:${remainingSeconds < 10 ? '0' + remainingSeconds : remainingSeconds}`;
};

// Format date and time
const formatDateTime = (dateTimeStr) => {
  if (!dateTimeStr) return 'Unknown';
  const date = new Date(dateTimeStr);
  return date.toLocaleString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

// Update when props change
if (props.modelValue !== showDialog.value) {
  updateDialog(props.modelValue);
}
</script> 