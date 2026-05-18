const BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

export async function register(email, password, firstName) {
  const response = await fetch(`${BASE_URL}/api/v1/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password, firstName })
  });

  const data = await response.json().catch(() => null);
  const errorMessage = typeof data?.error === 'string' ? data.error : 'Registration failed. Please try again.';

  if (response.status === 409) {
    throw new Error('duplicate email');
  }
  if (!response.ok) {
    throw new Error(errorMessage);
  }
}

export async function login(email, password) {
  const response = await fetch(`${BASE_URL}/api/v1/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });

  console.log('login: response status:', response.status);
  if (!response.ok) {
    throw new Error('invalid credentials');
  }

  const result = await response.json();
  console.log('login: response json:', result);
  return result;
}

export async function changePassword(currentPassword, newPassword) {
  const token = localStorage.getItem('token');
  const response = await fetch(`${BASE_URL}/api/v1/auth/password`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`
    },
    body: JSON.stringify({ current_password: currentPassword, new_password: newPassword })
  });

  if (response.status === 409) {
    throw new Error('current password is incorrect');
  }
  if (!response.ok) {
    throw new Error('password change failed');
  }
}
