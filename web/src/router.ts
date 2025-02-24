import { createRouter, createWebHistory } from 'vue-router';

import HomeView from './pages/HomeView.vue';
import JobView from './pages/JobView.vue';

const routes = [
  { path: '/', name: 'homeView', component: HomeView, meta: { title: 'GoCron' } },
  { path: '/jobs/:id', name: 'jobView', component: JobView, meta: { title: 'Job' } },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, _, next) => {
  const id = to.params.id;
  document.title = `${to.meta.title}`;
  if (id !== undefined) {
    const name = id + '';
    document.title += ` - ${name.toUpperCase()}`;
  }
  next();
});

export default router;
