import { createApp } from 'vue';
import App from './App.vue';
import { createPinia } from 'pinia';
import './style.css';
import router from './router';
import { OpenAPI } from './openapi';

export const BackendURL = import.meta.env.MODE === 'development' ? 'http://localhost:8156/' : '/';

OpenAPI.BASE = BackendURL + 'api';

createApp(App).use(router).use(createPinia()).mount('#app');
