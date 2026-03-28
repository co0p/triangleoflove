import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  plugins: [vue()],
  server: {
    host: true,
    allowedHosts: ['frontend', 'localhost'],
    proxy: {
      '/api': {
        target: process.env.BACKEND_URL || 'http://backend:8080',
        changeOrigin: true
      }
    }
  }
});
