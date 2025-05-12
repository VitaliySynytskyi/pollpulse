<template>
  <div class="min-h-screen bg-gray-50 py-6">
    <div class="container mx-auto px-4 max-w-3xl">
      <h1 class="text-3xl font-bold text-gray-900 mb-6">My Profile</h1>

      <div class="bg-white rounded-lg shadow-md overflow-hidden">
        <div class="p-6 border-b border-gray-200">
          <h2 class="text-xl font-semibold text-gray-800 mb-4">Profile Information</h2>
          <p class="text-gray-600 mb-6">Update your personal information and email address</p>

          <form @submit.prevent="updateProfile">
            <div v-if="profileError" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md mb-4">
              {{ profileError }}
            </div>

            <div v-if="profileSuccess" class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-md mb-4">
              {{ profileSuccess }}
            </div>

            <div class="mb-4">
              <label for="name" class="block text-sm font-medium text-gray-700 mb-1">
                Full Name
              </label>
              <input
                id="name"
                v-model="name"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                :class="{ 'border-red-500': v$.name.$error }"
              />
              <div v-if="v$.name.$error" class="text-red-500 text-xs mt-1">
                {{ v$.name.$errors[0].$message }}
              </div>
            </div>

            <div class="mb-4">
              <label for="email" class="block text-sm font-medium text-gray-700 mb-1">
                Email
              </label>
              <input
                id="email"
                v-model="email"
                type="email"
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                :class="{ 'border-red-500': v$.email.$error }"
              />
              <div v-if="v$.email.$error" class="text-red-500 text-xs mt-1">
                {{ v$.email.$errors[0].$message }}
              </div>
            </div>

            <div class="mb-4">
              <label for="company" class="block text-sm font-medium text-gray-700 mb-1">
                Company/Organization (Optional)
              </label>
              <input
                id="company"
                v-model="company"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>

            <div class="flex justify-end">
              <v-btn
                type="submit"
                color="primary"
                :loading="profileSubmitting"
                :disabled="profileSubmitting"
              >
                Save Changes
              </v-btn>
            </div>
          </form>
        </div>

        <div class="p-6 border-b border-gray-200">
          <h2 class="text-xl font-semibold text-gray-800 mb-4">Change Password</h2>
          <p class="text-gray-600 mb-6">Update your password to maintain a secure account</p>

          <form @submit.prevent="updatePassword">
            <div v-if="passwordError" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md mb-4">
              {{ passwordError }}
            </div>

            <div v-if="passwordSuccess" class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-md mb-4">
              {{ passwordSuccess }}
            </div>

            <div class="mb-4">
              <label for="currentPassword" class="block text-sm font-medium text-gray-700 mb-1">
                Current Password
              </label>
              <input
                id="currentPassword"
                v-model="currentPassword"
                type="password"
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                :class="{ 'border-red-500': v$.currentPassword.$error }"
              />
              <div v-if="v$.currentPassword.$error" class="text-red-500 text-xs mt-1">
                {{ v$.currentPassword.$errors[0].$message }}
              </div>
            </div>

            <div class="mb-4">
              <label for="newPassword" class="block text-sm font-medium text-gray-700 mb-1">
                New Password
              </label>
              <input
                id="newPassword"
                v-model="newPassword"
                type="password"
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                :class="{ 'border-red-500': v$.newPassword.$error }"
              />
              <div v-if="v$.newPassword.$error" class="text-red-500 text-xs mt-1">
                {{ v$.newPassword.$errors[0].$message }}
              </div>
            </div>

            <div class="mb-4">
              <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">
                Confirm New Password
              </label>
              <input
                id="confirmPassword"
                v-model="confirmPassword"
                type="password"
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                :class="{ 'border-red-500': v$.confirmPassword.$error }"
              />
              <div v-if="v$.confirmPassword.$error" class="text-red-500 text-xs mt-1">
                {{ v$.confirmPassword.$errors[0].$message }}
              </div>
            </div>

            <div class="flex justify-end">
              <v-btn
                type="submit"
                color="primary"
                :loading="passwordSubmitting"
                :disabled="passwordSubmitting"
              >
                Update Password
              </v-btn>
            </div>
          </form>
        </div>

        <div class="p-6">
          <h2 class="text-xl font-semibold text-gray-800 mb-4">Account Settings</h2>
          <p class="text-gray-600 mb-6">Manage your account preferences and notifications</p>

          <div class="mb-6">
            <h3 class="text-lg font-medium text-gray-800 mb-3">Email Notifications</h3>
            
            <div class="space-y-3">
              <div class="flex items-start">
                <div class="flex items-center h-5">
                  <input
                    id="surveyResponses"
                    v-model="notifySurveyResponses"
                    type="checkbox"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                </div>
                <div class="ml-3 text-sm">
                  <label for="surveyResponses" class="font-medium text-gray-700">Survey responses</label>
                  <p class="text-gray-500">Receive an email when someone completes your survey</p>
                </div>
              </div>
              
              <div class="flex items-start">
                <div class="flex items-center h-5">
                  <input
                    id="weeklyDigest"
                    v-model="notifyWeeklyDigest"
                    type="checkbox"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                </div>
                <div class="ml-3 text-sm">
                  <label for="weeklyDigest" class="font-medium text-gray-700">Weekly digest</label>
                  <p class="text-gray-500">Get a weekly summary of your survey performance</p>
                </div>
              </div>
              
              <div class="flex items-start">
                <div class="flex items-center h-5">
                  <input
                    id="productUpdates"
                    v-model="notifyProductUpdates"
                    type="checkbox"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                </div>
                <div class="ml-3 text-sm">
                  <label for="productUpdates" class="font-medium text-gray-700">Product updates</label>
                  <p class="text-gray-500">Stay informed about new features and improvements</p>
                </div>
              </div>
            </div>
          </div>

          <div class="flex justify-end">
            <v-btn
              color="primary"
              @click="saveNotificationSettings"
              :loading="notificationsSubmitting"
              :disabled="notificationsSubmitting"
            >
              Save Settings
            </v-btn>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useVuelidate } from '@vuelidate/core';
