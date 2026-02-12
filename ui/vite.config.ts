import tailwindcss from "@tailwindcss/vite";
import { devtools } from "@tanstack/devtools-vite";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import viteReact from "@vitejs/plugin-react";
import { nitro } from "nitro/vite";
import { defineConfig } from "vite";
import viteTsConfigPaths from "vite-tsconfig-paths";

const config = defineConfig({
  plugins: [
    devtools(),
    nitro({
      routeRules: {
        "/api/**": { proxy: "http://localhost:8080/api/**" },
      },
    }),
    viteTsConfigPaths({
      projects: ["./tsconfig.json"],
    }),
    tailwindcss(),
    tanstackStart({
      spa: {
        enabled: true,
        prerender: {
          outputPath: "/index.html",
        },
      },
    }),
    viteReact(),
  ],
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
        bypass: (req) => {
          // Bypass proxy for OAuth routes - let browser handle redirects
          if (req.url?.startsWith("/api/auth/google")) {
            return req.url.replace("/api", "http://localhost:8080/api");
          }
        },
      },
    },
  },
});

export default config;
