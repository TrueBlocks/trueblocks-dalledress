import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';
import tsconfigPaths from 'vite-tsconfig-paths';

export default defineConfig({
  plugins: [react(), tsconfigPaths()],
  resolve: {
    alias: {
      '@app': 'wailsjs/go/app/App',
      '@hooks': 'src/hooks',
      '@utils': 'src/utils',
      '@contexts': 'src/contexts',
    },
  },
  esbuild: {
    logOverride: { 'ignored-use-directive': 'silent' },
  },
});
