const SEED_EMAIL = 'river@triangleoflove.app';
const SEED_PASSWORD = 'lovecoach1';

const FRONTEND_BASE_URL = process.env.FRONTEND_BASE_URL || 'http://frontend:5173';

async function getToken(request) {
  const response = await request.post('/api/v1/auth/login', {
    data: { email: SEED_EMAIL, password: SEED_PASSWORD }
  });
  const { token } = await response.json();
  return token;
}

async function loginViaUI(page) {
  await page.goto(`${FRONTEND_BASE_URL}/login`);
  await page.fill('input[type="email"]', SEED_EMAIL);
  await page.fill('input[type="password"]', SEED_PASSWORD);
  await page.click('button[type="submit"]');
}

module.exports = { SEED_EMAIL, SEED_PASSWORD, FRONTEND_BASE_URL, getToken, loginViaUI };
