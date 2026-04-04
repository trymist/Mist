import type { User } from '@/types';

const API_BASE = '/api';

export interface LoginRequest {
  email: string;
  password: string;
}

export interface SignUpRequest {
  email: string;
  password: string;
  name: string;
}

export const authService = {
  /**
   * Login user
   */
  async login(credentials: LoginRequest): Promise<User> {
    const response = await fetch(`${API_BASE}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(credentials),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Login failed');
    }

    return data.data;
  },

  /**
   * Sign up new user
   */
  async signUp(userData: SignUpRequest): Promise<User> {
    const response = await fetch(`${API_BASE}/auth/signup`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(userData),
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Sign up failed');
    }

    return data.data;
  },

  /**
   * Get current user
   */
  async getCurrentUser(): Promise<User> {
    const response = await fetch(`${API_BASE}/auth/me`, {
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Failed to get user');
    }

    return data.data;
  },

  /**
   * Logout user
   */
  async logout(): Promise<void> {
    const response = await fetch(`${API_BASE}/auth/logout`, {
      method: 'POST',
      credentials: 'include',
    });

    const data = await response.json();
    if (!data.success) {
      throw new Error(data.error || 'Logout failed');
    }
  },
};
