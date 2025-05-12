<template>
  <div>
    <v-btn
      :color="buttonColor"
      :variant="buttonVariant"
      @click="showShareDialog = true"
      :class="buttonClass"
    >
      <v-icon icon="mdi-share-variant" class="mr-1"></v-icon> {{ buttonText }}
    </v-btn>
    
    <v-dialog v-model="showShareDialog" max-width="500px">
      <v-card>
        <v-card-title class="bg-indigo-50 py-4">
          <v-icon icon="mdi-share-variant" class="mr-2"></v-icon>
          Share Survey
        </v-card-title>
        
        <v-card-text class="py-6">
          <p class="mb-4">Share this survey with potential respondents:</p>
          
          <div class="mb-6">
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Survey Link
            </label>
            <div class="flex items-center">
              <input
                type="text"
                :value="surveyUrl"
                readonly
                class="flex-1 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
              <button
                @click="copyToClipboard"
                class="ml-2 px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 transition-colors"
              >
                {{ copied ? 'Copied!' : 'Copy' }}
              </button>
            </div>
          </div>
          
          <div class="mb-6" v-if="showQrCode">
            <label class="block text-sm font-medium text-gray-700 mb-1">
              QR Code
            </label>
            <div class="bg-white p-4 border border-gray-300 rounded-md shadow-sm flex justify-center">
              <img :src="qrCodeUrl" alt="QR Code" class="h-48 w-48" />
            </div>
            <div class="flex justify-center mt-2">
              <button
                @click="downloadQrCode"
                class="text-sm text-indigo-600 hover:text-indigo-800 flex items-center"
              >
                <v-icon icon="mdi-download" size="small" class="mr-1"></v-icon> Download QR Code
              </button>
            </div>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Share via
            </label>
            <div class="flex justify-center space-x-6">
              <a :href="`https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(surveyUrl)}`" 
                 target="_blank" 
                 class="text-blue-600 hover:text-blue-800"
                 title="Share on Facebook"
              >
                <v-icon icon="mdi-facebook" size="large"></v-icon>
              </a>
              <a :href="`https://twitter.com/intent/tweet?url=${encodeURIComponent(surveyUrl)}&text=${encodeURIComponent('Please take my survey: ' + surveyTitle)}`" 
                 target="_blank" 
                 class="text-blue-400 hover:text-blue-600"
                 title="Share on Twitter"
              >
                <v-icon icon="mdi-twitter" size="large"></v-icon>
              </a>
              <a :href="`https://www.linkedin.com/shareArticle?mini=true&url=${encodeURIComponent(surveyUrl)}&title=${encodeURIComponent('Survey: ' + surveyTitle)}`" 
                 target="_blank" 
                 class="text-blue-700 hover:text-blue-900"
                 title="Share on LinkedIn"
              >
                <v-icon icon="mdi-linkedin" size="large"></v-icon>
              </a>
              <a :href="`mailto:?subject=${encodeURIComponent('Please take my survey: ' + surveyTitle)}&body=${encodeURIComponent('I created a survey and would appreciate your feedback:\n\n' + surveyUrl)}`" 
                 class="text-red-600 hover:text-red-800"
                 title="Share via Email"
              >
                <v-icon icon="mdi-email" size="large"></v-icon>
              </a>
              <a :href="`https://api.whatsapp.com/send?text=${encodeURIComponent('Please take my survey: ' + surveyTitle + '\n\n' + surveyUrl)}`" 
                 target="_blank" 
                 class="text-green-600 hover:text-green-800"
                 title="Share on WhatsApp"
              >
                <v-icon icon="mdi-whatsapp" size="large"></v-icon>
              </a>
            </div>
          </div>
        </v-card-text>
        
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="gray" variant="text" @click="showShareDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useNotificationStore } from '@/stores/notification';

const props = defineProps({
  surveyId: {
    type: String,
    required: true
  },
  surveyTitle: {
    type: String,
    default: 'Survey'
  },
  buttonText: {
    type: String,
    default: 'Share'
  },
  buttonColor: {
    type: String,
    default: 'primary'
  },
  buttonVariant: {
    type: String,
    default: 'flat'
  },
  buttonClass: {
    type: String,
    default: ''
  },
  showQrCode: {
    type: Boolean,
    default: true
  }
});

const showShareDialog = ref(false);
const copied = ref(false);
const notificationStore = useNotificationStore();

// Calculate survey URL for sharing
const surveyUrl = computed(() => {
  const baseUrl = window.location.origin;
  return `${baseUrl}/surveys/${props.surveyId}/take`;
});

// Generate QR code URL using an API service
const qrCodeUrl = computed(() => {
  return `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(surveyUrl.value)}`;
});

// Copy survey URL to clipboard
const copyToClipboard = () => {
  navigator.clipboard.writeText(surveyUrl.value).then(() => {
    copied.value = true;
    
    // Reset after 2 seconds
    setTimeout(() => {
      copied.value = false;
    }, 2000);
    
    notificationStore.showNotification({
      message: 'URL copied to clipboard',
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

// Download QR code
const downloadQrCode = () => {
  const link = document.createElement('a');
  link.href = qrCodeUrl.value;
  link.download = `survey-${props.surveyId}-qrcode.png`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  
  notificationStore.showNotification({
    message: 'QR code downloaded',
    type: 'success'
  });
};

// Reset copied state when dialog closes
watch(showShareDialog, (newValue) => {
  if (!newValue) {
    copied.value = false;
  }
});
</script> 