const { test, expect } = require('@playwright/test');
const { getToken } = require('../helpers/auth');

test('TestUsersMe GivenValidJWT WhenRequested ThenFirstNameReturned', async ({ request }) => {
  const token = await getToken(request);

  const response = await request.get('/api/v1/users/me', {
    headers: { Authorization: `Bearer ${token}` }
  });

  expect(response.status(), await response.text()).toBe(200);

  const body = await response.json();
  expect(body).toEqual(
    expect.objectContaining({
      firstName: 'River'
    })
  );
});

test('TestUsersMe GivenNoJWT WhenRequested ThenReturns401', async ({ request }) => {
  const response = await request.get('/api/v1/users/me');

  expect(response.status()).toBe(401);
});
