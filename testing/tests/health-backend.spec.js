const { test, expect } = require('@playwright/test');

test('TestBootstrap GivenRunningStack WhenHealthEndpointCalled ThenHealthyResponseReturned', async ({ request }) => {
  const response = await request.get('/api/v1/health');

  expect(response.ok()).toBeTruthy();

  const body = await response.json();
  console.log('health endpoint response:', body);
  expect(body).toEqual(
    expect.objectContaining({
      status: 'healthy'
    })
  );
});

test('status endpoint returns HTTP status payload', async ({ request }) => {
  const response = await request.get('/api/v1/status');

  expect(response.status()).toBe(200);

  const body = await response.json();
  expect(body).toEqual({
    status: 'ok',
    code: 200
  });
});
