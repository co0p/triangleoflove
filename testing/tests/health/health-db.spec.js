const { test, expect } = require('@playwright/test');
const net = require('net');

function canConnectToDb(host, port, timeoutMs) {
  return new Promise((resolve) => {
    const socket = net.createConnection({ host, port: Number(port) });

    const cleanup = () => {
      socket.removeAllListeners();
      socket.destroy();
    };

    socket.setTimeout(timeoutMs);

    socket.on('connect', () => {
      cleanup();
      resolve(true);
    });

    socket.on('timeout', () => {
      cleanup();
      resolve(false);
    });

    socket.on('error', () => {
      cleanup();
      resolve(false);
    });
  });
}

test('TestHealth_GivenDBRunning_WhenTCPConnected_ThenPortReachable', async () => {
  const dbHost = process.env.DB_HOST || 'db';
  const dbPort = process.env.DB_PORT || '5432';
  const connected = await canConnectToDb(dbHost, dbPort, 1000);
  expect(connected).toBeTruthy();
});

test('db port is reachable from test container', async () => {
  const dbHost = process.env.DB_HOST || 'db';
  const dbPort = process.env.DB_PORT || '5432';

  let connected = false;
  for (let attempt = 0; attempt < 10; attempt += 1) {
    connected = await canConnectToDb(dbHost, dbPort, 1000);
    if (connected) {
      break;
    }
    await new Promise((resolve) => setTimeout(resolve, 300));
  }

  expect(connected).toBeTruthy();
});
