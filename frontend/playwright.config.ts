import { defineConfig, devices } from "@playwright/test";

// Visual regression runs locally only (no CI wiring). The dev server is started
// by Playwright; GraphQL is mocked per-test via page.route, so neither the Go
// backend nor Postgres is needed.
const PORT = 5180;

export default defineConfig({
  testDir: "./tests/visual",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: 0,
  reporter: "list",
  snapshotPathTemplate: "{testDir}/__screenshots__/{arg}{ext}",
  use: {
    baseURL: `http://localhost:${PORT}`,
    // Headless shell only; pin the rendering surface so snapshots are stable.
    deviceScaleFactor: 1,
  },
  expect: {
    toHaveScreenshot: {
      // Recharts SVG rendering has sub-pixel jitter across runs; a small
      // tolerance keeps the non-regression check meaningful without flaking.
      maxDiffPixelRatio: 0.01,
    },
  },
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
  ],
  webServer: {
    command: `pnpm exec vite --port ${PORT} --strictPort`,
    port: PORT,
    reuseExistingServer: !process.env.CI,
    timeout: 120_000,
  },
});
