import { defineConfig, configDefaults } from "vitest/config";
import react from "@vitejs/plugin-react";

// In dev, Vite serves the SPA and proxies GraphQL requests to the Go server so
// the frontend always talks to a same-origin `/query` (matching production,
// where the Go binary embeds frontend/dist and serves both on one origin).
export default defineConfig({
  plugins: [react({ jsxImportSource: "@emotion/react" })],
  server: {
    proxy: {
      "/query": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
  test: {
    // Unit tests live beside the source as *.test.ts(x). The Playwright visual
    // specs under tests/visual share the .spec.ts suffix but are driven by
    // `pnpm test:visual`, so keep them out of the vitest run.
    exclude: [...configDefaults.exclude, "tests/**"],
  },
});
