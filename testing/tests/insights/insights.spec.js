const { test, expect } = require('@playwright/test');
const { FRONTEND_BASE_URL, loginViaUI } = require('../helpers/auth');
const { resetDb } = require('../helpers/seed');

const SEEDED_INSIGHTS_DATE = '20260330';

test.beforeAll(async ({ request }) => {
  await resetDb(request);
});

test('TestInsights_GivenNoAuth_WhenNavigated_ThenRedirectsToLogin', async ({ page }) => {
  await page.goto(`${FRONTEND_BASE_URL}/insights/${SEEDED_INSIGHTS_DATE}`);

  await expect(page).toHaveURL(/\/login$/);
  await expect(page.getByRole('heading', { name: 'Sign in' })).toBeVisible();
});

test('TestInsights_GivenCheckinExists_WhenRequested_ThenReturnsThreeScores', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);

  await page.goto(`${FRONTEND_BASE_URL}/insights/${SEEDED_INSIGHTS_DATE}`);

  await expect(page).toHaveURL(new RegExp(`/insights/${SEEDED_INSIGHTS_DATE}$`));
  await expect(page.getByRole('heading', { name: 'Daily insights' })).toBeVisible();
  await expect(page.getByTestId('insight-intimacy')).toHaveText('100');
  await expect(page.getByTestId('insight-commitment')).toHaveText('100');
  await expect(page.getByTestId('insight-passion')).toHaveText('100');
});

test('TestInsights_GivenUnpairedUser_WhenRequested_ThenReturnsScores', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);
  await expect(page.locator('main')).toContainText('Not connected yet');

  await page.goto(`${FRONTEND_BASE_URL}/insights/${SEEDED_INSIGHTS_DATE}`);

  await expect(page).toHaveURL(new RegExp(`/insights/${SEEDED_INSIGHTS_DATE}$`));
  await expect(page.getByRole('heading', { name: 'Daily insights' })).toBeVisible();
  await expect(page.getByTestId('insight-intimacy')).toHaveText('100');
  await expect(page.getByTestId('insight-commitment')).toHaveText('100');
  await expect(page.getByTestId('insight-passion')).toHaveText('100');
});