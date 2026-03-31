import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'jsdom',
    globals: true,
  },
  server: {
    host: true,
    allowedHosts: ['frontend', 'localhost', '.railway.app'],
    proxy: {
      '/api': {
        target: process.env.BACKEND_URL || 'http://backend:8080',
        changeOrigin: true
      }
    }
  }
});
