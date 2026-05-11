const { test, expect } = require('@playwright/test');
const { SEED_EMAIL, SEED_PASSWORD } = require('../helpers/auth');
const { ADMIN_EMAIL, ADMIN_PASSWORD } = require('../helpers/users');

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

test('TestRegistration_GivenValidInput_WhenSubmitted_ThenAccountCreated', async ({ request }) => {
  const email = `reg-${Date.now()}@example.com`;
  const response = await request.post('/api/v1/register', {
    data: { email, password: 'securepass!', firstName: 'Tester' }
  });
  expect(response.status(), await response.text()).toBe(201);
});

test('TestRole_GivenAdminSeed_WhenInspected_ThenRoleIsAdmin', async ({ request }) => {
  const response = await request.post('/api/v1/auth/login', {
    data: { email: ADMIN_EMAIL, password: ADMIN_PASSWORD }
  });

  expect(response.status(), await response.text()).toBe(200);

  const { token } = await response.json();
  expect(token).toBeTruthy();

  // Decode the JWT payload (base64url middle segment) without verifying the signature.
  const payloadBase64 = token.split('.')[1];
  const payload = JSON.parse(Buffer.from(payloadBase64, 'base64url').toString('utf8'));

  expect(payload.role).toBe('admin');
});
