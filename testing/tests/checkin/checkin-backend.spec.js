const { test, expect } = require('@playwright/test');
const { getToken } = require('../helpers/auth');

const CHECKIN_PAYLOAD = {
  felt_close: 2,
  positive_energy: 3,
  supported: 1,
  communication_healthy: 4,
  stress_level: -1,
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

test('TestCheckin_GivenNegativeSliderValue_WhenSubmittedAndReloaded_ThenValuePreserved', async ({ request }) => {
  const token = await getToken(request);

  const payload = {
    felt_close: -3,
    positive_energy: -5,
    supported: null,
    communication_healthy: null,
    stress_level: null,
    note: ''
  };

  const putResponse = await request.put('/api/v1/checkins/today', {
    headers: { Authorization: `Bearer ${token}` },
    data: payload
  });
  expect(putResponse.status()).toBe(200);

  const getResponse = await request.get('/api/v1/checkins/today', {
    headers: { Authorization: `Bearer ${token}` }
  });
  expect(getResponse.status()).toBe(200);
  const body = await getResponse.json();

  expect(body.felt_close).toBe(-3);
  expect(body.positive_energy).toBe(-5);
  expect(body.supported).toBeNull();
});

test('TestCheckin_GivenUnsetSliders_WhenSubmittedAndReloaded_ThenUnsetsRemainDistinct', async ({ request }) => {
  const token = await getToken(request);

  // felt_close explicitly set to 0; all other rating fields left unset (null)
  const payload = {
    felt_close: 0,
    positive_energy: null,
    supported: null,
    communication_healthy: null,
    stress_level: null,
    note: ''
  };

  const putResponse = await request.put('/api/v1/checkins/today', {
    headers: { Authorization: `Bearer ${token}` },
    data: payload
  });
  expect(putResponse.status()).toBe(200);

  const getResponse = await request.get('/api/v1/checkins/today', {
    headers: { Authorization: `Bearer ${token}` }
  });
  expect(getResponse.status()).toBe(200);
  const body = await getResponse.json();

  expect(body.felt_close).toBe(0);              // explicitly set to 0 → not null
  expect(body.positive_energy).toBeNull();       // unset → null
  expect(body.supported).toBeNull();             // unset → null
  expect(body.communication_healthy).toBeNull(); // unset → null
  expect(body.stress_level).toBeNull();          // unset → null
});
