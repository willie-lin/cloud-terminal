import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
  ],
  server: {
    port: 3000,
    proxy: {
      // Proxy API requests to the Go backend (port 8443)
      '/api': {
        target: 'http://localhost:8443',
        changeOrigin: true,
        secure: false,
      },
      '/admin': {
        target: 'http://localhost:8443',
        changeOrigin: true,
        secure: false,
      },
      '/webhook': {
        target: 'http://localhost:8443',
        changeOrigin: true,
        secure: false,
      },
    },
  },
})
