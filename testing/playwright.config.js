const { defineConfig } = require('@playwright/test');

module.exports = defineConfig({
  testDir: './tests',
  reporter: 'line',
  workers: 1,
  globalSetup: './global-setup.js',
  use: {
    baseURL: process.env.API_BASE_URL || 'http://backend:8080'
  }
});
