import { APIResponse } from '../types';

const API_BASE = '/api/v1';

export async function fetchApi<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<APIResponse<T>> {
  const token = localStorage.getItem('access_token');

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(options.headers || {}),
  };

  if (token) {
    (headers as any)['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
  });

  const rawText = await response.text();
  let json: APIResponse<T>;

  try {
    json = JSON.parse(rawText);
  } catch (parseErr) {
    if (!response.ok) {
      throw new Error(`Server returned HTTP ${response.status} (${response.statusText || 'Error'}). API connection failed.`);
    }
    throw new Error('Respon dari server bukan format JSON yang valid.');
  }

  if (!response.ok || !json.success) {
    throw new Error(json.message || 'Terjadi kesalahan saat berkomunikasi dengan server');
  }

  return json;
}