import { required, email as emailValidator, minLength, sameAs, helpers } from '@vuelidate/validators';
import { useAuthStore } from '@/stores/auth';
import { useNotificationStore } from '@/stores/notification';

const authStore = useAuthStore();
const notificationStore = useNotificationStore();
const user = computed(() => authStore.user);

// Profile form data
const name = ref('');
const email = ref('');
const company = ref('');
const profileSubmitting = ref(false);
const profileError = ref('');
const profileSuccess = ref('');

// Password form data
const currentPassword = ref('');
const newPassword = ref('');
const confirmPassword = ref('');
const passwordSubmitting = ref(false);
const passwordError = ref('');
const passwordSuccess = ref('');

// Notification settings
const notifySurveyResponses = ref(false);
const notifyWeeklyDigest = ref(false);
const notifyProductUpdates = ref(false);
const notificationsSubmitting = ref(false);

// Load user data
onMounted(() => {
  if (user.value) {
    name.value = user.value.name || '';
    email.value = user.value.email || '';
    company.value = user.value.company || '';
    
    // Load notification preferences
    notifySurveyResponses.value = user.value.preferences?.notifySurveyResponses || false;
    notifyWeeklyDigest.value = user.value.preferences?.notifyWeeklyDigest || false;
    notifyProductUpdates.value = user.value.preferences?.notifyProductUpdates || false;
  }
});

// Validation rules
const profileRules = {
  name: { required: helpers.withMessage('Name is required', required) },
  email: { 
    required: helpers.withMessage('Email is required', required), 
    email: helpers.withMessage('Please enter a valid email address', emailValidator)
  }
};

const passwordRules = {
  currentPassword: { required: helpers.withMessage('Current password is required', required) },
  newPassword: { 
    required: helpers.withMessage('New password is required', required),
    minLength: helpers.withMessage(
      'Password must be at least 8 characters long', 
      minLength(8)
    ),
    containsUppercase: helpers.withMessage(
      'Password must contain at least one uppercase letter',
      (value) => /[A-Z]/.test(value)
    ),
    containsNumber: helpers.withMessage(
      'Password must contain at least one number',
      (value) => /[0-9]/.test(value)
    )
  },
  confirmPassword: { 
    required: helpers.withMessage('Please confirm your new password', required),
    sameAsPassword: helpers.withMessage(
      'Passwords must match', 
      sameAs(newPassword)
    )
  }
};

const v$ = useVuelidate(
  {
    ...profileRules,
    ...passwordRules
  }, 
  {
    name, 
    email, 
    currentPassword, 
    newPassword, 
    confirmPassword
  }
);

// Update profile information
const updateProfile = async () => {
  try {
    profileSuccess.value = '';
    profileError.value = '';
    
    // Validate only profile fields
    const profileV$ = useVuelidate(profileRules, { name, email });
    const isValid = await profileV$.value.$validate();
    if (!isValid) return;
    
    profileSubmitting.value = true;
    
    await authStore.updateProfile({
      name: name.value,
      email: email.value,
      company: company.value
    });
    
    profileSuccess.value = 'Profile updated successfully';
    
    // Update local user store
    authStore.updateUserData({
      ...user.value,
      name: name.value,
      email: email.value,
      company: company.value
    });
  } catch (error) {
    profileError.value = error.message || 'Failed to update profile';
  } finally {
    profileSubmitting.value = false;
  }
};

// Update password
const updatePassword = async () => {
  try {
    passwordSuccess.value = '';
    passwordError.value = '';
    
    // Validate only password fields
    const passwordV$ = useVuelidate(passwordRules, { 
      currentPassword, 
      newPassword, 
      confirmPassword 
    });
    
    const isValid = await passwordV$.value.$validate();
    if (!isValid) return;
    
    passwordSubmitting.value = true;
    
    await authStore.updatePassword({
      currentPassword: currentPassword.value,
      newPassword: newPassword.value
    });
    
    passwordSuccess.value = 'Password updated successfully';
    
    // Clear password fields
    currentPassword.value = '';
    newPassword.value = '';
    confirmPassword.value = '';
    
    // Clear validations
    passwordV$.value.$reset();
  } catch (error) {
    passwordError.value = error.message || 'Failed to update password';
  } finally {
    passwordSubmitting.value = false;
  }
};

// Save notification preferences
const saveNotificationSettings = async () => {
  try {
    notificationsSubmitting.value = true;
    
    const preferences = {
      notifySurveyResponses: notifySurveyResponses.value,
      notifyWeeklyDigest: notifyWeeklyDigest.value,
      notifyProductUpdates: notifyProductUpdates.value
    };
    
    await authStore.updatePreferences(preferences);
    
    // Update local user store
    authStore.updateUserData({
      ...user.value,
      preferences
    });
    
    notificationStore.showNotification({
      message: 'Notification settings saved successfully',
      type: 'success'
    });
  } catch (error) {
    notificationStore.showNotification({
      message: error.message || 'Failed to save notification settings',
      type: 'error'
    });
  } finally {
    notificationsSubmitting.value = false;
  }
};
</script> 