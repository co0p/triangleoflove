const { test, expect } = require('@playwright/test');
const { loginViaUI, FRONTEND_BASE_URL } = require('../helpers/auth');

const METRIC_IDS = [
  'felt_understood', 'meaningful_sharing', 'could_count_on_them',
  'effort_for_us', 'desire', 'spark', 'mood'
];

/**
 * Set a range input to a value by evaluating on the element directly.
 * Dispatches an 'input' event so Vue's @input handler fires.
 */
async function setSlider(page, id, value) {
  await page.locator(`#${id}`).evaluate((el, val) => {
    el.value = val;
    el.dispatchEvent(new Event('input', { bubbles: true }));
  }, value);
}

/**
 * Log in and navigate to the checkin page, waiting for the page to render.
 */
async function goToCheckin(page) {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);
  await page.goto(`${FRONTEND_BASE_URL}/checkin`);
  await expect(page.locator('h1')).toHaveText('Daily check-in');
}

test('GivenLoggedIn_WhenCheckinPageOpened_ThenNewMetricSlidersVisible', async ({ page }) => {
  await goToCheckin(page);

  const expectedLabels = [
    'Felt understood', 'Meaningful sharing', 'Could count on them',
    'Effort for us', 'Desire', 'Spark', 'My mood today'
  ];
  for (const label of expectedLabels) {
    await expect(page.getByText(label)).toBeVisible();
  }

  const sliders = page.locator('input[type="range"]');
  await expect(sliders).toHaveCount(7);
});

test('GivenLoggedIn_WhenSlidersAdjustedAndSaved_ThenConfirmationShown', async ({ page }) => {
  await goToCheckin(page);

  await setSlider(page, 'felt_understood', 4);
  await setSlider(page, 'mood', 3);

  await page.click('[data-testid="save-checkin"]');
  await expect(page.locator('[data-testid="checkin-confirmation"]')).toBeVisible();
});

test('GivenCheckinSaved_WhenPageReloaded_ThenSavedValuesPrePopulated', async ({ page }) => {
  await goToCheckin(page);

  await setSlider(page, 'felt_understood', 5);
  await setSlider(page, 'desire', 3);
  await setSlider(page, 'mood', 2);

  await page.click('[data-testid="save-checkin"]');
  await expect(page.locator('[data-testid="checkin-confirmation"]')).toBeVisible();

  await page.reload();
  await expect(page.locator('h1')).toHaveText('Daily check-in');
  await expect(page.locator('#felt_understood')).toHaveValue('5');
  await expect(page.locator('#desire')).toHaveValue('3');
  await expect(page.locator('#mood')).toHaveValue('2');
});
