const BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

export async function getPairing() {
  const token = localStorage.getItem('token');
  const response = await fetch(`${BASE_URL}/api/v1/pairing`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!response.ok) throw new Error('failed to load pairing');
  return response.json();
}

export async function getCoupleStatus() {
  const token = localStorage.getItem('token');
  const response = await fetch(`${BASE_URL}/api/v1/couples/me`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!response.ok) throw new Error('failed to load couple status');
  return response.json();
}

export async function regeneratePairing() {
  const token = localStorage.getItem('token');
  const response = await fetch(`${BASE_URL}/api/v1/pairing/regenerate`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!response.ok) throw new Error('failed to regenerate code');
  return response.json();
}

export async function unpairCouple() {
  const token = localStorage.getItem('token');
  const response = await fetch(`${BASE_URL}/api/v1/couples/me`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${token}` }
  });
  if (!response.ok) throw new Error('failed to unpair');
  return response.json();
}

export async function connectPairing(inviteCode) {
  const token = localStorage.getItem('token');
  const response = await fetch(`${BASE_URL}/api/v1/pairing/connect`, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ invite_code: inviteCode })
  });
  if (response.status === 422) throw Object.assign(new Error('invalid invite code'), { code: 'INVALID_CODE' });
  if (response.status === 409) throw Object.assign(new Error('already paired'), { code: 'ALREADY_PAIRED' });
  if (!response.ok) throw new Error('failed to connect');
  return response.json();
}
