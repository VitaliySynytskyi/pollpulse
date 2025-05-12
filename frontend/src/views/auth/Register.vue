<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 bg-white p-8 rounded-lg shadow-md">
      <div>
        <h2 class="mt-6 text-center text-3xl font-bold text-gray-900">
          Create your account
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Or
          <router-link to="/login" class="font-medium text-indigo-600 hover:text-indigo-500">
            sign in to your existing account
          </router-link>
        </p>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="register">
        <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
          {{ error }}
        </div>
        
        <div>
          <label for="name" class="block text-sm font-medium text-gray-700">
            Full Name
          </label>
          <div class="mt-1">
            <input
              id="name"
              v-model="name"
              name="name"
              type="text"
              autocomplete="name"
              required
              class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              :class="{ 'border-red-500': v$.name.$error }"
            />
            <div v-if="v$.name.$error" class="text-red-500 text-xs mt-1">
              {{ v$.name.$errors[0].$message }}
            </div>
          </div>
        </div>

        <div>
          <label for="email" class="block text-sm font-medium text-gray-700">
            Email address
          </label>
          <div class="mt-1">
            <input
              id="email"
              v-model="email"
              name="email"
              type="email"
              autocomplete="email"
              required
              class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              :class="{ 'border-red-500': v$.email.$error }"
            />
            <div v-if="v$.email.$error" class="text-red-500 text-xs mt-1">
              {{ v$.email.$errors[0].$message }}
            </div>
          </div>
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-gray-700">
            Password
          </label>
          <div class="mt-1">
            <input
              id="password"
              v-model="password"
              name="password"
              type="password"
              autocomplete="new-password"
              required
              class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              :class="{ 'border-red-500': v$.password.$error }"
            />
            <div v-if="v$.password.$error" class="text-red-500 text-xs mt-1">
              {{ v$.password.$errors[0].$message }}
            </div>
          </div>
        </div>

        <div>
          <label for="confirmPassword" class="block text-sm font-medium text-gray-700">
            Confirm Password
          </label>
          <div class="mt-1">
            <input
              id="confirmPassword"
              v-model="confirmPassword"
              name="confirmPassword"
              type="password"
              autocomplete="new-password"
              required
              class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              :class="{ 'border-red-500': v$.confirmPassword.$error }"
            />
            <div v-if="v$.confirmPassword.$error" class="text-red-500 text-xs mt-1">
              {{ v$.confirmPassword.$errors[0].$message }}
            </div>
          </div>
        </div>

        <div class="flex items-center">
          <input
            id="terms"
            name="terms"
            type="checkbox"
            v-model="agreeToTerms"
            required
            class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
            :class="{ 'border-red-500': v$.agreeToTerms.$error }"
          />
          <label for="terms" class="ml-2 block text-sm text-gray-900">
            I agree to the
            <a href="#" class="font-medium text-indigo-600 hover:text-indigo-500">Terms of Service</a>
            and
            <a href="#" class="font-medium text-indigo-600 hover:text-indigo-500">Privacy Policy</a>
          </label>
        </div>
        <div v-if="v$.agreeToTerms.$error" class="text-red-500 text-xs">
          {{ v$.agreeToTerms.$errors[0].$message }}
        </div>

        <div>
          <v-btn
            type="submit"
            block
            color="primary"
            :loading="loading"
            :disabled="loading"
            size="large"
            class="w-full"
          >
            Create account
          </v-btn>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useVuelidate } from '@vuelidate/core';
import { required, email as emailValidator, minLength, sameAs, helpers } from '@vuelidate/validators';
import { useAuthStore } from '@/stores/auth';
import { useNotificationStore } from '@/stores/notification';

const router = useRouter();
const authStore = useAuthStore();
const notificationStore = useNotificationStore();

// Form data
const name = ref('');
const email = ref('');
const password = ref('');
const confirmPassword = ref('');
const agreeToTerms = ref(false);
const loading = ref(false);
const error = ref('');

// Validation rules
const rules = {
  name: { required },
  email: { required, email: emailValidator },
  password: { 
    required, 
    minLength: minLength(8),
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
    required,
    sameAsPassword: helpers.withMessage(
      'Passwords must match', 
      sameAs(password)
    ) 
  },
  agreeToTerms: { 
    isChecked: helpers.withMessage(
      'You must agree to the terms and conditions', 
      (value) => value === true
    ) 
  }
};

const v$ = useVuelidate(rules, { name, email, password, confirmPassword, agreeToTerms });

// Handle registration form submission
const register = async () => {
  try {
    const isValid = await v$.value.$validate();
    if (!isValid) return;

    loading.value = true;
    error.value = '';

    await authStore.register({
      name: name.value,
      email: email.value,
      password: password.value
    });

    notificationStore.showNotification({
      message: 'Account created successfully! Welcome to PollPulse.',
      type: 'success'
    });

    router.push('/dashboard');
  } catch (err) {
    error.value = err.message || 'Failed to create account. Please try again later.';
  } finally {
    loading.value = false;
  }
};
</script>