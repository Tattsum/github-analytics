import { expect, test } from "@playwright/test";
import { mockGraphql } from "./mockGraphql";

const WIDTHS = [
  { name: "desktop-1280", width: 1280, height: 900 },
  { name: "tablet-768", width: 768, height: 1024 },
  { name: "mobile-375", width: 375, height: 812 },
] as const;

const ROUTES = [
  { name: "team-overview", path: "/" },
  { name: "repositories", path: "/repositories" },
  { name: "repository-detail", path: "/repositories/acme%2Fweb" },
  { name: "member-detail", path: "/members/octocat" },
] as const;

// Recharts renders into an SVG sized by ResponsiveContainer; wait for it so the
// screenshot is taken after layout settles.
async function waitForRender(page: import("@playwright/test").Page) {
  await page.locator(".recharts-surface").first().waitFor({ state: "visible" });
  await page.waitForFunction(() => document.fonts.ready.then(() => true));
}

for (const viewport of WIDTHS) {
  test.describe(viewport.name, () => {
    test.use({ viewport: { width: viewport.width, height: viewport.height } });

    for (const route of ROUTES) {
      test(route.name, async ({ page }) => {
        await mockGraphql(page);
        await page.goto(route.path);
        await waitForRender(page);
        await expect(page).toHaveScreenshot(`${route.name}-${viewport.name}.png`, {
          fullPage: true,
          animations: "disabled",
        });
      });
    }
  });
}

// Mobile-only: capture the header after the hamburger menu is opened. Skipped
// when the toggle is absent (e.g. pre-migration baseline run) so the same spec
// works before and after the emotion/responsive migration.
test.describe("mobile-375 interactions", () => {
  test.use({ viewport: { width: 375, height: 812 } });

  test("hamburger-open", async ({ page }) => {
    await mockGraphql(page);
    await page.goto("/");
    await waitForRender(page);

    const toggle = page.getByRole("button", { name: "メニュー" });
    test.skip((await toggle.count()) === 0, "hamburger toggle not present");

    await toggle.click();
    await expect(page).toHaveScreenshot("hamburger-open-mobile-375.png", {
      clip: { x: 0, y: 0, width: 375, height: 300 },
      animations: "disabled",
    });
  });
});
