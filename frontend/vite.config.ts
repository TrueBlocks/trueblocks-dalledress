import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@wailsjs': path.resolve(__dirname, './wailsjs'),
      '@trueblocks/ui': path.resolve(__dirname, '../../packages/ui/src'),
    },
    dedupe: ['react', 'react-dom', '@mantine/core', '@mantine/hooks'],
  },
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
});
