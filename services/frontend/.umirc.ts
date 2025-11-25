import { defineConfig } from 'umi';

export default defineConfig({
  npmClient: 'npm',
  history: {
    type: 'browser',
  },
  proxy: {
    '/api': {
      target: process.env.API_BASE_URL || 'http://localhost:8080',
      changeOrigin: true,
    },
  },
});
