<template>
  <div class="min-h-screen bg-gray-50 py-6">
    <div class="container mx-auto px-4">
      <!-- Loading state -->
      <div v-if="loading" class="bg-white rounded-lg shadow-md p-16 text-center">
        <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
        <p class="mt-4 text-gray-600">Loading survey results...</p>
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
          <h1 class="text-3xl font-bold text-gray-900">Survey Results</h1>
          <div class="flex space-x-2">
            <router-link 
              :to="`/surveys/${surveyId}`" 
              class="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition-colors"
            >
              <v-icon icon="mdi-arrow-left" class="mr-1"></v-icon> Back to Survey
            </router-link>
            
            <v-btn
              color="success"
              variant="outlined"
              @click="showExportOptions = true"
              :disabled="exporting"
            >
              <v-icon icon="mdi-file-export" class="mr-1"></v-icon> Export
            </v-btn>
          </div>
        </div>

        <!-- Survey info card -->
        <div class="bg-white rounded-lg shadow-md p-6 mb-6">
          <h2 class="text-2xl font-bold mb-2">{{ survey.title }}</h2>
          <p v-if="survey.description" class="text-gray-600 mb-4">{{ survey.description }}</p>
          
          <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div class="bg-gray-50 p-4 rounded-lg">
              <div class="text-gray-500 text-sm mb-1">Total Responses</div>
              <div class="text-2xl font-semibold">{{ resultsData.totalResponses || 0 }}</div>
            </div>
            
            <div class="bg-gray-50 p-4 rounded-lg">
              <div class="text-gray-500 text-sm mb-1">Completion Rate</div>
              <div class="text-2xl font-semibold">{{ resultsData.completionRate || 0 }}%</div>
            </div>
            
            <div class="bg-gray-50 p-4 rounded-lg">
              <div class="text-gray-500 text-sm mb-1">Avg. Time to Complete</div>
              <div class="text-2xl font-semibold">{{ formatTime(resultsData.avgTimeToComplete) }}</div>
            </div>
            
            <div class="bg-gray-50 p-4 rounded-lg">
              <div class="text-gray-500 text-sm mb-1">Status</div>
              <div class="text-2xl font-semibold flex items-center">
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
              </div>
            </div>
          </div>
        </div>

        <!-- Add this tab navigation after the survey info card -->
        <div class="bg-white rounded-lg shadow-md p-4 mb-6">
          <div class="flex border-b border-gray-200">
            <button
              @click="activeTab = 'summary'"
              class="px-4 py-2 font-medium text-sm"
              :class="activeTab === 'summary' ? 'text-indigo-600 border-b-2 border-indigo-600' : 'text-gray-500 hover:text-gray-700'"
            >
              Summary
            </button>
            <button
              @click="activeTab = 'responses'"
              class="px-4 py-2 font-medium text-sm"
              :class="activeTab === 'responses' ? 'text-indigo-600 border-b-2 border-indigo-600' : 'text-gray-500 hover:text-gray-700'"
            >
              Individual Responses
            </button>
          </div>
        </div>

        <!-- No responses yet -->
        <div v-if="resultsData.totalResponses === 0" class="bg-white rounded-lg shadow-md p-8 text-center">
          <v-icon icon="mdi-poll" size="x-large" color="gray" class="mb-4"></v-icon>
          <h2 class="text-xl font-medium text-gray-900 mb-4">No responses yet</h2>
          <p class="text-gray-500 mb-6">This survey hasn't received any responses yet.</p>
          <share-survey 
            :survey-id="surveyId"
            :survey-title="survey.title"
            button-text="Share Survey"
            button-color="primary"
            button-variant="flat"
          />
        </div>

        <!-- Results content -->
        <div v-else>
          <!-- Summary view -->
          <div v-if="activeTab === 'summary'">
            <!-- Response over time chart -->
            <div class="bg-white rounded-lg shadow-md p-6 mb-6">
              <h2 class="text-xl font-semibold text-gray-800 mb-4">Responses Over Time</h2>
              <div class="h-80">
                <line-chart 
                  :chart-data="responseTimeChartData"
                  :options="responseTimeChartOptions"
                />
              </div>
            </div>
            
            <!-- Questions and answers -->
            <div v-for="(question, index) in survey.questions" :key="index" class="bg-white rounded-lg shadow-md p-6 mb-6">
              <h2 class="text-xl font-semibold text-gray-800 mb-2">
                Question {{ index + 1 }}: {{ question.text }}
              </h2>
              
              <!-- Text responses -->
              <div v-if="question.type === 'text'" class="mt-4">
                <h3 class="text-md font-medium text-gray-700 mb-3">Text Responses ({{ getQuestionResponses(index).length }})</h3>
                
                <div v-if="getQuestionResponses(index).length === 0" class="text-gray-500 italic">
                  No responses yet
                </div>
                
                <div v-else class="space-y-2">
                  <div 
                    v-for="(response, respIndex) in getQuestionResponses(index)"
                    :key="respIndex"
                    class="bg-gray-50 p-3 rounded-md"
                  >
                    <p class="text-gray-800">{{ response }}</p>
                  </div>
                </div>
              </div>
              
              <!-- Charts for multiple choice, checkbox, dropdown, rating -->
              <div v-else class="mt-4">
                <div v-if="getQuestionResponses(index).length === 0" class="text-gray-500 italic">
                  No responses yet
                </div>
                <div v-else>
                  <div class="flex flex-col md:flex-row">
                    <div class="md:w-1/2">
                      <div class="h-64">
                        <pie-chart
                          v-if="question.type === 'multiple_choice' || question.type === 'dropdown'"
                          :chart-data="getChartData(index, question)"
                          :options="pieChartOptions"
                        />
                        <bar-chart
                          v-else-if="question.type === 'checkbox' || question.type === 'rating'"
                          :chart-data="getChartData(index, question)"
                          :options="barChartOptions"
                        />
                      </div>
                    </div>
                    <div class="md:w-1/2 mt-4 md:mt-0 md:pl-6">
                      <table class="min-w-full divide-y divide-gray-200">
                        <thead class="bg-gray-50">
                          <tr>
                            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                              Option
                            </th>
                            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                              Count
                            </th>
                            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                              Percentage
                            </th>
                          </tr>
                        </thead>
                        <tbody class="bg-white divide-y divide-gray-200">
                          <tr v-for="(count, option) in getResponseCounts(index, question)" :key="option">
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                              {{ option }}
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              {{ count }}
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              {{ calculatePercentage(count, getQuestionResponses(index).length) }}%
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          <!-- Individual responses view -->
          <div v-else-if="activeTab === 'responses'" class="bg-white rounded-lg shadow-md overflow-hidden">
            <div class="p-6 border-b border-gray-200">
              <h2 class="text-xl font-semibold text-gray-800 mb-4">Individual Responses</h2>
              
              <!-- Replace the filter and table with ResponsesTable component -->
              <responses-table 
                :responses="resultsData.responses"
                @view-response="viewResponseDetails"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Response Details Dialog -->
    <response-details 
      v-model="showResponseDetails"
      :response="selectedResponse"
    ></response-details>

    <!-- Export options dialog -->
    <v-dialog v-model="showExportOptions" max-width="400px">
      <v-card>
        <v-card-title class="bg-indigo-50 py-4">
          <v-icon icon="mdi-file-export" class="mr-2"></v-icon>
          Export Survey Results
        </v-card-title>
        
        <v-card-text class="py-6">
          <p class="mb-4">Select the data you want to export:</p>
          
          <div class="space-y-4">
            <div class="flex items-center p-3 border border-gray-200 rounded-md hover:bg-gray-50 cursor-pointer"
                 @click="exportResults('summary')">
              <v-icon icon="mdi-chart-bar" color="indigo" class="mr-3"></v-icon>
              <div>
                <h3 class="font-medium">Summary Data</h3>
                <p class="text-sm text-gray-500">Question statistics and aggregate results</p>
              </div>
            </div>
            
            <div class="flex items-center p-3 border border-gray-200 rounded-md hover:bg-gray-50 cursor-pointer"
                 @click="exportResults('responses')">
              <v-icon icon="mdi-format-list-bulleted" color="indigo" class="mr-3"></v-icon>
              <div>
                <h3 class="font-medium">Individual Responses</h3>
                <p class="text-sm text-gray-500">Detailed data from each respondent</p>
              </div>
            </div>
            
            <div class="flex items-center p-3 border border-gray-200 rounded-md hover:bg-gray-50 cursor-pointer"
                 @click="exportResults('all')">
              <v-icon icon="mdi-database" color="indigo" class="mr-3"></v-icon>
              <div>
                <h3 class="font-medium">Complete Dataset</h3>
                <p class="text-sm text-gray-500">Both summary and individual response data</p>
              </div>
            </div>
          </div>
        </v-card-text>
        
        <v-card-actions class="px-6 py-3 bg-gray-50">
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="showExportOptions = false">Cancel</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useSurveyStore } from '@/stores/survey';
import { useResponseStore } from '@/stores/response';
import { useNotificationStore } from '@/stores/notification';
import { Line as LineChart, Bar as BarChart, Pie as PieChart } from 'vue-chartjs';
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, PointElement, LineElement, ArcElement } from 'chart.js';
import ResponseDetails from '@/components/surveys/ResponseDetails.vue';
import ShareSurvey from '@/components/ShareSurvey.vue';
import ResponsesTable from '@/components/surveys/ResponsesTable.vue';

