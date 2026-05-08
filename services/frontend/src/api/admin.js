const BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

function authHeaders() {
  const token = localStorage.getItem('token');
  return { Authorization: `Bearer ${token}` };
}

export async function getAdminUsers() {
  const response = await fetch(`${BASE_URL}/api/v1/admin/users`, {
    headers: authHeaders(),
  });
  if (!response.ok) {
    throw new Error('failed to fetch users');
  }
  return response.json();
}

export async function activateUser(id) {
  const response = await fetch(`${BASE_URL}/api/v1/admin/users/${id}/activate`, {
    method: 'PUT',
    headers: authHeaders(),
  });
  if (!response.ok) {
    throw new Error('failed to activate user');
  }
}

export async function deactivateUser(id) {
  const response = await fetch(`${BASE_URL}/api/v1/admin/users/${id}/deactivate`, {
    method: 'PUT',
    headers: authHeaders(),
  });
  if (!response.ok) {
    throw new Error('failed to deactivate user');
  }
}
