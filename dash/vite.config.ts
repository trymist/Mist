import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react({
      babel: {
        plugins: [['babel-plugin-react-compiler']],
      },
    }),
    tailwindcss()
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    proxy: {
      '/api': {
        changeOrigin: true,
        target: process.env.vite_api_url || 'http://localhost:8080/',
        ws: true,
      },
      '/uploads': {
        changeOrigin: true,
        target: process.env.vite_api_url || 'http://localhost:8080/',
        ws: true,
      },


    }
  }
})
