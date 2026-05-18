const { test, expect } = require('@playwright/test');
const { loginViaUI, FRONTEND_BASE_URL } = require('../helpers/auth');
const { resetDb } = require('../helpers/seed');

test.beforeAll(async ({ request }) => {
  await resetDb(request);
});

test('GivenLoggedIn_WhenNavigatingBetweenPages_ThenNavBarGreetingPersists', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);

  const greeting = page.locator('.navbar-greeting');

  // Greeting visible on dashboard immediately after login; wait for async load to complete
  await expect(greeting).toContainText('Hello, River', { timeout: 10000 });

  // Navigate to check-in — greeting must still be visible
  await page.goto(`${FRONTEND_BASE_URL}/checkin`);
  await expect(greeting).toBeVisible();
  await expect(greeting).toContainText('Hello, River');

  // Navigate to pairing — greeting must still be visible
  await page.goto(`${FRONTEND_BASE_URL}/pairing`);
  await expect(greeting).toBeVisible();
  await expect(greeting).toContainText('Hello, River');

  // Navigate back to dashboard — greeting must still be visible
  await page.goto(`${FRONTEND_BASE_URL}/dashboard`);
  await expect(greeting).toBeVisible();
  await expect(greeting).toContainText('Hello, River');
});
