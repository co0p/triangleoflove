const BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

export async function getMe() {
  const token = localStorage.getItem('token');
  console.log('getMe: fetching from', `${BASE_URL}/api/v1/users/me`, 'with token:', token?.substring(0, 20) + '...');
  const response = await fetch(`${BASE_URL}/api/v1/users/me`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  console.log('getMe: response status:', response.status);
  if (response.status === 401) {
    throw new Error('unauthorized');
  }
  if (!response.ok) {
    const errorText = await response.text();
    console.error('getMe: error response:', errorText);
    throw new Error(`failed to load profile: ${response.status} ${errorText}`);
  }
  const json = await response.json();
  console.log('getMe: response json:', json);
  return json;
}