// Register ChartJS components
ChartJS.register(
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  ArcElement
);

const route = useRoute();
const surveyStore = useSurveyStore();
const responseStore = useResponseStore();
const notificationStore = useNotificationStore();

const surveyId = route.params.id;
const loading = ref(true);
const error = ref('');
const exporting = ref(false);
const activeTab = ref('summary');
const showResponseDetails = ref(false);
const selectedResponse = ref({});
const showExportOptions = ref(false);

// Survey and results data
const survey = ref({
  id: '',
  title: '',
  description: '',
  status: '',
  questions: []
});

const resultsData = ref({
  totalResponses: 0,
  completionRate: 0,
  avgTimeToComplete: 0,
  responses: [],
  responsesByDate: {}
});

// Chart configurations
const pieChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'bottom'
    }
  }
};

const barChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      ticks: {
        precision: 0
      }
    }
  }
};

const responseTimeChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      ticks: {
        precision: 0
      }
    }
  }
};

// Responses over time chart data
const responseTimeChartData = computed(() => {
  if (!resultsData.value.responsesByDate) return { labels: [], datasets: [] };
  
  const dates = Object.keys(resultsData.value.responsesByDate).sort();
  const counts = dates.map(date => resultsData.value.responsesByDate[date]);
  
  return {
    labels: dates.map(date => formatDate(date)),
    datasets: [
      {
        label: 'Responses',
        backgroundColor: 'rgba(99, 102, 241, 0.2)',
        borderColor: 'rgba(99, 102, 241, 1)',
        data: counts,
        fill: true,
        tension: 0.4
      }
    ]
  };
});

