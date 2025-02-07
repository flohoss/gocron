import { createApp } from 'vue';
import { createPinia } from 'pinia';
import { createMemoryHistory, createRouter } from 'vue-router';
import './style.css';
import App from './App.vue';

import HomeView from './pages/HomeView.vue';
import JobView from './pages/JobView.vue';

export const BackendURL = import.meta.env.MODE === 'development' ? 'http://localhost:8080/' : '/';

const pinia = createPinia();
const app = createApp(App);
const router = createRouter({
  history: createMemoryHistory(),
  routes: [
    { path: '/', component: HomeView },
    { path: '/jobs/:id', component: JobView },
  ],
});

app.use(pinia);
app.use(router);
app.mount('#app');
