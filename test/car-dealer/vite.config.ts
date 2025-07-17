import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import tailwindcss from '@tailwindcss/vite';
import path from 'path';

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    port: 8081,
    host: true, // Needed to allow custom domains
    allowedHosts: ['cardealer.local', 'tracker.local','zeals-tracker-api.ngrok.app', 'cardealerform.local', 'rental-car.ngrok.app', 'car-form.ngrok.app'],
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
});
