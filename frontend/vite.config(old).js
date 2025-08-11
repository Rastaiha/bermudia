import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import fs from 'fs'

export default defineConfig({
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  plugins: [
    vue(),
    {
      name: 'mock-api-plugin',
      enforce: 'pre',
      configureServer(server) {
        server.middlewares.use((req, res, next) => {
          if (req.url.startsWith('/api/v1/client/territory/')) {
            console.log(`[Mock API] Intercepted API call: ${req.url}`);

            const territoryId = req.url.split('/').pop();
            const filePath = path.resolve(__dirname, `public/data/territory${territoryId}.json`);

            try {
              const data = fs.readFileSync(filePath, 'utf-8');
              res.setHeader('Content-Type', 'application/json');
              res.end(data);
            } catch (error) {
              res.statusCode = 404;
              res.end(`File not found in public/data/territory${territoryId}.json`);
            }
            return;
          }

          next();
        });
      },
    },
  ],
})