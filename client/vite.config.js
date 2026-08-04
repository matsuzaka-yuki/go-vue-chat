import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,
    proxy: {
      '/ws': {
        target: 'ws://localhost:8086',
        ws: true,
      },
      '/api': {
        target: 'http://localhost:8086',
      },
      '/avatars': {
        target: 'http://localhost:8086',
      },
      '/media': {
        target: 'http://localhost:8086',
      },
    },
  },
})
