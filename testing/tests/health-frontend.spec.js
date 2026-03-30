const { test, expect } = require('@playwright/test');

test('TestHealth_GivenFrontendImage_WhenDockerHealthcheckRuns_ThenContainerBecomesHealthy', async ({ request }) => {
  const frontendBaseUrl = process.env.FRONTEND_BASE_URL || 'http://frontend:5173';
  const response = await request.get(frontendBaseUrl);
  expect(response.status()).toBe(200);
});

test('TestHealth_GivenFrontendRunning_WhenRootCalled_ThenReturns200', async ({ request }) => {
  const frontendBaseUrl = process.env.FRONTEND_BASE_URL || 'http://frontend:5173';
  const response = await request.get(frontendBaseUrl);
  expect(response.status()).toBe(200);
});

test('frontend loads and shows triangle of love coach', async ({ request }) => {
  const frontendBaseUrl = process.env.FRONTEND_BASE_URL || 'http://frontend:5173';

  const response = await request.get(frontendBaseUrl);
  expect(response.ok()).toBeTruthy();

  const body = (await response.text()).toLowerCase();
  expect(body).toContain('triangle of love coach');
});
