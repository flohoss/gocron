import { createApp } from 'vue';
import App from './App.vue';
import './style.css';
import router from './router';
import { client } from './client/client.gen';
import { BackendURL } from './backend';

client.setConfig({ baseUrl: BackendURL });

createApp(App).use(router).mount('#app');
