import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

// Lazy-loaded components for better performance
const Home = () => import('@/views/Home.vue');
const Login = () => import('@/views/auth/Login.vue');
const Register = () => import('@/views/auth/Register.vue');
const Dashboard = () => import('@/views/Dashboard.vue');
const SurveyList = () => import('@/views/surveys/SurveyList.vue');
const SurveyCreate = () => import('@/views/surveys/SurveyCreate.vue');
const SurveyEdit = () => import('@/views/surveys/SurveyEdit.vue');
const SurveyView = () => import('@/views/surveys/SurveyView.vue');
const SurveyPreview = () => import('@/views/surveys/SurveyPreview.vue');
const SurveyTake = () => import('@/views/surveys/SurveyTake.vue');
const SurveyResults = () => import('@/views/surveys/SurveyResults.vue');
const NotFound = () => import('@/views/errors/NotFound.vue');
const Profile = () => import('@/views/Profile.vue');

const routes = [
  {
    path: '/',
    name: 'home',
    component: Home,
    meta: { title: 'PollPulse - Create and Share Surveys' }
  },
  {
    path: '/login',
    name: 'login',
    component: Login,
    meta: { title: 'Login', public: true }
  },
  {
    path: '/register',
    name: 'register',
    component: Register,
    meta: { title: 'Register', public: true }
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: Dashboard,
    meta: { title: 'Dashboard', requiresAuth: true }
  },
  {
    path: '/surveys',
    name: 'survey-list',
    component: SurveyList,
    meta: { title: 'My Surveys', requiresAuth: true }
  },
  {
    path: '/surveys/create',
    name: 'survey-create',
    component: SurveyCreate,
    meta: { title: 'Create Survey', requiresAuth: true }
  },
  {
    path: '/surveys/:id',
    name: 'survey-view',
    component: SurveyView,
    meta: { title: 'Survey Details', requiresAuth: true }
  },
  {
    path: '/surveys/:id/edit',
    name: 'survey-edit',
    component: SurveyEdit,
    meta: { title: 'Edit Survey', requiresAuth: true }
  },
  {
    path: '/surveys/:id/preview',
    name: 'survey-preview',
    component: SurveyPreview,
    meta: { title: 'Survey Preview', requiresAuth: true }
  },
  {
    path: '/surveys/:id/results',
    name: 'survey-results',
    component: SurveyResults,
    meta: { title: 'Survey Results', requiresAuth: true }
  },
  {
    path: '/surveys/:id/take',
    name: 'survey-take',
    component: SurveyTake,
    meta: { title: 'Take Survey', public: true }
  },
  {
    path: '/profile',
    name: 'profile',
    component: Profile,
    meta: { title: 'My Profile', requiresAuth: true }
  },
  // 404 route
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFound,
    meta: { title: 'Page Not Found', public: true }
  }
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition;
    } else {
      return { top: 0 };
    }
  }
});

// Navigation guards
router.beforeEach((to, from, next) => {
  // Update page title
  document.title = to.meta.title || 'PollPulse';
  
  const authStore = useAuthStore();
  const isLoggedIn = authStore.isAuthenticated;
  
  // Check if the route requires authentication
  if (to.meta.requiresAuth && !isLoggedIn) {
    // Redirect to login page with return URL
    return next({ 
      name: 'login', 
      query: { redirect: to.fullPath } 
    });
  }
  
  // If user is logged in and tries to access login/register pages, redirect to dashboard
  if (isLoggedIn && (to.name === 'login' || to.name === 'register')) {
    return next({ name: 'dashboard' });
  }
  
  next();
});

export default router; 