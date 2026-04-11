import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
  },
  esbuild: {
    include: /\.js$/,
    loader: 'jsx',
  },
  optimizeDeps: {
    esbuild: {
      loader: { '.js': 'jsx' },
    },
  },
  publicDir: 'public',
  server: {
    port: 3000,
    proxy: {
      '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,  
      },
    },
  },
  build: {
    outDir: 'dist',
  }
});
  

