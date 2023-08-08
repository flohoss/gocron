import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import JobsFormView from '../views/Jobs/FormView.vue';
import JobView from '../views/Jobs/JobView.vue';
import RestoreView from '@/views/Jobs/RestoreView.vue';
import { useSystemStore } from '@/stores/system';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: {
        title: ' - System overview',
      },
    },
    {
      path: '/jobs/:id(\\d+)+',
      name: 'jobs',
      component: JobView,
      meta: {
        title: ' - Job details',
      },
    },
    {
      path: '/jobs/restore',
      name: 'jobsRestore',
      component: RestoreView,
      meta: {
        title: ' - Job restore',
      },
    },
    {
      path: '/jobs/form/:id(\\d+)*',
      name: 'jobsForm',
      component: JobsFormView,
      meta: {
        title: ' - Job form',
      },
    },
  ],
});

router.beforeEach((to, from, next) => {
  const store = useSystemStore();
  const title: string = to.meta.title as string;
  if (title) {
    document.title = (store.system.config.identifier || 'GoBackup') + title;
  }
  next();
});

export default router;
