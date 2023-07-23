import './assets/main.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import router from './router';
import { OpenAPI } from './openapi';
import { plugin, defaultConfig } from '@formkit/vue';
import { generateClasses } from '@formkit/themes';

OpenAPI.BASE = '/api';

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(
  plugin,
  defaultConfig({
    config: {
      classes: generateClasses({
        global: {
          outer: 'col-span-1',
          message: 'label-text-alt',
          messages: 'label-text-alt',
          help: 'label-text-alt opacity-50',
          label: 'label-text',
        },
        'family:actions': {
          input: 'btn btn-primary',
        },
        'family:button': {
          input: 'btn btn-primary',
        },
        'family:text': {
          input: 'input input-bordered w-full my-1',
        },
        'family:dropdown': {
          input: 'select',
        },
      }),
    },
  })
);

app.mount('#app');
