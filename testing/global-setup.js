'use strict';

const { createServer, runSeed } = require('./seed-server');

const PORT = parseInt(process.env.SEED_SERVER_PORT || '4000', 10);

module.exports = async function globalSetup() {
  const server = createServer();

  await new Promise((resolve, reject) => {
    server.once('error', reject);
    server.listen(PORT, resolve);
  });

  console.log(`[seed-server] Started on port ${PORT}`);

  // Seed once before any spec runs so all specs start from a known baseline.
  await runSeed();
  console.log('[seed-server] Initial seed complete');

  // Playwright (>=1.41) calls the returned function as globalTeardown.
  return async function globalTeardown() {
    await new Promise((resolve) => server.close(resolve));
    console.log('[seed-server] Stopped');
  };
};
