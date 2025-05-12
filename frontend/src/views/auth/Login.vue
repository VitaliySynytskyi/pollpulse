<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 bg-white p-8 rounded-lg shadow-md">
      <div>
        <h2 class="mt-6 text-center text-3xl font-bold text-gray-900">
          Sign in to your account
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Or
          <router-link to="/register" class="font-medium text-indigo-600 hover:text-indigo-500">
            create a new account
          </router-link>
        </p>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="login">
        <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
          {{ error }}
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
              autocomplete="current-password"
              required
              class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              :class="{ 'border-red-500': v$.password.$error }"
            />
            <div v-if="v$.password.$error" class="text-red-500 text-xs mt-1">
              {{ v$.password.$errors[0].$message }}
            </div>
          </div>
        </div>

        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <input
              id="remember_me"
              name="remember_me"
              type="checkbox"
              v-model="rememberMe"
              class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
            />
            <label for="remember_me" class="ml-2 block text-sm text-gray-900">
              Remember me
            </label>
          </div>

          <div class="text-sm">
            <a href="#" class="font-medium text-indigo-600 hover:text-indigo-500">
              Forgot your password?
            </a>
          </div>
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
            Sign in
          </v-btn>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useVuelidate } from '@vuelidate/core';
import { required, email as emailValidator, minLength } from '@vuelidate/validators';
import { useAuthStore } from '@/stores/auth';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

// Form data
const email = ref('');
const password = ref('');
const rememberMe = ref(false);
const loading = ref(false);
const error = ref('');

// Validation rules
const rules = {
  email: { required, email: emailValidator },
  password: { required, minLength: minLength(6) }
};

const v$ = useVuelidate(rules, { email, password });

// Handle login form submission
const login = async () => {
  try {
    const isValid = await v$.value.$validate();
    if (!isValid) return;

    loading.value = true;
    error.value = '';

    await authStore.login({
      email: email.value,
      password: password.value,
      remember: rememberMe.value
    });

    // Get redirect path from route query or default to dashboard
    const redirectPath = route.query.redirect || '/dashboard';
    router.push(redirectPath);
  } catch (err) {
    error.value = err.message || 'Failed to log in. Please check your credentials.';
  } finally {
    loading.value = false;
  }
};
</script> 