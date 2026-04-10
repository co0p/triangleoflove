'use strict';

const http = require('http');
const fs = require('fs');
const path = require('path');
const { Client } = require('pg');

const PORT = parseInt(process.env.SEED_SERVER_PORT || '4000', 10);
const DATABASE_URL = process.env.DATABASE_URL;

// Path to the canonical seed file shared with the db service.
const INIT_SQL_PATH = path.resolve(__dirname, 'services/db/init.sql');

async function runSeed() {
  if (!DATABASE_URL) {
    throw new Error('DATABASE_URL environment variable is not set');
  }
  const initSql = fs.readFileSync(INIT_SQL_PATH, 'utf8');
  const client = new Client({ connectionString: DATABASE_URL });
  await client.connect();
  try {
    // Drop in FK-safe order so CREATE TABLE in init.sql runs cleanly.
    await client.query('DROP TABLE IF EXISTS checkins, couples, accounts CASCADE');
    await client.query(initSql);
  } finally {
    await client.end();
  }
}

function createServer() {
  return http.createServer(async (req, res) => {
    if (req.method === 'GET' && req.url === '/health') {
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ status: 'ok' }));
      return;
    }

    if (req.method === 'POST' && req.url === '/seed') {
      try {
        await runSeed();
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ status: 'seeded' }));
      } catch (err) {
        console.error('[seed-server] Seed failed:', err.message);
        res.writeHead(500, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ status: 'error', message: err.message }));
      }
      return;
    }

    res.writeHead(404);
    res.end();
  });
}

module.exports = { createServer, runSeed };

// Allow running standalone: `node seed-server.js`
if (require.main === module) {
  const server = createServer();
  server.listen(PORT, () => {
    console.log(`[seed-server] Listening on port ${PORT}`);
  });
}
