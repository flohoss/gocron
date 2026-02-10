const { defineConfig } = require('cypress');

module.exports = defineConfig({
  allowCypressEnv: false,
  e2e: {
    supportFile: false,
    baseUrl: process.env.CYPRESS_baseUrl || 'http://localhost:5173',
  },
});
