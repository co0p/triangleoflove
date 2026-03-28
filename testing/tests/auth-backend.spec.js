const { test, expect } = require('@playwright/test');
const { SEED_EMAIL, SEED_PASSWORD } = require('./helpers/auth');

test('TestAuthLogin GivenValidCredentials WhenRequested ThenReturns200WithJWT', async ({ request }) => {
  const response = await request.post('/api/v1/auth/login', {
    data: { email: SEED_EMAIL, password: SEED_PASSWORD }
  });

  expect(response.status(), await response.text()).toBe(200);

  const body = await response.json();
  expect(body).toEqual(
    expect.objectContaining({
      token: expect.any(String)
    })
  );
});

test('TestAuthLogin GivenInvalidCredentials WhenRequested ThenReturns401', async ({ request }) => {
  const response = await request.post('/api/v1/auth/login', {
    data: { email: SEED_EMAIL, password: 'wrongpassword' }
  });

  expect(response.status()).toBe(401);
});
