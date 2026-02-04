import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  build: {
    lib: {
      entry: "./src/index.tsx",
      formats: ["iife"],
      name: "TexlyWidget",
      fileName: () => "texly-widget.js",
    },
    rollupOptions: {
      output: {
        inlineDynamicImports: true,
      },
    },
    cssCodeSplit: false,
    minify: "esbuild",
  },
  define: {
    "process.env": {
      NODE_ENV: JSON.stringify("production"),
    },
  },
});
