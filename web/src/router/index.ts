import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import JobsFormView from '../views/Jobs/FormView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/jobs/:id',
      name: 'jobs',
      component: HomeView,
    },
    {
      path: '/jobs/form',
      name: 'jobsForm',
      component: JobsFormView,
    },
  ],
});

export default router;
