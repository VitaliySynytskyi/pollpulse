<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full text-center">
      <div class="text-8xl font-bold text-indigo-600 mb-4">404</div>
      <h1 class="text-3xl font-bold text-gray-900 mb-2">Page Not Found</h1>
      <p class="text-gray-600 mb-8">
        We couldn't find the page you were looking for. The page might have been moved, deleted, or never existed.
      </p>
      <div class="flex flex-col sm:flex-row space-y-4 sm:space-y-0 sm:space-x-4 justify-center">
        <router-link
          to="/"
          class="px-6 py-3 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 transition-colors"
        >
          Go Home
        </router-link>
        <button
          @click="goBack"
          class="px-6 py-3 bg-gray-200 text-gray-700 font-medium rounded-lg hover:bg-gray-300 transition-colors"
        >
          Go Back
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';

const router = useRouter();

const goBack = () => {
  router.go(-1);
};
</script> 