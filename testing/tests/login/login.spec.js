const { test, expect } = require('@playwright/test');
const { FRONTEND_BASE_URL, SEED_EMAIL, loginViaUI } = require('../helpers/auth');

test('GivenNoJWT WhenAppOpened ThenLoginScreenShown', async ({ page }) => {
  await page.goto(FRONTEND_BASE_URL);
  await expect(page.locator('h1')).toHaveText('Sign in');
});

test('GivenValidCredentials WhenSubmitted ThenJWTStoredInBrowser', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);
  const token = await page.evaluate(() => localStorage.getItem('token'));
  expect(token).toBeTruthy();
});

test('GivenValidCredentials WhenSubmitted ThenDashboardShowsFirstName', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);
  // Wait for the firstName to load and be displayed in the header
  await expect(page.locator('header h1')).toContainText('Welcome back, River', { timeout: 10000 });
});

test('GivenInvalidCredentials WhenSubmitted ThenErrorShownOnLoginPage', async ({ page }) => {
  await page.goto(`${FRONTEND_BASE_URL}/login`);
  await page.fill('input[type="email"]', SEED_EMAIL);
  await page.fill('input[type="password"]', 'wrongpassword');
  await page.click('button[type="submit"]');
  await expect(page.locator('[role="alert"]')).toBeVisible();
  await expect(page).toHaveURL(/\/login/);
});

test('GivenStoredJWT WhenPageRefreshed ThenDashboardStillShown', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);
  await page.reload();
  await expect(page).toHaveURL(/\/dashboard/);
  // After reload, wait for firstName to load again
  await expect(page.locator('header h1')).toContainText('Welcome back, River', { timeout: 10000 });
});
