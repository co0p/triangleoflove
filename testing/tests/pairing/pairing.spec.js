const { test, expect } = require('@playwright/test');
const { FRONTEND_BASE_URL, SEED_EMAIL_2, SEED_PASSWORD_2, loginViaUI } = require('../helpers/auth');
const { resetDb } = require('../helpers/seed');

// Reset to unpaired baseline once before the suite.
// Tests in this file are ordered: the "GivenPaired" tests rely on the pair
// formed by GivenTwoUnpairedUsers, so per-test resets would break them.
test.beforeAll(async ({ request }) => {
  await resetDb(request);
});

test('GivenUnpaired_WhenVisitsPairingPage_ThenShowsOwnCodeAndInput', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);
  await page.goto(`${FRONTEND_BASE_URL}/pairing`);

  const code = page.getByTestId('invite-code');
  await expect(code).toBeVisible();

  await expect(page.getByTestId('partner-code-input')).toBeVisible();
});

test('GivenTwoUnpairedUsers_WhenCodeEntered_ThenCoupleFormed', async ({ browser }) => {
  // River logs in and gets their invite code
  const riverContext = await browser.newContext();
  const riverPage = await riverContext.newPage();
  await loginViaUI(riverPage);
  await expect(riverPage).toHaveURL(/\/dashboard/);
  await riverPage.goto(`${FRONTEND_BASE_URL}/pairing`);
  const riverCode = (await riverPage.getByTestId('invite-code').textContent()).trim();

  // Jordan logs in and enters River's code
  const jordanContext = await browser.newContext();
  const jordanPage = await jordanContext.newPage();
  await loginViaUI(jordanPage, SEED_EMAIL_2, SEED_PASSWORD_2);
  await expect(jordanPage).toHaveURL(/\/dashboard/);
  await jordanPage.goto(`${FRONTEND_BASE_URL}/pairing`);
  await jordanPage.getByTestId('partner-code-input').fill(riverCode);
  await jordanPage.getByTestId('connect-button').click();

  await expect(jordanPage.getByTestId('connect-success')).toBeVisible();

  await riverContext.close();
  await jordanContext.close();
});

// D4: these tests run after the D3 test above has paired River and Jordan
test('GivenPaired_WhenVisitsPairingPage_ThenShowsPartnerNameAndSince', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);
  await page.goto(`${FRONTEND_BASE_URL}/pairing`);

  await expect(page.getByTestId('partner-name')).toBeVisible();
  await expect(page.getByTestId('partner-name')).toContainText('Jordan');
  await expect(page.getByTestId('paired-since')).toBeVisible();
});

test('GivenPaired_WhenVisitsDashboard_ThenShowsPairingStatus', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);

  await expect(page.getByTestId('pairing-status')).toBeVisible();
  await expect(page.getByTestId('pairing-status')).toContainText('Jordan');
});

test('GivenPaired_WhenUnpairConfirmed_ThenInviteCodeVisible', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);
  await page.goto(`${FRONTEND_BASE_URL}/pairing`);

  // Confirm paired state is shown before unpairing.
  await expect(page.getByTestId('unpair-button')).toBeVisible();

  // Open confirmation modal and confirm.
  await page.getByTestId('unpair-button').click();
  await expect(page.getByTestId('unpair-modal')).toBeVisible();
  await page.getByTestId('unpair-modal-confirm').click();

  // After unpairing the invite code flow must be shown.
  await expect(page.getByTestId('invite-code')).toBeVisible();
  await expect(page.getByTestId('partner-name')).not.toBeVisible();
  await expect(page.getByTestId('unpair-button')).not.toBeVisible();
});
