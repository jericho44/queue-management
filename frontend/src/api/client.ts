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

  const json: APIResponse<T> = await response.json();

  if (!response.ok || !json.success) {
    throw new Error(json.message || 'An error occurred while communicating with the server');
  }

  return json;
}
