const BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

export async function login(email, password) {
  const response = await fetch(`${BASE_URL}/api/v1/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });

  if (!response.ok) {
    throw new Error('invalid credentials');
  }

  return response.json();
}
