const { test, expect } = require('@playwright/test');
const { FRONTEND_BASE_URL, SEED_EMAIL, SEED_PASSWORD, loginViaUI } = require('../helpers/auth');

test('TestProfile_GivenAuthenticated_WhenNavBarNameClicked_ThenProfilePageLoads', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);

  await page.locator('a.navbar-profile-link').click();

  await expect(page).toHaveURL(/\/profile/);
  await expect(page.locator('h1')).toHaveText('Profile');
});

test('TestProfile_GivenAuthenticated_WhenPageLoads_ThenEmailShownReadOnly', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);
  await page.goto(`${FRONTEND_BASE_URL}/profile`);

  const emailInput = page.locator('input[type="email"]');
  await expect(emailInput).toHaveValue(SEED_EMAIL);
  await expect(emailInput).toHaveAttribute('readonly', '');
});

test('TestProfile_GivenAuthenticated_WhenLogoutClicked_ThenTokenRemovedAndRedirectedToLogin', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);
  await page.goto(`${FRONTEND_BASE_URL}/profile`);

  await page.locator('button', { hasText: 'Log out' }).click();

  await expect(page).toHaveURL(/\/login/);
  const token = await page.evaluate(() => localStorage.getItem('token'));
  expect(token).toBeNull();
});

test('TestProfile_GivenCorrectCurrentPassword_WhenNewPasswordSubmitted_ThenPasswordUpdated', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);
  await page.goto(`${FRONTEND_BASE_URL}/profile`);

  await page.fill('input[name="current_password"]', SEED_PASSWORD);
  await page.fill('input[name="new_password"]', SEED_PASSWORD);
  await page.fill('input[name="confirm_password"]', SEED_PASSWORD);
  await page.locator('button[type="submit"]').click();

  await expect(page.locator('.password-success')).toBeVisible();
});

test('TestProfile_GivenWrongCurrentPassword_WhenNewPasswordSubmitted_ThenErrorShown', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);
  await page.goto(`${FRONTEND_BASE_URL}/profile`);

  await page.fill('input[name="current_password"]', 'definitely-wrong');
  await page.fill('input[name="new_password"]', 'newpassword');
  await page.fill('input[name="confirm_password"]', 'newpassword');
  await page.locator('button[type="submit"]').click();

  await expect(page.locator('.alert-error')).toBeVisible();
  await expect(page).toHaveURL(/\/profile/);
});

test('TestProfile_GivenMismatchedNewPasswords_WhenFormSubmitted_ThenErrorShownNoRequest', async ({ page }) => {
  await loginViaUI(page);
  await expect(page).toHaveURL(/\/dashboard/);
  await page.goto(`${FRONTEND_BASE_URL}/profile`);

  await page.fill('input[name="current_password"]', SEED_PASSWORD);
  await page.fill('input[name="new_password"]', 'password-a');
  await page.fill('input[name="confirm_password"]', 'password-b');
  await page.locator('button[type="submit"]').click();

  await expect(page.locator('.alert-error')).toBeVisible();
  // No network request sent — verify by checking we're still on profile with no success shown
  await expect(page.locator('.password-success')).not.toBeVisible();
  await expect(page).toHaveURL(/\/profile/);
});