// Get responses for a specific question
const getQuestionResponses = (questionIndex) => {
  if (!resultsData.value.responses) return [];
  
  return resultsData.value.responses
    .map(response => response.answers[questionIndex]?.answer)
    .filter(answer => answer !== undefined && answer !== null && answer !== '');
};

// Count responses for chart data
const getResponseCounts = (questionIndex, question) => {
  const responses = getQuestionResponses(questionIndex);
  const counts = {};
  
  if (question.type === 'checkbox') {
    // For checkbox, each response can have multiple selected options
    responses.forEach(response => {
      if (Array.isArray(response)) {
        response.forEach(option => {
          counts[option] = (counts[option] || 0) + 1;
        });
      }
    });
  } else if (question.type === 'rating') {
    // For rating, initialize all possible ratings
    const scale = parseInt(question.scale || 5);
    for (let i = 1; i <= scale; i++) {
      counts[i] = 0;
    }
    
    // Then count actual responses
    responses.forEach(response => {
      const rating = parseInt(response);
      if (!isNaN(rating)) {
        counts[rating] = (counts[rating] || 0) + 1;
      }
    });
  } else {
    // For multiple choice and dropdown
    responses.forEach(response => {
      counts[response] = (counts[response] || 0) + 1;
    });
  }
  
  return counts;
};

// Get chart data for a question
const getChartData = (questionIndex, question) => {
  const counts = getResponseCounts(questionIndex, question);
  const labels = Object.keys(counts);
  const data = Object.values(counts);
  
  // Generate random colors for chart segments
  const backgroundColors = labels.map(() => 
    `rgba(${Math.floor(Math.random() * 200)}, ${Math.floor(Math.random() * 200)}, ${Math.floor(Math.random() * 200)}, 0.6)`
  );
  
  return {
    labels,
    datasets: [
      {
        label: question.text,
        data,
        backgroundColor: backgroundColors,
        borderWidth: 1
      }
    ]
  };
};

// Calculate percentage
const calculatePercentage = (count, total) => {
  if (!total) return 0;
  return Math.round((count / total) * 100);
};

// Format time (seconds to minutes:seconds)
const formatTime = (seconds) => {
  if (!seconds) return '0:00';
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = Math.floor(seconds % 60);
  return `${minutes}:${remainingSeconds < 10 ? '0' + remainingSeconds : remainingSeconds}`;
};

// Format date with optional time
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

