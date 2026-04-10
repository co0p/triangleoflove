'use strict';

const SEED_URL = process.env.SEED_URL || 'http://localhost:4000';

/**
 * Truncates all tables and resets to the baseline seed data.
 * Call in beforeEach / beforeAll for tests that require a clean slate.
 *
 * @param {import('@playwright/test').APIRequestContext} request
 */
async function resetDb(request) {
  const response = await request.post(`${SEED_URL}/seed`);
  if (!response.ok()) {
    const body = await response.text();
    throw new Error(`Seed endpoint returned ${response.status()}: ${body}`);
  }
}

module.exports = { resetDb, SEED_URL };
