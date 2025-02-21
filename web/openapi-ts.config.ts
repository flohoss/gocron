import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  input: 'openapi.json',
  output: {
    indexFile: false,
    path: 'src/client',
  },
  plugins: ['@hey-api/client-fetch'],
});
