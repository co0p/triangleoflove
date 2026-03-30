const BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

export async function getMe() {
  const token = localStorage.getItem('token');
  const response = await fetch(`${BASE_URL}/api/v1/users/me`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!response.ok) {
    throw new Error('unauthorized');
  }
  return response.json();
}
