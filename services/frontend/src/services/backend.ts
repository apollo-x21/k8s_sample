export interface AuthPayload {
  username: string;
  password: string;
}

export interface AuthResponse {
  token?: string;
  message: string;
}

export interface UsersResponse {
  users: { username: string }[];
  me: string;
}

export interface MessageResponse {
  message: string;
}

const API_BASE = '/api';

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(options.headers || {}),
    },
    ...options,
  });

  if (!response.ok) {
    const text = await response.text();
    throw new Error(text || '请求失败');
  }

  return response.json() as Promise<T>;
}

export function register(payload: AuthPayload) {
  return request<AuthResponse>('/register', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export function login(payload: AuthPayload) {
  return request<AuthResponse>('/login', {
    method: 'POST',
    body: JSON.stringify(payload),
  });
}

export function getUsers(token: string) {
  return request<UsersResponse>('/users', {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
}

export function logout(token: string) {
  return request<MessageResponse>('/logout', {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
}
