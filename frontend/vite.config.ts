import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config
export default defineConfig({
  plugins: [react(), tailwindcss()],
  define: {
    // Build date, baked identically into the SSG and browser bundles. The
    // tournament calendar needs a "today" for its first render: using this
    // constant keeps the pre-rendered HTML and the hydration pass byte-identical
    // (a live `new Date()` would differ between them), and the real date takes
    // over in an effect right after mount.
    __BUILD_DATE__: JSON.stringify(new Date().toISOString().slice(0, 10)),
  },
  // Static files live at the repo root /public (per the project layout) and are
  // copied verbatim into dist/ at build.
  publicDir: '../public',
  build: {
    // Content-hashed filenames → the Go server can cache /assets/* forever.
    assetsInlineLimit: 2048,
    cssCodeSplit: true,
    target: 'es2020',
    rollupOptions: {
      output: {
        assetFileNames: 'assets/[name].[hash][extname]',
        chunkFileNames: 'assets/[name].[hash].js',
        entryFileNames: 'assets/[name].[hash].js',
      },
    },
  },
  server: {
    // Allow importing Markdown from ../content (outside the Vite root).
    fs: { allow: ['..'] },
  },
  ssr: {
    // Bundle these for the SSG render pass rather than externalizing.
    noExternal: ['react-helmet-async'],
  },
})
