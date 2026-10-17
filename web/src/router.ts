import { createRouter, createWebHistory, isNavigationFailure, NavigationFailureType } from 'vue-router';

import HomeView from './pages/HomeView.vue';
import JobView from './pages/JobView.vue';
import CommandView from './pages/CommandView.vue';
import LoginView from './pages/LoginView.vue';
import { useAuth } from './stores/useAuth';

const routes = [
  { path: '/', name: 'homeView', component: HomeView, meta: { title: 'GoCron' } },
  { path: '/jobs/:id', name: 'jobView', component: JobView, meta: { title: 'GoCron' } },
  { path: '/commands', name: 'commandView', component: CommandView, meta: { title: 'GoCron - Command' } },
  { path: '/login', name: 'loginView', component: LoginView, meta: { title: 'GoCron - Sign in', public: true } },
];

declare module 'vue-router' {
  interface RouteMeta {
    title: string;
    public?: boolean;
  }
}

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to) => {
  document.title = `${to.meta.title}`;

  if (to.meta.public) {
    return true;
  }

  const auth = useAuth();
  if (!auth.ready.value) await auth.fetchCurrentUser();

  if (auth.authEnabled.value && !auth.authenticated.value) {
    return { name: 'loginView', query: { redirect: to.fullPath } };
  }
});

router.onError((error) => {
  if (!isNavigationFailure(error, NavigationFailureType.aborted | NavigationFailureType.cancelled)) {
    throw error;
  }
});

export default router;
