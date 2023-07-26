import './assets/main.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import router from './router';
import { OpenAPI } from './openapi';
import { useJobStore } from './stores/jobs';

OpenAPI.BASE = '/api';

const app = createApp(App);

app.use(createPinia());
app.use(router);

const store = useJobStore();

store
  .getJobs()
  .then(() => app.mount('#app'))
  .catch((err) => console.log(err));
