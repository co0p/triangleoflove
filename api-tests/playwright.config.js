const { defineConfig } = require('@playwright/test');

module.exports = defineConfig({
  testDir: './tests',
  reporter: 'line',
  use: {
    baseURL: process.env.API_BASE_URL || 'http://backend:8080',
    extraHTTPHeaders: {
      Accept: 'application/json'
    }
  }
});
