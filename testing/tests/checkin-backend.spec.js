const { test, expect } = require('@playwright/test');
const { getToken } = require('./helpers/auth');

const CHECKIN_PAYLOAD = {
  felt_close: 7,
  positive_energy: 8,
  supported: 6,
  communication_healthy: 9,
  stress_level: 4,
  note: ''
};

test('TestCheckin_GivenNoAuth_WhenGetTodayRequested_ThenReturns401', async ({ request }) => {
  const response = await request.get('/api/v1/checkins/today');
  expect(response.status()).toBe(401);
});

test('TestCheckin_GivenNoAuth_WhenPutTodayRequested_ThenReturns401', async ({ request }) => {
  const response = await request.put('/api/v1/checkins/today', { data: CHECKIN_PAYLOAD });
  expect(response.status()).toBe(401);
});

test('TestCheckin_GivenValidPayload_WhenPutTodayRequested_ThenReturns200WithSavedBody', async ({ request }) => {
  const token = await getToken(request);
  const response = await request.put('/api/v1/checkins/today', {
    headers: { Authorization: `Bearer ${token}` },
    data: CHECKIN_PAYLOAD
  });
  expect(response.status()).toBe(200);
  const body = await response.json();
  expect(body).toMatchObject(CHECKIN_PAYLOAD);
});

test('TestCheckin_GivenEmptyNote_WhenSaved_ThenSucceeds', async ({ request }) => {
  const token = await getToken(request);
  const response = await request.put('/api/v1/checkins/today', {
    headers: { Authorization: `Bearer ${token}` },
    data: { ...CHECKIN_PAYLOAD, note: '' }
  });
  expect(response.status()).toBe(200);
  const body = await response.json();
  expect(body.note).toBe('');
});

test('TestCheckin_GivenExistingEntry_WhenGetTodayRequested_ThenReturns200WithValues', async ({ request }) => {
  const token = await getToken(request);
  await request.put('/api/v1/checkins/today', {
    headers: { Authorization: `Bearer ${token}` },
    data: CHECKIN_PAYLOAD
  });
  const response = await request.get('/api/v1/checkins/today', {
    headers: { Authorization: `Bearer ${token}` }
  });
  expect(response.status()).toBe(200);
  const body = await response.json();
  expect(body).toMatchObject(CHECKIN_PAYLOAD);
});
