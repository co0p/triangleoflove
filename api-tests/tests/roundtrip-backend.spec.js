const { test, expect } = require('@playwright/test');

test('TestBootstrap_GivenBackendAndDb_WhenRoundtripExecuted_ThenDataPersistedAndReturned', async ({ request }) => {
  const response = await request.post('/demo/roundtrip');

  expect(response.status(), await response.text()).toBe(200);

  const body = await response.json();
  expect(body).toEqual(
    expect.objectContaining({
      id: expect.any(Number),
      value: expect.any(String)
    })
  );
});