// Export results
const exportResults = async (exportType = 'summary') => {
  try {
    exporting.value = true;
    showExportOptions.value = false;
    
    // Prepare CSV data
    let csvContent = "data:text/csv;charset=utf-8,";
    
    // Add survey info as header
    csvContent += `Survey: ${survey.value.title}\r\n`;
    csvContent += `Total Responses: ${resultsData.value.totalResponses}\r\n\r\n`;
    
    // Export summary data
    if (exportType === 'summary' || exportType === 'all') {
      csvContent += "# SUMMARY DATA\r\n\r\n";
      
      // For each question, add the question and its responses
      survey.value.questions.forEach((question, qIndex) => {
        csvContent += `Question ${qIndex + 1}: ${question.text}\r\n`;
        
        if (question.type === 'text') {
          // For text questions, list all text responses
          csvContent += "Responses:\r\n";
          const responses = getQuestionResponses(qIndex);
          responses.forEach(response => {
            csvContent += `"${response.replace(/"/g, '""')}"\r\n`; // Escape quotes in CSV
          });
        } else {
          // For other question types, show counts and percentages
          csvContent += "Option,Count,Percentage\r\n";
          const counts = getResponseCounts(qIndex, question);
          const total = getQuestionResponses(qIndex).length;
          
          Object.entries(counts).forEach(([option, count]) => {
            const percentage = calculatePercentage(count, total);
            csvContent += `"${option.replace(/"/g, '""')}",${count},${percentage}%\r\n`;
          });
        }
        
        csvContent += "\r\n"; // Add blank line between questions
      });
    }
    
    // Export individual response data
    if (exportType === 'responses' || exportType === 'all') {
      if (exportType === 'all') {
        csvContent += "\r\n# INDIVIDUAL RESPONSES\r\n\r\n";
      }
      
      // Create headers for the responses table
      let headers = ['Respondent ID', 'Name', 'Email', 'Submitted At', 'Time to Complete (sec)'];
      
      // Add headers for each question
      survey.value.questions.forEach((question, idx) => {
        headers.push(`Q${idx + 1}: ${question.text}`);
      });
      
      csvContent += headers.join(',') + '\r\n';
      
      // Add each response as a row
      resultsData.value.responses.forEach(response => {
        let row = [
          response.id || '',
          `"${(response.respondentName || 'Anonymous').replace(/"/g, '""')}"`,
          `"${(response.respondentEmail || '').replace(/"/g, '""')}"`,
          response.submittedAt ? new Date(response.submittedAt).toISOString() : '',
          response.timeToComplete || ''
        ];
        
        // Add answer for each question
        survey.value.questions.forEach((question, idx) => {
          const answer = response.answers[idx]?.answer;
          
          if (Array.isArray(answer)) {
            // Handle array answers (checkboxes)
            row.push(`"${answer.join(', ').replace(/"/g, '""')}"`);
          } else if (answer !== null && answer !== undefined) {
            // Handle scalar answers
            row.push(`"${String(answer).replace(/"/g, '""')}"`);
          } else {
            // No answer
            row.push('');
          }
        });
        
        csvContent += row.join(',') + '\r\n';
      });
    }
    
    // Create download link
    const encodedUri = encodeURI(csvContent);
    const link = document.createElement("a");
    link.setAttribute("href", encodedUri);
    link.setAttribute("download", `survey-results-${surveyId}-${exportType}.csv`);
    document.body.appendChild(link);
    
    // Trigger download
    link.click();
    
    // Clean up
    document.body.removeChild(link);
    
    notificationStore.showNotification({
      message: 'Survey results exported successfully',
      type: 'success'
    });
  } catch (err) {
    console.error('Error exporting results:', err);
    notificationStore.showNotification({
      message: 'Failed to export survey results',
      type: 'error'
    });
  } finally {
    exporting.value = false;
  }
};

// View response details
const viewResponseDetails = (response) => {
  selectedResponse.value = response;
  showResponseDetails.value = true;
};

// Load survey and results data
onMounted(async () => {
  try {
    loading.value = true;
    error.value = '';
    
    // Load survey data
    const surveyData = await surveyStore.fetchSurveyById(surveyId);
    
    if (!surveyData) {
      throw new Error('Survey not found');
    }
    
    survey.value = surveyData;
    
    // Load results data
    const results = await responseStore.fetchSurveyResults(surveyId);
    
    if (results) {
      resultsData.value = results;
    }
    
  } catch (err) {
    error.value = err.message || 'Failed to load survey results';
    console.error('Error loading survey results:', err);
  } finally {
    loading.value = false;
  }
});
</script> 