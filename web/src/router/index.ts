import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import JobsFormView from '../views/Jobs/FormView.vue';
import JobView from '../views/Jobs/JobView.vue';
import RestoreView from '@/views/Jobs/RestoreView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/jobs/:id(\\d+)+',
      name: 'jobs',
      component: JobView,
    },
    {
      path: '/jobs/restore',
      name: 'jobsRestore',
      component: RestoreView,
    },
    {
      path: '/jobs/form/:id(\\d+)*',
      name: 'jobsForm',
      component: JobsFormView,
    },
  ],
});

export default router;
