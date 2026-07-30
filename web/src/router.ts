import { createRouter, createWebHistory } from 'vue-router';

import HomeView from './pages/HomeView.vue';
import JobView from './pages/JobView.vue';
import CommandView from './pages/CommandView.vue';

const routes = [
  { path: '/', name: 'homeView', component: HomeView, meta: { title: 'GoCron' } },
  { path: '/jobs/:id', name: 'jobView', component: JobView, meta: { title: 'Job' } },
  { path: '/commands', name: 'commandView', component: CommandView, meta: { title: 'Command' } },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, _, next) => {
  document.title = `${to.meta.title}`;
  next();
});

export default router;
