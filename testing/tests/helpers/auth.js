const { USERS } = require('./users');

const SEED_EMAIL = USERS[0].email;
const SEED_PASSWORD = USERS[0].password;

const SEED_EMAIL_2 = USERS[1].email;
const SEED_PASSWORD_2 = USERS[1].password;

const FRONTEND_BASE_URL = process.env.FRONTEND_BASE_URL || 'http://frontend:5173';

async function getToken(request) {
  const response = await request.post('/api/v1/auth/login', {
    data: { email: SEED_EMAIL, password: SEED_PASSWORD }
  });
  const { token } = await response.json();
  return token;
}

async function loginViaUI(page, email, password) {
  const loginEmail = email || SEED_EMAIL;
  const loginPassword = password || SEED_PASSWORD;
  await page.goto(`${FRONTEND_BASE_URL}/login`);
  await page.fill('input[type="email"]', loginEmail);
  await page.fill('input[type="password"]', loginPassword);
  await Promise.all([
    page.waitForURL('**/dashboard'),
    page.click('button[type="submit"]')
  ]);
}

module.exports = { SEED_EMAIL, SEED_PASSWORD, SEED_EMAIL_2, SEED_PASSWORD_2, USERS, FRONTEND_BASE_URL, getToken, loginViaUI };

