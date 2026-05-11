const BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

export async function getWeeklyInsights(past) {
  const token = localStorage.getItem('token');
  const response = await fetch(`${BASE_URL}/api/v1/insights?past=${past}`, {
    headers: { Authorization: `Bearer ${token}` }
  });

  if (response.status === 401) {
    throw new Error('unauthorized');
  }
  if (!response.ok) {
    throw new Error('failed to load weekly insights');
  }

  return response.json();
}

export async function getInsights(date) {
  const token = localStorage.getItem('token');
  const response = await fetch(`${BASE_URL}/api/v1/insights/${date}`, {
    headers: { Authorization: `Bearer ${token}` }
  });

  if (response.status === 401) {
    throw new Error('unauthorized');
  }
  if (response.status === 404) {
    throw new Error('not_found');
  }
  if (response.status === 400) {
    throw new Error('invalid_date');
  }
  if (!response.ok) {
    throw new Error('failed to load insights');
  }

  return response.json();
}