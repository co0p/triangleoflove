const { test, expect } = require('@playwright/test');
const { FRONTEND_BASE_URL, loginViaUI } = require('../helpers/auth');
const { resetDb } = require('../helpers/seed');

test.beforeAll(async ({ request }) => {
  await resetDb(request);
});

test('TestDashboardMatrix_GivenRecentCheckins_WhenDashboardLoads_ThenMonthlyMatrixShowsColoredCells', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);

  const matrix = page.getByTestId('monthly-matrix');
  await expect(matrix).toBeVisible();

  // At least one cell should have a color class (not unavailable) because
  // init.sql seeds 7 relative-date check-ins for the past week.
  const coloredCell = matrix.locator(
    '[data-testid="matrix-cell"]:not(.insight-value--unavailable)'
  ).first();
  await expect(coloredCell).toBeVisible();
});

test('TestDashboardMatrix_GivenRecentCheckins_WhenDashboardLoads_ThenAllThreeDimensionRowsVisible', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);

  const rows = page.getByTestId('monthly-matrix').getByTestId('matrix-row');
  await expect(rows).toHaveCount(3);
});
