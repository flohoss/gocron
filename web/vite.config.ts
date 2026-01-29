import { defineConfig } from 'vite';
import tailwindcss from '@tailwindcss/vite';
import vue from '@vitejs/plugin-vue';
import Icons from 'unplugin-icons/vite';

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss(), Icons({ compiler: 'vue3' })],
});
