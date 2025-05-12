<template>
  <div class="app-container">
    <app-header v-if="showHeader" />
    <main class="min-h-[calc(100vh-64px)]">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
    <app-footer v-if="showFooter" />
    <notifications />
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useAuthStore } from './stores/auth';
import AppHeader from './components/layout/AppHeader.vue';
import AppFooter from './components/layout/AppFooter.vue';
import Notifications from './components/common/Notifications.vue';

const route = useRoute();
const authStore = useAuthStore();

// Hide header and footer on certain routes like login and register
const showHeader = computed(() => {
  return !['login', 'register', 'survey-take'].includes(route.name);
});

const showFooter = computed(() => {
  return !['login', 'register', 'survey-take'].includes(route.name);
});

onMounted(() => {
  // Check if user is logged in on app start
  authStore.checkAuth();
});
</script>

<style>
.app-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}
</style> 