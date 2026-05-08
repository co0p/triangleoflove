const BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

export async function register(email, password, firstName) {
  const response = await fetch(`${BASE_URL}/api/v1/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password, firstName })
  });

  if (response.status === 409) {
    throw new Error('duplicate email');
  }
  if (!response.ok) {
    throw new Error('registration failed');
  }
}

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
