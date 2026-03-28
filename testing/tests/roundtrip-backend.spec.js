const { test, expect } = require('@playwright/test');

test('TestBootstrap GivenBackendAndDb WhenRoundtripExecuted ThenDataPersistedAndReturned', async ({ request }) => {
  const response = await request.post('/api/v1/demo/roundtrip');

  expect(response.status(), await response.text()).toBe(200);

  const body = await response.json();
  expect(body).toEqual(
    expect.objectContaining({
      id: expect.any(Number),
      value: expect.any(String)
    })
  );
});