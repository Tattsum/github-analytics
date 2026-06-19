import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// In dev, Vite serves the SPA and proxies GraphQL requests to the Go server so
// the frontend always talks to a same-origin `/query` (matching production,
// where the Go binary embeds frontend/dist and serves both on one origin).
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/query": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
