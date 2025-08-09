import { createApp } from 'vue';
import App from './App.vue';
import './style.css';
import router from './router';
import { client } from './client/client.gen';

export const BackendURL = import.meta.env.MODE === 'development' ? 'http://localhost:8156' : '';

client.setConfig({ baseUrl: BackendURL });

createApp(App).use(router).mount('#app');
