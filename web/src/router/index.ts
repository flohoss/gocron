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
      meta: {
        title: 'GoBackup - System overview',
      },
    },
    {
      path: '/jobs/:id(\\d+)+',
      name: 'jobs',
      component: JobView,
      meta: {
        title: 'GoBackup - Job details',
      },
    },
    {
      path: '/jobs/restore',
      name: 'jobsRestore',
      component: RestoreView,
      meta: {
        title: 'GoBackup - Job restore',
      },
    },
    {
      path: '/jobs/form/:id(\\d+)*',
      name: 'jobsForm',
      component: JobsFormView,
      meta: {
        title: 'GoBackup - Job form',
      },
    },
  ],
});

router.beforeEach((to, from, next) => {
  const title: string = to.meta.title as string;
  if (title) {
    document.title = title;
  }
  next();
});

export default router;
