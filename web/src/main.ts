import './assets/main.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import router from './router';
import { OpenAPI } from './openapi';
import { useSystemStore } from './stores/system';

OpenAPI.BASE = '/api';

const app = createApp(App);

app.use(createPinia());
app.use(router);

const systemStore = useSystemStore();
try {
  await systemStore.fetchSystem();
} catch (err) {
  console.log(err);
}
app.mount('#app');
