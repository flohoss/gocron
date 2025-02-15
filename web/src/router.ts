import { createRouter, createWebHistory } from 'vue-router';

import HomeView from './pages/HomeView.vue';
import JobView from './pages/JobView.vue';

const routes = [
  { path: '/', name: 'homeView', component: HomeView },
  { path: '/jobs/:id', name: 'jobView', component: JobView },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
