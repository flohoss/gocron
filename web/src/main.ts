import { createApp } from 'vue';
import App from './App.vue';
import { createPinia } from 'pinia';
import './style.css';
import router from './router';

export const BackendURL = import.meta.env.MODE === 'development' ? 'http://localhost:8156/' : '/';

createApp(App).use(router).use(createPinia()).mount('#app');
