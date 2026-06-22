import type { Page } from "@playwright/test";
import { graphqlResponses } from "./fixtures/graphql";

// Intercepts the same-origin /query endpoint and answers with fixed fixtures
// keyed by operation name, so the SPA renders without the Go backend.
export async function mockGraphql(page: Page): Promise<void> {
  await page.route("**/query", async (route) => {
    const postData = route.request().postData() ?? "{}";
    const { operationName } = JSON.parse(postData) as { operationName?: string };
    const data = operationName ? graphqlResponses[operationName] : undefined;

    if (!data) {
      await route.fulfill({
        status: 500,
        contentType: "application/json",
        body: JSON.stringify({ errors: [{ message: `no fixture for operation: ${operationName}` }] }),
      });
      return;
    }

    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ data }),
    });
  });
}
