const BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

export async function getTodayCheckin() {
  const token = localStorage.getItem('token');
  const response = await fetch(`${BASE_URL}/api/v1/checkins/today`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (response.status === 401) throw new Error('unauthorized');
  if (response.status === 404) return null;
  if (!response.ok) throw new Error('failed to load check-in');
  return response.json();
}

export async function saveTodayCheckin(data) {
  const token = localStorage.getItem('token');
  const response = await fetch(`${BASE_URL}/api/v1/checkins/today`, {
    method: 'PUT',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  });
  if (!response.ok) throw new Error('failed to save check-in');
  return response.json();
}
