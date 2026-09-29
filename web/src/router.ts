import { createRouter, createWebHistory, isNavigationFailure, NavigationFailureType } from 'vue-router';

import HomeView from './pages/HomeView.vue';
import JobView from './pages/JobView.vue';
import CommandView from './pages/CommandView.vue';
import { useAuth } from './stores/useAuth';

const routes = [
  { path: '/', name: 'homeView', component: HomeView, meta: { title: 'GoCron' } },
  { path: '/jobs/:id', name: 'jobView', component: JobView, meta: { title: 'Job' } },
  { path: '/commands', name: 'commandView', component: CommandView, meta: { title: 'Command' } },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to) => {
  document.title = `${to.meta.title}`;

  const auth = useAuth();
  if (!auth.ready.value) await auth.fetchCurrentUser();

  if (auth.authEnabled.value && !auth.authenticated.value) {
    window.location.href = '/api/auth/login';
    return false;
  }
});

router.onError((error) => {
  if (!isNavigationFailure(error, NavigationFailureType.aborted | NavigationFailureType.cancelled)) {
    throw error;
  }
});

export default router;
