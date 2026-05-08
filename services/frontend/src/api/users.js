const BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

export async function getMe() {
  const token = localStorage.getItem('token');
  const response = await fetch(`${BASE_URL}/api/v1/users/me`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (response.status === 401) {
    throw new Error('unauthorized');
  }
  if (!response.ok) {
    throw new Error('failed to load profile');
  }
  return response.json();
}
