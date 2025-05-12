<template>
  <div class="min-h-screen">
    <!-- Hero Section -->
    <section class="bg-gradient-to-r from-indigo-600 to-purple-600 text-white py-20">
      <div class="container mx-auto px-4">
        <div class="flex flex-col md:flex-row items-center justify-between">
          <div class="md:w-1/2 mb-10 md:mb-0">
            <h1 class="text-4xl md:text-5xl font-bold mb-4">Create Powerful Surveys with PollPulse</h1>
            <p class="text-xl mb-8">
              Collect data, gather insights, and make informed decisions with our easy-to-use survey platform.
            </p>
            <div class="flex flex-col sm:flex-row space-y-4 sm:space-y-0 sm:space-x-4">
              <router-link
                to="/register"
                class="px-6 py-3 bg-white text-indigo-600 font-medium rounded-lg hover:bg-gray-100 transition-colors"
              >
                Get Started Free
              </router-link>
              <router-link
                to="/login"
                class="px-6 py-3 bg-transparent border border-white text-white font-medium rounded-lg hover:bg-white/10 transition-colors"
              >
                Sign In
              </router-link>
            </div>
          </div>
          <div class="md:w-1/2">
            <img src="@/assets/images/hero-image.svg" alt="Survey Illustration" class="w-full max-w-md mx-auto" />
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="py-16 bg-gray-50">
      <div class="container mx-auto px-4">
        <h2 class="text-3xl font-bold text-center mb-12">Why Choose PollPulse?</h2>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <!-- Feature 1 -->
          <div class="bg-white p-6 rounded-lg shadow-md">
            <div class="w-14 h-14 bg-indigo-100 rounded-full flex items-center justify-center mb-4">
              <v-icon icon="mdi-form-select" class="text-indigo-600 text-2xl"></v-icon>
            </div>
            <h3 class="text-xl font-semibold mb-3">Easy Survey Creation</h3>
            <p class="text-gray-600">
              Build beautiful surveys in minutes with our intuitive drag-and-drop interface. No coding required.
            </p>
          </div>
          
          <!-- Feature 2 -->
          <div class="bg-white p-6 rounded-lg shadow-md">
            <div class="w-14 h-14 bg-indigo-100 rounded-full flex items-center justify-center mb-4">
              <v-icon icon="mdi-chart-box" class="text-indigo-600 text-2xl"></v-icon>
            </div>
            <h3 class="text-xl font-semibold mb-3">Real-time Analytics</h3>
            <p class="text-gray-600">
              Get instant insights with powerful analytics and visualization tools as responses come in.
            </p>
          </div>
          
          <!-- Feature 3 -->
          <div class="bg-white p-6 rounded-lg shadow-md">
            <div class="w-14 h-14 bg-indigo-100 rounded-full flex items-center justify-center mb-4">
              <v-icon icon="mdi-shield-check" class="text-indigo-600 text-2xl"></v-icon>
            </div>
            <h3 class="text-xl font-semibold mb-3">Secure & Reliable</h3>
            <p class="text-gray-600">
              Your data is always secure with our enterprise-grade encryption and privacy-first approach.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="py-16 bg-indigo-900 text-white">
      <div class="container mx-auto px-4 text-center">
        <h2 class="text-3xl font-bold mb-6">Ready to get started?</h2>
        <p class="text-xl max-w-2xl mx-auto mb-8">
          Join thousands of organizations using PollPulse to collect feedback and make data-driven decisions.
        </p>
        <router-link
          to="/register"
          class="px-8 py-4 bg-white text-indigo-900 font-medium rounded-lg hover:bg-gray-100 transition-colors inline-block"
        >
          Create Your First Survey
        </router-link>
      </div>
    </section>
  </div>
</template>

<script setup>
// No additional logic needed for this component
</script> 